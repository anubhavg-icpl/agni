<script lang="ts">
	import { onMount } from 'svelte';
	import { vms } from '$lib/stores/vms';
	import VMCard from '$lib/components/Dashboard/VMCard.svelte';

	onMount(() => {
		vms.fetch();
	});

	// Refresh every 5 seconds
	let interval: ReturnType<typeof setInterval>;
	onMount(() => {
		interval = setInterval(() => {
			vms.fetch();
		}, 5000);
		return () => clearInterval(interval);
	});
</script>

<svelte:head>
	<title>Virtual Machines | Agni</title>
</svelte:head>

<div class="space-y-6 lg:space-y-8">
	<!-- Header -->
	<div class="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
		<div>
			<h1 class="text-2xl sm:text-3xl font-bold text-white mb-1">Virtual Machines</h1>
			<p class="text-gray-400 text-sm sm:text-base">Manage your Firecracker microVMs</p>
		</div>
		<a
			href="/vms/new"
			class="inline-flex items-center gap-2 px-5 py-3 bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white font-medium rounded-xl transition-all duration-300 transform hover:scale-[1.02] active:scale-[0.98] shadow-lg shadow-orange-500/20"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			<span>Create VM</span>
		</a>
	</div>

	<!-- VM List -->
	{#if $vms.loading && $vms.vms.length === 0}
		<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 border border-gray-700/50 rounded-2xl p-12">
			<div class="flex flex-col items-center justify-center text-center">
				<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-orange-500/20 to-red-500/20 flex items-center justify-center mb-4">
					<svg class="w-6 h-6 text-orange-400 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
					</svg>
				</div>
				<p class="text-gray-400">Loading virtual machines...</p>
			</div>
		</div>
	{:else if $vms.vms.length === 0}
		<div class="bg-gradient-to-br from-gray-800/50 to-gray-900/50 border border-gray-700/50 rounded-2xl p-12">
			<div class="flex flex-col items-center justify-center text-center max-w-md mx-auto">
				<div class="w-16 h-16 rounded-2xl bg-gradient-to-br from-gray-700/50 to-gray-800/50 flex items-center justify-center mb-5">
					<svg class="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
					</svg>
				</div>
				<h3 class="text-xl font-semibold text-white mb-2">No Virtual Machines</h3>
				<p class="text-gray-400 mb-6">Create your first microVM to get started.</p>
				<a
					href="/vms/new"
					class="inline-flex items-center gap-2 px-5 py-3 bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white font-medium rounded-xl transition-all duration-300 transform hover:scale-[1.02] active:scale-[0.98]"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
					</svg>
					<span>Create Your First VM</span>
				</a>
			</div>
		</div>
	{:else}
		<div>
			<div class="flex items-center justify-between mb-4">
				<h2 class="text-lg sm:text-xl font-semibold text-white">All VMs</h2>
				<span class="text-sm text-gray-500">{$vms.vms.length} total</span>
			</div>
			<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4">
				{#each $vms.vms as vm (vm.id)}
					<VMCard {vm} />
				{/each}
			</div>
		</div>
	{/if}

	{#if $vms.error}
		<div class="bg-gradient-to-r from-red-500/10 to-red-600/10 border border-red-500/30 rounded-2xl p-4 flex items-center gap-3">
			<div class="w-10 h-10 rounded-xl bg-red-500/20 flex items-center justify-center flex-shrink-0">
				<svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
			</div>
			<div>
				<p class="text-red-400 font-medium">Error Loading VMs</p>
				<p class="text-red-400/70 text-sm">{$vms.error}</p>
			</div>
		</div>
	{/if}
</div>
