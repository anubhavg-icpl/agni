<script lang="ts">
	import type { VM } from '$lib/api/client';
	import { vms } from '$lib/stores/vms';
	import RiveAnimation from '$lib/components/RiveAnimation.svelte';

	export let vm: VM;

	const statusConfig: Record<string, { class: string; label: string; color: string }> = {
		running: { class: 'status-running', label: 'Running', color: 'text-green-400' },
		stopped: { class: 'status-stopped', label: 'Stopped', color: 'text-gray-400' },
		starting: { class: 'status-starting', label: 'Starting', color: 'text-yellow-400' },
		stopping: { class: 'status-stopping', label: 'Stopping', color: 'text-yellow-400' },
		error: { class: 'status-error', label: 'Error', color: 'text-red-400' }
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
		if (confirm(`Delete "${vm.name}"? This action cannot be undone.`)) {
			await vms.delete(vm.id);
		}
	}
</script>

<div class="card hover:border-orange-500/30 transition-all duration-200 hover:shadow-lg hover:shadow-orange-500/5 p-3 sm:p-4">
	<!-- Header -->
	<div class="flex items-start justify-between mb-3 sm:mb-4">
		<div class="flex items-center gap-2 flex-1 min-w-0">
			<div class="w-8 h-8 flex-shrink-0">
				{#if vm.status === 'running'}
					<RiveAnimation src="/fire-skull.riv" width={32} height={32} />
				{:else if vm.status === 'error'}
					<div class="w-8 h-8 rounded-full bg-red-500/20 flex items-center justify-center">
						<svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
						</svg>
					</div>
				{:else}
					<div class="w-8 h-8 rounded-full bg-gray-700/50 flex items-center justify-center">
						<svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2" />
						</svg>
					</div>
				{/if}
			</div>
			<div class="min-w-0">
				<h3 class="font-semibold text-base sm:text-lg text-gray-100 truncate">{vm.name}</h3>
				<p class="text-xs text-gray-500 font-mono">{vm.id.slice(0, 8)}...</p>
			</div>
		</div>
		<span class="flex items-center gap-1.5 px-2 sm:px-2.5 py-0.5 sm:py-1 text-xs font-medium rounded-full {status.class}">
			{#if vm.status === 'running'}
				<span class="w-2 h-2 rounded-full bg-green-400 animate-pulse"></span>
			{:else if vm.status === 'starting' || vm.status === 'stopping'}
				<svg class="w-3 h-3 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
				</svg>
			{:else if vm.status === 'error'}
				<span class="w-2 h-2 rounded-full bg-red-400"></span>
			{:else}
				<span class="w-2 h-2 rounded-full bg-gray-500"></span>
			{/if}
			<span class="hidden xs:inline">{status.label}</span>
		</span>
	</div>

	<!-- Specs -->
	<div class="grid grid-cols-3 gap-2 sm:gap-3 mb-3 sm:mb-4 p-2 sm:p-3 bg-black/20 rounded-lg">
		<div class="text-center">
			<div class="text-xs text-gray-500 mb-0.5">vCPU</div>
			<div class="text-xs sm:text-sm font-medium text-gray-200">{vm.config.cpus}</div>
		</div>
		<div class="text-center border-x border-gray-700/50">
			<div class="text-xs text-gray-500 mb-0.5">Memory</div>
			<div class="text-xs sm:text-sm font-medium text-gray-200">{vm.config.memory_mb} MB</div>
		</div>
		<div class="text-center">
			<div class="text-xs text-gray-500 mb-0.5">Kernel</div>
			<div class="text-xs sm:text-sm font-medium text-gray-200 truncate" title={vm.config.kernel_path}>
				{vm.config.kernel_path.split('/').pop()?.slice(0, 8) || 'N/A'}
			</div>
		</div>
	</div>

	<!-- Error Message -->
	{#if vm.error}
		<div class="text-xs sm:text-sm text-red-400 mb-3 sm:mb-4 p-2 sm:p-2.5 bg-red-500/10 rounded-lg border border-red-500/20">
			<span class="font-medium">Error:</span> {vm.error}
		</div>
	{/if}

	<!-- Actions -->
	<div class="flex gap-2">
		{#if vm.status === 'stopped' || vm.status === 'error'}
			<button on:click={handleStart} class="btn btn-success flex-1 text-xs sm:text-sm py-1.5 sm:py-2 flex items-center justify-center gap-1.5">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
				<span>Start</span>
			</button>
			<button on:click={handleDelete} class="btn btn-danger text-xs sm:text-sm py-1.5 sm:py-2 px-2 sm:px-3" title="Delete VM">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
				</svg>
			</button>
		{:else if vm.status === 'running'}
			<button on:click={handleShutdown} class="btn btn-secondary flex-1 text-xs sm:text-sm py-1.5 sm:py-2 flex items-center justify-center gap-1.5">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z" />
				</svg>
				<span>Stop</span>
			</button>
			<button on:click={handleStop} class="btn btn-danger text-xs sm:text-sm py-1.5 sm:py-2 px-2 sm:px-3" title="Force kill">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
				</svg>
			</button>
		{:else}
			<button disabled class="btn btn-secondary flex-1 text-xs sm:text-sm py-1.5 sm:py-2 opacity-50 cursor-wait flex items-center justify-center gap-1.5">
				<svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
				</svg>
				<span>{status.label}</span>
			</button>
		{/if}
	</div>
</div>
