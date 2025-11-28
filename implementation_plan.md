# Implementation Plan: Keystroke Capture & Replay System (Phase 1)

## Overview
Single-page HTML application that captures keystroke dynamics for one exam question, compresses the event log, and outputs JSON to both console and page display.

---

## File Structure

```
project/
├── index.html              (HTML structure with Tailwind CSS, no inline JS)
├── main.js                 (Event capture, form handling & question loading)
├── process_and_pack.js     (Processing, compression & JSON generation)
└── questions.json          (Question bank)
```

---

## 1. questions.json

### Purpose
Store a bank of questions from which one will be randomly selected on page load.

### Structure
The file contains an array of question objects. Each object has:
- `question_title`: Title/name of the question (not displayed, used for tracking)
- `question`: The actual question text to display to the user
- `variables`: Metadata about variables (not used in Phase 1)

### Usage
- Only the `question` field is displayed in the HTML
- The array index or `question_title` is used as identifier in JSON output

---

## 2. index.html

### Purpose
Provide the form interface with Shadcn-inspired styling using Tailwind CSS.

### Key Elements

#### External Resources
- Tailwind CSS via CDN
- Script tags for `main.js` and `process_and_pack.js`

#### Custom CSS
```css
<style>
  /* Monospace font for code/question/answer areas */
  .code-font {
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    letter-spacing: 0.5px;
  }

  /* Ensure JSON output also uses monospace */
  #json-output {
    font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
    font-size: 0.875rem;
  }
</style>
```

#### Form Structure
1. **Container**: Centered card with max-width, padding, shadow
2. **Name Input**:
   - Label: "Enter name:"
   - Input type: text
   - ID: `student-name`
   - Not keystroke-tracked (metadata only)
3. **Question Display**:
   - Label: "Question"
   - Div with ID: `question-text`
   - Initially shows: "Loading question..."
   - Class: `code-font` (Consolas)
   - Background: `bg-gray-50`
4. **Answer Textarea**:
   - Label: "Answer"
   - Textarea ID: `answer-field`
   - Class: `code-font` (Consolas)
   - Rows: 10-12
   - Placeholder: "Type your answer here..."
   - **This field is keystroke-tracked**
5. **Submit Button**:
   - Text: "Submit"
   - Styled with Shadcn-like blue button
   - ID: `submit-btn`

#### JSON Output Section
- Initially hidden (class: `hidden`)
- ID: `output-section`
- Contains:
  - Heading: "Submission JSON:"
  - `<pre>` with `<code>` block
  - ID: `json-output`
  - Scrollable, max height ~400px

#### Visual Design (Shadcn-like)
- **Background**: `bg-slate-50`
- **Card**: `bg-white border border-gray-200 rounded-lg shadow-sm`
- **Inputs/Textarea**: `border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-500`
- **Button**: `bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-md`
- **Question Box**: `bg-gray-50 border border-gray-200 rounded-lg p-4`

---

## 3. main.js

### Purpose
Handle event capture, question loading, and form submission.

### Global State
```javascript
const captureData = {
  rawEvents: [],           // Array of captured events
  startTime_ms: null,      // First interaction timestamp
  lastSelection: { start: 0, end: 0 }  // Track last selection state
}

let selectedQuestion = null  // Currently displayed question
```

### A. Question Loading

#### Function: `loadRandomQuestion()`
```javascript
async function loadRandomQuestion() {
  try {
    const response = await fetch('questions.json')
    const questions = await response.json()  // Array of questions

    // Select random question
    const randomIndex = Math.floor(Math.random() * questions.length)
    selectedQuestion = {
      index: randomIndex,
      title: questions[randomIndex].question_title,
      text: questions[randomIndex].question  // Use 'question' field only
    }

    // Display only the question text
    document.getElementById('question-text').textContent = selectedQuestion.text

  } catch (error) {
    console.error('Error loading questions:', error)
    // Fallback question
    selectedQuestion = {
      index: -1,
      title: 'Fallback Question',
      text: 'Write a function to reverse a string'
    }
    document.getElementById('question-text').textContent = selectedQuestion.text
  }
}

// Call on page load
window.addEventListener('DOMContentLoaded', loadRandomQuestion)
```

**Behavior:**
- Fetches `questions.json` on page load (returns array directly)
- Randomly selects one question
- Stores index and question_title for tracking
- **Only displays the `question` field** in the HTML
- Falls back to default question if fetch fails

### B. Event Listeners Setup

#### Initialize on DOMContentLoaded
```javascript
window.addEventListener('DOMContentLoaded', () => {
  const answerField = document.getElementById('answer-field')
  const submitBtn = document.getElementById('submit-btn')

  // Track first interaction
  answerField.addEventListener('focus', setStartTime, { once: true })
  answerField.addEventListener('click', setStartTime, { once: true })
  answerField.addEventListener('keydown', handleFirstKeydown, { once: true })

  // Ongoing event capture
  answerField.addEventListener('keydown', handleKeydown)
  answerField.addEventListener('paste', handlePaste)
  answerField.addEventListener('mouseup', handleSelection)
  answerField.addEventListener('select', handleSelection)

  // Submit
  submitBtn.addEventListener('click', handleSubmit)

  // Load question
  loadRandomQuestion()
})
```

### C. Start Time Tracking

#### Function: `setStartTime()`
```javascript
function setStartTime() {
  if (!captureData.startTime_ms) {
    captureData.startTime_ms = performance.now()
  }
}
```

#### Function: `handleFirstKeydown(e)`
```javascript
function handleFirstKeydown(e) {
  setStartTime()
  handleKeydown(e)  // Also capture this first keystroke
}
```

**Behavior:**
- `startTime_ms` is set on the FIRST interaction:
  - Focus event
  - Click event
  - First keydown event
- Uses `performance.now()` for high-resolution timestamp

### D. Event Capture Handlers

#### Function: `handleKeydown(e)`
```javascript
function handleKeydown(e) {
  const now = performance.now()
  const key = e.key

  if (key.length === 1) {
    // Regular character (a, b, 1, !, etc.)
    captureData.rawEvents.push({
      type: 'key',
      key: key,
      timestamp: now
    })
  } else if (['Backspace', 'Delete', 'Enter', 'ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown'].includes(key)) {
    // Special keys
    captureData.rawEvents.push({
      type: 'special',
      key: key,
      timestamp: now
    })
  }
  // Ignore other keys (Shift, Ctrl, Alt, etc.)
}
```

**Captured Special Keys:**
- Backspace
- Delete
- Enter
- Arrow keys (Left, Right, Up, Down)

**Ignored Keys:**
- Modifiers (Shift, Ctrl, Alt, Meta)
- Function keys (F1, F2, etc.)
- Other non-printable keys

#### Function: `handlePaste(e)`
```javascript
function handlePaste(e) {
  const now = performance.now()
  const content = e.clipboardData.getData('text')

  captureData.rawEvents.push({
    type: 'paste',
    content: content,
    timestamp: now
  })
}
```

**Behavior:**
- Captures full pasted text content
- Records timestamp

#### Function: `handleSelection(e)`
```javascript
function handleSelection(e) {
  const start = e.target.selectionStart
  const end = e.target.selectionEnd

  // Only log if selection changed
  if (start !== captureData.lastSelection.start ||
      end !== captureData.lastSelection.end) {
    const now = performance.now()

    captureData.rawEvents.push({
      type: 'selection',
      start: start,
      end: end,
      timestamp: now
    })

    captureData.lastSelection = { start, end }
  }
}
```

**Behavior:**
- Captures cursor position changes (when start === end)
- Captures text selections (when start !== end)
- Only logs when selection actually changes
- Triggered by mouse clicks, drags, or keyboard selection

### E. Submit Handler

#### Function: `handleSubmit(e)`
```javascript
function handleSubmit(e) {
  e.preventDefault()

  const answerField = document.getElementById('answer-field')
  const nameInput = document.getElementById('student-name')

  // Call process_and_pack.js function
  const payload = processAndPack({
    rawEvents: captureData.rawEvents,
    startTime_ms: captureData.startTime_ms,
    finalAnswer: answerField.value,
    questionIndex: selectedQuestion.index,
    questionTitle: selectedQuestion.title,
    questionText: selectedQuestion.text,
    metadata: {
      studentName: nameInput.value
    }
  })

  // Console output
  console.log(JSON.stringify(payload, null, 2))

  // Display on page
  displayJSON(payload)
}
```

#### Function: `displayJSON(payload)`
```javascript
function displayJSON(payload) {
  const outputSection = document.getElementById('output-section')
  const jsonOutput = document.getElementById('json-output')

  // Format JSON with indentation
  jsonOutput.textContent = JSON.stringify(payload, null, 2)

  // Show output section
  outputSection.classList.remove('hidden')

  // Scroll to output
  outputSection.scrollIntoView({ behavior: 'smooth' })
}
```

**Behavior:**
- Prevents default form submission
- Collects all data (events, metadata, question)
- Calls `processAndPack()` from process_and_pack.js
- Logs JSON to console
- Displays formatted JSON on page
- Scrolls to output section

---

## 4. process_and_pack.js

### Purpose
Process raw events, apply compression algorithm, and generate final JSON payload.

### Main Function: `processAndPack(data)`

```javascript
function processAndPack(data) {
  // data = {
  //   rawEvents: [],
  //   startTime_ms: number,
  //   finalAnswer: string,
  //   questionIndex: number,
  //   questionTitle: string,
  //   questionText: string,
  //   metadata: { studentName: string }
  // }

  // 1. Calculate endTime_ms from last event or now
  const endTime_ms = calculateEndTime(data.rawEvents, data.startTime_ms)

  // 2. Compress event log
  const compressedLog = compressEvents(data.rawEvents)

  // 3. Build final JSON
  return {
    examId: "EXAM-DEMO-001",
    studentId: crypto.randomUUID(),
    submissionTime: new Date().toISOString(),
    metadata: data.metadata,
    q1: {
      questionIndex: data.questionIndex,
      questionTitle: data.questionTitle,
      question: data.questionText,
      finalAnswer: data.finalAnswer,
      startTime_ms: data.startTime_ms,
      endTime_ms: endTime_ms,
      eventLog: compressedLog
    }
  }
}
```

### A. Calculate End Time

#### Function: `calculateEndTime(rawEvents, startTime_ms)`
```javascript
function calculateEndTime(rawEvents, startTime_ms) {
  if (rawEvents.length === 0) {
    return startTime_ms || performance.now()
  }

  // Return timestamp of last event
  const lastEvent = rawEvents[rawEvents.length - 1]
  return lastEvent.timestamp
}
```

**Behavior:**
- If no events: use startTime or current time
- Otherwise: use timestamp of last captured event

### B. Compression Algorithm

#### Function: `compressEvents(rawEvents)`
```javascript
function compressEvents(rawEvents) {
  const compressed = []
  let i = 0

  while (i < rawEvents.length) {
    const event = rawEvents[i]

    // Calculate latency from previous event
    const latency_ms = i === 0 ? 0 :
      Math.round(event.timestamp - rawEvents[i-1].timestamp)

    if (event.type === 'key') {
      // Try to build a compressible segment
      const segment = extractSegment(rawEvents, i)

      if (segment.canCompress) {
        compressed.push({
          type: 'COMPRESSED',
          string: segment.string,
          interval_ms: segment.meanInterval
        })
        i += segment.length
      } else {
        compressed.push({
          type: 'RAW_KEY',
          key: event.key,
          latency_ms: latency_ms
        })
        i++
      }
    }
    else if (event.type === 'special') {
      compressed.push({
        type: 'RAW_SPECIAL',
        key: event.key,
        latency_ms: latency_ms
      })
      i++
    }
    else if (event.type === 'paste') {
      compressed.push({
        type: 'RAW_PASTE',
        content: event.content,
        latency_ms: latency_ms
      })
      i++
    }
    else if (event.type === 'selection') {
      compressed.push({
        type: 'SELECTION_CHANGE',
        start: event.start,
        end: event.end,
        latency_ms: latency_ms
      })
      i++
    }
  }

  return compressed
}
```

**Key Points:**
- Processes events sequentially
- Calculates latency relative to previous event
- First event always has latency_ms = 0
- Attempts compression for 'key' events only
- Special keys, paste, selection always stay raw

#### Function: `extractSegment(rawEvents, startIdx)`
```javascript
function extractSegment(rawEvents, startIdx) {
  const THRESHOLD_STDDEV = 30  // milliseconds
  const MIN_SEGMENT_LENGTH = 3

  const segment = []
  let i = startIdx

  // Collect consecutive 'key' events
  while (i < rawEvents.length && rawEvents[i].type === 'key') {
    segment.push(rawEvents[i])
    i++
  }

  // Need at least MIN_SEGMENT_LENGTH to compress
  if (segment.length < MIN_SEGMENT_LENGTH) {
    return { canCompress: false }
  }

  // Calculate inter-key intervals
  const intervals = []
  for (let j = 1; j < segment.length; j++) {
    intervals.push(segment[j].timestamp - segment[j-1].timestamp)
  }

  // Calculate standard deviation
  const stddev = calculateStdDev(intervals)

  // Compress if consistent timing
  if (stddev <= THRESHOLD_STDDEV) {
    return {
      canCompress: true,
      length: segment.length,
      string: segment.map(e => e.key).join(''),
      meanInterval: Math.round(mean(intervals))
    }
  }

  return { canCompress: false }
}
```

**Algorithm:**
1. Collect consecutive 'key' events into a segment
2. If segment < 3 characters: don't compress
3. Calculate inter-key intervals (time differences)
4. Calculate standard deviation of intervals
5. If StdDev ≤ 30ms: compress into COMPRESSED event
6. Otherwise: return individual RAW_KEY events

**Compression Criteria:**
- Minimum segment length: 3 characters
- Maximum standard deviation: 30ms
- Only applies to regular character keys

#### Helper Functions

##### `calculateStdDev(values)`
```javascript
function calculateStdDev(values) {
  if (values.length === 0) return 0

  const avg = mean(values)
  const squareDiffs = values.map(value => Math.pow(value - avg, 2))
  const avgSquareDiff = mean(squareDiffs)
  return Math.sqrt(avgSquareDiff)
}
```

##### `mean(values)`
```javascript
function mean(values) {
  if (values.length === 0) return 0
  return values.reduce((sum, val) => sum + val, 0) / values.length
}
```

---

## 5. Output JSON Schema

### Complete Payload Structure

```json
{
  "examId": "EXAM-DEMO-001",
  "studentId": "550e8400-e29b-41d4-a716-446655440000",
  "submissionTime": "2025-11-28T10:30:00.000Z",
  "metadata": {
    "studentName": "John Doe"
  },
  "q1": {
    "questionIndex": 2,
    "questionTitle": "E-commerce Package Delivery Notification",
    "question": "Your order OD20231208001 containing 2 units of Samsung Mobile (Rs.45,999 each) has been dispatched. Tracking ID: TRK987654321, Expected delivery: 12-Dec-2023 by 6:00 PM, Delivery address: Flat 5, Maple Apartments, Hyderabad 500081.",
    "finalAnswer": "OrderID: OD20231208001\nQuantity: 2\nProductName: Samsung Mobile\nUnitPrice: Rs.45,999",
    "startTime_ms": 1732790000000,
    "endTime_ms": 1732790045000,
    "eventLog": [
      { "type": "RAW_KEY", "key": "O", "latency_ms": 0 },
      { "type": "COMPRESSED", "string": "rderID: OD", "latency_ms": 145, "interval_ms": 145 },
      { "type": "RAW_SPECIAL", "key": "Enter", "latency_ms": 200 },
      { "type": "RAW_PASTE", "content": "Quantity: 2", "latency_ms": 1500 },
      { "type": "SELECTION_CHANGE", "start": 35, "end": 35, "latency_ms": 600 },
      { "type": "COMPRESSED", "string": "ProductName", "latency_ms": 150, "interval_ms": 150 }
    ]
  }
}
```

### Field Descriptions

#### Top Level
| Field | Type | Description |
|-------|------|-------------|
| `examId` | string | Hardcoded: "EXAM-DEMO-001" |
| `studentId` | string | Auto-generated UUID |
| `submissionTime` | string | ISO 8601 timestamp of submission |
| `metadata` | object | Non-tracked form fields |
| `q1` | object | Question 1 data and event log |

#### metadata Object
| Field | Type | Description |
|-------|------|-------------|
| `studentName` | string | Value from name input field |

#### q1 Object
| Field | Type | Description |
|-------|------|-------------|
| `questionIndex` | number | Index in questions.json array |
| `questionTitle` | string | question_title from questions.json |
| `question` | string | Question text displayed (from 'question' field) |
| `finalAnswer` | string | Final textarea content |
| `startTime_ms` | number | Timestamp of first interaction |
| `endTime_ms` | number | Timestamp of last event |
| `eventLog` | array | Compressed event log |

#### Event Types in eventLog

##### COMPRESSED Event
```json
{
  "type": "COMPRESSED",
  "string": "hello world",
  "latency_ms": 500,
  "interval_ms": 145
}
```

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Must be "COMPRESSED" |
| `string` | string | The compressed string segment |
| `latency_ms` | number | Time from previous event to first character |
| `interval_ms` | number | Mean inter-key interval within segment |

##### RAW_KEY Event
```json
{
  "type": "RAW_KEY",
  "key": "a",
  "latency_ms": 250
}
```

##### RAW_SPECIAL Event
```json
{
  "type": "RAW_SPECIAL",
  "key": "Backspace",
  "latency_ms": 300
}
```

##### RAW_PASTE Event
```json
{
  "type": "RAW_PASTE",
  "content": "print('Hello')",
  "latency_ms": 1200
}
```

##### SELECTION_CHANGE Event
```json
{
  "type": "SELECTION_CHANGE",
  "start": 10,
  "end": 15,
  "latency_ms": 650
}
```

---

## 6. Application Flow

### Page Load
1. Browser loads `index.html`
2. Tailwind CSS loads from CDN
3. `main.js` and `process_and_pack.js` load
4. `DOMContentLoaded` event fires
5. `loadRandomQuestion()` fetches and displays random question
6. Event listeners attached to answer textarea and submit button

### User Interaction
1. User types name (not tracked)
2. User clicks/focuses on answer textarea
   - **First interaction triggers `startTime_ms` recording**
3. User types answer:
   - Each keystroke captured with timestamp
   - Special keys (Backspace, Enter, etc.) captured
   - Paste events captured with full content
   - Mouse selections captured
4. All events stored in `captureData.rawEvents` array

### Submission
1. User clicks Submit button
2. `handleSubmit()` called:
   - Prevents default form submission
   - Collects all data
   - Calls `processAndPack(data)`
3. `processAndPack()` executes:
   - Calculates `endTime_ms` from last event
   - Compresses event log using standard deviation algorithm
   - Builds final JSON payload
   - Returns payload
4. `handleSubmit()` continues:
   - Logs JSON to console
   - Displays JSON on page
   - Scrolls to output section

---

## 7. Compression Algorithm Details

### Goal
Reduce payload size by identifying and compressing segments of consistent typing.

### Criteria for Compression
- **Type**: Only regular character keys (not special keys, paste, or selection)
- **Minimum Length**: At least 3 consecutive characters
- **Consistency**: Standard deviation of inter-key intervals ≤ 30ms

### Example

**Input (Raw Events):**
```javascript
[
  { type: 'key', key: 'h', timestamp: 1000 },
  { type: 'key', key: 'e', timestamp: 1150 },  // 150ms
  { type: 'key', key: 'l', timestamp: 1300 },  // 150ms
  { type: 'key', key: 'l', timestamp: 1445 },  // 145ms
  { type: 'key', key: 'o', timestamp: 1595 },  // 150ms
]
```

**Intervals:** [150, 150, 145, 150]
**Mean:** 148.75ms
**StdDev:** ~2.17ms ✓ (< 30ms)
**Length:** 5 characters ✓ (>= 3)

**Output (Compressed):**
```json
{
  "type": "COMPRESSED",
  "string": "hello",
  "interval_ms": 149
}
```

### Non-Compressible Events

**Special Keys:**
```json
{ "type": "RAW_SPECIAL", "key": "Backspace", "latency_ms": 300 }
```

**Paste:**
```json
{ "type": "RAW_PASTE", "content": "pasted text", "latency_ms": 1200 }
```

**Selection:**
```json
{ "type": "SELECTION_CHANGE", "start": 5, "end": 10, "latency_ms": 650 }
```

**Variable Timing Keys:**
If inter-key intervals have high variance (StdDev > 30ms), each character becomes:
```json
{ "type": "RAW_KEY", "key": "a", "latency_ms": 250 }
```

---

## 8. Implementation Checklist

### Phase 1: Setup
- [ ] Create `questions.json` with 5+ sample questions
- [ ] Create `index.html` with form structure
- [ ] Add Tailwind CSS CDN
- [ ] Add custom CSS for Consolas font
- [ ] Create empty `main.js` file
- [ ] Create empty `process_and_pack.js` file
- [ ] Link JS files in HTML

### Phase 2: Question Loading (main.js)
- [ ] Implement `loadRandomQuestion()` function
- [ ] Add fetch logic for `questions.json`
- [ ] Add random selection logic
- [ ] Update question display div
- [ ] Add error handling with fallback question

### Phase 3: Event Capture (main.js)
- [ ] Initialize `captureData` object
- [ ] Implement `setStartTime()` function
- [ ] Attach focus/click/keydown listeners for first interaction
- [ ] Implement `handleKeydown()` for regular and special keys
- [ ] Implement `handlePaste()` for paste events
- [ ] Implement `handleSelection()` for cursor/selection changes
- [ ] Test: Verify all events captured in console

### Phase 4: Compression Algorithm (process_and_pack.js)
- [ ] Implement `mean()` helper function
- [ ] Implement `calculateStdDev()` helper function
- [ ] Implement `extractSegment()` function
- [ ] Implement `compressEvents()` function
- [ ] Test: Verify compression with sample data

### Phase 5: JSON Generation (process_and_pack.js)
- [ ] Implement `calculateEndTime()` function
- [ ] Implement `processAndPack()` main function
- [ ] Generate proper JSON structure
- [ ] Test: Verify JSON matches schema

### Phase 6: Submit & Display (main.js)
- [ ] Implement `handleSubmit()` function
- [ ] Call `processAndPack()` with collected data
- [ ] Implement `displayJSON()` function
- [ ] Add console.log output
- [ ] Add page display in output section
- [ ] Add smooth scroll to output

### Phase 7: Testing
- [ ] Test: Type continuously (should compress)
- [ ] Test: Type with long pauses (should not compress)
- [ ] Test: Backspace and Delete keys
- [ ] Test: Arrow key navigation
- [ ] Test: Paste text
- [ ] Test: Mouse selection
- [ ] Test: Mixed typing patterns
- [ ] Test: Empty answer submission
- [ ] Test: Question loading failure
- [ ] Verify JSON output format matches schema

### Phase 8: Polish
- [ ] Verify Consolas font applied correctly
- [ ] Test responsive layout
- [ ] Add loading state for question
- [ ] Verify Shadcn-like styling
- [ ] Test in different browsers (Chrome, Firefox, Safari)
- [ ] Validate JSON structure
- [ ] Add comments to code

---

## 9. Requirements Summary

✅ **File Structure**
- Separate HTML, JS files (no inline JS)
- `main.js` for event capture
- `process_and_pack.js` for compression

✅ **Question Loading**
- Random selection from `questions.json`
- Question displayed on page load
- Question ID and text in JSON output

✅ **Timing**
- `startTime_ms` recorded on first interaction (focus/click/type)
- `endTime_ms` calculated in `process_and_pack.js`

✅ **Event Capture**
- All keystrokes with timestamps
- Special keys (Backspace, Delete, Enter, Arrows)
- Paste events with full content
- Mouse selections and cursor positioning

✅ **Compression**
- Detect consistent typing segments
- Standard deviation threshold: 30ms
- Minimum segment length: 3 characters
- Compress into COMPRESSED events with mean interval

✅ **Output**
- Console.log JSON payload
- Display formatted JSON on page
- Follow schema from json_schema.md

✅ **Styling**
- Tailwind CSS via CDN
- Shadcn-inspired components
- **Consolas font for question and answer areas**
- **Consolas font for JSON output**

---

## 10. Future Enhancements (Not in Phase 1)

- Multiple questions support
- Backend API integration
- Data persistence
- Replay functionality
- Export JSON as file
- Visual timeline of events
- Real-time event preview
- Pause detection indicators
- Copy-paste highlighting
- Mobile responsiveness improvements

---

## End of Implementation Plan
