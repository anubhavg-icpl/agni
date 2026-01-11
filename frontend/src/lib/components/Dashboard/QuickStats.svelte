<script lang="ts">
	import type { VM } from '$lib/api/client';
	import RiveAnimation from '$lib/components/RiveAnimation.svelte';

	export let vms: VM[] = [];

	$: totalVMs = vms.length;
	$: runningVMs = vms.filter((vm) => vm.status === 'running').length;
	$: stoppedVMs = vms.filter((vm) => vm.status === 'stopped').length;
	$: errorVMs = vms.filter((vm) => vm.status === 'error').length;
</script>

<div class="grid grid-cols-2 lg:grid-cols-4 gap-2 sm:gap-4">
	<!-- Total VMs -->
	<div class="card hover:border-orange-500/30 transition-all duration-200 group p-3 sm:p-4">
		<div class="flex items-center gap-2 mb-1 sm:mb-2">
			<div class="w-6 h-6 sm:w-8 sm:h-8 flex items-center justify-center">
				<RiveAnimation src="/fire-logo.riv" width={32} height={32} />
			</div>
			<span class="text-xs sm:text-sm text-gray-400 truncate">Total VMs</span>
		</div>
		<div class="text-xl sm:text-3xl font-bold text-orange-400">
			{totalVMs}
		</div>
	</div>

	<!-- Running VMs -->
	<div class="card hover:border-green-500/30 transition-all duration-200 group p-3 sm:p-4">
		<div class="flex items-center gap-2 mb-1 sm:mb-2">
			<div class="w-6 h-6 sm:w-8 sm:h-8 flex items-center justify-center relative">
				{#if runningVMs > 0}
					<RiveAnimation src="/fire-skull.riv" width={32} height={32} />
				{:else}
					<div class="w-4 h-4 rounded-full bg-gray-600"></div>
				{/if}
			</div>
			<span class="text-xs sm:text-sm text-gray-400 truncate">Running</span>
		</div>
		<div class="text-xl sm:text-3xl font-bold text-green-400">
			{runningVMs}
		</div>
	</div>

	<!-- Stopped VMs -->
	<div class="card hover:border-gray-500/30 transition-all duration-200 group p-3 sm:p-4">
		<div class="flex items-center gap-2 mb-1 sm:mb-2">
			<div class="w-6 h-6 sm:w-8 sm:h-8 flex items-center justify-center">
				<svg class="w-5 h-5 sm:w-6 sm:h-6 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z" />
				</svg>
			</div>
			<span class="text-xs sm:text-sm text-gray-400 truncate">Stopped</span>
		</div>
		<div class="text-xl sm:text-3xl font-bold text-gray-400">
			{stoppedVMs}
		</div>
	</div>

	<!-- Error VMs -->
	<div class="card hover:border-red-500/30 transition-all duration-200 group p-3 sm:p-4 {errorVMs > 0 ? 'border-red-500/30 bg-red-500/5' : ''}">
		<div class="flex items-center gap-2 mb-1 sm:mb-2">
			<div class="w-6 h-6 sm:w-8 sm:h-8 flex items-center justify-center">
				{#if errorVMs > 0}
					<svg class="w-5 h-5 sm:w-6 sm:h-6 text-red-500 animate-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
					</svg>
				{:else}
					<svg class="w-5 h-5 sm:w-6 sm:h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				{/if}
			</div>
			<span class="text-xs sm:text-sm text-gray-400 truncate">Errors</span>
		</div>
		<div class="text-xl sm:text-3xl font-bold {errorVMs > 0 ? 'text-red-400' : 'text-gray-500'}">
			{errorVMs}
		</div>
	</div>
</div>
