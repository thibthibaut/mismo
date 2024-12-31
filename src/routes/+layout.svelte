<!-- +layout.svelte -->
<script>
	import '../app.css';
	let { children } = $props();
</script>

<!-- Add font-outfit to the outer div to make Outfit the default font -->
<div class="min-h-screen bg-gradient-to-br from-rose-400 via-violet-400 to-teal-400 font-outfit">
    <div class="container mx-auto px-4 py-8">
        <!-- Add font-barriecito to the title -->
        <h1 class="text-6xl font-bold text-center mb-8 text-white 
                   tracking-wider
                   animate-pulse
                   drop-shadow-[0_0_15px_rgba(255,255,255,0.5)]
                   font-barriecito">
            MISMO
        </h1>

        <div class="max-w-3xl mx-auto 
                    bg-white/10 backdrop-blur-lg 
                    rounded-xl shadow-2xl 
                    p-8 
                    border border-white/20">
            {@render children()}
        </div>
    </div>
</div>
