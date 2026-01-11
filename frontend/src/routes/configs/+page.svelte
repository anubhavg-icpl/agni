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
	<title>Cookie Cutters | Agni</title>
</svelte:head>

<div class="space-y-4 sm:space-y-6">
	<!-- Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
		<div>
			<h1 class="text-xl sm:text-2xl font-bold">Cookie Cutters</h1>
			<p class="text-gray-500 text-xs sm:text-sm">Reusable configs for the lazy (that's a compliment)</p>
		</div>
	</div>

	<!-- Content -->
	{#if loading}
		<div class="text-center py-16">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-orange-500 border-t-transparent"></div>
			<p class="text-gray-400 mt-4">Rummaging through the drawer...</p>
		</div>
	{:else if configs.length === 0}
		<div class="card text-center py-16">
			<div class="text-6xl mb-4">🍪</div>
			<h3 class="text-xl font-medium text-gray-300 mb-2">No Cookie Cutters Yet</h3>
			<p class="text-gray-500 mb-2">Save a VM configuration as a template to reuse it.</p>
			<p class="text-gray-600 text-sm">Because copy-paste is beneath you.</p>
		</div>
	{:else}
		<div class="grid gap-3 sm:gap-4">
			{#each configs as config}
				<div class="card hover:border-orange-500/30 transition-all duration-200">
					<div class="flex flex-col sm:flex-row sm:items-start justify-between gap-3 sm:gap-4">
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2">
								<span class="text-base sm:text-lg">📋</span>
								<h3 class="font-semibold text-base sm:text-lg text-gray-100 truncate">{config.name}</h3>
							</div>
							{#if config.description}
								<p class="text-xs sm:text-sm text-gray-400 mt-1 ml-6 sm:ml-7 line-clamp-2">{config.description}</p>
							{/if}
							<div class="flex flex-wrap gap-2 sm:gap-3 mt-2 sm:mt-3 ml-6 sm:ml-7">
								<span class="px-2 py-0.5 sm:px-2.5 sm:py-1 text-xs bg-gray-700/50 rounded-full text-gray-300">
									🧠 {config.config.cpus} CPUs
								</span>
								<span class="px-2 py-0.5 sm:px-2.5 sm:py-1 text-xs bg-gray-700/50 rounded-full text-gray-300">
									💾 {config.config.memory_mb} MB
								</span>
								<span class="px-2 py-0.5 sm:px-2.5 sm:py-1 text-xs bg-gray-700/50 rounded-full text-gray-300 truncate max-w-[150px] sm:max-w-[200px]" title={config.config.kernel_path}>
									🔧 {config.config.kernel_path.split('/').pop()}
								</span>
							</div>
						</div>
						<button
							on:click={() => handleDelete(config.id, config.name)}
							class="btn btn-danger text-sm py-1.5 px-3 opacity-60 hover:opacity-100 transition-opacity self-end sm:self-start"
							title="Delete template"
						>
							🗑️
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	{#if error}
		<div class="bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg p-4">
			<span class="font-medium">Oops:</span> {error}
		</div>
	{/if}
</div>
