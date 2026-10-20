<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type ConfigTemplate } from '$lib/api/client';
	import { formatDate, getFilename } from '$lib/utils/formatters';
	import PageHeader from '$lib/components/Common/PageHeader.svelte';
	import EmptyState from '$lib/components/Common/EmptyState.svelte';
	import TableSkeleton from '$lib/components/Common/TableSkeleton.svelte';
	import Button from '$lib/components/Common/Button.svelte';
	import Modal from '$lib/components/Common/Modal.svelte';
	import Icon from '$lib/components/Common/Icon.svelte';

	// SvelteKit passes params to all route components
	export let params: Record<string, string> = {};

	let configs: ConfigTemplate[] = [];
	let loading = true;
	let error = '';
	let deleteTarget: ConfigTemplate | null = null;
	let deleting = false;

	onMount(() => {
		loadConfigs();
	});

	async function loadConfigs() {
		loading = true;
		error = '';
		try {
			configs = await api.listConfigs();
		} catch (e) {
			error = (e as Error).message || 'Failed to load templates.';
			configs = [];
		} finally {
			loading = false;
		}
	}

	async function confirmDelete() {
		if (!deleteTarget) return;
		deleting = true;
		try {
			await api.deleteConfig(deleteTarget.id);
			configs = configs.filter((c) => c.id !== deleteTarget!.id);
			deleteTarget = null;
		} catch (e) {
			error = `Failed to delete template: ${(e as Error).message}`;
		}
		deleting = false;
	}
</script>

<svelte:head>
	<title>Templates · Agni</title>
</svelte:head>

<div class="space-y-6">
	<PageHeader
		title="Templates"
		description="Saved VM configurations to reuse when creating a VM."
	/>

	<div class="card">
		{#if loading}
			<TableSkeleton rows={4} />
		{:else if error && configs.length === 0}
			<EmptyState
				icon="ph-plug-charging"
				title="Failed to load templates"
				description={error}
				actionLabel="Try again"
				onAction={loadConfigs}
			/>
		{:else if configs.length === 0}
			<EmptyState
				icon="ph-files"
				title="No templates"
				description="Saved VM configurations appear here and can be reused when creating a VM."
			/>
		{:else}
			{#if error}
				<div class="p-4 pb-0 text-sm text-status-warn-text">{error}</div>
			{/if}
			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="border-b text-left">
							<th class="py-2.5 px-4 text-xs font-medium text-muted-foreground" scope="col">Name</th
							>
							<th class="py-2.5 px-4 text-xs font-medium text-muted-foreground" scope="col"
								>Description</th
							>
							<th
								class="py-2.5 px-4 text-xs font-medium text-muted-foreground text-right"
								scope="col">vCPUs</th
							>
							<th
								class="py-2.5 px-4 text-xs font-medium text-muted-foreground text-right"
								scope="col">Memory</th
							>
							<th class="py-2.5 px-4 text-xs font-medium text-muted-foreground" scope="col"
								>Kernel</th
							>
							<th class="py-2.5 px-4 text-xs font-medium text-muted-foreground" scope="col"
								>Updated</th
							>
							<th class="py-2.5 px-4 w-0" scope="col"><span class="sr-only">Actions</span></th>
						</tr>
					</thead>
					<tbody>
						{#each configs as config (config.id)}
							<tr class="border-b last:border-b-0 hover:bg-accent/50 transition-colors">
								<td class="py-2.5 px-4">
									<div class="font-medium text-foreground">{config.name}</div>
									<div class="text-xs text-muted-foreground font-mono" title={config.id}>
										{config.id.slice(0, 8)}
									</div>
								</td>
								<td class="py-2.5 px-4 text-muted-foreground max-w-[240px]">
									{#if config.description}
										<span class="block truncate" title={config.description}
											>{config.description}</span
										>
									{:else}
										<span class="text-muted-foreground/60">-</span>
									{/if}
								</td>
								<td class="py-2.5 px-4 text-right tabular-nums text-muted-foreground"
									>{config.config.cpus}</td
								>
								<td class="py-2.5 px-4 text-right tabular-nums text-muted-foreground"
									>{config.config.memory_mb} MB</td
								>
								<td class="py-2.5 px-4 max-w-[180px]">
									<span
										class="font-mono text-[13px] text-muted-foreground block truncate"
										title={config.config.kernel_path}
									>
										{getFilename(config.config.kernel_path) || 'N/A'}
									</span>
								</td>
								<td class="py-2.5 px-4 text-muted-foreground whitespace-nowrap"
									>{formatDate(config.updated_at)}</td
								>
								<td class="py-2.5 px-4">
									<div class="flex items-center justify-end">
										<button
											class="h-8 w-8 inline-flex items-center justify-center rounded-md text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors"
											on:click={() => (deleteTarget = config)}
											aria-label="Delete template {config.name}"
											title="Delete"
										>
											<Icon name="ph-trash" size={16} />
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>

<Modal
	open={deleteTarget !== null}
	alert
	title="Delete {deleteTarget?.name ?? ''}?"
	description="This permanently removes the template. VMs already created from it are not affected. This cannot be undone."
	on:close={() => (deleteTarget = null)}
>
	<div slot="footer">
		<Button variant="ghost" on:click={() => (deleteTarget = null)}>Cancel</Button>
		<Button variant="destructive" loading={deleting} on:click={confirmDelete}
			>Delete template</Button
		>
	</div>
</Modal>
