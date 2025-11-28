// ============================================
// CONSTANTS
// ============================================

// Compression algorithm parameters
const THRESHOLD_STDDEV = 300000  // Maximum standard deviation (ms) for compression
const MIN_SEGMENT_LENGTH = 3  // Minimum characters required for compression

// Default exam configuration
const DEFAULT_EXAM_ID = "EXAM-DEMO-001"

// ============================================
// HELPER FUNCTIONS
// ============================================

// Helper function: Generate UUID with fallback for older browsers
function generateUUID() {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    return crypto.randomUUID()
  }
  // Fallback UUID v4 generator
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

// Helper function: Calculate mean of array
function mean(values) {
  if (values.length === 0) return 0
  return values.reduce((sum, val) => sum + val, 0) / values.length
}

// Helper function: Calculate standard deviation
function calculateStdDev(values) {
  if (values.length === 0) return 0

  const avg = mean(values)
  const squareDiffs = values.map(value => Math.pow(value - avg, 2))
  const avgSquareDiff = mean(squareDiffs)
  return Math.sqrt(avgSquareDiff)
}

// Calculate end time from last event
function calculateEndTime(rawEvents, startTime_ms) {
  if (rawEvents.length === 0) {
    return startTime_ms || performance.now()
  }

  // Return timestamp of last event
  const lastEvent = rawEvents[rawEvents.length - 1]
  return lastEvent.timestamp
}

// Extract and analyze segment for compression
function extractSegment(rawEvents, startIdx) {
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

// Compress events based on timing consistency
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
          latency_ms: latency_ms,
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
    else {
      // Unknown event type - log warning and skip
      console.warn('Unknown event type encountered:', event.type, event)
      i++
    }
  }

  return compressed
}

// Main function: Process and pack data into final JSON
function processAndPack(data) {
  // Validate input
  if (!data) {
    throw new Error('No data provided to processAndPack')
  }

  if (!Array.isArray(data.rawEvents)) {
    throw new Error('rawEvents must be an array')
  }

  if (!data.metadata || typeof data.metadata !== 'object') {
    throw new Error('metadata is required and must be an object')
  }

  if (!data.metadata.studentName || !data.metadata.studentName.trim()) {
    throw new Error('Student name is required in metadata')
  }

  if (data.finalAnswer === undefined || data.finalAnswer === null) {
    throw new Error('finalAnswer is required')
  }

  // 1. Calculate endTime_ms from last event or now
  const endTime_ms = calculateEndTime(data.rawEvents, data.startTime_ms)

  // 2. Compress event log
  const compressedLog = compressEvents(data.rawEvents)

  // 3. Build final JSON
  return {
    examId: DEFAULT_EXAM_ID,
    studentId: generateUUID(),
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
