<script>
    export let name;
    export let lives = 7;
    export let isHost = false;
    export let hasPlayed = false;
    export let number = 0;
</script>

<div class="p-4 bg-white rounded-lg shadow-sm border border-gray-200">
    <div class="flex items-center justify-between">
        <div class="flex items-center space-x-2">
            <span class="font-medium">{name}</span>
            {#if isHost}
                <span class="text-xs bg-purple-100 text-purple-800 px-2 py-1 rounded">Host</span>
            {/if}
        </div>
        <div class="flex items-center space-x-2">
            {#if hasPlayed}
                <span class="text-green-500">✓</span>
            {/if}
            <span class="text-sm text-gray-600">♥ {lives}</span>
        </div>
    </div>
    {#if number !== 0}
        <div class="mt-2 text-sm text-gray-600">
            Numéro: {number}
        </div>
    {/if}
</div>
