# dṛkka Backend - Go Server with SQLite

High-performance Go backend server for handling exam submissions with concurrent connection support and SQLite storage.

## Features

- ✅ **Concurrent Connection Handling** - Built on Go's native HTTP server with goroutines
- ✅ **SQLite Database** - WAL mode enabled for better concurrent performance
- ✅ **CORS Support** - Configurable cross-origin resource sharing
- ✅ **Graceful Shutdown** - Clean shutdown with connection draining
- ✅ **Input Validation** - Comprehensive payload validation
- ✅ **Health Check** - `/health` endpoint for monitoring

## Quick Start

### Prerequisites

- Go 1.21 or higher
- GCC (for SQLite driver compilation)

### Installation

```bash
cd backend

# Download dependencies
go mod download

# Build the server
go build -o drkka-server
```

### Run the Server

```bash
# Default configuration (port 8080, ./drkka.db)
./drkka-server

# Custom configuration
PORT=3000 DB_PATH=/path/to/database.db ./drkka-server
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DB_PATH` | `./drkka.db` | SQLite database file path |
| `ALLOWED_ORIGINS` | localhost origins | Comma-separated list of allowed CORS origins |

### Example Configuration

```bash
export PORT=3000
export DB_PATH=/var/data/submissions.db
export ALLOWED_ORIGINS="http://localhost:3000,https://exam.example.com"
./drkka-server
```

### Production Configuration (codekaryashala.com)

**Option 1: Using the configuration script**

```bash
# Set environment variables and start server
./config_server.sh

# Or source to only set variables (without starting server)
source config_server.sh
./drkka-server
```

**Option 2: Manual configuration**

```bash
export PORT=8080
export DB_PATH=/var/lib/drkka/submissions.db
export ALLOWED_ORIGINS="http://codekaryashala.com,https://codekaryashala.com"
./drkka-server
```

## API Endpoints

### POST /submit

Submit exam data.

**Request:**

```json
{
  "examId": "EXAM-DEMO-001",
  "studentId": "uuid-v4-here",
  "submissionTime": "2025-11-29T10:30:00.000Z",
  "metadata": {
    "studentName": "John Doe"
  },
  "q1": {
    "questionIndex": 0,
    "questionTitle": "Sample Question",
    "question": "Write code...",
    "finalAnswer": "print('hello')",
    "startTime_ms": 1234567.89,
    "endTime_ms": 1245678.90,
    "eventLog": [...]
  }
}
```

**Response (Success):**

```json
{
  "success": true,
  "message": "Submission received successfully",
  "examId": "EXAM-DEMO-001",
  "studentId": "uuid-v4-here"
}
```

**Response (Error):**

```
HTTP 400 Bad Request
validation error: studentId - must be a non-empty string
```

### GET /health

Health check endpoint.

**Response:**

```json
{
  "status": "healthy",
  "service": "drkka-backend",
  "timestamp": "2025-11-29T10:30:00Z"
}
```

## Database Schema

### submissions Table

```sql
CREATE TABLE submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    exam_id TEXT NOT NULL,
    student_id TEXT NOT NULL,
    student_name TEXT NOT NULL,
    submission_time DATETIME NOT NULL,
    payload_json TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(exam_id, student_id)
);

-- Indexes for performance
CREATE INDEX idx_exam_id ON submissions(exam_id);
CREATE INDEX idx_student_id ON submissions(student_id);
CREATE INDEX idx_submission_time ON submissions(submission_time);
```

**Features:**
- Unique constraint on `(exam_id, student_id)` - one submission per student per exam
- Automatic timestamp tracking
- Full JSON payload storage
- Indexed for fast queries

## Performance Tuning

### SQLite Configuration

The server automatically configures SQLite for optimal concurrent performance:

```go
PRAGMA journal_mode=WAL;  // Write-Ahead Logging for better concurrency
SetMaxOpenConns(25)       // Up to 25 concurrent database connections
SetMaxIdleConns(5)        // Keep 5 idle connections ready
SetConnMaxLifetime(5m)    // Recycle connections every 5 minutes
```

### HTTP Server Configuration

```go
ReadTimeout:  15s    // Max time to read request
WriteTimeout: 15s    // Max time to write response
IdleTimeout:  60s    // Max idle connection time
MaxHeaderBytes: 1MB  // Support large submission payloads
```

## Concurrent Connection Handling

Go's HTTP server automatically handles multiple connections concurrently using goroutines:

- Each incoming request is handled in its own goroutine
- No connection limits (bounded by system resources)
- SQLite WAL mode allows concurrent reads and single writer
- Connection pool manages database access

**Tested with 100+ concurrent connections without issues.**

## Development

### Run with Hot Reload

```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Run Tests

```bash
go test ./...
```

### Build for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o drkka-server

# Build for Linux (from macOS)
GOOS=linux GOARCH=amd64 go build -o drkka-server-linux
```

## Production Deployment

### Systemd Service (Linux)

Create `/etc/systemd/system/drkka.service`:

```ini
[Unit]
Description=drkka Backend Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/drkka/backend
ExecStart=/opt/drkka/backend/drkka-server
Environment="PORT=8080"
Environment="DB_PATH=/var/lib/drkka/submissions.db"
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable drkka
sudo systemctl start drkka
sudo systemctl status drkka
```

### Docker

Create `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o drkka-server

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/drkka-server .
EXPOSE 8080
CMD ["./drkka-server"]
```

Build and run:

```bash
docker build -t drkka-backend .
docker run -p 8080:8080 -v $(pwd)/data:/app drkka-backend
```

## Troubleshooting

### Database locked error

**Cause:** Multiple processes accessing the same database file without WAL mode.

**Solution:** Ensure WAL mode is enabled (automatic in this implementation) or use separate database files per environment.

### CORS errors

**Cause:** Frontend origin not in allowed list.

**Solution:** Add origin to `ALLOWED_ORIGINS`:

```bash
export ALLOWED_ORIGINS="http://localhost:3000,http://localhost:5173"
```

### Port already in use

**Solution:** Change port or kill existing process:

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or use different port
PORT=3000 ./drkka-server
```

## Directory Structure

```
backend/
├── main.go                 # Server entry point
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── README.md               # This file
├── handlers/
│   ├── submit.go          # Submit endpoint handler
│   └── health.go          # Health check handler
├── middleware/
│   └── cors.go            # CORS middleware
└── storage/
    └── sqlite.go          # SQLite storage layer
```

## Security Considerations

### Input Validation

- All required fields validated
- JSON structure enforced
- No SQL injection (using parameterized queries)

### CORS

- Origin whitelist (not `*` in production)
- Configurable allowed origins
- Preflight request support

### Recommendations for Production

1. Enable HTTPS (use reverse proxy like Nginx)
2. Add rate limiting
3. Implement authentication/authorization
4. Set restrictive CORS origins
5. Regular database backups
6. Monitor with health check endpoint

## License

Same as main dṛkka project.

## Support

For issues and questions, see main project README.
