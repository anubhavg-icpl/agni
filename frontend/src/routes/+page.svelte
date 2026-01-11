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
	<title>Dashboard - Firectl</title>
</svelte:head>

{#if $auth.setupRequired}
	<!-- Setup screen -->
	<div class="min-h-screen bg-gray-900 flex items-center justify-center">
		<div class="card max-w-md w-full p-8">
			<h1 class="text-2xl font-bold text-center mb-6">Welcome to Firectl</h1>
			<p class="text-gray-400 text-center mb-8">Create your admin account to get started.</p>

			<form on:submit|preventDefault={handleSetup} class="space-y-4">
				<div>
					<label class="label" for="username">Username</label>
					<input
						type="text"
						id="username"
						name="username"
						bind:value={setupUsername}
						class="input w-full"
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
						required
						minlength="8"
					/>
					<p class="text-xs text-gray-500 mt-1">Minimum 8 characters</p>
				</div>
				{#if $auth.error}
					<p class="text-red-400 text-sm">{$auth.error}</p>
				{/if}
				<button type="submit" class="btn btn-primary w-full" disabled={$auth.loading}>
					{$auth.loading ? 'Creating...' : 'Create Admin Account'}
				</button>
			</form>
		</div>
	</div>
{:else}
	<div class="space-y-8">
		<!-- Header -->
		<div class="flex items-center justify-between">
			<h1 class="text-2xl font-bold">Mission Control</h1>
			<a href="/vms/new" class="btn btn-primary"> + Birth a VM </a>
		</div>

		<!-- Quick Stats -->
		<QuickStats vms={$vms.vms} />

		<!-- VM Grid -->
		{#if $vms.loading && $vms.vms.length === 0}
			<div class="text-center py-12 text-gray-400">Summoning demons...</div>
		{:else if $vms.vms.length === 0}
			<div class="card text-center py-12">
				<h3 class="text-lg font-medium text-gray-300 mb-2">It's quiet... too quiet.</h3>
				<p class="text-gray-500 mb-4">Start a fire. Create a VM.</p>
				<a href="/vms/new" class="btn btn-primary">Ignite</a>
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
				{$vms.error}
			</div>
		{/if}
	</div>
{/if}
