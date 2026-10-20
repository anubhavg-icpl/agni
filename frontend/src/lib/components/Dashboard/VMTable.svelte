<script lang="ts">
	import type { VM } from '$lib/api/client';
	import { vms as vmStore } from '$lib/stores/vms';
	import { formatDate, getFilename } from '$lib/utils/formatters';
	import StatusBadge from '../Common/StatusBadge.svelte';
	import EmptyState from '../Common/EmptyState.svelte';
	import TableSkeleton from '../Common/TableSkeleton.svelte';
	import StatusBanner from '../Common/StatusBanner.svelte';
	import Button from '../Common/Button.svelte';
	import Modal from '../Common/Modal.svelte';
	import Icon from '../Common/Icon.svelte';

	export let vms: VM[] = [];
	export let loading = false;
	export let error: string | null = null;
	export let onRetry: (() => void) | undefined = undefined;

	let deleteTarget: VM | null = null;
	let deleting = false;

	async function handleStart(vm: VM) {
		await vmStore.start(vm.id);
	}

	async function handleShutdown(vm: VM) {
		await vmStore.shutdown(vm.id);
	}

	async function handleStop(vm: VM) {
		await vmStore.stop(vm.id);
	}

	async function confirmDelete() {
		if (!deleteTarget) return;
		deleting = true;
		await vmStore.delete(deleteTarget.id);
		deleting = false;
		deleteTarget = null;
	}
</script>

<div class="card">
	{#if loading && vms.length === 0}
		<TableSkeleton rows={5} />
	{:else if error && vms.length === 0}
		<EmptyState
			icon="ph-plug-charging"
			title="Failed to load virtual machines"
			description={error}
			actionLabel="Try again"
			onAction={onRetry}
		/>
	{:else if vms.length === 0}
		<EmptyState
			icon="ph-cube-transparent"
			title="No virtual machines"
			description="Create a VM to launch a Firecracker microVM on this machine."
			actionLabel="Create VM"
			actionHref="/vms/new"
		/>
	{:else}
		{#if error}
			<div class="p-4 pb-0">
				<StatusBanner variant="warning" message="Showing the last loaded VMs. {error}" />
			</div>
		{/if}
		<div class="overflow-x-auto">
			<table class="w-full text-sm table-fixed">
				<colgroup>
					<col class="w-[26%]" />
					<col class="w-[11%]" />
					<col class="w-[7%]" />
					<col class="w-[10%]" />
					<col class="w-[22%]" />
					<col class="w-[13%]" />
					<col class="w-[11%]" />
				</colgroup>
				<thead class="sticky top-0 bg-card">
					<tr class="border-b text-left">
						<th class="th" scope="col">Name</th>
						<th class="th" scope="col">Status</th>
						<th class="th text-right" scope="col">vCPUs</th>
						<th class="th text-right" scope="col">Memory</th>
						<th class="th" scope="col">Kernel</th>
						<th class="th" scope="col">Created</th>
						<th class="th w-0" scope="col"><span class="sr-only">Actions</span></th>
					</tr>
				</thead>
				<tbody>
					{#each vms as vm (vm.id)}
						<tr class="border-b last:border-b-0 hover:bg-accent/50 transition-colors">
							<td class="td align-middle">
								<div class="min-w-0">
									<a
										href="/vms/{vm.id}"
										class="font-medium text-foreground hover:underline leading-tight block truncate"
										title={vm.name}
									>
										{vm.name}
									</a>
									<span class="text-xs text-muted-foreground font-mono" title={vm.id}>
										{vm.id.slice(0, 8)}
									</span>
								</div>
							</td>
							<td class="td align-middle"><StatusBadge status={vm.status} /></td>
							<td class="td align-middle text-right tabular-nums text-muted-foreground">
								{vm.config.cpus}
							</td>
							<td class="td align-middle text-right tabular-nums text-muted-foreground whitespace-nowrap">
								{vm.config.memory_mb} MB
							</td>
							<td class="td align-middle">
								<span
									class="font-mono text-[13px] text-muted-foreground block truncate"
									title={vm.config.kernel_path}
								>
									{getFilename(vm.config.kernel_path) || 'N/A'}
								</span>
							</td>
							<td class="td align-middle text-muted-foreground whitespace-nowrap">
								{formatDate(vm.created_at)}
							</td>
							<td class="td align-middle">
								<div class="flex items-center justify-end gap-1">
									<a
										href="/vms/{vm.id}"
										class="h-8 w-8 inline-flex items-center justify-center rounded-md text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
										aria-label="Open console for {vm.name}"
										title="Console"
									>
										<Icon name="ph-terminal-window" size={16} />
									</a>
									{#if vm.status === 'stopped' || vm.status === 'error'}
										<button
											class="icon-btn"
											on:click={() => handleStart(vm)}
											aria-label="Start {vm.name}"
											title="Start"
										>
											<Icon name="ph-play" size={16} />
										</button>
										<button
											class="icon-btn icon-btn-danger"
											on:click={() => (deleteTarget = vm)}
											aria-label="Delete {vm.name}"
											title="Delete"
										>
											<Icon name="ph-trash" size={16} />
										</button>
									{:else if vm.status === 'running'}
										<button
											class="icon-btn"
											on:click={() => handleShutdown(vm)}
											aria-label="Shut down {vm.name}"
											title="Shut down"
										>
											<Icon name="ph-power" size={16} />
										</button>
										<button
											class="icon-btn icon-btn-danger"
											on:click={() => handleStop(vm)}
											aria-label="Force stop {vm.name}"
											title="Force stop"
										>
											<Icon name="ph-stop" size={16} />
										</button>
									{:else}
										<span class="text-muted-foreground inline-flex items-center px-2" title={vm.status}>
											<svg class="animate-spin h-4 w-4" viewBox="0 0 24 24" aria-hidden="true">
												<circle
													class="opacity-25"
													cx="12"
													cy="12"
													r="10"
													stroke="currentColor"
													stroke-width="4"
													fill="none"
												/>
												<path
													class="opacity-75"
													fill="currentColor"
													d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
												/>
											</svg>
										</span>
									{/if}
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<Modal
	open={deleteTarget !== null}
	alert
	title="Delete {deleteTarget?.name ?? ''}?"
	description="This permanently removes the VM and its configuration from Agni. This cannot be undone."
	on:close={() => (deleteTarget = null)}
>
	<div slot="footer">
		<Button variant="ghost" on:click={() => (deleteTarget = null)}>Cancel</Button>
		<Button variant="destructive" loading={deleting} on:click={confirmDelete}>Delete VM</Button>
	</div>
</Modal>

<style>
	.th {
		@apply py-2.5 px-4 text-xs font-medium text-muted-foreground;
	}
	.td {
		@apply py-2.5 px-4;
	}
	.icon-btn {
		@apply h-8 w-8 inline-flex items-center justify-center rounded-md text-muted-foreground hover:text-foreground hover:bg-accent transition-colors;
	}
	.icon-btn-danger {
		@apply hover:text-destructive hover:bg-destructive/10;
	}
</style>
