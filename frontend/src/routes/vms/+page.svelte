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
	<title>The Zoo | Agni</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold">The Zoo</h1>
			<p class="text-gray-500 text-sm">All your little monsters in one place</p>
		</div>
		<a href="/vms/new" class="btn btn-primary">+ Summon New Demon</a>
	</div>

	<!-- VM List -->
	{#if $vms.loading && $vms.vms.length === 0}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-orange-500 border-t-transparent"></div>
			<p class="text-gray-400 mt-4">Counting sheep... wait, wrong animals</p>
		</div>
	{:else if $vms.vms.length === 0}
		<div class="card text-center py-16">
			<div class="text-6xl mb-4">🦗</div>
			<h3 class="text-xl font-medium text-gray-300 mb-2">*crickets*</h3>
			<p class="text-gray-500 mb-6">The zoo is empty. No demons have been summoned yet.</p>
			<a href="/vms/new" class="btn btn-primary">Summon Your First Demon</a>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each $vms.vms as vm (vm.id)}
				<VMCard {vm} />
			{/each}
		</div>
	{/if}

	{#if $vms.error}
		<div class="bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg p-4">
			<span class="font-medium">Oops:</span> {$vms.error}
		</div>
	{/if}
</div>
