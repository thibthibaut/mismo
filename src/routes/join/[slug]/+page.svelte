<script>
    import { fade } from 'svelte/transition';
    import WaitingRoom from '$lib/WaitingRoom.svelte';
    import { onMount } from 'svelte';
    
    let playerName = "";
    let gameID = "";
    let showJoinForm = true;
    let showWaitingRoom = false;
    let error = "";
    let socket;
    let gameState = null;

    function joinGame() {
        if (!playerName.trim()) {
            error = "Please enter your name";
            return;
        }

        if (!gameID.trim()) {
            error = "Please enter the game ID";
            return;
        }

        showWaitingRoom = true
      }

    //     try {
    //         const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    //         socket = new WebSocket(`${protocol}://${window.location.host}/ws/game/${gameID}`);

    //         socket.onopen = () => {
    //             // Send join message when connection is established
    //             socket.send(JSON.stringify({ type: "join", name: playerName }));
    //         };

    //         socket.onmessage = (event) => {
    //             const data = JSON.parse(event.data);
    //             gameState = data;
    //             console.log('Game State Updated:', gameState);
                
    //             // Hide join form and show waiting room once we receive game state
    //             showJoinForm = false;
    //             showWaitingRoom = true;
    //             error = "";
    //         };

    //         socket.onclose = () => {
    //             console.log('WebSocket connection closed.');
    //             error = 'Connection closed.';
    //             showJoinForm = true;
    //             showWaitingRoom = false;
    //         };

    //         socket.onerror = (error) => {
    //             console.error('WebSocket error:', error);
    //             error = 'Failed to connect to game.';
    //             showJoinForm = true;
    //             showWaitingRoom = false;
    //         };

    //     } catch (err) {
    //         console.error('Join Game Error:', err);
    //         error = err.message;
    //     }
    // }

    // // Cleanup WebSocket connection when component is destroyed
    // onMount(() => {
    //     return () => {
    //         if (socket) {
    //             socket.close();
    //         }
    //     };
    // });
</script>

<div class="flex flex-col space-y-4 items-center">
    {#if showJoinForm}
        <div transition:fade class="flex flex-col space-y-4 items-center">
            {#if error}
                <div class="text-red-500 bg-red-50 p-3 rounded-lg">
                    {error}
                </div>
            {/if}

            <input 
                bind:value={gameID}
                name="gameID" 
                placeholder="Game ID" 
                class="px-6 py-3
                       bg-white
                       border-2 border-violet-200
                       focus:border-violet-400 focus:ring-2 focus:ring-violet-200
                       rounded-lg
                       shadow-sm
                       placeholder-violet-300
                       text-violet-600
                       min-w-[200px]
                       transition-all duration-200
                       outline-none"
            />

            <input 
                bind:value={playerName}
                name="name" 
                placeholder="Ton Prénom" 
                class="px-6 py-3
                       bg-white
                       border-2 border-violet-200
                       focus:border-violet-400 focus:ring-2 focus:ring-violet-200
                       rounded-lg
                       shadow-sm
                       placeholder-violet-300
                       text-violet-600
                       min-w-[200px]
                       transition-all duration-200
                       outline-none"
            />

            <button 
                on:click={joinGame}
                class="px-6 py-3 
                       bg-gradient-to-r from-violet-600 to-purple-600 
                       hover:from-violet-700 hover:to-purple-700
                       text-white font-medium rounded-lg 
                       shadow-lg hover:shadow-xl
                       transition-all duration-200 
                       min-w-[200px]">
                Rejoindre la partie
            </button>
        </div>
    {/if}

    {#if showWaitingRoom}
        <div transition:fade>
            <WaitingRoom 
                {gameID}
                {playerName}
            />
        </div>
    {/if}
</div>
