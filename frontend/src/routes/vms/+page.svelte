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

<div class="space-y-4 sm:space-y-6">
	<!-- Header -->
	<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
		<div>
			<h1 class="text-xl sm:text-2xl font-bold">Virtual Machines</h1>
			<p class="text-gray-500 text-xs sm:text-sm">Manage your Firecracker microVMs</p>
		</div>
		<a href="/vms/new" class="btn btn-primary text-sm sm:text-base whitespace-nowrap flex items-center gap-1.5">
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
			</svg>
			<span>New VM</span>
		</a>
	</div>

	<!-- VM List -->
	{#if $vms.loading && $vms.vms.length === 0}
		<div class="text-center py-12 flex flex-col items-center gap-3">
			<svg class="w-8 h-8 animate-spin text-orange-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
			</svg>
			<p class="text-gray-400">Loading virtual machines...</p>
		</div>
	{:else if $vms.vms.length === 0}
		<div class="card text-center py-12 sm:py-16">
			<div class="w-20 h-20 mx-auto mb-4 rounded-full bg-gray-800 flex items-center justify-center">
				<svg class="w-10 h-10 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2" />
				</svg>
			</div>
			<h3 class="text-lg sm:text-xl font-medium text-gray-300 mb-2">No Virtual Machines</h3>
			<p class="text-gray-500 mb-6 text-sm">Create your first microVM to get started.</p>
			<a href="/vms/new" class="btn btn-primary">Create VM</a>
		</div>
	{:else}
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
			{#each $vms.vms as vm (vm.id)}
				<VMCard {vm} />
			{/each}
		</div>
	{/if}

	{#if $vms.error}
		<div class="bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg p-4 flex items-center gap-2">
			<svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
			</svg>
			<span>{$vms.error}</span>
		</div>
	{/if}
</div>
