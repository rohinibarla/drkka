// Global state
let selectedQuestion = null
let isSubmitted = false  // Track if already submitted

// Global state for event capture
const captureData = {
  rawEvents: [],           // Array of captured events
  startTime_ms: null,      // First interaction timestamp
  lastSelection: { start: 0, end: 0 }  // Track last selection state
}

// Load and select random question on page load
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
    // Show error message with retry option
    const questionDiv = document.getElementById('question-text')
    questionDiv.innerHTML = `
      <div class="text-red-600 mb-4">
        ⚠️ Failed to load question: ${error.message}
      </div>
      <button
        onclick="loadRandomQuestion()"
        class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
      >
        🔄 Retry Loading Question
      </button>
    `
  }
}

// Set start time on first interaction
function setStartTime() {
  if (!captureData.startTime_ms) {
    captureData.startTime_ms = performance.now()
  }
}

// Handle first keydown event (sets start time and captures keystroke)
function handleFirstKeydown(e) {
  setStartTime()
  handleKeydown(e)  // Also capture this first keystroke

  // Now attach the ongoing keydown listener
  const answerField = document.getElementById('answer-field')
  answerField.addEventListener('keydown', handleKeydown)
}

// Handle keydown events
function handleKeydown(e) {
  const now = performance.now()
  const key = e.key

  // Ignore keys pressed with modifiers (Ctrl+V, Cmd+C, etc.) - these are shortcuts
  if (e.ctrlKey || e.metaKey || e.altKey) {
    return
  }

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

// Handle paste events
function handlePaste(e) {
  const now = performance.now()
  const content = e.clipboardData.getData('text')

  captureData.rawEvents.push({
    type: 'paste',
    content: content,
    timestamp: now
  })
}

// Handle selection changes (mouse clicks, drags)
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

// Handle form submission
async function handleSubmit(e) {
  e.preventDefault()

  // Prevent double submission
  if (isSubmitted) {
    alert('Already submitted. Refresh the page to submit again.')
    return
  }

  // Validate question is loaded
  if (!selectedQuestion) {
    alert('Please wait for the question to load.')
    return
  }

  const answerField = document.getElementById('answer-field')
  const nameInput = document.getElementById('student-name')

  // Validate student name
  if (!nameInput.value.trim()) {
    alert('Please enter your name.')
    nameInput.focus()
    return
  }

  // Validate answer has been typed
  if (!answerField.value.trim()) {
    alert('Please type an answer before submitting.')
    answerField.focus()
    return
  }

  // Validate has events
  if (captureData.rawEvents.length === 0) {
    alert('No typing activity detected. Please type your answer.')
    answerField.focus()
    return
  }

  const submitBtn = document.getElementById('submit-btn')

  // Disable submit button
  submitBtn.disabled = true
  submitBtn.textContent = 'Submitting...'
  submitBtn.classList.add('opacity-50', 'cursor-not-allowed')

  try {
    // Call process_and_pack.js function
    const payload = processAndPack({
      rawEvents: captureData.rawEvents,
      startTime_ms: captureData.startTime_ms,
      finalAnswer: answerField.value,
      questionIndex: selectedQuestion.index,
      questionTitle: selectedQuestion.title,
      questionText: selectedQuestion.text,
      metadata: {
        studentName: nameInput.value.trim()
      }
    })

    // Send to server
    const response = await fetch('/submit', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    })

    if (!response.ok) {
      throw new Error('Server returned error: ' + response.status)
    }

    const result = await response.json()

    // Mark as submitted
    isSubmitted = true

    // Console output (for debugging)
    // console.log(JSON.stringify(payload, null, 2))

    // Display on page (commented out for production, uncomment for debugging)
    // displayJSON(payload)

    // Update button
    submitBtn.textContent = 'Submitted ✓'

    // Show success message
    alert('Submission successful! Your answer has been recorded.')

  } catch (error) {
    console.error('Error processing submission:', error)
    alert('Error creating submission: ' + error.message + '\n\nPlease try again.')

    // Re-enable button on error
    submitBtn.disabled = false
    submitBtn.textContent = 'Submit'
    submitBtn.classList.remove('opacity-50', 'cursor-not-allowed')
  }
}

// Display JSON on page
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

// Initialize event listeners on page load
window.addEventListener('DOMContentLoaded', () => {
  const answerField = document.getElementById('answer-field')
  const submitBtn = document.getElementById('submit-btn')

  // Track first interaction
  answerField.addEventListener('focus', setStartTime, { once: true })
  answerField.addEventListener('click', setStartTime, { once: true })
  answerField.addEventListener('keydown', handleFirstKeydown, { once: true })

  // Ongoing event capture (keydown is added by handleFirstKeydown)
  answerField.addEventListener('paste', handlePaste)
  answerField.addEventListener('mouseup', handleSelection)
  answerField.addEventListener('select', handleSelection)

  // Submit
  submitBtn.addEventListener('click', handleSubmit)

  // Load question
  loadRandomQuestion()
})

// Warn user before leaving if they have unsaved work
window.addEventListener('beforeunload', (e) => {
  if (captureData.rawEvents.length > 0 && !isSubmitted) {
    e.preventDefault()
    e.returnValue = ''
    return ''
  }
})
