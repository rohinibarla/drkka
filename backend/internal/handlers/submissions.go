package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"backend/internal/storage"
)

// SubmissionsHandler handles submission listing requests
type SubmissionsHandler struct {
	storage *storage.SQLiteStorage
}

// NewSubmissionsHandler creates a new submissions handler
func NewSubmissionsHandler(storage *storage.SQLiteStorage) *SubmissionsHandler {
	return &SubmissionsHandler{storage: storage}
}

// SubmissionSummary represents a simplified submission for listing
type SubmissionSummary struct {
	ExamID         string `json:"examId"`
	StudentID      string `json:"studentId"`
	StudentName    string `json:"studentName"`
	SubmissionTime string `json:"submissionTime"`
}

// HandleListSubmissions handles GET /submissions requests
func (h *SubmissionsHandler) HandleListSubmissions(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if summary=true query parameter is set
	summaryOnly := r.URL.Query().Get("summary") == "true"

	// Get all submissions from database
	submissions, err := h.storage.GetAllSubmissions()
	if err != nil {
		log.Printf("Error retrieving submissions: %v", err)
		http.Error(w, "Failed to retrieve submissions", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if summaryOnly {
		// Convert to summary format for listing
		summaries := make([]SubmissionSummary, 0, len(submissions))
		for _, submission := range submissions {
			summary := SubmissionSummary{
				ExamID:         getStringField(submission, "examId"),
				StudentID:      getStringField(submission, "studentId"),
				SubmissionTime: getStringField(submission, "submissionTime"),
			}

			// Extract student name from metadata
			if metadata, ok := submission["metadata"].(map[string]interface{}); ok {
				summary.StudentName = getStringField(metadata, "studentName")
			}

			summaries = append(summaries, summary)
		}
		json.NewEncoder(w).Encode(summaries)
		log.Printf("📋 Listed %d submission summaries", len(summaries))
	} else {
		// Return full submissions
		json.NewEncoder(w).Encode(submissions)
		log.Printf("📋 Listed %d full submissions", len(submissions))
	}
}

// getStringField safely extracts a string field from a map
func getStringField(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
