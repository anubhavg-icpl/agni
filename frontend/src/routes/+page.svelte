<script lang="ts">
	import { onMount } from 'svelte';
	import { vms } from '$lib/stores/vms';
	import { auth } from '$lib/stores/auth';
	import VMCard from '$lib/components/Dashboard/VMCard.svelte';
	import QuickStats from '$lib/components/Dashboard/QuickStats.svelte';

	let setupUsername = '';
	let setupPassword = '';

	onMount(() => {
		if (!$auth.setupRequired) {
			vms.fetch();
		}
	});

	// Refresh every 5 seconds
	let interval: ReturnType<typeof setInterval>;
	onMount(() => {
		interval = setInterval(() => {
			if (!$auth.setupRequired) {
				vms.fetch();
			}
		}, 5000);
		return () => clearInterval(interval);
	});

	async function handleSetup() {
		await auth.setup(setupUsername, setupPassword);
	}
</script>

<svelte:head>
	<title>Command Center | Agni</title>
</svelte:head>

{#if $auth.setupRequired}
	<!-- Setup screen -->
	<div class="min-h-screen bg-gray-900 flex items-center justify-center">
		<div class="card max-w-md w-full p-8">
			<h1 class="text-2xl font-bold text-center mb-6">Fresh Meat Detected</h1>
			<p class="text-gray-400 text-center mb-8">Someone needs to run this circus. Might as well be you.</p>

			<form on:submit|preventDefault={handleSetup} class="space-y-4">
				<div>
					<label class="label" for="username">Username</label>
					<input
						type="text"
						id="username"
						name="username"
						bind:value={setupUsername}
						class="input w-full"
						placeholder="Pick wisely"
						required
					/>
				</div>
				<div>
					<label class="label" for="password">Password</label>
					<input
						type="password"
						id="password"
						name="password"
						bind:value={setupPassword}
						class="input w-full"
						placeholder="Make it memorable"
						required
						minlength="8"
					/>
					<p class="text-xs text-gray-500 mt-1">8+ characters. Impress us.</p>
				</div>
				{#if $auth.error}
					<p class="text-red-400 text-sm">{$auth.error}</p>
				{/if}
				<button type="submit" class="btn btn-primary w-full" disabled={$auth.loading}>
					{$auth.loading ? 'Forging credentials...' : 'Claim Your Throne'}
				</button>
			</form>
		</div>
	</div>
{:else}
	<div class="space-y-4 sm:space-y-6 lg:space-y-8">
		<!-- Header -->
		<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
			<div>
				<h1 class="text-xl sm:text-2xl font-bold">Mission Control</h1>
				<p class="text-gray-500 text-xs sm:text-sm">Where chaos is managed (barely)</p>
			</div>
			<a href="/vms/new" class="btn btn-primary text-sm sm:text-base whitespace-nowrap"> + Spawn Another </a>
		</div>

		<!-- Quick Stats -->
		<QuickStats vms={$vms.vms} />

		<!-- VM Grid -->
		{#if $vms.loading && $vms.vms.length === 0}
			<div class="text-center py-8 sm:py-12 text-gray-400 text-sm sm:text-base">Waking up the minions...</div>
		{:else if $vms.vms.length === 0}
			<div class="card text-center py-8 sm:py-12">
				<h3 class="text-base sm:text-lg font-medium text-gray-300 mb-2">Tumbleweeds...</h3>
				<p class="text-gray-500 mb-4 text-sm">No VMs yet. This emptiness is judging you.</p>
				<a href="/vms/new" class="btn btn-primary">Create Something</a>
			</div>
		{:else}
			<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
				{#each $vms.vms as vm (vm.id)}
					<VMCard {vm} />
				{/each}
			</div>
		{/if}

		{#if $vms.error}
			<div class="bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg p-3 sm:p-4 text-sm">
				Well, that didn't work: {$vms.error}
			</div>
		{/if}
	</div>
{/if}
