document.addEventListener('htmx:afterOnLoad', function() {
	if (window.gameSocket) {
		return;
	}

	// Connect to WebSocket for real-time updates
	const socket = new WebSocket(`ws://${window.location.host}/ws/admin`);
	window.gameSocket = socket;
	console.log('Connected to admin WebSocket');
	socket.onmessage = function(event) {
		const data = JSON.parse(event.data);
		console.log('Received message:', data);
		switch (data.type) {
			case 'gameStatus':
				htmx.ajax('GET', '/admin/game/status', { target: '#gameStatus' });
				break;
			case 'playerList':
				htmx.ajax('GET', '/admin/game/players', { target: '#playerList' });
				break;
		}
	};

	socket.onclose = function() {
		console.log('WebSocket connection closed');
		window.gameSocket = null;
	};
});

