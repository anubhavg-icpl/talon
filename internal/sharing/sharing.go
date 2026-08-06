// Package sharing implements engagement sharing with capability-based
// roles, adapted from Cloudflare OS's sharing pattern. Users can share
// pentest engagements with team members using build (write) or use
// (read-only) roles. Revocation is lazy and non-destructive.
package sharing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Role determines what a collaborator can do with a shared engagement.
type Role string

const (
	RoleOwner Role = "owner" // full control, creator
	RoleBuild Role = "build" // can modify runs, add findings, execute tools
	RoleUse   Role = "use"   // read-only access to results
)

// ShareLink represents a shareable link with an embedded capability token.
type ShareLink struct {
	ID          string    `json:"id"`
	EngagementID string  `json:"engagement_id"`
	Role        Role      `json:"role"`
	Token       string    `json:"token"`       // HMAC-signed capability token
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Revoked     bool      `json:"revoked"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
	Label       string    `json:"label,omitempty"` // human-friendly name
}

// Collaborator represents a user who has access to an engagement.
type Collaborator struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	Role         Role      `json:"role"`
	GrantedAt    time.Time `json:"granted_at"`
	GrantedBy    string    `json:"granted_by"`
	ShareLinkID  string    `json:"share_link_id,omitempty"`
}

// Engagement represents a shared pentest engagement.
type Engagement struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	OwnerID     string            `json:"owner_id"`
	RunIDs      []string          `json:"run_ids"`
	CreatedAt   time.Time         `json:"created_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Store manages engagement sharing, backed by an in-memory store
// with JSON serialization support. Can be adapted to *sql.DB.
type Store struct {
	mu            sync.RWMutex
	engagements   map[string]*Engagement
	collaborators map[string][]*Collaborator // engagementID -> collaborators
	links         map[string]*ShareLink
	secret        []byte // HMAC signing key
}

// New creates a sharing store with the given HMAC secret.
func New(secret string) *Store {
	return &Store{
		engagements:   make(map[string]*Engagement),
		collaborators: make(map[string][]*Collaborator),
		links:         make(map[string]*ShareLink),
		secret:        []byte(secret),
	}
}

// CreateEngagement registers a new engagement.
func (s *Store) CreateEngagement(engagement *Engagement) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if engagement.ID == "" {
		engagement.ID = uuid.NewString()
	}
	if engagement.CreatedAt.IsZero() {
		engagement.CreatedAt = time.Now()
	}
	s.engagements[engagement.ID] = engagement

	// Owner is always a collaborator with owner role
	s.collaborators[engagement.ID] = []*Collaborator{
		{
			UserID:    engagement.OwnerID,
			Role:      RoleOwner,
			GrantedAt: time.Now(),
			GrantedBy: engagement.OwnerID,
		},
	}
	return nil
}

// GetEngagement returns an engagement by ID.
func (s *Store) GetEngagement(id string) (*Engagement, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	eng, ok := s.engagements[id]
	if !ok {
		return nil, fmt.Errorf("sharing: engagement %s not found", id)
	}
	return eng, nil
}

// ListEngagements returns all engagements for a user.
func (s *Store) ListEngagements(userID string) ([]*Engagement, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*Engagement
	for engID, collabs := range s.collaborators {
		for _, c := range collabs {
			if c.UserID == userID {
				if eng, ok := s.engagements[engID]; ok {
					result = append(result, eng)
				}
				break
			}
		}
	}
	return result, nil
}

// CreateShareLink generates a shareable link with a capability token.
func (s *Store) CreateShareLink(engagementID, createdBy string, role Role, label string, expiresAt *time.Time) (*ShareLink, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.engagements[engagementID]; !ok {
		return nil, fmt.Errorf("sharing: engagement %s not found", engagementID)
	}

	linkID := uuid.NewString()
	token := s.generateToken(engagementID, role, linkID)

	link := &ShareLink{
		ID:           linkID,
		EngagementID: engagementID,
		Role:         role,
		Token:        token,
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
		ExpiresAt:    expiresAt,
		Label:        label,
	}
	s.links[linkID] = link
	return link, nil
}

// AcceptShareLink validates a share token and adds the user as a collaborator.
func (s *Store) AcceptShareLink(token, userID, username string) (*Collaborator, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find link by token
	var link *ShareLink
	for _, l := range s.links {
		if l.Token == token {
			link = l
			break
		}
	}
	if link == nil {
		return nil, fmt.Errorf("sharing: invalid or expired share token")
	}
	if link.Revoked {
		return nil, fmt.Errorf("sharing: share link has been revoked")
	}
	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return nil, fmt.Errorf("sharing: share link has expired")
	}

	collab := &Collaborator{
		UserID:      userID,
		Username:    username,
		Role:        link.Role,
		GrantedAt:   time.Now(),
		GrantedBy:   link.CreatedBy,
		ShareLinkID: link.ID,
	}
	s.collaborators[link.EngagementID] = append(s.collaborators[link.EngagementID], collab)
	return collab, nil
}

// RevokeShareLink revokes a share link. Existing collaborators who used
// the link keep their access (lazy non-destructive revocation, matching
// CF-OS pattern) but no new users can join.
func (s *Store) RevokeShareLink(linkID, revokedBy string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	link, ok := s.links[linkID]
	if !ok {
		return fmt.Errorf("sharing: link %s not found", linkID)
	}
	link.Revoked = true
	now := time.Now()
	link.RevokedAt = &now
	return nil
}

// RemoveCollaborator removes a user's access to an engagement.
func (s *Store) RemoveCollaborator(engagementID, userID, removedBy string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	collabs := s.collaborators[engagementID]
	for i, c := range collabs {
		if c.UserID == userID {
			if c.Role == RoleOwner {
				return fmt.Errorf("sharing: cannot remove owner")
			}
			s.collaborators[engagementID] = append(collabs[:i], collabs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("sharing: user %s is not a collaborator", userID)
}

// CheckPermission verifies that a user has at least the required role
// for the given engagement. Role hierarchy: owner > build > use.
func (s *Store) CheckPermission(engagementID, userID string, required Role) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	collabs := s.collaborators[engagementID]
	for _, c := range collabs {
		if c.UserID == userID {
			return roleSatisfies(c.Role, required)
		}
	}
	return false
}

// ListCollaborators returns all collaborators for an engagement.
func (s *Store) ListCollaborators(engagementID string) ([]*Collaborator, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.engagements[engagementID]; !ok {
		return nil, fmt.Errorf("sharing: engagement %s not found", engagementID)
	}
	return s.collaborators[engagementID], nil
}

// ListShareLinks returns all share links for an engagement.
func (s *Store) ListShareLinks(engagementID string) ([]*ShareLink, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*ShareLink
	for _, link := range s.links {
		if link.EngagementID == engagementID {
			result = append(result, link)
		}
	}
	return result, nil
}

// AddRunToEngagement associates a run with an engagement.
func (s *Store) AddRunToEngagement(engagementID, runID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.checkPermissionLocked(engagementID, userID, RoleBuild) {
		return fmt.Errorf("sharing: insufficient permissions")
	}
	eng, ok := s.engagements[engagementID]
	if !ok {
		return fmt.Errorf("sharing: engagement not found")
	}
	// Check for duplicate
	for _, rid := range eng.RunIDs {
		if rid == runID {
			return nil // already added
		}
	}
	eng.RunIDs = append(eng.RunIDs, runID)
	return nil
}

// --- Helpers ---

func (s *Store) generateToken(engagementID string, role Role, linkID string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(engagementID))
	mac.Write([]byte(role))
	mac.Write([]byte(linkID))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *Store) checkPermissionLocked(engagementID, userID string, required Role) bool {
	for _, c := range s.collaborators[engagementID] {
		if c.UserID == userID {
			return roleSatisfies(c.Role, required)
		}
	}
	return false
}

// roleSatisfies checks if the granted role satisfies the required role.
// Hierarchy: owner > build > use
func roleSatisfies(granted, required Role) bool {
	rank := map[Role]int{RoleUse: 1, RoleBuild: 2, RoleOwner: 3}
	return rank[granted] >= rank[required]
}

// ShareURL generates a URL for accepting a share link.
func (s *Store) ShareURL(link *ShareLink, baseURL string) string {
	u, _ := url.Parse(baseURL)
	u.Path = "/share/accept"
	q := u.Query()
	q.Set("token", link.Token)
	q.Set("engagement", link.EngagementID)
	u.RawQuery = q.Encode()
	return u.String()
}

// Serialize exports the entire sharing state as JSON (for persistence).
func (s *Store) Serialize() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state := struct {
		Engagements   map[string]*Engagement   `json:"engagements"`
		Collaborators map[string][]*Collaborator `json:"collaborators"`
		Links         map[string]*ShareLink    `json:"links"`
	}{
		Engagements:   s.engagements,
		Collaborators: s.collaborators,
		Links:         s.links,
	}
	return json.MarshalIndent(state, "", "  ")
}

// Deserialize imports sharing state from JSON.
func (s *Store) Deserialize(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var state struct {
		Engagements   map[string]*Engagement    `json:"engagements"`
		Collaborators map[string][]*Collaborator `json:"collaborators"`
		Links         map[string]*ShareLink     `json:"links"`
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	s.engagements = state.Engagements
	s.collaborators = state.Collaborators
	s.links = state.Links
	if s.engagements == nil {
		s.engagements = make(map[string]*Engagement)
	}
	if s.collaborators == nil {
		s.collaborators = make(map[string][]*Collaborator)
	}
	if s.links == nil {
		s.links = make(map[string]*ShareLink)
	}
	return nil
}

// Ensure strings import is used
var _ = strings.Contains
