<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type ConfigTemplate } from '$lib/api/client';

	// SvelteKit passes params to all route components
	export let params: Record<string, string> = {};

	let configs: ConfigTemplate[] = [];
	let loading = true;
	let error = '';

	onMount(async () => {
		await loadConfigs();
	});

	async function loadConfigs() {
		loading = true;
		error = '';
		try {
			configs = await api.listConfigs();
		} catch (e) {
			error = (e as Error).message || 'Failed to load templates';
			configs = [];
		} finally {
			loading = false;
		}
	}

	async function handleDelete(id: string, name: string) {
		if (!confirm(`Delete template "${name}"? This cannot be undone.`)) return;
		try {
			await api.deleteConfig(id);
			configs = configs.filter((c) => c.id !== id);
		} catch (e) {
			error = (e as Error).message;
		}
	}
</script>

<svelte:head>
	<title>Config Templates | Agni</title>
</svelte:head>

<div class="space-y-6 lg:space-y-8">
	<!-- Header -->
	<div class="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
		<div>
			<h1 class="text-2xl sm:text-3xl font-bold text-white mb-1">The Recipe Book</h1>
			<p class="text-gray-400 text-sm sm:text-base">Pre-made suffering configurations. Work smarter, not harder.</p>
		</div>
	</div>

	<!-- Content -->
	{#if loading}
		<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 border border-gray-700/50 rounded-2xl p-12">
			<div class="flex flex-col items-center justify-center text-center">
				<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-orange-500/20 to-red-500/20 flex items-center justify-center mb-4">
					<svg class="w-6 h-6 text-orange-400 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
					</svg>
				</div>
				<p class="text-gray-400">Dusting off the ancient scrolls...</p>
			</div>
		</div>
	{:else if error}
		<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 border border-gray-700/50 rounded-2xl p-12">
			<div class="flex flex-col items-center justify-center text-center max-w-md mx-auto">
				<div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-red-500/20 to-red-600/20 flex items-center justify-center mb-5">
					<svg class="w-8 h-8 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</div>
				<h3 class="text-xl font-semibold text-white mb-2">The Scrolls Caught Fire</h3>
				<p class="text-gray-400 mb-6">{error}</p>
				<button
					on:click={loadConfigs}
					class="inline-flex items-center gap-2 px-5 py-3 bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white font-medium rounded-xl transition-all duration-300 transform hover:scale-[1.02] active:scale-[0.98]"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
					</svg>
					<span>Try Again, Mortal</span>
				</button>
			</div>
		</div>
	{:else if configs.length === 0}
		<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 border border-gray-700/50 rounded-2xl p-12">
			<div class="flex flex-col items-center justify-center text-center max-w-md mx-auto">
				<div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-gray-700/50 to-gray-800/50 flex items-center justify-center mb-5">
					<svg class="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
					</svg>
				</div>
				<h3 class="text-xl font-semibold text-white mb-2">The Recipe Book is Empty</h3>
				<p class="text-gray-400">No evil recipes yet? Save a VM config to reuse your dark creations.</p>
			</div>
		</div>
	{:else}
		<div>
			<div class="flex items-center justify-between mb-4">
				<h2 class="text-lg sm:text-xl font-semibold text-white">Your Evil Recipes</h2>
				<span class="text-sm text-gray-500">{configs.length} {configs.length === 1 ? 'recipe' : 'recipes'} of doom</span>
			</div>
			<div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4">
				{#each configs as config}
					<div class="group relative bg-gradient-to-br from-gray-800/50 to-gray-900/50 border border-gray-700/50 rounded-2xl p-4 sm:p-5 hover:border-orange-500/30 transition-all duration-300 hover:shadow-xl hover:shadow-orange-500/5 overflow-hidden">
						<div class="absolute inset-0 bg-gradient-to-br from-orange-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"></div>

						<div class="relative">
							<div class="flex items-start justify-between gap-3 mb-3">
								<div class="flex items-center gap-3 min-w-0">
									<div class="w-10 h-10 rounded-xl bg-gradient-to-br from-orange-500/20 to-red-500/20 flex items-center justify-center flex-shrink-0">
										<svg class="w-5 h-5 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
										</svg>
									</div>
									<div class="min-w-0">
										<h3 class="font-semibold text-base text-gray-100 truncate group-hover:text-white transition-colors">{config.name}</h3>
										<p class="text-xs text-gray-500 font-mono">{config.id.slice(0, 8)}...</p>
									</div>
								</div>
								<button
									on:click={() => handleDelete(config.id, config.name)}
									class="p-2 rounded-lg bg-red-500/10 text-red-400 hover:bg-red-500/20 transition-all opacity-0 group-hover:opacity-100"
									title="Delete template"
								>
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
									</svg>
								</button>
							</div>

							{#if config.description}
								<p class="text-sm text-gray-400 mb-4 line-clamp-2">{config.description}</p>
							{/if}

							<div class="grid grid-cols-3 gap-2 p-3 bg-black/30 rounded-xl border border-gray-700/30">
								<div class="text-center">
									<div class="flex items-center justify-center gap-1 mb-1">
										<svg class="w-3.5 h-3.5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
										</svg>
										<span class="text-xs text-gray-500">vCPU</span>
									</div>
									<div class="text-sm font-semibold text-gray-200">{config.config.cpus}</div>
								</div>
								<div class="text-center border-x border-gray-700/30">
									<div class="flex items-center justify-center gap-1 mb-1">
										<svg class="w-3.5 h-3.5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
										</svg>
										<span class="text-xs text-gray-500">RAM</span>
									</div>
									<div class="text-sm font-semibold text-gray-200">{config.config.memory_mb} MB</div>
								</div>
								<div class="text-center">
									<div class="flex items-center justify-center gap-1 mb-1">
										<svg class="w-3.5 h-3.5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
										</svg>
										<span class="text-xs text-gray-500">Kernel</span>
									</div>
									<div class="text-sm font-semibold text-gray-200 truncate" title={config.config.kernel_path}>
										{config.config.kernel_path.split('/').pop()?.slice(0, 8) || 'N/A'}
									</div>
								</div>
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</div>
