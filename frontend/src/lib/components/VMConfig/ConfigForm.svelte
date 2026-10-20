<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { VMConfig } from '$lib/api/client';
	import ResourcesTab from './ResourcesTab.svelte';
	import StorageTab from './StorageTab.svelte';
	import NetworkTab from './NetworkTab.svelte';
	import KernelTab from './KernelTab.svelte';
	import SecurityTab from './SecurityTab.svelte';
	import Button from '../Common/Button.svelte';
	import Icon from '../Common/Icon.svelte';

	export let config: Partial<VMConfig> = {};
	export let loading = false;
	export let submitLabel = 'Create VM';

	let activeTab = 'resources';
	let errors: Record<string, string> = {};

	const dispatch = createEventDispatcher();

	const tabs = [
		{ id: 'resources', label: 'Resources', icon: 'ph-cpu' },
		{ id: 'storage', label: 'Storage', icon: 'ph-hard-drives' },
		{ id: 'network', label: 'Network', icon: 'ph-wifi-high' },
		{ id: 'kernel', label: 'Kernel', icon: 'ph-code' },
		{ id: 'security', label: 'Security', icon: 'ph-shield-check' }
	];

	function handleSubmit() {
		errors = {};

		if (!config.name) {
			errors.name = 'VM name is required.';
		}
		if (!config.root_drive?.path) {
			errors.root_drive = 'Root drive path is required.';
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

<form on:submit|preventDefault={handleSubmit} class="space-y-6" novalidate>
	<div class="card p-6">
		<label for="vm-name" class="label">VM name</label>
		<input
			type="text"
			id="vm-name"
			bind:value={config.name}
			class="input w-full {errors.name ? 'border-destructive focus:ring-destructive/30' : ''}"
			placeholder="my-microvm"
			required
		/>
		{#if errors.name}
			<p class="text-[13px] text-destructive mt-1.5">{errors.name}</p>
		{/if}
	</div>

	<div class="border-b">
		<nav class="flex gap-1 -mb-px" aria-label="Configuration sections">
			{#each tabs as tab}
				<button
					type="button"
					on:click={() => (activeTab = tab.id)}
					class="px-4 py-2 text-sm font-medium rounded-t-md transition-colors inline-flex items-center gap-2 border-b-2
						{activeTab === tab.id
						? 'border-primary text-primary-accent'
						: 'border-transparent text-muted-foreground hover:text-foreground'}"
					aria-current={activeTab === tab.id ? 'true' : undefined}
				>
					<Icon name={tab.icon} size={16} />
					{tab.label}
				</button>
			{/each}
		</nav>
	</div>

	<div class="card p-6">
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

	<div class="flex justify-end gap-2">
		<Button variant="ghost" href="/">Cancel</Button>
		<Button type="submit" variant="primary" {loading}>{submitLabel}</Button>
	</div>
</form>
