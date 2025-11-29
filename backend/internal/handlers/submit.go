package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"backend/internal/storage"
)

// SubmitHandler handles submission requests
type SubmitHandler struct {
	storage *storage.SQLiteStorage
}

// NewSubmitHandler creates a new submit handler
func NewSubmitHandler(storage *storage.SQLiteStorage) *SubmitHandler {
	return &SubmitHandler{storage: storage}
}

// HandleSubmit handles POST /submit requests
func (h *SubmitHandler) HandleSubmit(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON payload
	var payload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if err := validatePayload(payload); err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save to database
	if err := h.storage.SaveSubmission(payload); err != nil {
		log.Printf("Error saving submission: %v", err)
		http.Error(w, "Failed to save submission", http.StatusInternalServerError)
		return
	}

	// Get student info for logging
	examID, _ := payload["examId"].(string)
	studentID, _ := payload["studentId"].(string)
	var studentName string
	if metadata, ok := payload["metadata"].(map[string]interface{}); ok {
		studentName, _ = metadata["studentName"].(string)
	}

	log.Printf("✅ Submission saved: exam=%s, student=%s (%s)", examID, studentID, studentName)

	// Return success response
	response := map[string]interface{}{
		"success": true,
		"message": "Submission received successfully",
		"examId":  examID,
		"studentId": studentID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// validatePayload validates the submission payload
func validatePayload(payload map[string]interface{}) error {
	// Check required top-level fields
	requiredFields := []string{"examId", "studentId", "submissionTime", "metadata"}
	for _, field := range requiredFields {
		if _, ok := payload[field]; !ok {
			return &ValidationError{Field: field, Message: "required field missing"}
		}
	}

	// Validate examId
	examID, ok := payload["examId"].(string)
	if !ok || examID == "" {
		return &ValidationError{Field: "examId", Message: "must be a non-empty string"}
	}

	// Validate studentId
	studentID, ok := payload["studentId"].(string)
	if !ok || studentID == "" {
		return &ValidationError{Field: "studentId", Message: "must be a non-empty string"}
	}

	// Validate submissionTime
	submissionTime, ok := payload["submissionTime"].(string)
	if !ok || submissionTime == "" {
		return &ValidationError{Field: "submissionTime", Message: "must be a non-empty string"}
	}

	// Validate metadata
	metadata, ok := payload["metadata"].(map[string]interface{})
	if !ok {
		return &ValidationError{Field: "metadata", Message: "must be an object"}
	}

	// Validate studentName in metadata
	studentName, ok := metadata["studentName"].(string)
	if !ok || studentName == "" {
		return &ValidationError{Field: "metadata.studentName", Message: "must be a non-empty string"}
	}

	// Check for at least one question (q1, q2, etc.)
	hasQuestion := false
	for key := range payload {
		if len(key) == 2 && key[0] == 'q' && key[1] >= '1' && key[1] <= '9' {
			hasQuestion = true
			break
		}
	}
	if !hasQuestion {
		return &ValidationError{Field: "questions", Message: "at least one question (q1, q2, etc.) is required"}
	}

	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return "validation error: " + e.Field + " - " + e.Message
}
