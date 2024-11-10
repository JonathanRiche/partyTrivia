//@ts-check

/**
 * @typedef {Object} Player
 * @property {string} id - Player's unique identifier
 * @property {string} name - Player's display name
 * @property {number} score - Player's current score
 */

/**
 * @typedef {Object} Question
 * @property {string} id - Question identifier
 * @property {string} text - Question text
 * @property {string[]} options - Available answer options
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
 * @property {Question} payload
 */

/** @typedef {PlayerJoinedMessage | GameStartedMessage | QuestionMessage} GameMessage */


/**
 * Initializes the game client
 */
function main() {
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
 * Handles incoming game messages
 * @param {GameMessage} message
 */
function handleGameMessage(message) {
	console.log('Received message:', message);
	// Handle different message types
	switch (message.type) {
		case 'playerJoined':
			updatePlayersList(message.payload);
			break;
		case 'gameStarted':
			handleGameStart(message.payload);
			break;
		case 'question':
			showQuestion(message.payload);
			break;
	}
}

/**
 * Updates the players list in the UI
 * @param {Object.<string, Player>} players - Map of player IDs to Player objects
 */
function updatePlayersList(players) {
	const playersListElement = document.getElementById('players-list');
	if (playersListElement) {
		console.log('Updating players list:', players.players);
		// Convert the players object to an array and map over it
		playersListElement.innerHTML = Object.values(players).map(player =>
			`<div class="p-2 border-b">${player.name}</div>`
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

/**
 * Displays a question to the user
 * @param {Question} question
 */
function showQuestion(question) {
	console.log('New question:', question);
	// Implement question display logic
}

export { main };
