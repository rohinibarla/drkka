# dṛkka

```
దృక్క
దృక్ + క
చూడగలవాడు; చూసే శక్తి గల వ్యక్తి.
```

## Goal 
Form-based exams with a precise replay of how each answer was typed.

## System goal
The system captures a complete, time-stamped record of all user input—keystrokes, edits, cursor movements, selections, mouse actions, and paste events—for each answer, augmented with clear start and end time boundaries. It detects segments of consistent typing speed, compresses them to reduce payload size, and enables an accurate, step-by-step replay of how every answer was written.


---

# ✅ **Complete Feature List of the Keystroke Dynamics + Replay + Evaluation System**

## **A. Core Data Capture Features**

1. **Capture every keystroke** with high-resolution timestamps (keydown events).
2. **Record inter-key latency** (time difference between consecutive keystrokes).
3. **Capture special keys** such as Backspace, Enter, Delete, Arrow keys as `RAW_SPECIAL`.
4. **Capture paste events** with full pasted content as `RAW_PASTE`.
5. **Capture mouse interactions** affecting cursor/selection:
   * Cursor repositioning
   * Text selection ranges
     Logged as `SELECTION_CHANGE`.

6. **Capture selection changes caused by mouse click or drag** (`selectionStart`, `selectionEnd`).
7. **Track absolute start time** when the student begins typing an answer.
8. **Track absolute end time** at submission click.
9. **Capture final text content** of each answer.
10. **Capture all non-tracked fields** (name, roll number, etc.) as simple metadata.

---

## **B. Compression & Processing Features**

11. **Detect segments of consistent typing speed** using standard deviation threshold.
12. **Compress consistent segments** into `COMPRESSED` events:
    * `string` (full typed segment)
    * `interval_ms` (mean latency)
13. **Store variable-timing events as raw events** (`RAW_KEY`, `RAW_SPECIAL`, `RAW_PASTE`).
14. **Preserve selection changes as non-compressible critical events**.
15. **Client-side compression** before sending payload to backend.
16. **Construct a single JSON payload** combining metadata + compressed logs + timestamps.

---

## **C. Submission & Transmission Features**

17. **Disable submit button and show “Submitting…” UI feedback** during compression/transmission.
18. **Transmit a unified payload** containing:
    * exam info
    * metadata
    * answers
    * event logs
    * start/end times

---

## **D. Storage & Backend Features**

19. **Store the complete JSON payload**, including compressed and raw events.
20. **Compute “PASTE used?” flag** for evaluator convenience.
21. **Offer an API endpoint** for Evaluation Interface to fetch submissions.

---

## **E. Replay System Features**

22. **Accurate step-by-step replay** of the student’s entire typing process.
    * Keystrokes
    * Backspaces
    * Cursor moves
    * Selections
    * Paste events
23. **Interpret compressed vs raw events** during replay.
24. **Correctly apply cursor and selection changes** before applying next action.
25. **Real-time timing reconstruction** using `latency_ms` and `interval_ms`.
26. **Replay paste events instantly** and optionally highlight them.

---

## **F. Evaluation Interface Features**

27. **Submission dashboard** with paginated, filterable table.
28. **Show all metadata fields** (name, ID, etc.).
29. **Show final answers** for each question.
30. **Show “Paste Used?” indicator**.
31. **Evaluator mark field** with options:
    * COPIED
    * WRONG
    * CORRECT
    * OK
32. **Replay button for each answer** launching replay module.
33. **Replay module UI** with:
    * Real-time simulation
    * Speed control (0.5x to 2x)
    * Time scrubber
34. **Optional event timeline** showing paste events, selection changes, long pauses.

---

## **G. Additional Behaviour & Safeguards**

35.  **Keystrokes immediately after a paste** are logged seperately (just take care, the replay works).
36. **Non-tracked fields are never keystroke-logged**—only final values stored.

---

## **H. Compression Algorithm Discussion & Evolution**

### **Current Implementation: Standard Deviation Method**

The initial compression algorithm uses **standard deviation** to detect consistent typing segments:

- Calculate standard deviation of inter-key intervals
- Compress when: `stddev ≤ 30ms AND length ≥ 3`
- Works well but computationally intensive

**Configuration:**
- `THRESHOLD_STDDEV = 30` (maximum standard deviation in ms)
- `MIN_SEGMENT_LENGTH = 3` (minimum characters to compress)

---

### **Proposed Enhancement: Threshold-Based Compression**

**Goal:** Simplify algorithm for better maintainability while achieving same compression benefits with clearer visibility of thinking/copying pauses.

#### **New Algorithm**

**Principle:** Compress consecutive keys where inter-key interval is below a maximum threshold (normal typing speed).

**Steps:**
1. Start from first key event
2. Look ahead at next key
3. If `interval < THRESHOLD_MAX_INTERVAL_MS` → add to current segment
4. If `interval >= THRESHOLD_MAX_INTERVAL_MS` → break segment, start new one
5. When segment completes (minimum 3 characters):
   - Calculate **average interval** of all keys in segment
   - Store as `COMPRESSED` event with `interval_ms = average`
6. First segment has `latency_ms = 0`
7. Subsequent segments have `latency_ms = pause since last event`

**Configuration:**
- `THRESHOLD_MAX_INTERVAL_MS` - configurable threshold for experimentation
- `MIN_SEGMENT_LENGTH = 3` - same as current

**Implementation Note:** Keep old standard deviation method in code (commented out) for reference.

---

#### **Special Keys Handling in Compression**

**Can be included in COMPRESSED segments:**
- **Backspace** → represented as `\b` in string
- **Enter** → represented as `\n` in string
- **Delete** → represented as `\x7F` in string
  - Note: Backspace removes character **before** cursor
  - Delete removes character **after** cursor

**Always kept as RAW_SPECIAL (never compressed):**
- **Arrow keys** (ArrowLeft, ArrowRight, ArrowUp, ArrowDown)
- **Selection events** (mouse clicks to select text)

---

#### **Benefits of Threshold-Based Approach**

1. **Simpler computation** - no standard deviation calculation needed
2. **Better maintainability** - easier to understand and debug
3. **Clear pause visibility** - thinking/copying pauses clearly visible in `latency_ms` of next segment
4. **Good compression** - consecutive fast typing still compressed effectively
5. **Configurable** - can experiment with different thresholds to optimize

---

#### **Example Comparison**

**Student typing behavior:**
- Types "print" fast (~120ms between keys)
- Pauses 5000ms (thinking or copying)
- Types "hello" fast (~115ms between keys)

**Compressed Output:**
```json
[
  {
    "type": "COMPRESSED",
    "string": "print",
    "latency_ms": 0,
    "interval_ms": 120    // average interval
  },
  {
    "type": "COMPRESSED",
    "string": "hello",
    "latency_ms": 5000,   // ← Thinking pause clearly visible!
    "interval_ms": 115
  }
]
```

**With Backspace/Enter:**
```json
{
  "type": "COMPRESSED",
  "string": "hello\b\blo\nworld",  // hello, backspace twice, lo, enter, world
  "latency_ms": 0,
  "interval_ms": 125
}
```

---

#### **JSON Schema**

No changes to schema - maintains backward compatibility:

```json
{
  "type": "COMPRESSED",
  "string": "text with \b and \n and \x7F",
  "latency_ms": 0,
  "interval_ms": 120.5
}
```

---

**Status:** Proposed enhancement - implementation pending after threshold experimentation.
