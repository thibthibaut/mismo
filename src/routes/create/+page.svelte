<script>

    import WaitingRoom from '$lib/WaitingRoom.svelte';

    import { onMount } from 'svelte';
    import { fade } from 'svelte/transition';

    let playerName = '';
    let gameId = '';
    let joinUrl = '';
    let showGameInfo = false;
    let showCreateGame = true;
    $: console.log('showCreateGame:', showCreateGame);
    
    async function createGame() {
        if (!playerName.trim()) {
            alert('Please enter your name first');
            return;
        }

        try {
            const response = await fetch('/create-game', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ playerName })
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || 'Failed to create game.');
            }

            const data = await response.json();
            gameId = data.gameId;
            
            // Construct the join URL (adjust the base URL according to your setup)
            joinUrl = `${window.location.origin}/join/${gameId}`;
            showCreateGame = false;
            showGameInfo = true;

        } catch (error) {
            console.error('Create Game Error:', error);
            alert(`Error: ${error.message}`);
        }
    }

    async function copyToClipboard() {
        try {
            await navigator.clipboard.writeText(joinUrl);
            // alert('Link copied to clipboard!');
        } catch (err) {
            console.error('Failed to copy text: ', err);
        }
    }
</script>

<div class="flex flex-col space-y-4 items-center">

    {#if showCreateGame}
        <input 
            bind:value={playerName}
            name="name" 
            placeholder="Ton Blaze" 
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
            on:click={createGame}
            class="px-6 py-3 
                   bg-gradient-to-r from-violet-600 to-purple-600 
                   hover:from-violet-700 hover:to-purple-700
                   text-white font-medium rounded-lg 
                   shadow-lg hover:shadow-xl
                   transition-all duration-200 
                   min-w-[200px]">
            Créer la partie
        </button>
    {/if}

    {#if showGameInfo && gameId && playerName }
        <div 
            transition:fade
            class="mt-4 p-4 bg-violet-50 rounded-lg border-2 border-violet-200"
        >
            <p class="text-violet-600 mb-2">Game ID: {gameId}</p>
            <div class="flex items-center space-x-2">
                <input 
                    readonly 
                    value={joinUrl}
                    class="px-4 py-2 bg-white rounded border border-violet-200 text-violet-600"
                />
                <button 
                    on:click={copyToClipboard}
                    class="p-2 bg-violet-100 hover:bg-violet-200 rounded-lg transition-colors"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-violet-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
                    </svg>
                </button>
            </div>
        </div>

        <WaitingRoom gameID={gameId} playerName={playerName}/>

    {/if}
</div>

