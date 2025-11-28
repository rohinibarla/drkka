# Implementation Plan: Replay & Review System

## Overview
A standalone HTML page that accepts submission JSON from the exam system and provides an accurate, step-by-step replay of how the student typed their answer, with playback controls and timing visualization.

---

## File Structure

```
project/
├── review.html            (Replay interface - standalone file)
└── review.js              (Replay engine and event handlers)
```

---

## 1. review.html

### Purpose
Provide a user interface for loading submission JSON and replaying the typing process with full fidelity.

### Key Elements

#### External Resources
- Tailwind CSS via CDN (same as index.html)
- Script tag for `review.js`

#### Custom CSS
```css
<style>
  /* Monospace font for question/answer/JSON areas */
  .code-font {
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    letter-spacing: 0.5px;
  }

  /* Cursor indicator in replay */
  .cursor-indicator {
    display: inline-block;
    width: 2px;
    height: 1.2em;
    background-color: #3b82f6;
    animation: blink 1s infinite;
  }

  @keyframes blink {
    0%, 50% { opacity: 1; }
    51%, 100% { opacity: 0; }
  }
</style>
```

#### Layout Structure

**1. JSON Input Section**
- Heading: "Load Submission"
- Textarea for pasting JSON (ID: `json-input`)
- "Load Submission" button (ID: `load-btn`)
- Error display area (ID: `error-message`, hidden by default)

**2. Submission Info Panel** (ID: `info-panel`, hidden initially)
- Student name display
- Exam ID and submission time
- Question display (in code-font)
- Duration calculation (endTime - startTime)
- Final answer display (read-only, for reference)

**3. Replay Section** (ID: `replay-section`, hidden initially)
- Section heading: "Typing Replay"
- Replay textarea (ID: `replay-field`)
  - Read-only
  - Consolas font
  - Same size as original answer field
- Playback controls:
  - Play/Pause button (ID: `play-pause-btn`)
  - Speed selector dropdown (ID: `speed-select`)
    - Options: 0.5x, 1x, 2x, 4x
  - Reset button (ID: `reset-btn`)
- Progress bar (ID: `progress-bar`)
- Event counter: "Event X of Y" (ID: `event-counter`)

**4. Event Timeline** (ID: `timeline-section`, optional, initially hidden)
- Visual timeline showing:
  - Paste events (highlighted)
  - Long pauses (>2 seconds)
  - Selection changes
  - Special key usage

### Visual Design
- Same Shadcn-inspired styling as index.html
- Clean, professional layout
- Clear visual separation between sections
- Blue accent color for interactive elements

---

## 2. review.js

### Purpose
Handle JSON parsing, replay engine logic, and UI controls.

### Global State

```javascript
const reviewState = {
  submission: null,           // Parsed submission JSON
  events: [],                 // Flattened event log
  currentEventIndex: 0,       // Current position in replay
  isPlaying: false,           // Playback state
  isPaused: false,            // Pause state
  speed: 1.0,                 // Playback speed multiplier
  replayField: null,          // Reference to replay textarea
  timeoutId: null,            // Current timeout ID
  currentText: '',            // Current text in replay field
  cursorPosition: 0           // Current cursor position
}
```

---

## A. JSON Loading & Validation

### Function: `loadSubmission()`

```javascript
function loadSubmission() {
  try {
    const jsonInput = document.getElementById('json-input').value
    const submission = JSON.parse(jsonInput)

    // Validate required fields
    if (!submission.q1 || !submission.q1.eventLog) {
      throw new Error('Invalid submission format')
    }

    reviewState.submission = submission
    reviewState.events = submission.q1.eventLog

    // Display submission info
    displaySubmissionInfo(submission)

    // Show replay section
    showReplaySection()

    // Clear any previous replay
    resetReplay()

    // Hide error
    hideError()

  } catch (error) {
    showError('Invalid JSON: ' + error.message)
  }
}
```

**Validation Checks:**
- Valid JSON format
- Contains `q1` object
- Contains `eventLog` array
- Has required metadata fields

### Function: `displaySubmissionInfo(submission)`

```javascript
function displaySubmissionInfo(submission) {
  // Display metadata
  document.getElementById('student-name').textContent =
    submission.metadata.studentName

  document.getElementById('exam-id').textContent =
    submission.examId

  document.getElementById('question-text').textContent =
    submission.q1.question

  document.getElementById('final-answer').textContent =
    submission.q1.finalAnswer

  // Calculate duration
  const duration = submission.q1.endTime_ms - submission.q1.startTime_ms
  const seconds = Math.round(duration / 1000)
  document.getElementById('duration').textContent =
    `${seconds} seconds`

  // Show info panel
  document.getElementById('info-panel').classList.remove('hidden')
}
```

---

## B. Replay Engine

### Core Replay Logic

#### Function: `startReplay()`

```javascript
function startReplay() {
  if (reviewState.events.length === 0) return

  reviewState.isPlaying = true
  reviewState.isPaused = false
  updatePlayPauseButton()

  playNextEvent()
}
```

#### Function: `playNextEvent()`

```javascript
function playNextEvent() {
  // Check if we've finished
  if (reviewState.currentEventIndex >= reviewState.events.length) {
    stopReplay()
    return
  }

  // Check if paused
  if (reviewState.isPaused) {
    return
  }

  const event = reviewState.events[reviewState.currentEventIndex]

  // Calculate delay based on speed
  const delay = event.latency_ms / reviewState.speed

  // Schedule event application
  reviewState.timeoutId = setTimeout(() => {
    applyEvent(event)
    reviewState.currentEventIndex++
    updateProgress()

    // Continue to next event
    playNextEvent()
  }, delay)
}
```

#### Function: `applyEvent(event)`

```javascript
function applyEvent(event) {
  switch(event.type) {
    case 'COMPRESSED':
      applyCompressedEvent(event)
      break

    case 'RAW_KEY':
      insertCharacter(event.key)
      break

    case 'RAW_SPECIAL':
      applySpecialKey(event.key)
      break

    case 'RAW_PASTE':
      insertText(event.content)
      break

    case 'SELECTION_CHANGE':
      updateCursorPosition(event.start, event.end)
      break
  }

  // Update textarea
  updateReplayField()
}
```

---

## C. Event Type Handlers

### 1. COMPRESSED Event Handler

```javascript
function applyCompressedEvent(event) {
  const chars = event.string.split('')
  const intervalDelay = event.interval_ms / reviewState.speed

  let charIndex = 0

  function typeNextChar() {
    if (charIndex >= chars.length) return

    insertCharacter(chars[charIndex])
    updateReplayField()
    charIndex++

    if (charIndex < chars.length) {
      setTimeout(typeNextChar, intervalDelay)
    }
  }

  // Start typing the first character immediately
  // (latency_ms already handled by playNextEvent)
  typeNextChar()
}
```

**Note:** This function needs to be modified to work with the async event loop. Alternative approach:

```javascript
function applyCompressedEvent(event) {
  // For simplicity in initial implementation:
  // Type entire string at once, but we'll enhance later for char-by-char
  insertText(event.string)
}
```

**Enhanced version (Phase 2):**
- Type character by character with interval_ms delays
- Requires recursive setTimeout or async/await approach

### 2. RAW_KEY Event Handler

```javascript
function insertCharacter(char) {
  const text = reviewState.currentText
  const pos = reviewState.cursorPosition

  reviewState.currentText =
    text.slice(0, pos) + char + text.slice(pos)

  reviewState.cursorPosition = pos + 1
}
```

### 3. RAW_SPECIAL Event Handler

```javascript
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
      // Move to previous line (complex, optional for v1)
      break

    case 'ArrowDown':
      // Move to next line (complex, optional for v1)
      break
  }
}
```

### 4. RAW_PASTE Event Handler

```javascript
function insertText(text) {
  const currentText = reviewState.currentText
  const pos = reviewState.cursorPosition

  reviewState.currentText =
    currentText.slice(0, pos) + text + currentText.slice(pos)

  reviewState.cursorPosition = pos + text.length
}
```

### 5. SELECTION_CHANGE Event Handler

```javascript
function updateCursorPosition(start, end) {
  // For now, just update cursor position
  // Selection highlighting can be added in Phase 2
  reviewState.cursorPosition = start

  // If start !== end, we have a selection
  if (start !== end) {
    reviewState.selectionStart = start
    reviewState.selectionEnd = end
  }
}
```

---

## D. Replay Field Update

### Function: `updateReplayField()`

```javascript
function updateReplayField() {
  const field = document.getElementById('replay-field')
  field.value = reviewState.currentText

  // Set cursor position
  field.selectionStart = reviewState.cursorPosition
  field.selectionEnd = reviewState.cursorPosition

  // Scroll to cursor
  field.focus()
}
```

---

## E. Playback Controls

### Function: `pauseReplay()`

```javascript
function pauseReplay() {
  reviewState.isPaused = true
  clearTimeout(reviewState.timeoutId)
  updatePlayPauseButton()
}
```

### Function: `resumeReplay()`

```javascript
function resumeReplay() {
  reviewState.isPaused = false
  updatePlayPauseButton()
  playNextEvent()
}
```

### Function: `togglePlayPause()`

```javascript
function togglePlayPause() {
  if (!reviewState.isPlaying) {
    startReplay()
  } else if (reviewState.isPaused) {
    resumeReplay()
  } else {
    pauseReplay()
  }
}
```

### Function: `resetReplay()`

```javascript
function resetReplay() {
  // Clear any ongoing replay
  clearTimeout(reviewState.timeoutId)

  // Reset state
  reviewState.currentEventIndex = 0
  reviewState.isPlaying = false
  reviewState.isPaused = false
  reviewState.currentText = ''
  reviewState.cursorPosition = 0

  // Clear replay field
  const field = document.getElementById('replay-field')
  field.value = ''

  // Update UI
  updatePlayPauseButton()
  updateProgress()
}
```

### Function: `changeSpeed()`

```javascript
function changeSpeed() {
  const speedSelect = document.getElementById('speed-select')
  reviewState.speed = parseFloat(speedSelect.value)

  // If currently playing, the new speed will apply to next event
}
```

---

## F. UI Updates

### Function: `updateProgress()`

```javascript
function updateProgress() {
  const total = reviewState.events.length
  const current = reviewState.currentEventIndex
  const percentage = (current / total) * 100

  // Update progress bar
  const progressBar = document.getElementById('progress-bar')
  progressBar.style.width = percentage + '%'

  // Update counter
  const counter = document.getElementById('event-counter')
  counter.textContent = `Event ${current} of ${total}`
}
```

### Function: `updatePlayPauseButton()`

```javascript
function updatePlayPauseButton() {
  const btn = document.getElementById('play-pause-btn')

  if (reviewState.isPlaying && !reviewState.isPaused) {
    btn.textContent = '⏸ Pause'
  } else {
    btn.textContent = '▶ Play'
  }
}
```

### Function: `showError(message)`

```javascript
function showError(message) {
  const errorDiv = document.getElementById('error-message')
  errorDiv.textContent = message
  errorDiv.classList.remove('hidden')
}
```

### Function: `hideError()`

```javascript
function hideError() {
  const errorDiv = document.getElementById('error-message')
  errorDiv.classList.add('hidden')
}
```

---

## G. Event Listeners Setup

```javascript
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
})
```

---

## 3. HTML Structure Details

### JSON Input Section

```html
<div class="bg-white border border-gray-200 rounded-lg shadow-sm p-8">
  <h2 class="text-xl font-semibold text-gray-900 mb-4">Load Submission</h2>

  <textarea
    id="json-input"
    rows="12"
    class="code-font w-full px-4 py-3 border border-gray-300 rounded-md"
    placeholder="Paste submission JSON here..."
  ></textarea>

  <button
    id="load-btn"
    class="mt-4 bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-md"
  >
    Load Submission
  </button>

  <div id="error-message" class="hidden mt-4 p-4 bg-red-50 border border-red-200 text-red-700 rounded-md">
  </div>
</div>
```

### Submission Info Panel

```html
<div id="info-panel" class="hidden mt-8 bg-white border border-gray-200 rounded-lg shadow-sm p-8">
  <h2 class="text-xl font-semibold text-gray-900 mb-4">Submission Details</h2>

  <div class="grid grid-cols-2 gap-4 mb-4">
    <div>
      <span class="text-sm font-medium text-gray-700">Student:</span>
      <span id="student-name" class="ml-2 text-gray-900"></span>
    </div>
    <div>
      <span class="text-sm font-medium text-gray-700">Exam ID:</span>
      <span id="exam-id" class="ml-2 text-gray-900"></span>
    </div>
    <div>
      <span class="text-sm font-medium text-gray-700">Duration:</span>
      <span id="duration" class="ml-2 text-gray-900"></span>
    </div>
  </div>

  <div class="mb-4">
    <label class="block text-sm font-medium text-gray-700 mb-2">Question</label>
    <div id="question-text" class="code-font text-sm p-4 bg-gray-50 rounded-lg border">
    </div>
  </div>

  <div>
    <label class="block text-sm font-medium text-gray-700 mb-2">Final Answer (Reference)</label>
    <pre id="final-answer" class="code-font text-sm p-4 bg-gray-50 rounded-lg border overflow-auto max-h-48">
    </pre>
  </div>
</div>
```

### Replay Section

```html
<div id="replay-section" class="hidden mt-8 bg-white border border-gray-200 rounded-lg shadow-sm p-8">
  <h2 class="text-xl font-semibold text-gray-900 mb-4">Typing Replay</h2>

  <textarea
    id="replay-field"
    rows="12"
    class="code-font w-full px-4 py-3 border border-gray-300 rounded-md mb-4"
    readonly
  ></textarea>

  <div class="flex items-center gap-4 mb-4">
    <button
      id="play-pause-btn"
      class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-md"
    >
      ▶ Play
    </button>

    <select
      id="speed-select"
      class="border border-gray-300 rounded-md px-4 py-2"
    >
      <option value="0.5">0.5x</option>
      <option value="1" selected>1x</option>
      <option value="2">2x</option>
      <option value="4">4x</option>
    </select>

    <button
      id="reset-btn"
      class="border border-gray-300 hover:bg-gray-50 px-6 py-2 rounded-md"
    >
      ⟲ Reset
    </button>

    <div class="flex-1"></div>

    <span id="event-counter" class="text-sm text-gray-600">
      Event 0 of 0
    </span>
  </div>

  <div class="w-full bg-gray-200 rounded-full h-2">
    <div
      id="progress-bar"
      class="bg-blue-600 h-2 rounded-full transition-all duration-200"
      style="width: 0%"
    ></div>
  </div>
</div>
```

---

## 4. Character-by-Character Implementation Strategy

### Overview
We'll use async/await to properly handle character-by-character typing in COMPRESSED events while maintaining timing accuracy and supporting pause/resume.

### Async Replay Architecture

```javascript
// Main replay loop (async)
async function playNextEvent() {
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
  await sleep(event.latency_ms / reviewState.speed)

  // Apply the event (may involve char-by-char typing)
  await applyEvent(event)

  // Move to next event
  reviewState.currentEventIndex++
  updateProgress()

  // Continue replay
  playNextEvent()
}

// Sleep utility
function sleep(ms) {
  return new Promise(resolve => {
    reviewState.timeoutId = setTimeout(resolve, ms)
  })
}
```

### COMPRESSED Event Handler (Character-by-Character)

```javascript
async function applyCompressedEvent(event) {
  const chars = event.string.split('')
  const intervalDelay = event.interval_ms / reviewState.speed

  for (let i = 0; i < chars.length; i++) {
    // Check for pause before each character
    if (reviewState.isPaused) {
      // Wait until resumed
      await waitForResume()
    }

    // Type the character
    insertCharacter(chars[i])
    updateReplayField()

    // Wait before next character (except after last char)
    if (i < chars.length - 1) {
      await sleep(intervalDelay)
    }
  }
}

// Wait for resume helper
async function waitForResume() {
  return new Promise(resolve => {
    reviewState.resumeCallback = resolve
  })
}
```

### Resume Logic

```javascript
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
```

### Updated State

```javascript
const reviewState = {
  submission: null,
  events: [],
  currentEventIndex: 0,
  isPlaying: false,
  isPaused: false,
  speed: 1.0,
  replayField: null,
  timeoutId: null,
  currentText: '',
  cursorPosition: 0,
  resumeCallback: null        // ← NEW: For resuming mid-character typing
}
```

### Speed Change Handling

```javascript
function changeSpeed() {
  const speedSelect = document.getElementById('speed-select')
  const oldSpeed = reviewState.speed
  const newSpeed = parseFloat(speedSelect.value)

  reviewState.speed = newSpeed

  // If paused, just update the speed
  // The new speed will be used when resumed

  // If playing, current timeouts will complete at old speed
  // New speed applies to next sleep() call
}
```

### Complete applyEvent Function (Async)

```javascript
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
  }
}
```

### Pause Logic (Updated)

```javascript
function pauseReplay() {
  reviewState.isPaused = true

  // Clear any pending timeout
  if (reviewState.timeoutId) {
    clearTimeout(reviewState.timeoutId)
    reviewState.timeoutId = null
  }

  updatePlayPauseButton()
}
```

### Reset Logic (Updated)

```javascript
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
  field.value = ''

  // Update UI
  updatePlayPauseButton()
  updateProgress()
}
```

---

## 5. Implementation Phases

### Phase 1: Full Character-by-Character Replay
- ✅ JSON loading and validation
- ✅ Display submission info
- ✅ **Character-by-character COMPRESSED events with interval_ms**
- ✅ All event types (RAW_KEY, RAW_SPECIAL, RAW_PASTE, SELECTION_CHANGE)
- ✅ Play/Pause/Reset controls **with pause during char typing**
- ✅ Speed control **affecting char-by-char timing**
- ✅ Progress tracking
- ✅ Async/await architecture

### Phase 2: Visual Enhancements (Future)
- ⏭ Visual cursor indicator
- ⏭ Text selection highlighting
- ⏭ Event timeline visualization
- ⏭ Long pause indicators

### Phase 3: Advanced Features (Future)
- ⏭ Scrubbing (click progress bar to jump)
- ⏭ Event-by-event stepping
- ⏭ Highlight paste events
- ⏭ Export replay
- ⏭ Side-by-side comparison

---

## 6. Key Implementation Details

### Cursor Position Visualization

**Problem:** Textarea doesn't show visible cursor when read-only.

**Solution:** Remove readonly and prevent user input:

```javascript
// In HTML: Remove readonly attribute
<textarea id="replay-field" ...></textarea>

// In JavaScript: Prevent user input
document.getElementById('replay-field').addEventListener('keydown', (e) => {
  e.preventDefault()  // Block all keyboard input
})

document.getElementById('replay-field').addEventListener('paste', (e) => {
  e.preventDefault()  // Block paste
})

document.getElementById('replay-field').addEventListener('cut', (e) => {
  e.preventDefault()  // Block cut
})
```

### Handling Arrow Keys (ArrowUp/ArrowDown)

**Simple Implementation:**
```javascript
case 'ArrowUp':
case 'ArrowDown':
  // For Phase 1, ignore these
  // They're complex to implement accurately
  break
```

**Future Enhancement:**
- Calculate line positions
- Move cursor to same column on previous/next line
- Requires analyzing newline positions in text

### Stop Replay

```javascript
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
```

---

## 7. Testing Checklist

### Functionality Tests
- [ ] Load valid JSON successfully
- [ ] Show error for invalid JSON
- [ ] Display all submission info correctly
- [ ] Replay RAW_KEY events
- [ ] Replay RAW_SPECIAL events (Backspace, Delete, Enter)
- [ ] Replay RAW_PASTE events
- [ ] Replay SELECTION_CHANGE events
- [ ] Replay COMPRESSED events
- [ ] Play/Pause works correctly
- [ ] Reset clears replay field
- [ ] Speed control affects playback
- [ ] Progress bar updates correctly
- [ ] Final text matches finalAnswer

### Edge Cases
- [ ] Empty event log
- [ ] Only COMPRESSED events
- [ ] Only special keys
- [ ] Multiple paste events
- [ ] Very long answers
- [ ] Very fast typing (low interval_ms)
- [ ] Pause during COMPRESSED event

---

## 8. Sample Test JSON

```json
{
  "examId": "EXAM-DEMO-001",
  "studentId": "550e8400-e29b-41d4-a716-446655440000",
  "submissionTime": "2025-11-28T10:30:00.000Z",
  "metadata": {
    "studentName": "John Doe"
  },
  "q1": {
    "questionIndex": 0,
    "questionTitle": "Flight Booking Confirmation",
    "question": "Extract passenger name and flight details.",
    "finalAnswer": "Name: Rohini\nFlight: AI-205",
    "startTime_ms": 1732790000000,
    "endTime_ms": 1732790045000,
    "eventLog": [
      { "type": "RAW_KEY", "key": "N", "latency_ms": 0 },
      { "type": "COMPRESSED", "string": "ame: ", "latency_ms": 150, "interval_ms": 145 },
      { "type": "RAW_PASTE", "content": "Rohini", "latency_ms": 1200 },
      { "type": "RAW_SPECIAL", "key": "Enter", "latency_ms": 300 },
      { "type": "COMPRESSED", "string": "Flight: AI-205", "latency_ms": 200, "interval_ms": 140 }
    ]
  }
}
```

---

## 9. Requirements Summary

✅ **Core Features**
- Load and validate submission JSON
- Display submission metadata and question
- **Character-by-character replay of COMPRESSED events with interval_ms timing**
- Accurate replay of all 5 event types (COMPRESSED, RAW_KEY, RAW_SPECIAL, RAW_PASTE, SELECTION_CHANGE)
- Play/Pause/Reset controls **with pause during character typing**
- Speed adjustment (0.5x, 1x, 2x, 4x) **affecting char-by-char delays**
- Progress tracking
- Async/await architecture for timing accuracy

✅ **Styling**
- Tailwind CSS
- Shadcn-inspired design
- Consolas font for code areas
- Responsive layout

✅ **User Experience**
- Clear error messages
- Smooth character-by-character playback
- Visual feedback (progress bar, event counter)
- Visible cursor in replay field
- Reference final answer visible

---

## End of Review Implementation Plan
