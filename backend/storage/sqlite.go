package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Submission represents the exam submission data
type Submission struct {
	ExamID         string    `json:"examId"`
	StudentID      string    `json:"studentId"`
	SubmissionTime time.Time `json:"submissionTime"`
	StudentName    string    `json:"-"`
	PayloadJSON    string    `json:"-"`
}

// SQLiteStorage handles SQLite database operations
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrent performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Set connection pool settings for better concurrency
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	storage := &SQLiteStorage{db: db}

	if err := storage.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return storage, nil
}

// createTables creates the necessary database tables
func (s *SQLiteStorage) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS submissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		exam_id TEXT NOT NULL,
		student_id TEXT NOT NULL,
		student_name TEXT NOT NULL,
		submission_time DATETIME NOT NULL,
		payload_json TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(exam_id, student_id)
	);

	CREATE INDEX IF NOT EXISTS idx_exam_id ON submissions(exam_id);
	CREATE INDEX IF NOT EXISTS idx_student_id ON submissions(student_id);
	CREATE INDEX IF NOT EXISTS idx_submission_time ON submissions(submission_time);
	`

	_, err := s.db.Exec(query)
	return err
}

// SaveSubmission saves a submission to the database
func (s *SQLiteStorage) SaveSubmission(payload map[string]interface{}) error {
	// Extract metadata
	examID, _ := payload["examId"].(string)
	studentID, _ := payload["studentId"].(string)
	submissionTimeStr, _ := payload["submissionTime"].(string)

	var studentName string
	if metadata, ok := payload["metadata"].(map[string]interface{}); ok {
		studentName, _ = metadata["studentName"].(string)
	}

	// Parse submission time
	submissionTime, err := time.Parse(time.RFC3339, submissionTimeStr)
	if err != nil {
		submissionTime = time.Now()
	}

	// Convert payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Insert into database (replace if exists)
	query := `
	INSERT INTO submissions (exam_id, student_id, student_name, submission_time, payload_json)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(exam_id, student_id) DO UPDATE SET
		student_name = excluded.student_name,
		submission_time = excluded.submission_time,
		payload_json = excluded.payload_json,
		created_at = CURRENT_TIMESTAMP
	`

	_, err = s.db.Exec(query, examID, studentID, studentName, submissionTime, string(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to save submission: %w", err)
	}

	return nil
}

// GetSubmission retrieves a submission by exam ID and student ID
func (s *SQLiteStorage) GetSubmission(examID, studentID string) (map[string]interface{}, error) {
	query := `
	SELECT payload_json FROM submissions
	WHERE exam_id = ? AND student_id = ?
	`

	var payloadJSON string
	err := s.db.QueryRow(query, examID, studentID).Scan(&payloadJSON)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("submission not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve submission: %w", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return payload, nil
}

// GetSubmissionsByExam retrieves all submissions for an exam
func (s *SQLiteStorage) GetSubmissionsByExam(examID string) ([]map[string]interface{}, error) {
	query := `
	SELECT payload_json FROM submissions
	WHERE exam_id = ?
	ORDER BY submission_time DESC
	`

	rows, err := s.db.Query(query, examID)
	if err != nil {
		return nil, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []map[string]interface{}
	for rows.Next() {
		var payloadJSON string
		if err := rows.Scan(&payloadJSON); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		submissions = append(submissions, payload)
	}

	return submissions, nil
}

// GetAllSubmissions retrieves all submissions
func (s *SQLiteStorage) GetAllSubmissions() ([]map[string]interface{}, error) {
	query := `
	SELECT payload_json FROM submissions
	ORDER BY submission_time DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []map[string]interface{}
	for rows.Next() {
		var payloadJSON string
		if err := rows.Scan(&payloadJSON); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		submissions = append(submissions, payload)
	}

	return submissions, nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
