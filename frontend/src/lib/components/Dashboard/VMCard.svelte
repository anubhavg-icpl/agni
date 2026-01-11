<script lang="ts">
	import type { VM } from '$lib/api/client';
	import { vms } from '$lib/stores/vms';

	export let vm: VM;

	const statusConfig: Record<string, { class: string; icon: string; label: string }> = {
		running: { class: 'status-running', icon: '🔥', label: 'Alive' },
		stopped: { class: 'status-stopped', icon: '💀', label: 'Dead' },
		starting: { class: 'status-starting', icon: '⚡', label: 'Waking...' },
		stopping: { class: 'status-stopping', icon: '😴', label: 'Dying...' },
		error: { class: 'status-error', icon: '💥', label: 'Broken' }
	};

	$: status = statusConfig[vm.status] || statusConfig.stopped;

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
		if (confirm(`Yeet "${vm.name}" into the void? This is permanent.`)) {
			await vms.delete(vm.id);
		}
	}
</script>

<div class="card hover:border-orange-500/30 transition-all duration-200 hover:shadow-lg hover:shadow-orange-500/5">
	<!-- Header -->
	<div class="flex items-start justify-between mb-4">
		<div class="flex-1 min-w-0">
			<h3 class="font-semibold text-lg text-gray-100 truncate">{vm.name}</h3>
			<p class="text-xs text-gray-500 font-mono mt-0.5">{vm.id.slice(0, 8)}...</p>
		</div>
		<span class="flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium rounded-full {status.class}">
			<span>{status.icon}</span>
			{status.label}
		</span>
	</div>

	<!-- Specs -->
	<div class="grid grid-cols-3 gap-3 mb-4 p-3 bg-black/20 rounded-lg">
		<div class="text-center">
			<div class="text-xs text-gray-500 mb-0.5">Brains</div>
			<div class="text-sm font-medium text-gray-200">{vm.config.cpus} CPU</div>
		</div>
		<div class="text-center border-x border-gray-700/50">
			<div class="text-xs text-gray-500 mb-0.5">Memory</div>
			<div class="text-sm font-medium text-gray-200">{vm.config.memory_mb} MB</div>
		</div>
		<div class="text-center">
			<div class="text-xs text-gray-500 mb-0.5">Kernel</div>
			<div class="text-sm font-medium text-gray-200 truncate" title={vm.config.kernel_path}>
				{vm.config.kernel_path.split('/').pop()?.slice(0, 10) || 'N/A'}
			</div>
		</div>
	</div>

	<!-- Error Message -->
	{#if vm.error}
		<div class="text-sm text-red-400 mb-4 p-2.5 bg-red-500/10 rounded-lg border border-red-500/20">
			<span class="font-medium">Uh oh:</span> {vm.error}
		</div>
	{/if}

	<!-- Actions -->
	<div class="flex gap-2">
		{#if vm.status === 'stopped' || vm.status === 'error'}
			<button on:click={handleStart} class="btn btn-success flex-1 text-sm py-2">
				⚡ Resurrect
			</button>
			<button on:click={handleDelete} class="btn btn-danger text-sm py-2 px-3" title="Delete VM">
				🗑️
			</button>
		{:else if vm.status === 'running'}
			<button on:click={handleShutdown} class="btn btn-secondary flex-1 text-sm py-2">
				😴 Sleep
			</button>
			<button on:click={handleStop} class="btn btn-danger text-sm py-2 px-3" title="Force kill">
				💀
			</button>
		{:else}
			<button disabled class="btn btn-secondary flex-1 text-sm py-2 opacity-50 cursor-wait">
				{status.icon} {status.label}
			</button>
		{/if}
	</div>
</div>
