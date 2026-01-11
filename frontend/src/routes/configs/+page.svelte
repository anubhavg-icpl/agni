<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type ConfigTemplate } from '$lib/api/client';

	let configs: ConfigTemplate[] = [];
	let loading = true;
	let error = '';

	onMount(async () => {
		try {
			configs = await api.listConfigs();
		} catch (e) {
			error = (e as Error).message;
		}
		loading = false;
	});

	async function handleDelete(id: string, name: string) {
		if (!confirm(`Shred template "${name}"? Gone forever. No backsies.`)) return;
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

<div class="space-y-4 sm:space-y-6">
	<!-- Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
		<div>
			<h1 class="text-xl sm:text-2xl font-bold">Config Templates</h1>
			<p class="text-gray-500 text-xs sm:text-sm">Reusable VM configurations</p>
		</div>
	</div>

	<!-- Content -->
	{#if loading}
		<div class="text-center py-16 flex flex-col items-center gap-3">
			<svg class="w-8 h-8 animate-spin text-orange-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
			</svg>
			<p class="text-gray-400">Loading templates...</p>
		</div>
	{:else if configs.length === 0}
		<div class="card text-center py-12 sm:py-16">
			<div class="w-20 h-20 mx-auto mb-4 rounded-full bg-gray-800 flex items-center justify-center">
				<svg class="w-10 h-10 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
				</svg>
			</div>
			<h3 class="text-lg sm:text-xl font-medium text-gray-300 mb-2">No Templates</h3>
			<p class="text-gray-500 text-sm">Save a VM configuration as a template to reuse it.</p>
		</div>
	{:else}
		<div class="grid gap-3 sm:gap-4">
			{#each configs as config}
				<div class="card hover:border-orange-500/30 transition-all duration-200">
					<div class="flex flex-col sm:flex-row sm:items-start justify-between gap-3 sm:gap-4">
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2">
								<svg class="w-5 h-5 text-orange-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
								</svg>
								<h3 class="font-semibold text-base sm:text-lg text-gray-100 truncate">{config.name}</h3>
							</div>
							{#if config.description}
								<p class="text-xs sm:text-sm text-gray-400 mt-1 ml-7 line-clamp-2">{config.description}</p>
							{/if}
							<div class="flex flex-wrap gap-2 sm:gap-3 mt-2 sm:mt-3 ml-7">
								<span class="px-2 py-0.5 sm:px-2.5 sm:py-1 text-xs bg-gray-700/50 rounded-full text-gray-300">
									{config.config.cpus} vCPU
								</span>
								<span class="px-2 py-0.5 sm:px-2.5 sm:py-1 text-xs bg-gray-700/50 rounded-full text-gray-300">
									{config.config.memory_mb} MB
								</span>
								<span class="px-2 py-0.5 sm:px-2.5 sm:py-1 text-xs bg-gray-700/50 rounded-full text-gray-300 truncate max-w-[150px] sm:max-w-[200px]" title={config.config.kernel_path}>
									{config.config.kernel_path.split('/').pop()}
								</span>
							</div>
						</div>
						<button
							on:click={() => handleDelete(config.id, config.name)}
							class="btn btn-danger text-sm py-1.5 px-3 opacity-60 hover:opacity-100 transition-opacity self-end sm:self-start"
							title="Delete template"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
							</svg>
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	{#if error}
		<div class="bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg p-4 flex items-center gap-2">
			<svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
			</svg>
			<span>{error}</span>
		</div>
	{/if}
</div>
