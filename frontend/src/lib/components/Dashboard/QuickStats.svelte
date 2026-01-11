<script lang="ts">
	import type { VM } from '$lib/api/client';

	export let vms: VM[] = [];

	$: totalVMs = vms.length;
	$: runningVMs = vms.filter((vm) => vm.status === 'running').length;
	$: stoppedVMs = vms.filter((vm) => vm.status === 'stopped').length;
	$: errorVMs = vms.filter((vm) => vm.status === 'error').length;
</script>

<div class="grid grid-cols-2 md:grid-cols-4 gap-3 sm:gap-4">
	<!-- Total VMs -->
	<div class="group relative bg-gradient-to-br from-gray-800/50 to-gray-900/50 border border-gray-700/50 rounded-2xl p-4 sm:p-5 hover:border-orange-500/30 transition-all duration-300 hover:shadow-lg hover:shadow-orange-500/5 overflow-hidden">
		<div class="absolute inset-0 bg-gradient-to-br from-orange-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
		<div class="relative">
			<div class="flex items-center justify-between mb-3">
				<div class="w-10 h-10 sm:w-12 sm:h-12 rounded-xl bg-gradient-to-br from-orange-500/20 to-red-500/20 flex items-center justify-center">
					<svg class="w-5 h-5 sm:w-6 sm:h-6 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
					</svg>
				</div>
				<span class="text-xs text-gray-500 uppercase tracking-wider font-medium">Total</span>
			</div>
			<div class="text-3xl sm:text-4xl font-bold text-white mb-1">
				{totalVMs}
			</div>
			<p class="text-xs sm:text-sm text-gray-400">Virtual Machines</p>
		</div>
	</div>

	<!-- Running VMs -->
	<div class="group relative bg-gradient-to-br from-gray-800/50 to-gray-900/50 border border-gray-700/50 rounded-2xl p-4 sm:p-5 hover:border-green-500/30 transition-all duration-300 hover:shadow-lg hover:shadow-green-500/5 overflow-hidden">
		<div class="absolute inset-0 bg-gradient-to-br from-green-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
		<div class="relative">
			<div class="flex items-center justify-between mb-3">
				<div class="w-10 h-10 sm:w-12 sm:h-12 rounded-xl bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center relative">
					{#if runningVMs > 0}
						<span class="absolute top-0 right-0 w-3 h-3 bg-green-400 rounded-full animate-pulse"></span>
					{/if}
					<svg class="w-5 h-5 sm:w-6 sm:h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</div>
				<span class="text-xs text-gray-500 uppercase tracking-wider font-medium">Active</span>
			</div>
			<div class="text-3xl sm:text-4xl font-bold text-green-400 mb-1">
				{runningVMs}
			</div>
			<p class="text-xs sm:text-sm text-gray-400">Running Now</p>
		</div>
	</div>

	<!-- Stopped VMs -->
	<div class="group relative bg-gradient-to-br from-gray-800/50 to-gray-900/50 border border-gray-700/50 rounded-2xl p-4 sm:p-5 hover:border-gray-500/30 transition-all duration-300 hover:shadow-lg hover:shadow-gray-500/5 overflow-hidden">
		<div class="absolute inset-0 bg-gradient-to-br from-gray-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
		<div class="relative">
			<div class="flex items-center justify-between mb-3">
				<div class="w-10 h-10 sm:w-12 sm:h-12 rounded-xl bg-gradient-to-br from-gray-600/30 to-gray-700/30 flex items-center justify-center">
					<svg class="w-5 h-5 sm:w-6 sm:h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z" />
					</svg>
				</div>
				<span class="text-xs text-gray-500 uppercase tracking-wider font-medium">Idle</span>
			</div>
			<div class="text-3xl sm:text-4xl font-bold text-gray-300 mb-1">
				{stoppedVMs}
			</div>
			<p class="text-xs sm:text-sm text-gray-400">Stopped</p>
		</div>
	</div>

	<!-- Error VMs -->
	<div class="group relative bg-gradient-to-br from-gray-800/50 to-gray-900/50 border rounded-2xl p-4 sm:p-5 transition-all duration-300 hover:shadow-lg overflow-hidden {errorVMs > 0 ? 'border-red-500/30 hover:border-red-500/50 hover:shadow-red-500/10' : 'border-gray-700/50 hover:border-red-500/30 hover:shadow-red-500/5'}">
		<div class="absolute inset-0 bg-gradient-to-br from-red-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 {errorVMs > 0 ? '!opacity-100' : ''}"></div>
		<div class="relative">
			<div class="flex items-center justify-between mb-3">
				<div class="w-10 h-10 sm:w-12 sm:h-12 rounded-xl bg-gradient-to-br from-red-500/20 to-red-600/20 flex items-center justify-center relative">
					{#if errorVMs > 0}
						<span class="absolute -top-1 -right-1 w-5 h-5 bg-red-500 rounded-full flex items-center justify-center text-xs font-bold text-white">{errorVMs}</span>
					{/if}
					<svg class="w-5 h-5 sm:w-6 sm:h-6 {errorVMs > 0 ? 'text-red-400' : 'text-gray-500'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
					</svg>
				</div>
				<span class="text-xs text-gray-500 uppercase tracking-wider font-medium">Issues</span>
			</div>
			<div class="text-3xl sm:text-4xl font-bold mb-1 {errorVMs > 0 ? 'text-red-400' : 'text-gray-500'}">
				{errorVMs}
			</div>
			<p class="text-xs sm:text-sm text-gray-400">{errorVMs > 0 ? 'Need Attention' : 'All Clear'}</p>
		</div>
	</div>
</div>
