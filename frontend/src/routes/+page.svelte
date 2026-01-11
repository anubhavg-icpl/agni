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
				<h1 class="text-xl sm:text-2xl font-bold">Dashboard</h1>
				<p class="text-gray-500 text-xs sm:text-sm">VM overview and quick actions</p>
			</div>
			<a href="/vms/new" class="btn btn-primary text-sm sm:text-base whitespace-nowrap flex items-center gap-1.5">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
				</svg>
				<span>New VM</span>
			</a>
		</div>

		<!-- Quick Stats -->
		<QuickStats vms={$vms.vms} />

		<!-- VM Grid -->
		{#if $vms.loading && $vms.vms.length === 0}
			<div class="text-center py-8 sm:py-12 text-gray-400 text-sm sm:text-base flex flex-col items-center gap-3">
				<svg class="w-8 h-8 animate-spin text-orange-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
				</svg>
				<span>Loading VMs...</span>
			</div>
		{:else if $vms.vms.length === 0}
			<div class="card text-center py-8 sm:py-12">
				<div class="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-800 flex items-center justify-center">
					<svg class="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2" />
					</svg>
				</div>
				<h3 class="text-base sm:text-lg font-medium text-gray-300 mb-2">No VMs Found</h3>
				<p class="text-gray-500 mb-4 text-sm">Get started by creating your first virtual machine.</p>
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
			<div class="bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg p-3 sm:p-4 text-sm flex items-center gap-2">
				<svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
				<span>{$vms.error}</span>
			</div>
		{/if}
	</div>
{/if}
