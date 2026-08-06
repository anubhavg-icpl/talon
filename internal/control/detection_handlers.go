package control

import (
	"net/http"

	"github.com/anubhavg-icpl/talon/internal/core"
)

// handleDetectionListCases returns detection cases for the SOC pipeline view.
// Since CaseStore is run-scoped, this returns an empty list for completed runs
// unless cases were persisted.
func (s *Server) handleDetectionListCases(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"cases": []any{},
		"note":  "Detection cases are managed per-run via agent tools (detection_create_case, detection_triage, detection_investigate, detection_tune). Use the run evidence/traffic endpoints to inspect case data for completed runs.",
	})
}

// handleDetectionSkills returns the detection skill catalog.
func (s *Server) handleDetectionSkills(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	stage := r.URL.Query().Get("stage")
	cat := r.URL.Query().Get("category")

	// Search detection skills
	res := core.QuerySkills(core.SkillQuery{
		Q:        "detection " + query,
		Category: cat,
		Stage:    stage,
		Brief:    true,
		Limit:    50,
	})

	type skill struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Category string `json:"category"`
		Stage    string `json:"stage"`
	}

	skills := make([]skill, 0, len(res.Skills))
	for _, sk := range res.Skills {
		// Filter to detection skills only
		if sk.Category == "triage" || sk.Category == "investigation" || sk.Category == "tuning" {
			skills = append(skills, skill{
				ID: sk.ID, Name: sk.Name, Category: sk.Category, Stage: sk.Stage,
			})
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"total": len(skills),
		"skills": skills,
	})
}

// handleDetectionSkillsByType returns detection skills filtered by type.
func (s *Server) handleDetectionSkillsByType(w http.ResponseWriter, r *http.Request) {
	skillType := r.PathValue("type")
	if skillType != "triage" && skillType != "investigation" && skillType != "tuning" {
		writeError(w, http.StatusBadRequest, "type must be triage, investigation, or tuning")
		return
	}

	res := core.QuerySkills(core.SkillQuery{
		Q:        "detection",
		Category: skillType,
		Brief:    true,
		Limit:    50,
	})

	type skill struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Category string `json:"category"`
		Stage    string `json:"stage"`
		Path     string `json:"path"`
	}

	skills := make([]skill, 0, len(res.Skills))
	for _, sk := range res.Skills {
		skills = append(skills, skill{
			ID: sk.ID, Name: sk.Name, Category: sk.Category, Stage: sk.Stage, Path: sk.Path,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"type":   skillType,
		"total":  len(skills),
		"skills": skills,
	})
}
