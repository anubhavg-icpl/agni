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
		if (!confirm(`Delete template "${name}"?`)) return;
		try {
			await api.deleteConfig(id);
			configs = configs.filter(c => c.id !== id);
		} catch (e) {
			error = (e as Error).message;
		}
	}
</script>

<svelte:head>
	<title>Templates - Firectl</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<h1 class="text-2xl font-bold">Configuration Templates</h1>
	</div>

	{#if loading}
		<div class="text-center py-12 text-gray-400">Loading templates...</div>
	{:else if configs.length === 0}
		<div class="card text-center py-12">
			<h3 class="text-lg font-medium text-gray-300 mb-2">No templates yet</h3>
			<p class="text-gray-500">Save a VM configuration as a template to reuse it later</p>
		</div>
	{:else}
		<div class="grid gap-4">
			{#each configs as config}
				<div class="card">
					<div class="flex items-start justify-between">
						<div>
							<h3 class="font-medium text-lg">{config.name}</h3>
							{#if config.description}
								<p class="text-sm text-gray-400 mt-1">{config.description}</p>
							{/if}
							<div class="flex gap-4 mt-2 text-sm text-gray-500">
								<span>{config.config.cpus} CPUs</span>
								<span>{config.config.memory_mb} MB</span>
								<span>{config.config.kernel_path.split('/').pop()}</span>
							</div>
						</div>
						<button
							on:click={() => handleDelete(config.id, config.name)}
							class="text-red-400 hover:text-red-300 text-sm"
						>
							Delete
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	{#if error}
		<div class="bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg p-4">
			{error}
		</div>
	{/if}
</div>
