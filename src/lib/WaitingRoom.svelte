<script>
    import Player from '$lib/Player.svelte';
    import { onMount } from 'svelte';

    export let gameID = "";
    export let playerName = "";
    
    // Store game state
    let gameState = {
        players: {},
        round: 0,
        state: "waiting"
    };

    let socket;

    onMount(() => {
        const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
        socket = new WebSocket(`${protocol}://${window.location.host}/ws/game/${gameID}`);

        socket.onopen = () => {
            socket.send(JSON.stringify({ type: "join", name: playerName }));
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            gameState = data; // Update the game state
            console.log('Game State Updated:', gameState);
        };

        socket.onclose = () => {
            console.log('WebSocket connection closed.');
            alert('Connection closed.');
        };

        socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        return () => {
            if (socket) {
                socket.close();
            }
        };
    });

    // Helper function to get players array from object
    $: players = Object.values(gameState.players);
</script>

<div class="container mx-auto">
    <h2 class="text-xl font-bold mb-4">Salle d'attente</h2>
    
    {#if players.length === 0}
        <p class="text-gray-500">En attente de joueurs...</p>
    {:else}
        <div class="space-y-2">
            {#each players as player (player.id)}
                <Player 
                    name={player.name}
                    lives={player.lives}
                    isHost={player.isHost}
                    hasPlayed={player.hasPlayed}
                    number={player.number}
                />
            {/each}
        </div>
        
        <div class="mt-4 text-sm text-gray-600">
            {players.length} joueur{players.length > 1 ? 's' : ''} connecté{players.length > 1 ? 's' : ''}
        </div>
    {/if}
</div>
