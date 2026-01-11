<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { VMConfig } from '$lib/api/client';
	import ResourcesTab from './ResourcesTab.svelte';
	import StorageTab from './StorageTab.svelte';
	import NetworkTab from './NetworkTab.svelte';
	import KernelTab from './KernelTab.svelte';
	import SecurityTab from './SecurityTab.svelte';
	import Button from '../Common/Button.svelte';

	export let config: Partial<VMConfig> = {};
	export let loading = false;
	export let submitLabel = 'Create VM';

	let activeTab = 'resources';
	let errors: Record<string, string> = {};

	const dispatch = createEventDispatcher();

	const tabs = [
		{ id: 'resources', label: 'Resources', icon: 'M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z' },
		{ id: 'storage', label: 'Storage', icon: 'M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4' },
		{ id: 'network', label: 'Network', icon: 'M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9' },
		{ id: 'kernel', label: 'Kernel', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z' },
		{ id: 'security', label: 'Security', icon: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z' }
	];

	function handleSubmit() {
		// Basic validation
		errors = {};

		if (!config.name) {
			errors.name = 'VM name is required';
		}
		if (!config.root_drive?.path) {
			errors.root_drive = 'Root drive path is required';
		}

		if (Object.keys(errors).length > 0) {
			return;
		}

		dispatch('submit', config);
	}

	function handleConfigChange(event: CustomEvent<Partial<VMConfig>>) {
		config = { ...config, ...event.detail };
	}
</script>

<form on:submit|preventDefault={handleSubmit} class="space-y-6">
	<!-- VM Name -->
	<div class="card">
		<label for="vm-name" class="label">VM Name</label>
		<input
			type="text"
			id="vm-name"
			bind:value={config.name}
			class="input w-full {errors.name ? 'border-red-500' : ''}"
			placeholder="my-microvm"
			required
		/>
		{#if errors.name}
			<p class="text-sm text-red-400 mt-1">{errors.name}</p>
		{/if}
	</div>

	<!-- Tab Navigation -->
	<div class="border-b border-gray-700">
		<nav class="flex gap-1 -mb-px">
			{#each tabs as tab}
				<button
					type="button"
					on:click={() => (activeTab = tab.id)}
					class="px-4 py-2 text-sm font-medium rounded-t-lg transition-colors flex items-center gap-2
						{activeTab === tab.id
							? 'bg-gray-800 text-white border-b-2 border-blue-500'
							: 'text-gray-400 hover:text-white hover:bg-gray-800/50'}"
				>
					<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={tab.icon} />
					</svg>
					{tab.label}
				</button>
			{/each}
		</nav>
	</div>

	<!-- Tab Content -->
	<div class="card">
		{#if activeTab === 'resources'}
			<ResourcesTab {config} on:change={handleConfigChange} />
		{:else if activeTab === 'storage'}
			<StorageTab {config} {errors} on:change={handleConfigChange} />
		{:else if activeTab === 'network'}
			<NetworkTab {config} on:change={handleConfigChange} />
		{:else if activeTab === 'kernel'}
			<KernelTab {config} on:change={handleConfigChange} />
		{:else if activeTab === 'security'}
			<SecurityTab {config} on:change={handleConfigChange} />
		{/if}
	</div>

	<!-- Actions -->
	<div class="flex gap-4">
		<Button type="submit" variant="primary" fullWidth {loading}>
			{loading ? 'Creating...' : submitLabel}
		</Button>
		<a href="/" class="btn btn-secondary">Cancel</a>
	</div>
</form>
