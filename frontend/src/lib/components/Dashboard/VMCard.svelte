<script lang="ts">
	import type { VM } from '$lib/api/client';
	import { vms } from '$lib/stores/vms';

	export let vm: VM;

	const statusConfig: Record<string, { class: string; label: string; iconBg: string; iconColor: string }> = {
		running: { class: 'bg-green-500/10 text-green-400 border-green-500/30', label: 'Running', iconBg: 'from-green-500/20 to-emerald-500/20', iconColor: 'text-green-400' },
		stopped: { class: 'bg-gray-500/10 text-gray-400 border-gray-500/30', label: 'Stopped', iconBg: 'from-gray-600/20 to-gray-700/20', iconColor: 'text-gray-400' },
		starting: { class: 'bg-yellow-500/10 text-yellow-400 border-yellow-500/30', label: 'Starting', iconBg: 'from-yellow-500/20 to-amber-500/20', iconColor: 'text-yellow-400' },
		stopping: { class: 'bg-yellow-500/10 text-yellow-400 border-yellow-500/30', label: 'Stopping', iconBg: 'from-yellow-500/20 to-amber-500/20', iconColor: 'text-yellow-400' },
		error: { class: 'bg-red-500/10 text-red-400 border-red-500/30', label: 'Error', iconBg: 'from-red-500/20 to-red-600/20', iconColor: 'text-red-400' }
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

<div class="group relative bg-gradient-to-br from-gray-800/50 to-gray-900/50 border border-gray-700/50 rounded-2xl p-4 sm:p-5 hover:border-orange-500/30 transition-all duration-300 hover:shadow-xl hover:shadow-orange-500/5 overflow-hidden">
	<!-- Hover gradient overlay -->
	<div class="absolute inset-0 bg-gradient-to-br from-orange-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"></div>

	<!-- Header -->
	<div class="relative flex items-start justify-between mb-4">
		<div class="flex items-center gap-3 flex-1 min-w-0">
			<!-- Status Icon -->
			<div class="w-10 h-10 sm:w-12 sm:h-12 rounded-xl bg-gradient-to-br {status.iconBg} flex items-center justify-center flex-shrink-0 relative">
				{#if vm.status === 'running'}
					<span class="absolute top-0 right-0 w-3 h-3 bg-green-400 rounded-full animate-pulse"></span>
					<svg class="w-5 h-5 sm:w-6 sm:h-6 {status.iconColor}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				{:else if vm.status === 'error'}
					<svg class="w-5 h-5 sm:w-6 sm:h-6 {status.iconColor}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
					</svg>
				{:else if vm.status === 'starting' || vm.status === 'stopping'}
					<svg class="w-5 h-5 sm:w-6 sm:h-6 {status.iconColor} animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
					</svg>
				{:else}
					<svg class="w-5 h-5 sm:w-6 sm:h-6 {status.iconColor}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z" />
					</svg>
				{/if}
			</div>
			<div class="min-w-0">
				<h3 class="font-semibold text-base sm:text-lg text-gray-100 truncate group-hover:text-white transition-colors">{vm.name}</h3>
				<p class="text-xs text-gray-500 font-mono">{vm.id.slice(0, 8)}...</p>
			</div>
		</div>
		<span class="flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium rounded-lg border {status.class}">
			{#if vm.status === 'running'}
				<span class="w-1.5 h-1.5 rounded-full bg-green-400 animate-pulse"></span>
			{:else if vm.status === 'starting' || vm.status === 'stopping'}
				<svg class="w-3 h-3 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
				</svg>
			{:else if vm.status === 'error'}
				<span class="w-1.5 h-1.5 rounded-full bg-red-400"></span>
			{:else}
				<span class="w-1.5 h-1.5 rounded-full bg-gray-500"></span>
			{/if}
			<span class="hidden xs:inline">{status.label}</span>
		</span>
	</div>

	<!-- Specs -->
	<div class="relative grid grid-cols-3 gap-2 sm:gap-3 mb-4 p-3 bg-black/30 rounded-xl border border-gray-700/30">
		<div class="text-center">
			<div class="flex items-center justify-center gap-1 mb-1">
				<svg class="w-3.5 h-3.5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
				</svg>
				<span class="text-xs text-gray-500">vCPU</span>
			</div>
			<div class="text-sm sm:text-base font-semibold text-gray-200">{vm.config.cpus}</div>
		</div>
		<div class="text-center border-x border-gray-700/30">
			<div class="flex items-center justify-center gap-1 mb-1">
				<svg class="w-3.5 h-3.5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
				</svg>
				<span class="text-xs text-gray-500">RAM</span>
			</div>
			<div class="text-sm sm:text-base font-semibold text-gray-200">{vm.config.memory_mb} MB</div>
		</div>
		<div class="text-center">
			<div class="flex items-center justify-center gap-1 mb-1">
				<svg class="w-3.5 h-3.5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
				</svg>
				<span class="text-xs text-gray-500">Kernel</span>
			</div>
			<div class="text-sm sm:text-base font-semibold text-gray-200 truncate" title={vm.config.kernel_path}>
				{vm.config.kernel_path.split('/').pop()?.slice(0, 8) || 'N/A'}
			</div>
		</div>
	</div>

	<!-- Error Message -->
	{#if vm.error}
		<div class="relative text-xs sm:text-sm text-red-400 mb-4 p-3 bg-red-500/10 rounded-xl border border-red-500/20 flex items-start gap-2">
			<svg class="w-4 h-4 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
			</svg>
			<span>{vm.error}</span>
		</div>
	{/if}

	<!-- Actions -->
	<div class="relative flex gap-2">
		{#if vm.status === 'stopped' || vm.status === 'error'}
			<button on:click={handleStart} class="flex-1 py-2.5 px-4 bg-gradient-to-r from-green-500/20 to-emerald-500/20 hover:from-green-500/30 hover:to-emerald-500/30 text-green-400 text-xs sm:text-sm font-medium rounded-xl border border-green-500/30 transition-all duration-200 flex items-center justify-center gap-2">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
				<span>Start</span>
			</button>
			<button on:click={handleDelete} class="py-2.5 px-3 bg-gradient-to-r from-red-500/20 to-red-600/20 hover:from-red-500/30 hover:to-red-600/30 text-red-400 rounded-xl border border-red-500/30 transition-all duration-200" title="Delete VM">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
				</svg>
			</button>
		{:else if vm.status === 'running'}
			<button on:click={handleShutdown} class="flex-1 py-2.5 px-4 bg-gradient-to-r from-gray-600/20 to-gray-700/20 hover:from-gray-600/30 hover:to-gray-700/30 text-gray-300 text-xs sm:text-sm font-medium rounded-xl border border-gray-600/30 transition-all duration-200 flex items-center justify-center gap-2">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z" />
				</svg>
				<span>Stop</span>
			</button>
			<button on:click={handleStop} class="py-2.5 px-3 bg-gradient-to-r from-red-500/20 to-red-600/20 hover:from-red-500/30 hover:to-red-600/30 text-red-400 rounded-xl border border-red-500/30 transition-all duration-200" title="Force kill">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
				</svg>
			</button>
		{:else}
			<button disabled class="flex-1 py-2.5 px-4 bg-gray-700/20 text-gray-500 text-xs sm:text-sm font-medium rounded-xl border border-gray-700/30 cursor-wait flex items-center justify-center gap-2">
				<svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
				</svg>
				<span>{status.label}...</span>
			</button>
		{/if}
	</div>
</div>
