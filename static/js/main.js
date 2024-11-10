//@ts-check

/**
 * @typedef {Object} Player
 * @property {string} id - Player's unique identifier
 * @property {string} name - Player's display name
 * @property {number} score - Player's current score
 * @property {Object.<number, string>} answers - Player's answers to questions
 */
/**
 * @typedef {Object} Question
 * @property {string} id - Question identifier
 * @property {string} text - Question text
 * @property {string[]} options - Available answer options
 * @property {string} type - Question type ('single' or 'multiple')
 */
/**
 * @typedef {Object} GameData
 * @property {string} gameId - Game identifier
 * @property {string} status - Current game status
 */

/**
 * @typedef {Object} PlayerJoinedMessage
 * @property {'playerJoined'} type
 * @property {Object.<string, Player>} payload
 */

/**
 * @typedef {Object} GameStartedMessage
 * @property {'gameStarted'} type
 * @property {GameData} payload
 */

/**
 * @typedef {Object} QuestionMessage
 * @property {'question'} type
 * @property {Object} payload
 * @property {string} payload.state
 * @property {string} payload.message
 * @property {string} payload.gameId
 * @property {Array<Question>} payload.questions
 */
/**
 * @typedef {Object} GameStateMessage
 * @property {'gameState'} type
 * @property {Object} payload
 * @property {string} payload.state
 * @property {string} payload.message
 * @property {Object.<string, Player>} payload.players
 */

/** @typedef {PlayerJoinedMessage | GameStartedMessage | QuestionMessage | GameStateMessage} GameMessage */



/**
 * Initializes the game client
 */
function main() {
	window.currentQuestion = 0;
	// Check if we're on the game lobby page
	const playersListElement = document.getElementById('players-list');
	if (playersListElement) {
		connectGameWebSocket();
	}
}

/**
 * Establishes WebSocket connection to the game server
 */
function connectGameWebSocket() {
	const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
	const wsURL = `${wsProtocol}//${window.location.host}/ws/game`;

	const socket = new WebSocket(wsURL);

	socket.onopen = () => {
		console.log('Connected to game server');
	};

	/** @param {MessageEvent} event */
	socket.onmessage = (event) => {
		const data = JSON.parse(event.data);
		console.log('Received message:', data);
		handleGameMessage(data);
	};

	/** @param {Event} error */
	socket.onerror = (error) => {
		console.error('WebSocket error:', error);
	};

	socket.onclose = () => {
		console.log('Disconnected from game server');
		// Try to reconnect after 5 seconds
		setTimeout(connectGameWebSocket, 5000);
	};
}

/**
 * Handles game state changes from the server
 * @param {Object} state - The game state object
 * @param {string} state.state - The current game state ('waiting', 'active', 'ended')
 * @param {string} state.message - Status message to display
 */
function handleGameState(state) {
	console.log('Handling game state:', state);

	// Update status message
	const statusElement = document.getElementById('gameStatus');
	if (statusElement) {
		statusElement.textContent = state.message || 'Game status updated';

		// Update styling based on game state
		statusElement.classList.remove('text-blue-600', 'text-green-600', 'text-red-600');
		switch (state.state) {
			case 'active':
				statusElement.classList.add('text-green-600');
				break;
			case 'ended':
				statusElement.classList.add('text-red-600');
				break;
			default:
				statusElement.classList.add('text-blue-600');
		}
	}

	// Handle game container visibility
	const questionContainer = document.getElementById('question-container');
	if (questionContainer) {
		questionContainer.classList.toggle('hidden', state.state !== 'active');
	}

	// Handle additional game state-specific UI updates
	if (state.state === 'active') {
		console.log('Game is active');
		// Request fresh game view when game becomes active
		// htmx.ajax('GET', '/game/view', {
		// 	target: '#game-container',
		// 	swap: 'innerHTML'
		// });
		updateGameStatus(state.state);
		// Update player list if available
		if (state.players) {
			updatePlayersList(state.players);
		}
	} else if (state.state === 'ended') {
		// Handle game end state
		const gameContainer = document.getElementById('game-container');
		if (gameContainer) {
			gameContainer.innerHTML = '<div class="text-center p-4">Game has ended. Thank you for playing!</div>';
		}
	}
}

/**
 * Updates the game status display
 * @param {Object} status
 */
function updateGameStatus(status) {
	const statusElement = document.getElementById('gameStatus');
	if (statusElement) {
		statusElement.textContent = status.message;
		if (status.state === 'active') {
			statusElement.classList.remove('text-blue-600');
			statusElement.classList.add('text-green-600');
		}
	}

	// Show/hide question container based on game state
	const questionContainer = document.getElementById('question-container');
	if (questionContainer) {
		questionContainer.classList.toggle('hidden', status.state !== 'active');
	}
}

/**
 * Handles incoming game messages
 * @param {GameMessage} message
 */
function handleGameMessage(message) {
	console.log('Received message:', message);
	switch (message.type) {
		case 'playerJoined':
			updatePlayersList(message.payload);
			break;
		case 'gameStarted':
			handleGameStart(message.payload);
			break;
		case 'gameState':
			handleGameState(message.payload);
			break;
		case 'question':
			if (window.currentQuestion) {
				showQuestion(message.payload.questions, window.currentQuestion, message.payload.gameId);

				window.currentQuestion = window.currentQuestion + 1;
			} else {
				showQuestion(message.payload.questions, 0, message.payload.gameId);
				window.currentQuestion = 1;
			}

			break;
	}
}


/**
 * Updates the players list in the UI
 * @param {Object.<string, Player>} players
 */
function updatePlayersList(players) {
	console.log('Updating players list:', players);
	const playersListElement = document.getElementById('players-list');
	if (playersListElement) {
		playersListElement.innerHTML = Object.values(players).map(player =>
			`<div class="p-2 border-b">${player.name} (${player.score})</div>`
		).join('');
	}
}

/**
 * Handles game start event
 * @param {GameData} gameData
 */
function handleGameStart(gameData) {
	console.log('Game started:', gameData);
	// Implement game start logic
}
// {
// 	ID      int      `json:"id"`
// 	Text    string   `json:"text"`
// 	Options []string `json:"options"`
// 	Correct string   `json:"correct"`
// }

/**
 * Displays a question to the user
 * @param {Array<Question>} questions
 * @param {number} index
 * @param {string} gameID
 */
function showQuestion(questions, index, gameID) {
	const container = document.getElementById('question-container');
	if (!container) return;

	const question = questions[index];
	if (!question) return;

	// Determine if it's multiple choice or single choice
	const isMultiple = question.type === 'multiple';
	const inputType = isMultiple ? 'checkbox' : 'radio';
	const inputName = `question-${question.id}`;

	container.innerHTML = `
        <div class="border p-4 rounded-lg">
            <h3 class="text-lg font-semibold mb-4">${question.text}</h3>
            <form
                hx-post="/game/submit-answer"
                hx-target="#answer-status"
                class="space-y-2"
            >
                <input type="hidden" name="gameID" value="${gameID}">
                <input type="hidden" name="questionID" value="${question.id}">

                ${question.options.map((option, idx) => `
                    <div class="flex items-center p-2 border rounded hover:bg-blue-50 transition-colors">
                        <input
                            type="${inputType}"
                            id="option-${idx}"
                            name="${inputName}"
                            value="${option}"
                            class="mr-2"
                        >
                        <label
                            for="option-${idx}"
                            class="flex-grow cursor-pointer"
                        >
                            ${option}
                        </label>
                    </div>
                `).join('')}

                <button
                    type="submit"
                    class="w-full mt-4 p-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors
                           disabled:bg-gray-300 disabled:cursor-not-allowed"
                    hx-indicator="#submit-indicator"
                >
                    Submit Answer
                </button>
            </form>
            <div id="submit-indicator" class="htmx-indicator">
                <div class="text-center text-gray-600">
                    Submitting answer...
                </div>
            </div>
            <div id="answer-status" class="mt-4 text-center"></div>
        </div>
    `;

	// Add event listener for form submission
	const form = container.querySelector('form');
	if (form) {
		form.addEventListener('submit', function(e) {
			const submitButton = form.querySelector('button[type="submit"]');
			if (submitButton) {
				submitButton.disabled = true;
			}

			// For multiple choice, collect all selected answers
			if (isMultiple) {
				e.preventDefault();
				const selectedOptions = Array.from(form.querySelectorAll(`input[name="${inputName}"]:checked`))
					.map(input => input.value);

				// Use HTMX to submit the form with multiple answers
				htmx.ajax('POST', '/game/submit-answer', {
					target: '#answer-status',
					values: {
						gameID: gameID,
						questionID: question.id,
						answer: selectedOptions.join(',')
					}
				});
			}
		});
	}

	container.classList.remove('hidden');
}


export { main };
