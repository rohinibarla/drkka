// ============================================
// GLOBAL STATE
// ============================================

const reviewState = {
  submission: null,           // Parsed submission JSON
  events: [],                 // Flattened event log
  currentEventIndex: 0,       // Current position in replay
  isPlaying: false,           // Playback state
  isPaused: false,            // Pause state
  speed: 1.0,                 // Playback speed multiplier
  timeoutId: null,            // Current timeout ID
  currentText: '',            // Current text in replay field
  cursorPosition: 0,          // Current cursor position
  resumeCallback: null        // For resuming mid-character typing
}

// ============================================
// JSON LOADING & VALIDATION
// ============================================

function loadSubmission() {
  try {
    const jsonInput = document.getElementById('json-input').value.trim()

    if (!jsonInput) {
      throw new Error('Please paste a submission JSON')
    }

    const submission = JSON.parse(jsonInput)

    // Validate required fields
    if (!submission.q1) {
      throw new Error('Invalid submission format: missing q1')
    }

    if (!submission.q1.eventLog || !Array.isArray(submission.q1.eventLog)) {
      throw new Error('Invalid submission format: missing or invalid eventLog')
    }

    if (!submission.metadata) {
      throw new Error('Invalid submission format: missing metadata')
    }

    reviewState.submission = submission
    reviewState.events = submission.q1.eventLog

    // Display submission info
    displaySubmissionInfo(submission)

    // Show replay section
    document.getElementById('replay-section').classList.remove('hidden')

    // Clear any previous replay
    resetReplay()

    // Hide error
    hideError()

  } catch (error) {
    showError('Invalid JSON: ' + error.message)
  }
}

function displaySubmissionInfo(submission) {
  // Display metadata
  document.getElementById('student-name').textContent =
    submission.metadata.studentName || 'N/A'

  document.getElementById('exam-id').textContent =
    submission.examId || 'N/A'

  document.getElementById('submission-time').textContent =
    new Date(submission.submissionTime).toLocaleString() || 'N/A'

  document.getElementById('question-text').textContent =
    submission.q1.question || 'N/A'

  document.getElementById('final-answer').textContent =
    submission.q1.finalAnswer || ''

  // Calculate duration
  if (submission.q1.startTime_ms && submission.q1.endTime_ms) {
    const duration = submission.q1.endTime_ms - submission.q1.startTime_ms
    const seconds = Math.round(duration / 1000)
    document.getElementById('duration').textContent = `${seconds} seconds`
  } else {
    document.getElementById('duration').textContent = 'N/A'
  }

  // Show info panel
  document.getElementById('info-panel').classList.remove('hidden')
}

function showError(message) {
  const errorDiv = document.getElementById('error-message')
  errorDiv.textContent = message
  errorDiv.classList.remove('hidden')
}

function hideError() {
  const errorDiv = document.getElementById('error-message')
  errorDiv.classList.add('hidden')
}

// ============================================
// ASYNC REPLAY ENGINE
// ============================================

function sleep(ms) {
  return new Promise(resolve => {
    reviewState.timeoutId = setTimeout(resolve, ms)
  })
}

async function waitForResume() {
  return new Promise(resolve => {
    reviewState.resumeCallback = resolve
  })
}

async function playNextEvent() {
  // Check if we've finished
  if (reviewState.currentEventIndex >= reviewState.events.length) {
    stopReplay()
    return
  }

  // Check if paused
  if (reviewState.isPaused) {
    return  // Will resume when user clicks play
  }

  const event = reviewState.events[reviewState.currentEventIndex]

  // Wait for latency_ms before applying event
  const speed = reviewState.speed > 0 ? reviewState.speed : 1.0
  const delay = (event.latency_ms || 0) / speed
  await sleep(delay)

  // Apply the event (may involve char-by-char typing)
  await applyEvent(event)

  // Move to next event
  reviewState.currentEventIndex++
  updateProgress()

  // Continue replay
  await playNextEvent()
}

async function applyEvent(event) {
  switch(event.type) {
    case 'COMPRESSED':
      await applyCompressedEvent(event)
      break

    case 'RAW_KEY':
      insertCharacter(event.key)
      updateReplayField()
      break

    case 'RAW_SPECIAL':
      applySpecialKey(event.key)
      updateReplayField()
      break

    case 'RAW_PASTE':
      insertText(event.content)
      updateReplayField()
      break

    case 'SELECTION_CHANGE':
      updateCursorPosition(event.start, event.end)
      updateReplayField()
      break

    default:
      console.warn('Unknown event type:', event.type)
      break
  }
}

// ============================================
// EVENT HANDLERS
// ============================================

// Character-by-character COMPRESSED event handling
async function applyCompressedEvent(event) {
  const chars = event.string.split('')
  const speed = reviewState.speed > 0 ? reviewState.speed : 1.0
  const intervalDelay = (event.interval_ms || 0) / speed

  for (let i = 0; i < chars.length; i++) {
    // Check for pause before each character
    if (reviewState.isPaused) {
      // Wait until resumed
      await waitForResume()
    }

    const char = chars[i]

    // Handle escape sequences (backspace, enter, delete)
    switch (char) {
      case '\b':
        // Backspace
        applySpecialKey('Backspace')
        break

      case '\n':
        // Enter
        applySpecialKey('Enter')
        break

      case '\x7F':
        // Delete
        applySpecialKey('Delete')
        break

      default:
        // Regular character
        insertCharacter(char)
        break
    }

    updateReplayField()

    // Wait before next character (except after last char)
    if (i < chars.length - 1) {
      await sleep(intervalDelay)
    }
  }
}

// Insert a single character at cursor position
function insertCharacter(char) {
  const text = reviewState.currentText
  const pos = reviewState.cursorPosition

  reviewState.currentText =
    text.slice(0, pos) + char + text.slice(pos)

  reviewState.cursorPosition = pos + 1
}

// Insert text (for paste events)
function insertText(text) {
  const currentText = reviewState.currentText
  const pos = reviewState.cursorPosition

  reviewState.currentText =
    currentText.slice(0, pos) + text + currentText.slice(pos)

  reviewState.cursorPosition = pos + text.length
}

// Handle special keys
function applySpecialKey(key) {
  const text = reviewState.currentText
  const pos = reviewState.cursorPosition

  switch(key) {
    case 'Backspace':
      if (pos > 0) {
        reviewState.currentText =
          text.slice(0, pos - 1) + text.slice(pos)
        reviewState.cursorPosition = pos - 1
      }
      break

    case 'Delete':
      if (pos < text.length) {
        reviewState.currentText =
          text.slice(0, pos) + text.slice(pos + 1)
        // Cursor stays at same position
      }
      break

    case 'Enter':
      reviewState.currentText =
        text.slice(0, pos) + '\n' + text.slice(pos)
      reviewState.cursorPosition = pos + 1
      break

    case 'ArrowLeft':
      if (pos > 0) {
        reviewState.cursorPosition = pos - 1
      }
      break

    case 'ArrowRight':
      if (pos < text.length) {
        reviewState.cursorPosition = pos + 1
      }
      break

    case 'ArrowUp':
    case 'ArrowDown':
      // Complex to implement accurately, skip for now
      break
  }
}

// Update cursor position (for selection changes)
function updateCursorPosition(start, _end) {
  reviewState.cursorPosition = start

  // If start !== _end, we have a selection (not implemented in Phase 1)
  // For now, just move cursor to start position
}

// Update the replay textarea
function updateReplayField() {
  const field = document.getElementById('replay-field')
  field.value = reviewState.currentText

  // Set cursor position
  field.selectionStart = reviewState.cursorPosition
  field.selectionEnd = reviewState.cursorPosition

  // Refocus if replay is playing (to keep cursor visible)
  if (reviewState.isPlaying && !reviewState.isPaused) {
    field.focus()
  }
}

// ============================================
// PLAYBACK CONTROLS
// ============================================

function startReplay() {
  if (reviewState.events.length === 0) {
    showError('No events to replay')
    return
  }

  reviewState.isPlaying = true
  reviewState.isPaused = false
  updatePlayPauseButton()

  // Focus replay field once at start for visible cursor
  const field = document.getElementById('replay-field')
  field.focus()

  playNextEvent()
}

function pauseReplay() {
  reviewState.isPaused = true

  // Clear any pending timeout
  if (reviewState.timeoutId) {
    clearTimeout(reviewState.timeoutId)
    reviewState.timeoutId = null
  }

  updatePlayPauseButton()
}

function resumeReplay() {
  reviewState.isPaused = false
  updatePlayPauseButton()

  // If we were waiting during char-by-char typing, resume it
  if (reviewState.resumeCallback) {
    reviewState.resumeCallback()
    reviewState.resumeCallback = null
  } else {
    // Otherwise continue to next event
    playNextEvent()
  }
}

function togglePlayPause() {
  if (!reviewState.isPlaying) {
    startReplay()
  } else if (reviewState.isPaused) {
    resumeReplay()
  } else {
    pauseReplay()
  }
}

function resetReplay() {
  // Clear any ongoing replay
  if (reviewState.timeoutId) {
    clearTimeout(reviewState.timeoutId)
    reviewState.timeoutId = null
  }

  // Clear resume callback
  if (reviewState.resumeCallback) {
    reviewState.resumeCallback()
    reviewState.resumeCallback = null
  }

  // Reset state
  reviewState.currentEventIndex = 0
  reviewState.isPlaying = false
  reviewState.isPaused = false
  reviewState.currentText = ''
  reviewState.cursorPosition = 0

  // Clear replay field
  const field = document.getElementById('replay-field')
  if (field) {
    field.value = ''
  }

  // Update UI
  updatePlayPauseButton()
  updateProgress()
}

function stopReplay() {
  reviewState.isPlaying = false
  reviewState.isPaused = false

  // Clear timeouts
  if (reviewState.timeoutId) {
    clearTimeout(reviewState.timeoutId)
    reviewState.timeoutId = null
  }

  // Clear resume callback
  if (reviewState.resumeCallback) {
    reviewState.resumeCallback()
    reviewState.resumeCallback = null
  }

  updatePlayPauseButton()
}

function changeSpeed() {
  const speedSelect = document.getElementById('speed-select')
  const newSpeed = parseFloat(speedSelect.value)

  // Validate speed (must be positive)
  if (newSpeed > 0) {
    reviewState.speed = newSpeed
  } else {
    reviewState.speed = 1.0
    speedSelect.value = '1'
  }

  // New speed applies to next sleep() call
}

// ============================================
// UI UPDATES
// ============================================

function updateProgress() {
  const total = reviewState.events.length
  const current = reviewState.currentEventIndex
  const percentage = total > 0 ? (current / total) * 100 : 0

  // Update progress bar
  const progressBar = document.getElementById('progress-bar')
  progressBar.style.width = percentage + '%'

  // Update counter
  const counter = document.getElementById('event-counter')
  counter.textContent = `Event ${current} of ${total}`
}

function updatePlayPauseButton() {
  const btn = document.getElementById('play-pause-btn')

  if (reviewState.isPlaying && !reviewState.isPaused) {
    btn.textContent = '⏸ Pause'
  } else {
    btn.textContent = '▶ Play'
  }
}

// ============================================
// EVENT LISTENERS
// ============================================

window.addEventListener('DOMContentLoaded', () => {
  // Load button
  document.getElementById('load-btn')
    .addEventListener('click', loadSubmission)

  // Play/Pause button
  document.getElementById('play-pause-btn')
    .addEventListener('click', togglePlayPause)

  // Reset button
  document.getElementById('reset-btn')
    .addEventListener('click', resetReplay)

  // Speed selector
  document.getElementById('speed-select')
    .addEventListener('change', changeSpeed)

  // Prevent user input in replay field
  const replayField = document.getElementById('replay-field')

  replayField.addEventListener('keydown', (e) => {
    e.preventDefault()  // Block all keyboard input
  })

  replayField.addEventListener('paste', (e) => {
    e.preventDefault()  // Block paste
  })

  replayField.addEventListener('cut', (e) => {
    e.preventDefault()  // Block cut
  })
})
