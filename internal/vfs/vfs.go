// Package vfs implements a SQLite-backed virtual filesystem with
// content-addressed blob storage, adapted from Cloudflare Computer's
// dofs design. Each engagement gets its own in-process VFS for
// managing pentest artifacts (loot, evidence, notes, exploit code)
// with dedup, versioning, and sync support.
package vfs

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const chunkSize = 512 * 1024 // 512 KiB chunks, matching CF Computer's dofs

// NodeType enumerates filesystem entry kinds.
type NodeType string

const (
	NodeFile    NodeType = "file"
	NodeDir     NodeType = "dir"
	NodeSymlink NodeType = "symlink"
)

// VFS is a virtual filesystem backed by SQLite with content-addressed blobs.
type VFS struct {
	db   *sql.DB
	mu   sync.Mutex
	path string
}

// Entry represents a filesystem entry returned by List/Stat.
type Entry struct {
	Inode   int64     `json:"inode"`
	Name    string    `json:"name"`
	Type    NodeType  `json:"type"`
	Mode    uint32    `json:"mode"`
	Size    int64     `json:"size"`
	Mtime   time.Time `json:"mtime"`
	Rev     int64     `json:"rev"`
	Target  string    `json:"target,omitempty"` // for symlinks
	Content string    `json:"content,omitempty"` // for small files on read
}

// Open creates or opens a VFS at the given SQLite path.
func Open(dbPath string) (*VFS, error) {
	db, err := sql.Open("sqlite", dbPath+"?_journal=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("vfs: open %s: %w", dbPath, err)
	}
	v := &VFS{db: db, path: dbPath}
	if err := v.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return v, nil
}

// Close closes the underlying database.
func (v *VFS) Close() error {
	return v.db.Close()
}

func (v *VFS) migrate() error {
	_, err := v.db.Exec(schema)
	return err
}

const schema = `
CREATE TABLE IF NOT EXISTS vfs_nodes (
	inode INTEGER PRIMARY KEY AUTOINCREMENT,
	type  TEXT NOT NULL CHECK(type IN ('file','dir','symlink')),
	mode  INTEGER NOT NULL DEFAULT 493,
	mtime INTEGER NOT NULL,
	rev   INTEGER NOT NULL DEFAULT 0,
	mount_root TEXT,
	manifest_hash BLOB,
	link_target TEXT,
	size INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS vfs_dirents (
	parent_inode INTEGER NOT NULL,
	name TEXT NOT NULL,
	child_inode INTEGER NOT NULL,
	PRIMARY KEY (parent_inode, name)
) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS vfs_blobs (
	hash TEXT PRIMARY KEY,
	size INTEGER NOT NULL,
	last_seen INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS vfs_blob_bytes (
	hash TEXT PRIMARY KEY REFERENCES vfs_blobs(hash) ON DELETE CASCADE,
	bytes BLOB NOT NULL
);
CREATE TABLE IF NOT EXISTS vfs_chunks (
	inode INTEGER NOT NULL,
	idx INTEGER NOT NULL,
	hash TEXT NOT NULL,
	size INTEGER NOT NULL,
	PRIMARY KEY (inode, idx)
) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS vfs_rev_counter (
	id INTEGER PRIMARY KEY CHECK(id = 1),
	rev INTEGER NOT NULL DEFAULT 0
);
INSERT OR IGNORE INTO vfs_rev_counter (id, rev) VALUES (1, 0);
CREATE INDEX IF NOT EXISTS idx_vfs_nodes_rev ON vfs_nodes(rev);
CREATE INDEX IF NOT EXISTS idx_vfs_chunks_hash ON vfs_chunks(hash);
`

// nextRev bumps the global revision counter and returns the new value.
func (v *VFS) nextRev() (int64, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	var rev int64
	err := v.db.QueryRow("UPDATE vfs_rev_counter SET rev = rev + 1 WHERE id = 1 RETURNING rev").Scan(&rev)
	return rev, err
}

// currentRev returns the latest revision without incrementing.
func (v *VFS) currentRev() (int64, error) {
	var rev int64
	err := v.db.QueryRow("SELECT rev FROM vfs_rev_counter WHERE id = 1").Scan(&rev)
	return rev, err
}

// Mkdir creates a directory at the given path.
func (v *VFS) Mkdir(p string) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.mkdirLocked(p)
}

func (v *VFS) mkdirLocked(p string) error {
	clean := normalize(p)
	if clean == "/" || clean == "" {
		return nil // root always exists (inode 1)
	}
	parent, name := splitPath(clean)
	parentInode, err := v.lookupInode(parent)
	if err != nil {
		return err
	}
	if parentInode == 0 {
		// Create parent dirs recursively
		if err := v.mkdirLocked(parent); err != nil {
			return err
		}
		parentInode, err = v.lookupInode(parent)
		if err != nil {
			return err
		}
	}
	// Check if already exists
	existing, _ := v.lookupInode(clean)
	if existing != 0 {
		return nil // idempotent
	}
	rev, err := v.nextRev()
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	res, err := v.db.Exec(
		"INSERT INTO vfs_nodes (type, mode, mtime, rev) VALUES ('dir', ?, ?, ?)",
		os.ModeDir|0o755, now, rev)
	if err != nil {
		return fmt.Errorf("vfs: mkdir %s: %w", p, err)
	}
	inode, _ := res.LastInsertId()
	_, err = v.db.Exec(
		"INSERT INTO vfs_dirents (parent_inode, name, child_inode) VALUES (?, ?, ?)",
		parentInode, name, inode)
	return err
}

// WriteFile writes content to a file, chunking into content-addressed blobs.
func (v *VFS) WriteFile(p string, content []byte) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.writeFileLocked(p, content)
}

func (v *VFS) writeFileLocked(p string, content []byte) error {
	clean := normalize(p)
	parent, name := splitPath(clean)

	// Ensure parent dir exists
	parentInode, err := v.lookupInode(parent)
	if err != nil {
		return err
	}
	if parentInode == 0 {
		if err := v.mkdirLocked(parent); err != nil {
			return err
		}
		parentInode, err = v.lookupInode(parent)
		if err != nil {
			return err
		}
	}

	// Check if file already exists
	existingInode, _ := v.lookupInode(clean)
	now := time.Now().Unix()

	// Chunk the content and store blobs
	chunks := chunkContent(content)
	manifestHash := computeManifestHash(chunks)

	// Store blobs
	for _, chunk := range chunks {
		v.storeBlob(chunk.hash, chunk.data)
	}

	rev, err := v.nextRev()
	if err != nil {
		return err
	}

	if existingInode != 0 {
		// Update existing file: clear old chunks
		_, err = v.db.Exec("DELETE FROM vfs_chunks WHERE inode = ?", existingInode)
		if err != nil {
			return err
		}
		// Insert new chunks
		for i, chunk := range chunks {
			_, err = v.db.Exec(
				"INSERT INTO vfs_chunks (inode, idx, hash, size) VALUES (?, ?, ?, ?)",
				existingInode, i, chunk.hash, len(chunk.data))
			if err != nil {
				return err
			}
		}
		_, err = v.db.Exec(
			"UPDATE vfs_nodes SET mode = ?, mtime = ?, rev = ?, manifest_hash = ?, size = ? WHERE inode = ?",
			os.FileMode(0o644), now, rev, manifestHash, len(content), existingInode)
		return err
	}

	// Create new file node
	res, err := v.db.Exec(
		"INSERT INTO vfs_nodes (type, mode, mtime, rev, manifest_hash, size) VALUES ('file', ?, ?, ?, ?, ?)",
		os.FileMode(0o644), now, rev, manifestHash, len(content))
	if err != nil {
		return err
	}
	inode, _ := res.LastInsertId()
	for i, chunk := range chunks {
		_, err = v.db.Exec(
			"INSERT INTO vfs_chunks (inode, idx, hash, size) VALUES (?, ?, ?, ?)",
			inode, i, chunk.hash, len(chunk.data))
		if err != nil {
			return err
		}
	}
	_, err = v.db.Exec(
		"INSERT INTO vfs_dirents (parent_inode, name, child_inode) VALUES (?, ?, ?)",
		parentInode, name, inode)
	return err
}

// ReadFile reads the entire content of a file.
func (v *VFS) ReadFile(p string) ([]byte, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	inode, err := v.lookupInode(p)
	if err != nil {
		return nil, err
	}
	if inode == 0 {
		return nil, fmt.Errorf("vfs: %s: not found", p)
	}

	rows, err := v.db.Query(
		"SELECT hash, size FROM vfs_chunks WHERE inode = ? ORDER BY idx", inode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buf []byte
	for rows.Next() {
		var hash string
		var size int
		if err := rows.Scan(&hash, &size); err != nil {
			return nil, err
		}
		var data []byte
		err := v.db.QueryRow("SELECT bytes FROM vfs_blob_bytes WHERE hash = ?", hash).Scan(&data)
		if err != nil {
			return nil, fmt.Errorf("vfs: missing blob %s: %w", hash, err)
		}
		buf = append(buf, data...)
	}
	return buf, rows.Err()
}

// Stat returns metadata for a path.
func (v *VFS) Stat(p string) (*Entry, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	inode, err := v.lookupInode(p)
	if err != nil {
		return nil, err
	}
	if inode == 0 {
		return nil, fmt.Errorf("vfs: %s: not found", p)
	}
	return v.statInode(inode, path.Base(p))
}

// List lists entries in a directory.
func (v *VFS) List(dir string) ([]Entry, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	clean := normalize(dir)
	inode, err := v.lookupInode(clean)
	if err != nil {
		return nil, err
	}
	if inode == 0 {
		return nil, fmt.Errorf("vfs: %s: not found", dir)
	}

	rows, err := v.db.Query(
		`SELECT d.name, n.inode, n.type, n.mode, n.size, n.mtime, n.rev, n.link_target
		 FROM vfs_dirents d
		 JOIN vfs_nodes n ON d.child_inode = n.inode
		 WHERE d.parent_inode = ?
		 ORDER BY d.name`, inode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		var mtime int64
		var linkTarget sql.NullString
		if err := rows.Scan(&e.Name, &e.Inode, &e.Type, &e.Mode, &e.Size, &mtime, &e.Rev, &linkTarget); err != nil {
			return nil, err
		}
		e.Mtime = time.Unix(mtime, 0)
		if linkTarget.Valid {
			e.Target = linkTarget.String
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

// Delete removes a file or directory (recursively).
func (v *VFS) Delete(p string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	clean := normalize(p)
	inode, err := v.lookupInode(clean)
	if err != nil {
		return err
	}
	if inode == 0 {
		return fmt.Errorf("vfs: %s: not found", p)
	}
	return v.deleteInode(inode, clean)
}

func (v *VFS) deleteInode(inode int64, p string) error {
	// Check if directory - delete children first
	rows, err := v.db.Query("SELECT child_inode FROM vfs_dirents WHERE parent_inode = ?", inode)
	if err != nil {
		return err
	}
	var childInodes []int64
	for rows.Next() {
		var child int64
		rows.Scan(&child)
		childInodes = append(childInodes, child)
	}
	rows.Close()

	for _, child := range childInodes {
		v.deleteInode(child, path.Join(p, "?"))
	}

	_, err = v.db.Exec("DELETE FROM vfs_chunks WHERE inode = ?", inode)
	if err != nil {
		return err
	}
	_, err = v.db.Exec("DELETE FROM vfs_dirents WHERE child_inode = ?", inode)
	if err != nil {
		return err
	}
	_, err = v.db.Exec("DELETE FROM vfs_nodes WHERE inode = ?", inode)
	return err
}

// Symlink creates a symbolic link.
func (v *VFS) Symlink(target, link string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	clean := normalize(link)
	parent, name := splitPath(clean)
	parentInode, err := v.lookupInode(parent)
	if err != nil {
		return err
	}
	if parentInode == 0 {
		if err := v.mkdirLocked(parent); err != nil {
			return err
		}
		parentInode, _ = v.lookupInode(parent)
	}

	rev, err := v.nextRev()
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	res, err := v.db.Exec(
		"INSERT INTO vfs_nodes (type, mode, mtime, rev, link_target) VALUES ('symlink', ?, ?, ?, ?)",
		os.ModeSymlink|0o777, now, rev, target)
	if err != nil {
		return err
	}
	inode, _ := res.LastInsertId()
	_, err = v.db.Exec(
		"INSERT INTO vfs_dirents (parent_inode, name, child_inode) VALUES (?, ?, ?)",
		parentInode, name, inode)
	return err
}

// Walk traverses the VFS calling fn for each entry. If fn returns false, walking stops.
func (v *VFS) Walk(root string, fn func(path string, entry Entry) bool) error {
	entries, err := v.List(root)
	if err != nil {
		return err
	}
	for _, e := range entries {
		full := path.Join(root, e.Name)
		if !fn(full, e) {
			return nil
		}
		if e.Type == NodeDir {
			if err := v.Walk(full, fn); err != nil {
				return err
			}
		}
	}
	return nil
}

// --- Sync Support (from CF Computer sync protocol) ---

// ChangeEntry represents a state-based sync entry.
type ChangeEntry struct {
	Kind  string `json:"kind"` // "file", "dir", "symlink", "delete"
	Rev   int64  `json:"rev"`
	Path  string `json:"path"`
	Mode  uint32 `json:"mode,omitempty"`
	Size  int64  `json:"size,omitempty"`
	Mtime int64  `json:"mtime,omitempty"`
	Hash  string `json:"hash,omitempty"` // manifest hash for dedup check
}

// ChangesSince returns all changes since the given revision.
func (v *VFS) ChangesSince(sinceRev int64) ([]ChangeEntry, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	rows, err := v.db.Query(
		`SELECT n.inode, n.type, n.mode, n.size, n.mtime, n.rev, n.link_target,
		        (SELECT group_concat(parent_inode || '/' || name, '/') FROM vfs_dirents WHERE child_inode = n.inode)
		 FROM vfs_nodes n WHERE n.rev > ?
		 ORDER BY n.rev`, sinceRev)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []ChangeEntry
	for rows.Next() {
		var inode int64
		var nodeType string
		var mode uint32
		var size int64
		var mtime int64
		var rev int64
		var linkTarget sql.NullString
		var pathStr sql.NullString
		if err := rows.Scan(&inode, &nodeType, &mode, &size, &mtime, &rev, &linkTarget, &pathStr); err != nil {
			return nil, err
		}
		// Resolve path from dirents
		p := v.resolvePath(inode)
		ce := ChangeEntry{
			Kind:  nodeType,
			Rev:   rev,
			Path:  p,
			Mode:  mode,
			Size:  size,
			Mtime: mtime,
		}
		if linkTarget.Valid && linkTarget.String != "" {
			ce.Kind = "symlink"
		}
		entries = append(entries, ce)
	}
	return entries, rows.Err()
}

// --- Internal helpers ---

type chunk struct {
	hash string
	data []byte
}

func chunkContent(content []byte) []chunk {
	var chunks []chunk
	for i := 0; i < len(content); i += chunkSize {
		end := i + chunkSize
		if end > len(content) {
			end = len(content)
		}
		data := content[i:end]
		h := sha256.Sum256(data)
		chunks = append(chunks, chunk{
			hash: hex.EncodeToString(h[:]),
			data: data,
		})
	}
	if len(chunks) == 0 {
		// Empty file gets a single empty chunk
		h := sha256.Sum256(nil)
		chunks = append(chunks, chunk{
			hash: hex.EncodeToString(h[:]),
			data: nil,
		})
	}
	return chunks
}

func computeManifestHash(chunks []chunk) []byte {
	h := sha256.New()
	for _, c := range chunks {
		h.Write([]byte(c.hash))
	}
	return h.Sum(nil)
}

func (v *VFS) storeBlob(hash string, data []byte) error {
	var exists bool
	v.db.QueryRow("SELECT 1 FROM vfs_blobs WHERE hash = ?", hash).Scan(&exists)
	if exists {
		// Update last_seen
		_, err := v.db.Exec("UPDATE vfs_blobs SET last_seen = ? WHERE hash = ?", time.Now().Unix(), hash)
		return err
	}
	_, err := v.db.Exec("INSERT INTO vfs_blobs (hash, size, last_seen) VALUES (?, ?, ?)",
		hash, len(data), time.Now().Unix())
	if err != nil {
		return err
	}
	_, err = v.db.Exec("INSERT INTO vfs_blob_bytes (hash, bytes) VALUES (?, ?)", hash, data)
	return err
}

func (v *VFS) lookupInode(p string) (int64, error) {
	clean := normalize(p)
	if clean == "/" || clean == "" {
		return 1, nil // root inode is always 1
	}

	// Ensure root exists
	v.ensureRoot()

	parts := strings.Split(strings.Trim(clean, "/"), "/")
	currentInode := int64(1)

	for _, part := range parts {
		var childInode int64
		err := v.db.QueryRow(
			"SELECT child_inode FROM vfs_dirents WHERE parent_inode = ? AND name = ?",
			currentInode, part).Scan(&childInode)
		if err == sql.ErrNoRows {
			return 0, nil // not found, return 0 (not an error)
		}
		if err != nil {
			return 0, err
		}
		currentInode = childInode
	}
	return currentInode, nil
}

func (v *VFS) ensureRoot() {
	var count int
	v.db.QueryRow("SELECT COUNT(*) FROM vfs_nodes WHERE inode = 1").Scan(&count)
	if count == 0 {
		v.db.Exec("INSERT INTO vfs_nodes (inode, type, mode, mtime, rev) VALUES (1, 'dir', ?, ?, 0)",
			os.ModeDir|0o755, time.Now().Unix())
	}
}

func (v *VFS) resolvePath(inode int64) string {
	if inode == 1 {
		return "/"
	}
	var parentInode int64
	var name string
	err := v.db.QueryRow(
		"SELECT parent_inode, name FROM vfs_dirents WHERE child_inode = ?", inode).Scan(&parentInode, &name)
	if err != nil {
		return "?"
	}
	if parentInode == 1 {
		return "/" + name
	}
	return path.Join(v.resolvePath(parentInode), name)
}

func (v *VFS) statInode(inode int64, name string) (*Entry, error) {
	var entry Entry
	var mtime int64
	var nodeType string
	var linkTarget sql.NullString
	var manifestHash []byte
	err := v.db.QueryRow(
		"SELECT type, mode, size, mtime, rev, link_target, manifest_hash FROM vfs_nodes WHERE inode = ?",
		inode).Scan(&nodeType, &entry.Mode, &entry.Size, &mtime, &entry.Rev, &linkTarget, &manifestHash)
	if err != nil {
		return nil, err
	}
	entry.Inode = inode
	entry.Name = name
	entry.Type = NodeType(nodeType)
	entry.Mtime = time.Unix(mtime, 0)
	if linkTarget.Valid {
		entry.Target = linkTarget.String
	}
	return &entry, nil
}

// --- Path utilities ---

func normalize(p string) string {
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return path.Clean(p)
}

func splitPath(p string) (parent, name string) {
	clean := normalize(p)
	if clean == "/" {
		return "/", ""
	}
	idx := strings.LastIndex(clean, "/")
	if idx == 0 {
		return "/", clean[1:]
	}
	return clean[:idx], clean[idx+1:]
}

// ExportToWriter writes a file's content to the given writer.
func (v *VFS) ExportToWriter(p string, w io.Writer) error {
	data, err := v.ReadFile(p)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
