<script lang="ts">
	import RiveAnimation from '$lib/components/RiveAnimation.svelte';
</script>

<svelte:head>
	<title>404 - Page Not Found | Agni</title>
</svelte:head>

<div class="min-h-screen bg-gray-900 flex items-center justify-center relative overflow-hidden">
	<!-- Fire gradient overlay -->
	<div class="absolute inset-0 bg-fire-pattern"></div>
	
	<div class="text-center relative z-10">
		<!-- Fire Skull Animation -->
		<div class="flex justify-center mb-6">
			<RiveAnimation 
				src="/fire-skull.riv" 
				width={200} 
				height={200}
			/>
		</div>

		<h1 class="text-8xl font-bold text-fire mb-4">404</h1>
		<h2 class="text-2xl font-semibold text-gray-300 mb-2">Page Not Found</h2>
		<p class="text-gray-500 mb-8 max-w-md mx-auto">
			The page you're looking for has gone up in flames. 
			It might have been moved or doesn't exist.
		</p>

		<a 
			href="/" 
			class="btn btn-primary px-8 py-3 text-lg inline-flex items-center gap-2"
		>
			<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
				<path fill-rule="evenodd" d="M9.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L7.414 9H15a1 1 0 110 2H7.414l2.293 2.293a1 1 0 010 1.414z" clip-rule="evenodd" />
			</svg>
			Return Home
		</a>
	</div>
</div>
