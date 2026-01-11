<script lang="ts">
	import type { VM } from '$lib/api/client';
	import { vms } from '$lib/stores/vms';

	export let vm: VM;

	const statusColors: Record<string, string> = {
		running: 'status-running',
		stopped: 'status-stopped',
		starting: 'status-starting',
		stopping: 'status-stopping',
		error: 'status-error'
	};

	async function handleStart() {
		await vms.start(vm.id);
	}

	async function handleStop() {
		await vms.stop(vm.id);
	}

	async function handleShutdown() {
		await vms.shutdown(vm.id);
	}

	async function handleDelete() {
		if (confirm(`Delete VM "${vm.name}"? This cannot be undone.`)) {
			await vms.delete(vm.id);
		}
	}
</script>

<div class="card hover:border-gray-600 transition-colors">
	<div class="flex items-start justify-between mb-4">
		<div>
			<h3 class="font-medium text-lg">{vm.name}</h3>
			<p class="text-sm text-gray-500 font-mono">{vm.id.slice(0, 8)}</p>
		</div>
		<span class="px-2 py-1 text-xs rounded border {statusColors[vm.status] || 'status-stopped'}">
			{vm.status}
		</span>
	</div>

	<div class="space-y-2 text-sm text-gray-400 mb-4">
		<div class="flex justify-between">
			<span>CPUs:</span>
			<span class="text-gray-300">{vm.config.cpus}</span>
		</div>
		<div class="flex justify-between">
			<span>Memory:</span>
			<span class="text-gray-300">{vm.config.memory_mb} MB</span>
		</div>
		<div class="flex justify-between">
			<span>Kernel:</span>
			<span class="text-gray-300 truncate max-w-[150px]" title={vm.config.kernel_path}>
				{vm.config.kernel_path.split('/').pop()}
			</span>
		</div>
	</div>

	{#if vm.error}
		<div class="text-sm text-red-400 mb-4 p-2 bg-red-500/10 rounded">
			{vm.error}
		</div>
	{/if}

	<div class="flex gap-2">
		{#if vm.status === 'stopped' || vm.status === 'error'}
			<button on:click={handleStart} class="btn btn-success flex-1 text-sm">
				Start
			</button>
			<button on:click={handleDelete} class="btn btn-danger text-sm">
				Delete
			</button>
		{:else if vm.status === 'running'}
			<button on:click={handleShutdown} class="btn btn-secondary flex-1 text-sm">
				Shutdown
			</button>
			<button on:click={handleStop} class="btn btn-danger text-sm">
				Force Stop
			</button>
		{:else}
			<button disabled class="btn btn-secondary flex-1 text-sm opacity-50">
				{vm.status}...
			</button>
		{/if}
	</div>
</div>
