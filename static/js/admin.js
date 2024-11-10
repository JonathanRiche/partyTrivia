function updatePlayerList(playerListHtml) {
	const playerList = document.getElementById('playerList');
	if (playerList) {
		// Create a temporary container
		const temp = document.createElement('div');
		temp.innerHTML = playerListHtml;

		// Find new and existing players
		const newList = temp.querySelector('.space-y-2');
		const currentList = playerList.querySelector('.space-y-2');

		if (newList && currentList) {
			// Animate new players
			const newPlayers = Array.from(newList.children);
			const currentPlayers = Array.from(currentList.children);

			newPlayers.forEach(player => {
				if (!currentPlayers.find(p => p.querySelector('p').textContent === player.querySelector('p').textContent)) {
					player.classList.add('animate-slide-in');
				}
			});
		}

		playerList.innerHTML = playerListHtml;
	}
}
document.addEventListener('htmx:afterOnLoad', function() {
	if (window.gameSocket) {
		return;
	}

	// Connect to WebSocket for real-time updates
	const socket = new WebSocket(`ws://${window.location.host}/ws/admin`);
	window.gameSocket = socket;
	console.log('Connected to admin WebSocket');
	// Update the WebSocket message handler
	socket.onmessage = function(event) {
		const data = JSON.parse(event.data);
		console.log('Received admin message:', data);

		switch (data.type) {
			case 'gameStatus':
				htmx.ajax('GET', '/admin/game/status', { target: '#gameStatus' });
				break;
			case 'playerAnswered':
				console.log('Player answered:', data.payload);
				break;
			case 'playerList':
				// Only update player list if game is active
				if (data.payload.isActive) {
					htmx.ajax('GET', `/admin/game/players?gameID=${data.payload.gameId}`, {
						target: '#playerList',
						swap: 'innerHTML',
						afterSwap: () => {
							updatePlayerList(document.getElementById('playerList').innerHTML);
						}
					});
				} else {
					updatePlayerList(document.getElementById('playerList').innerHTML);
				}
				break;
		}
	};





	socket.onclose = function() {
		console.log('WebSocket connection closed');
		window.gameSocket = null;
	};
	const startButton = document.getElementById('startButton');
	const endButton = document.getElementById('endButton');
	/** @type {HTMLSelectElement} */
	const gameIDSelect = document.getElementById('gameIDSelect');
	function updateStartButton() {
		const gameID = gameIDSelect.value;
		if (!gameID || gameID === '') {
			console.log('No game selected');
			return;
		}
		startButton.setAttribute('hx-vals', `{"gameID": "${gameID}"}`);
		endButton.setAttribute('hx-vals', `{"gameID": "${gameID}"}`);
	}
	gameIDSelect.addEventListener('change', updateStartButton);


});

document.addEventListener('showMessage', (evt) => {
	alert(evt.detail.value);
});

