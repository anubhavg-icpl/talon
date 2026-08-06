// Package blueprint implements a pentest playbook template system. A
// Blueprint is a reusable, ordered sequence of tool invocations for a
// recurring pentest task (recon, exploitation, post-exploitation, reporting).
//
// This is adapted from Cloudflare's OS blueprints concept — declarative,
// shareable runbooks — reinterpreted for offensive-security engagements.
// Blueprints are stored in a SQL database so operators can browse, clone, and
// adapt them when spinning up new engagements.
package blueprint

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Category groups blueprints by the high-level phase of work they support.
type Category string

const (
	// CategoryRecon covers discovery and enumeration playbooks.
	CategoryRecon Category = "recon"
	// CategoryExploit covers active exploitation playbooks.
	CategoryExploit Category = "exploit"
	// CategoryPostExploit covers post-exploitation: privesc, lateral movement, persistence.
	CategoryPostExploit Category = "post-exploit"
	// CategoryReporting covers report generation and evidence packaging.
	CategoryReporting Category = "reporting"
)

// BlueprintStep is a single ordered action inside a blueprint. Each step names
// a tool, the arguments it should be invoked with, what a successful result
// looks like, and how to react when the tool fails or returns nothing useful.
type BlueprintStep struct {
	// Order is the 1-based execution position of this step within its blueprint.
	Order int `json:"order"`
	// Tool is the command or tool name to run (e.g. "nmap", "sqlmap", "ffuf").
	Tool string `json:"tool"`
	// Description is a human-readable explanation of what this step does and why.
	Description string `json:"description"`
	// Args holds named arguments passed to the tool. Values may contain
	// {{placeholders}} that the orchestrator substitutes at run time.
	Args map[string]string `json:"args,omitempty"`
	// ExpectedResult describes the observable outcome that indicates success,
	// so the agent can self-assess whether the step achieved its goal.
	ExpectedResult string `json:"expected_result,omitempty"`
	// OnFailure is the fallback strategy: a free-form instruction such as
	// "retry with -W 1", "skip", or "escalate to operator".
	OnFailure string `json:"on_failure,omitempty"`
}

// Blueprint is a reusable pentest playbook: a named, categorized, ordered set
// of steps that an operator or agent can execute against a target.
type Blueprint struct {
	// ID is the unique, stable identifier for the blueprint (UUID or slug).
	ID string `json:"id"`
	// Name is the human-readable title, e.g. "Web App Recon".
	Name string `json:"name"`
	// Description summarizes what the blueprint accomplishes and when to use it.
	Description string `json:"description"`
	// Category is one of recon/exploit/post-exploit/reporting.
	Category Category `json:"category"`
	// Phase is a finer-grained engagement-phase label (e.g. "reconnaissance",
	// "exploitation", "post-exploitation", "reporting").
	Phase string `json:"phase"`
	// Steps is the ordered list of tool actions that make up the playbook.
	Steps []BlueprintStep `json:"steps"`
	// Tags are free-form labels used for search and filtering
	// (e.g. "web", "windows", "cloud", "sqli").
	Tags []string `json:"tags,omitempty"`
	// Version is the blueprint schema/content version (e.g. "1.0").
	Version string `json:"version"`
	// Author identifies who created or last maintained the blueprint.
	Author string `json:"author"`
	// CreatedAt is the UTC creation timestamp.
	CreatedAt time.Time `json:"created_at"`
}

// BlueprintStore is a SQL-backed repository of pentest blueprints. It is safe
// for concurrent use: every mutating or read-multiple operation takes the store
// mutex. Steps and tags are stored as JSON text columns; everything else maps
// directly to scalar columns.
type BlueprintStore struct {
	db *sql.DB
	mu sync.Mutex
}

// NewBlueprintStore wraps an already-open *sql.DB. The caller is responsible
// for registering the SQL driver (e.g. sqlite3) and for the connection's
// lifetime. Call Migrate once before first use to create the schema.
func NewBlueprintStore(db *sql.DB) *BlueprintStore {
	return &BlueprintStore{db: db}
}

// blueprintSchema creates the blueprints table. Steps and tags are serialized
// as JSON text; category and phase are validated by CHECK constraints. Uses
// only portable SQLite-compatible DDL so it works with the same driver as the
// VFS layer.
const blueprintSchema = `
CREATE TABLE IF NOT EXISTS blueprints (
	id          TEXT PRIMARY KEY,
	name        TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	category    TEXT NOT NULL CHECK(category IN ('recon','exploit','post-exploit','reporting')),
	phase       TEXT NOT NULL DEFAULT '',
	steps_json  TEXT NOT NULL DEFAULT '[]',
	tags_json   TEXT NOT NULL DEFAULT '[]',
	version     TEXT NOT NULL DEFAULT '1.0',
	author      TEXT NOT NULL DEFAULT '',
	created_at  TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_blueprints_category ON blueprints(category);
`

// Migrate creates the blueprints table and indexes if they do not already
// exist. It is idempotent and safe to call on every startup.
func (s *BlueprintStore) Migrate() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(blueprintSchema)
	if err != nil {
		return fmt.Errorf("blueprint: migrate: %w", err)
	}
	return nil
}

// Create inserts a new blueprint. If ID is empty a UUID is generated; if
// CreatedAt is zero it is set to the current UTC time. Duplicate IDs fail at
// the database level (PRIMARY KEY violation).
func (s *BlueprintStore) Create(b Blueprint) error {
	if strings.TrimSpace(b.ID) == "" {
		b.ID = uuid.NewString()
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now().UTC()
	}
	if b.Steps == nil {
		b.Steps = []BlueprintStep{}
	}
	if b.Tags == nil {
		b.Tags = []string{}
	}

	stepsJSON, err := json.Marshal(b.Steps)
	if err != nil {
		return fmt.Errorf("blueprint: marshal steps: %w", err)
	}
	tagsJSON, err := json.Marshal(b.Tags)
	if err != nil {
		return fmt.Errorf("blueprint: marshal tags: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, err = s.db.Exec(
		`INSERT INTO blueprints
		    (id, name, description, category, phase, steps_json, tags_json, version, author, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		b.ID, b.Name, b.Description, string(b.Category), b.Phase,
		string(stepsJSON), string(tagsJSON), b.Version, b.Author,
		b.CreatedAt.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("blueprint: create %q: %w", b.ID, err)
	}
	return nil
}

// Get returns a single blueprint by ID. Returns an error wrapping sql.ErrNoRows
// when the blueprint does not exist.
func (s *BlueprintStore) Get(id string) (Blueprint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row := s.db.QueryRow(
		`SELECT id, name, description, category, phase, steps_json, tags_json, version, author, created_at
		   FROM blueprints WHERE id = ?`, id)
	b, err := scanBlueprint(row)
	if err != nil {
		return Blueprint{}, fmt.Errorf("blueprint: get %q: %w", id, err)
	}
	return b, nil
}

// List returns blueprints, optionally filtered by category. When categoryFilter
// is empty all blueprints are returned, ordered by name.
func (s *BlueprintStore) List(categoryFilter string) ([]Blueprint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		rows *sql.Rows
		err  error
	)
	categoryFilter = strings.TrimSpace(categoryFilter)
	if categoryFilter == "" {
		rows, err = s.db.Query(
			`SELECT id, name, description, category, phase, steps_json, tags_json, version, author, created_at
			   FROM blueprints ORDER BY name`)
	} else {
		rows, err = s.db.Query(
			`SELECT id, name, description, category, phase, steps_json, tags_json, version, author, created_at
			   FROM blueprints WHERE category = ? ORDER BY name`, categoryFilter)
	}
	if err != nil {
		return nil, fmt.Errorf("blueprint: list: %w", err)
	}
	defer rows.Close()

	out, err := scanBlueprints(rows)
	if err != nil {
		return nil, fmt.Errorf("blueprint: list: %w", err)
	}
	return out, nil
}

// Update replaces all mutable fields of an existing blueprint identified by ID.
// Returns an error if no blueprint with that ID exists.
func (s *BlueprintStore) Update(b Blueprint) error {
	if strings.TrimSpace(b.ID) == "" {
		return fmt.Errorf("blueprint: update: id required")
	}
	if b.Steps == nil {
		b.Steps = []BlueprintStep{}
	}
	if b.Tags == nil {
		b.Tags = []string{}
	}

	stepsJSON, err := json.Marshal(b.Steps)
	if err != nil {
		return fmt.Errorf("blueprint: marshal steps: %w", err)
	}
	tagsJSON, err := json.Marshal(b.Tags)
	if err != nil {
		return fmt.Errorf("blueprint: marshal tags: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	res, err := s.db.Exec(
		`UPDATE blueprints
		    SET name = ?, description = ?, category = ?, phase = ?, steps_json = ?,
		        tags_json = ?, version = ?, author = ?, created_at = ?
		  WHERE id = ?`,
		b.Name, b.Description, string(b.Category), b.Phase,
		string(stepsJSON), string(tagsJSON), b.Version, b.Author,
		b.CreatedAt.UTC().Format(time.RFC3339), b.ID,
	)
	if err != nil {
		return fmt.Errorf("blueprint: update %q: %w", b.ID, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("blueprint: update %q: %w", b.ID, sql.ErrNoRows)
	}
	return nil
}

// Delete removes a blueprint by ID. It is a no-op (returns nil) if the
// blueprint does not exist.
func (s *BlueprintStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(`DELETE FROM blueprints WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("blueprint: delete %q: %w", id, err)
	}
	return nil
}

// ListByTag returns every blueprint whose Tags slice contains tag. Because
// tags are stored as a JSON array, filtering is performed in Go after loading
// the matching rows; this avoids depending on a JSON1 extension and keeps the
// query portable across SQL drivers.
func (s *BlueprintStore) ListByTag(tag string) ([]Blueprint, error) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil, fmt.Errorf("blueprint: list by tag: tag required")
	}
	all, err := s.List("")
	if err != nil {
		return nil, err
	}
	var out []Blueprint
	for _, b := range all {
		for _, t := range b.Tags {
			if strings.EqualFold(t, tag) {
				out = append(out, b)
				break
			}
		}
	}
	return out, nil
}

// scanBlueprint decodes one row into a Blueprint. Used for the single-row
// (*sql.Row) path in Get.
func scanBlueprint(row *sql.Row) (Blueprint, error) {
	var (
		b          Blueprint
		category   string
		stepsJSON  string
		tagsJSON   string
		createdStr string
	)
	err := row.Scan(
		&b.ID, &b.Name, &b.Description, &category, &b.Phase,
		&stepsJSON, &tagsJSON, &b.Version, &b.Author, &createdStr,
	)
	if err != nil {
		return Blueprint{}, err
	}
	b.Category = Category(category)
	if err := decodeSteps(stepsJSON, &b.Steps); err != nil {
		return Blueprint{}, err
	}
	if err := json.Unmarshal([]byte(tagsJSON), &b.Tags); err != nil {
		return Blueprint{}, fmt.Errorf("decode tags: %w", err)
	}
	b.CreatedAt, err = time.Parse(time.RFC3339, createdStr)
	if err != nil {
		return Blueprint{}, fmt.Errorf("decode created_at: %w", err)
	}
	return b, nil
}

// scanBlueprints decodes all rows from a *sql.Rows iterator into a slice.
func scanBlueprints(rows *sql.Rows) ([]Blueprint, error) {
	var out []Blueprint
	for rows.Next() {
		var (
			b          Blueprint
			category   string
			stepsJSON  string
			tagsJSON   string
			createdStr string
		)
		if err := rows.Scan(
			&b.ID, &b.Name, &b.Description, &category, &b.Phase,
			&stepsJSON, &tagsJSON, &b.Version, &b.Author, &createdStr,
		); err != nil {
			return nil, err
		}
		b.Category = Category(category)
		if err := decodeSteps(stepsJSON, &b.Steps); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(tagsJSON), &b.Tags); err != nil {
			return nil, fmt.Errorf("decode tags: %w", err)
		}
		created, err := time.Parse(time.RFC3339, createdStr)
		if err != nil {
			return nil, fmt.Errorf("decode created_at: %w", err)
		}
		b.CreatedAt = created
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// decodeSteps unmarshals a steps JSON array, always returning a non-nil slice
// (json.Unmarshal into a nil slice works, but we normalize empties to keep the
// API consistent for callers).
func decodeSteps(data string, out *[]BlueprintStep) error {
	if strings.TrimSpace(data) == "" {
		*out = []BlueprintStep{}
		return nil
	}
	if err := json.Unmarshal([]byte(data), out); err != nil {
		return fmt.Errorf("decode steps: %w", err)
	}
	if *out == nil {
		*out = []BlueprintStep{}
	}
	return nil
}

// SeedDefaults populates the store with a starter library of useful pentest
// playbooks spanning all four categories. It is idempotent: each blueprint uses
// a stable ID and is inserted with INSERT OR IGNORE, so re-running SeedDefaults
// never duplicates or overwrites operator edits.
func (s *BlueprintStore) SeedDefaults() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	defaults := defaultBlueprints()

	for _, b := range defaults {
		if b.CreatedAt.IsZero() {
			b.CreatedAt = time.Now().UTC()
		}
		stepsJSON, err := json.Marshal(b.Steps)
		if err != nil {
			return fmt.Errorf("blueprint: seed %q: marshal steps: %w", b.ID, err)
		}
		tagsJSON, err := json.Marshal(b.Tags)
		if err != nil {
			return fmt.Errorf("blueprint: seed %q: marshal tags: %w", b.ID, err)
		}
		_, err = s.db.Exec(
			`INSERT OR IGNORE INTO blueprints
			    (id, name, description, category, phase, steps_json, tags_json, version, author, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			b.ID, b.Name, b.Description, string(b.Category), b.Phase,
			string(stepsJSON), string(tagsJSON), b.Version, b.Author, now,
		)
		if err != nil {
			return fmt.Errorf("blueprint: seed %q: %w", b.ID, err)
		}
	}
	return nil
}

// defaultBlueprints returns the built-in starter library. Each blueprint has a
// stable ID (bp-<slug>) so SeedDefaults is safely re-runnable.
func defaultBlueprints() []Blueprint {
	return []Blueprint{
		{
			ID:          "bp-web-app-recon",
			Name:        "Web App Recon",
			Description: "Comprehensive reconnaissance of a web application: port scan, tech fingerprinting, directory brute force, and vulnerability scanning.",
			Category:    CategoryRecon,
			Phase:       "reconnaissance",
			Version:     "1.0",
			Author:      "talon",
			Tags:        []string{"web", "recon", "http"},
			Steps: []BlueprintStep{
				{Order: 1, Tool: "nmap", Description: "TCP port scan to discover open services on the target.", Args: map[string]string{"target": "{{target_ip}}", "flags": "-sV -sC -p-"}, ExpectedResult: "List of open ports and identified service banners.", OnFailure: "Retry with -Pn if host appears down."},
				{Order: 2, Tool: "whatweb", Description: "Fingerprint web technologies, frameworks, and CMS in use.", Args: map[string]string{"url": "{{target_url}}"}, ExpectedResult: "Detected stack: server, language, CMS, JS libraries.", OnFailure: "Fall back to manual HTTP headers inspection."},
				{Order: 3, Tool: "gobuster", Description: "Brute-force common directories and files to map the application surface.", Args: map[string]string{"url": "{{target_url}}", "wordlist": "/usr/share/wordlists/dirb/common.txt", "mode": "dir"}, ExpectedResult: "Discoverable paths with their HTTP status codes.", OnFailure: "Try alternate wordlist or add extensions (-x php,asp,jsp)."},
				{Order: 4, Tool: "nuclei", Description: "Run template-based checks for known misconfigurations and CVEs.", Args: map[string]string{"target": "{{target_url}}", "severity": "low,medium,high,critical"}, ExpectedResult: "Matched vulnerability templates with evidence.", OnFailure: "Update nuclei templates (nuclei -update) and retry."},
			},
		},
		{
			ID:          "bp-network-port-scan",
			Name:        "Network Port Scan",
			Description: "Full TCP and top-UDP port discovery across a host or subnet to enumerate the attack surface.",
			Category:    CategoryRecon,
			Phase:       "discovery",
			Version:     "1.0",
			Author:      "talon",
			Tags:        []string{"network", "recon", "nmap"},
			Steps: []BlueprintStep{
				{Order: 1, Tool: "nmap", Description: "Fast SYN scan of all TCP ports.", Args: map[string]string{"target": "{{target_ip}}", "flags": "-sS -p- --min-rate 5000"}, ExpectedResult: "Complete list of open TCP ports.", OnFailure: "Fall back to -sT (connect scan) if raw sockets unavailable."},
				{Order: 2, Tool: "nmap", Description: "Service and script detection on discovered open ports.", Args: map[string]string{"target": "{{target_ip}}", "ports": "{{open_ports}}", "flags": "-sV -sC"}, ExpectedResult: "Service versions and default-script findings.", OnFailure: "Run individual service-specific tools manually."},
				{Order: 3, Tool: "nmap", Description: "Top common UDP ports probe for stateless services.", Args: map[string]string{"target": "{{target_ip}}", "flags": "-sU --top-ports 50"}, ExpectedResult: "Open/filtered UDP ports (DNS, SNMP, NTP, etc.).", OnFailure: "UDP scanning is lossy; re-run or target specific ports."},
			},
		},
		{
			ID:          "bp-subdomain-enumeration",
			Name:        "Subdomain Enumeration",
			Description: "Discover subdomains for a target domain using passive sources, active brute force, and DNS resolution.",
			Category:    CategoryRecon,
			Phase:       "reconnaissance",
			Version:     "1.0",
			Author:      "talon",
			Tags:        []string{"recon", "dns", "external", "osint"},
			Steps: []BlueprintStep{
				{Order: 1, Tool: "subfinder", Description: "Passive subdomain enumeration from public data sources.", Args: map[string]string{"domain": "{{domain}}"}, ExpectedResult: "A deduplicated list of resolved subdomain names.", OnFailure: "Add API keys to provider config and retry."},
				{Order: 2, Tool: "amass", Description: "Deeper passive + active enumeration to supplement subfinder.", Args: map[string]string{"domain": "{{domain}}", "mode": "enum"}, ExpectedResult: "Additional subdomains and associated IPs/ASNs.", OnFailure: "Run in passive-only mode if active methods are noisy."},
				{Order: 3, Tool: "dnsx", Description: "Resolve discovered names and probe for live HTTP services.", Args: map[string]string{"input": "{{subdomains_file}}", "flags": "-resp -status-code"}, ExpectedResult: "Resolved hosts with A records and HTTP status.", OnFailure: "Verify resolver config; try -wd for wildcard detection."},
			},
		},
		{
			ID:          "bp-sql-injection",
			Name:        "SQL Injection",
			Description: "Detect and exploit SQL injection vulnerabilities in a target parameter, validating with automated and manual techniques.",
			Category:    CategoryExploit,
			Phase:       "exploitation",
			Version:     "1.0",
			Author:      "talon",
			Tags:        []string{"web", "sqli", "exploit", "owasp"},
			Steps: []BlueprintStep{
				{Order: 1, Tool: "manual", Description: "Probe the target parameter with single-quote and boolean payloads to confirm injection.", Args: map[string]string{"param": "{{target_param}}", "payload_quote": "'", "payload_comment": "5; --", "payload_bool": "true"}, ExpectedResult: "Observable error, timing, or boolean difference indicating injection.", OnFailure: "Try alternate injection points (headers, JSON body)."},
				{Order: 2, Tool: "sqlmap", Description: "Automate fingerprinting and data extraction from the confirmed injection point.", Args: map[string]string{"url": "{{target_url}}", "param": "{{target_param}}", "flags": "--batch --level 3 --risk 2"}, ExpectedResult: "DBMS fingerprint and extracted database/table/column data.", OnFailure: "Increase --level/--risk or supply a custom tamper script."},
				{Order: 3, Tool: "sqlmap", Description: "Attempt to read sensitive data and, where applicable, escalate via OS commands.", Args: map[string]string{"url": "{{target_url}}", "flags": "--dump --os-shell"}, ExpectedResult: "Credential dumps or command execution evidence.", OnFailure: "Cap privileges: try --privileges and current-user checks first."},
			},
		},
		{
			ID:          "bp-xss-detection",
			Name:        "XSS Detection",
			Description: "Identify reflected and stored cross-site scripting vectors using automated fuzzing and manual payload verification.",
			Category:    CategoryExploit,
			Phase:       "exploitation",
			Version:     "1.0",
			Author:      "talon",
			Tags:        []string{"web", "xss", "exploit", "owasp"},
			Steps: []BlueprintStep{
				{Order: 1, Tool: "dalfox", Description: "Fuzz the target URL and its parameters for XSS vectors.", Args: map[string]string{"url": "{{target_url}}", "flags": "--silence"}, ExpectedResult: "List of confirmed XSS vectors with payload and parameter.", OnFailure: "Supply a custom payload list via --payload."},
				{Order: 2, Tool: "manual", Description: "Verify a flagged vector in a clean browser session to rule out false positives.", Args: map[string]string{"payload": "{{xss_payload}}"}, ExpectedResult: "Script executes in the page context (e.g. alert fires).", OnFailure: "Try bypass encodings: double-encoding, SVG, template literals."},
				{Order: 3, Tool: "manual", Description: "Assess exploitability: cookie theft, DOM keylogging, or CSRF-chain potential.", Args: map[string]string{"context": "stored"}, ExpectedResult: "Documented impact and proof-of-concept chain.", OnFailure: "Downgrade severity if impact is context-limited."},
			},
		},
		{
			ID:          "bp-web-fuzzing",
			Name:        "Web App Fuzzing",
			Description: "Content and parameter fuzzing to uncover hidden endpoints, virtual hosts, and injection-prone parameters.",
			Category:    CategoryExploit,
			Phase:       "exploitation",
			Version:     "1.0",
			Author:      "talon",
			Tags:        []string{"web", "fuzzing", "exploit"},
			Steps: []BlueprintStep{
				{Order: 1, Tool: "ffuf", Description: "Fuzz directories and file extensions beyond the common wordlist.", Args: map[string]string{"url": "{{target_url}}/FUZZ", "wordlist": "/usr/share/seclists/Discovery/Web-Content/raft-medium-directories.txt", "extensions": ".php,.asp,.aspx,.jsp,.html,.txt"}, ExpectedResult: "Additional discoverable resources and status codes.", OnFailure: "Adjust -mc (match codes) and add -recursion for depth."},
				{Order: 2, Tool: "ffuf", Description: "VHOST discovery to find host-based routing on the same IP.", Args: map[string]string{"url": "{{target_url}}", "mode": "vhost", "wordlist": "/usr/share/seclists/Discovery/DNS/subdomains-top1million-5000.txt", "header": "Host: FUZZ.{{domain}}"}, ExpectedResult: "Virtual hosts returning distinct content/sizes.", OnFailure: "Confirm target uses name-based virtual hosting first."},
				{Order: 3, Tool: "wfuzz", Description: "Parameter fuzzing to find hidden or vulnerable input names.", Args: map[string]string{"url": "{{target_url}}/?FUZZ=test", "wordlist": "/usr/share/seclists/Discovery/Web-Content/burp-parameter-names.txt"}, ExpectedResult: "Parameters that alter the response.", OnFailure: "Switch to POST body fuzzing with -d."},
			},
		},
		{
			ID:          "bp-privilege-escalation",
			Name:        "Privilege Escalation",
			Description: "Linux/Windows local privilege escalation enumeration and exploitation from an initial low-privilege foothold.",
			Category:    CategoryPostExploit,
			Phase:       "post-exploitation",
			Version:     "1.0",
			Author:      "talon",
			Tags:        []string{"linux", "windows", "privesc", "post-exploit"},
			Steps: []BlueprintStep{
				{Order: 1, Tool: "linpeas", Description: "Run comprehensive privilege-escalation enumeration script.", Args: map[string]string{"target": "{{shell}}", "mode": "auto"}, ExpectedResult: "Color-highlighted list of misconfigurations and privesc paths.", OnFailure: "Run manually section-by-section if output is truncated."},
				{Order: 2, Tool: "manual", Description: "Check sudo rights, SUID binaries, and writable cron/system paths.", Args: map[string]string{"checks": "sudo -l, find / -perm -4000 2>/dev/null, writable cron"}, ExpectedResult: "At least one exploitable misconfiguration or binary.", OnFailure: "Review kernel version for known CVE exploits (e.g. via searchsploit)."},
				{Order: 3, Tool: "manual", Description: "Attempt the identified escalation vector to obtain root/SYSTEM.", Args: map[string]string{"vector": "{{privesc_vector}}"}, ExpectedResult: "New shell/session with elevated privileges (id=0 / Administrator).", OnFailure: "Try alternative vector or escalate to operator for review."},
			},
		},
		{
			ID:          "bp-lateral-movement",
			Name:        "Lateral Movement",
			Description: "Harvest credentials from a compromised host and pivot to additional systems within the network.",
			Category:    CategoryPostExploit,
			Phase:       "post-exploitation",
			Version:     "1.0",
			Author:      "talon",
			Tags:        []string{"windows", "ad", "lateral", "post-exploit"},
			Steps: []BlueprintStep{
				{Order: 1, Tool: "manual", Description: "Harvest credentials from memory, registry, and common stores.", Args: map[string]string{"sources": "lsass, registry, browser, ssh keys"}, ExpectedResult: "Reusable credentials or password hashes.", OnFailure: "Attempt Kerberoasting/AS-REP if direct dump fails."},
				{Order: 2, Tool: "crackmapexec", Description: "Validate harvested credentials across the subnet to find reusable logins.", Args: map[string]string{"target": "{{subnet}}", "creds": "{{creds}}"}, ExpectedResult: "Hosts where the credentials authenticate (Pwn3d!).", OnFailure: "Try NTLM hash auth (-H) instead of plaintext."},
				{Order: 3, Tool: "manual", Description: "Pivot to a newly accessible host via WinRM/SMB/SSH using valid creds.", Args: map[string]string{"host": "{{reachable_host}}", "method": "winrm"}, ExpectedResult: "Interactive access on the new host; repeat enumeration.", OnFailure: "Check for admin-required protocols or LAPS rotation."},
			},
		},
		{
			ID:          "bp-initial-access",
			Name:        "Initial Access",
			Description: "Develop and deliver an initial-access payload to obtain a foothold on an exposed or phished target.",
			Category:    CategoryExploit,
			Phase:       "exploitation",
			Version:     "1.0",
			Author:      "talon",
			Tags:        []string{"exploit", "payload", "foothold"},
			Steps: []BlueprintStep{
				{Order: 1, Tool: "searchsploit", Description: "Search for public exploits matching the discovered service versions.", Args: map[string]string{"query": "{{service_name}} {{service_version}}"}, ExpectedResult: "At least one applicable exploit module or PoC.", OnFailure: "Search vendor advisories and GitHub for newer PoCs."},
				{Order: 2, Tool: "msfvenom", Description: "Generate a payload tailored to the target OS and delivery channel.", Args: map[string]string{"payload": "{{payload_name}}", "lhost": "{{lhost}}", "lhost_port": "{{lport}}", "format": "raw"}, ExpectedResult: "A working payload artifact staged for delivery.", OnFailure: "Try alternate formats/encoders or a different staged payload."},
				{Order: 3, Tool: "manual", Description: "Deliver the payload via the identified vector and confirm a callback.", Args: map[string]string{"vector": "{{delivery_vector}}"}, ExpectedResult: "Active session/callback observed on the listener.", OnFailure: "Verify listener reachability and egress filtering; try alternate port."},
			},
		},
		{
			ID:          "bp-engagement-reporting",
			Name:        "Engagement Reporting",
			Description: "Consolidate findings, evidence, and remediation guidance into a structured engagement report.",
			Category:    CategoryReporting,
			Phase:       "reporting",
			Version:     "1.0",
			Author:      "talon",
			Tags:        []string{"report", "reporting", "deliverable"},
			Steps: []BlueprintStep{
				{Order: 1, Tool: "manual", Description: "Triage and deduplicate findings, assigning final severity and status.", Args: map[string]string{"source": "{{findings}}"}, ExpectedResult: "Clean, deduplicated finding set with CVSS/impact ratings.", OnFailure: "Escalate ambiguous findings to operator for triage."},
				{Order: 2, Tool: "manual", Description: "Link each finding to its supporting evidence and reproduction steps.", Args: map[string]string{"evidence_store": "{{run_id}}"}, ExpectedResult: "Every finding has PoC + evidence references attached.", OnFailure: "Re-run the confirmatory steps to regenerate missing evidence."},
				{Order: 3, Tool: "report", Description: "Render the structured report (executive summary + technical detail).", Args: map[string]string{"format": "markdown", "run_id": "{{run_id}}"}, ExpectedResult: "A complete, review-ready report document.", OnFailure: "Fall back to a raw findings dump and flag for manual write-up."},
			},
		},
	}
}
