# JSON Schema for Compressed Keystroke Log Submission

This document details the specific JSON schema for the final payload sent from the client (exam environment) to the server upon submission. This structure is designed to be comprehensive, including all required metadata, non-tracked form data, and the compressed, high-fidelity event logs for the two tracked questions.

## 1. Overall Submission Payload Structure

The entire submission is a single JSON object.

| Field Name | Type | Description |
| :--- | :--- | :--- |
| `examId` | `string` | Unique identifier for the exam instance. |
| `studentId` | `string` | Unique identifier for the student. |
| `submissionTime` | `string` | Server-side timestamp of when the submission was received (for final record-keeping). |
| `metadata` | `object` | Contains the final values of all non-tracked form fields. |
| `q1` | `object` | Data structure for Question 1's answer and event log. |
| `q2` | `object` | Data structure for Question 2's answer and event log. |

## 2. Metadata Object Structure

The `metadata` object is a simple key-value map for all non-tracked form fields.

| Field Name | Type | Description |
| :--- | :--- | :--- |
| `[FieldName]` | `string` | The final text content of the non-tracked form field (e.g., `studentName`, `studentIdNumber`). |

## 3. Question Data Structure (`q1` and `q2`)

The structure for each tracked question (`q1` and `q2`) is identical.

| Field Name | Type | Description |
| :--- | :--- | :--- |
| `finalAnswer` | `string` | The final text content of the code-writing field at the moment of submission. |
| `startTime_ms` | `number` | Absolute high-resolution timestamp (e.g., milliseconds since epoch) of the first recorded input event for this question. |
| `endTime_ms` | `number` | Absolute high-resolution timestamp of the form submission click. |
| `eventLog` | `array` | The compressed and latency-based array of all user interaction events. |

## 4. Event Log Structure (`eventLog`)

The `eventLog` is an array of event objects. Each object represents a user action and includes a `type` field to distinguish its structure. The timing for all events is relative, using `latency_ms` (the time in milliseconds since the *previous* event).

### 4.1. Event Types

The log contains five possible event types:

| Event Type | Purpose | Key Fields |
| :--- | :--- | :--- |
| `COMPRESSED` | Represents a segment of consistent, evenly-spaced typing. | `string`, `interval_ms` |
| `RAW_KEY` | Represents a single character typed with variable timing. | `key`, `latency_ms` |
| `RAW_SPECIAL` | Represents a non-character key press (e.g., Backspace, Enter). | `key`, `latency_ms` |
| `RAW_PASTE` | Represents a text paste action. | `content`, `latency_ms` |
| `SELECTION_CHANGE` | Represents a mouse-driven cursor repositioning or text selection. | `start`, `end`, `latency_ms` |

### 4.2. Detailed Event Schemas

#### A. `COMPRESSED` Event (Compressed Typing Segment)

```json
{
  "type": "COMPRESSED",
  "string": "The quick brown fox",
  "interval_ms": 145
}
```

| Field Name | Type | Description |
| :--- | :--- | :--- |
| `type` | `string` | Must be `"COMPRESSED"`. |
| `string` | `string` | The sequence of characters typed in the consistent segment. |
| `interval_ms` | `number` | The mean inter-key latency (in milliseconds) for all characters in the `string`. |

#### B. `RAW_KEY` Event (Single Character Keystroke)

```json
{
  "type": "RAW_KEY",
  "key": "a",
  "latency_ms": 150
}
```

| Field Name | Type | Description |
| :--- | :--- | :--- |
| `type` | `string` | Must be `"RAW_KEY"`. |
| `key` | `string` | The single character typed. |
| `latency_ms` | `number` | Time (in milliseconds) since the previous event in the log. |

#### C. `RAW_SPECIAL` Event (Special Key Press)

```json
{
  "type": "RAW_SPECIAL",
  "key": "Backspace",
  "latency_ms": 300
}
```

| Field Name | Type | Description |
| :--- | :--- | :--- |
| `type` | `string` | Must be `"RAW_SPECIAL"`. |
| `key` | `string` | The name of the special key (e.g., `"Backspace"`, `"Enter"`, `"ArrowLeft"`). |
| `latency_ms` | `number` | Time (in milliseconds) since the previous event in the log. |

#### D. `RAW_PASTE` Event (Paste Action)

```json
{
  "type": "RAW_PASTE",
  "content": "print('Hello')",
  "latency_ms": 1200
}
```

| Field Name | Type | Description |
| :--- | :--- | :--- |
| `type` | `string` | Must be `"RAW_PASTE"`. |
| `content` | `string` | The full text content that was pasted. |
| `latency_ms` | `number` | Time (in milliseconds) since the previous event in the log. |

#### E. `SELECTION_CHANGE` Event (Mouse Interaction)

```json
{
  "type": "SELECTION_CHANGE",
  "start": 4,
  "end": 9,
  "latency_ms": 650
}
```

| Field Name | Type | Description |
| :--- | :--- | :--- |
| `type` | `string` | Must be `"SELECTION_CHANGE"`. |
| `start` | `number` | The new `selectionStart` index (0-based). |
| `end` | `number` | The new `selectionEnd` index (0-based). If `start` equals `end`, it represents a cursor repositioning. |
| `latency_ms` | `number` | Time (in milliseconds) since the previous event in the log. |

## 5. Example Final Payload

```json
{
  "examId": "EXAM-2025-001",
  "studentId": "STU-98765",
  "submissionTime": "2025-11-27T12:30:00.000Z",
  "metadata": {
    "studentName": "Jane Doe",
    "studentIdNumber": "123456789"
  },
  "q1": {
    "finalAnswer": "print('Hello World')",
    "startTime_ms": 1732700000000,
    "endTime_ms": 1732700060000,
    "eventLog": [
      { "type": "RAW_KEY", "key": "p", "latency_ms": 0 },
      { "type": "RAW_KEY", "key": "r", "latency_ms": 120 },
      { "type": "COMPRESSED", "string": "int('Hello", "interval_ms": 150 },
      { "type": "RAW_PASTE", "content": " World", "latency_ms": 2500 },
      { "type": "RAW_KEY", "key": ")", "latency_ms": 100 }
    ]
  },
  "q2": {
    "finalAnswer": "for i in range(10): print(i)",
    "startTime_ms": 1732700065000,
    "endTime_ms": 1732700090000,
    "eventLog": [
      { "type": "RAW_PASTE", "content": "for i in range(10):", "latency_ms": 0 },
      { "type": "SELECTION_CHANGE", "start": 19, "end": 19, "latency_ms": 500 },
      { "type": "COMPRESSED", "string": " print(i)", "interval_ms": 130 }
    ]
  }
}
```
