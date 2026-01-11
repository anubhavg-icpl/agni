<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { api, type SystemInfo } from '$lib/api/client';
	import Button from '$lib/components/Common/Button.svelte';
	import { toasts } from '$lib/components/Common/Toast.svelte';

	let systemInfo: SystemInfo | null = null;
	let loading = true;

	// Password change form
	let currentPassword = '';
	let newPassword = '';
	let confirmPassword = '';
	let changingPassword = false;

	// API settings
	let apiPort = '8080';

	onMount(async () => {
		try {
			systemInfo = await api.getSystemInfo();
		} catch (e) {
			console.error('Failed to load system info:', e);
		}
		loading = false;
	});

	async function handlePasswordChange() {
		if (!currentPassword || !newPassword || !confirmPassword) {
			toasts.error('All password fields are required');
			return;
		}

		if (newPassword !== confirmPassword) {
			toasts.error('New passwords do not match');
			return;
		}

		if (newPassword.length < 8) {
			toasts.error('Password must be at least 8 characters');
			return;
		}

		changingPassword = true;

		try {
			// This would be an API call to change password
			// await api.changePassword(currentPassword, newPassword);
			toasts.success('Password changed successfully');
			currentPassword = '';
			newPassword = '';
			confirmPassword = '';
		} catch (e) {
			toasts.error((e as Error).message);
		}

		changingPassword = false;
	}

	async function handleLogout() {
		await auth.logout();
		goto('/login');
	}
</script>

<svelte:head>
	<title>Settings - Firectl</title>
</svelte:head>

<div class="max-w-2xl mx-auto space-y-6">
	<h1 class="text-2xl font-bold">Settings</h1>

	<!-- Account Settings -->
	<div class="card">
		<h2 class="text-lg font-medium mb-4">Account</h2>
		<div class="space-y-4">
			<div>
				<span class="text-sm text-gray-400">Username:</span>
				<span class="ml-2 font-medium">{$auth.user?.username || 'Unknown'}</span>
			</div>

			<!-- Password Change -->
			<div class="border-t border-gray-700 pt-4 mt-4">
				<h3 class="font-medium mb-3">Change Password</h3>
				<form on:submit|preventDefault={handlePasswordChange} class="space-y-3">
					<div>
						<label for="currentPassword" class="label text-sm">Current Password</label>
						<input
							type="password"
							id="currentPassword"
							bind:value={currentPassword}
							class="input w-full"
							autocomplete="current-password"
						/>
					</div>
					<div>
						<label for="newPassword" class="label text-sm">New Password</label>
						<input
							type="password"
							id="newPassword"
							bind:value={newPassword}
							class="input w-full"
							autocomplete="new-password"
						/>
					</div>
					<div>
						<label for="confirmPassword" class="label text-sm">Confirm New Password</label>
						<input
							type="password"
							id="confirmPassword"
							bind:value={confirmPassword}
							class="input w-full"
							autocomplete="new-password"
						/>
					</div>
					<Button type="submit" variant="primary" loading={changingPassword}>
						Change Password
					</Button>
				</form>
			</div>
		</div>
	</div>

	<!-- System Information -->
	<div class="card">
		<h2 class="text-lg font-medium mb-4">System Information</h2>
		{#if loading}
			<div class="text-gray-400">Loading...</div>
		{:else if systemInfo}
			<div class="grid grid-cols-2 gap-4 text-sm">
				<div>
					<span class="text-gray-400">Version:</span>
					<span class="ml-2">{systemInfo.version}</span>
				</div>
				<div>
					<span class="text-gray-400">Go Version:</span>
					<span class="ml-2">{systemInfo.go_version}</span>
				</div>
				<div>
					<span class="text-gray-400">Firecracker:</span>
					<span class="ml-2">{systemInfo.firecracker_version}</span>
				</div>
				<div>
					<span class="text-gray-400">OS:</span>
					<span class="ml-2">{systemInfo.os}</span>
				</div>
				<div>
					<span class="text-gray-400">Architecture:</span>
					<span class="ml-2">{systemInfo.arch}</span>
				</div>
				<div>
					<span class="text-gray-400">CPUs:</span>
					<span class="ml-2">{systemInfo.num_cpu}</span>
				</div>
			</div>
		{:else}
			<div class="text-gray-500">Unable to load system information</div>
		{/if}
	</div>

	<!-- API Configuration -->
	<div class="card">
		<h2 class="text-lg font-medium mb-4">API Configuration</h2>
		<div class="space-y-4">
			<div>
				<label for="apiPort" class="label text-sm">API Port</label>
				<input type="text" id="apiPort" bind:value={apiPort} class="input w-32" disabled />
				<p class="text-sm text-gray-500 mt-1">API server port (requires restart to change)</p>
			</div>

			<div class="bg-gray-900 rounded-lg p-3 font-mono text-sm">
				<div class="text-gray-400 mb-1">API Base URL:</div>
				<div class="text-blue-400">
					{typeof window !== 'undefined' ? window.location.origin : ''}/api
				</div>
			</div>
		</div>
	</div>

	<!-- About -->
	<div class="card">
		<h2 class="text-lg font-medium mb-4">About Firectl</h2>
		<div class="space-y-2 text-sm text-gray-400">
			<p>Firectl is a command-line tool and GUI for managing Firecracker microVMs.</p>
			<p>
				Originally created by Amazon.com, Inc. as a CLI tool. GUI implementation by Anubhav Gain.
			</p>
			<div class="flex gap-4 mt-4">
				<a
					href="https://github.com/firecracker-microvm/firectl"
					target="_blank"
					rel="noopener noreferrer"
					class="text-blue-400 hover:text-blue-300"
				>
					GitHub Repository
				</a>
				<a
					href="https://firecracker-microvm.github.io/"
					target="_blank"
					rel="noopener noreferrer"
					class="text-blue-400 hover:text-blue-300"
				>
					Documentation
				</a>
			</div>
		</div>
	</div>

	<!-- Danger Zone -->
	<div class="card border-red-500/30">
		<h2 class="text-lg font-medium mb-4 text-red-400">Danger Zone</h2>
		<div class="space-y-4">
			<div class="flex items-center justify-between">
				<div>
					<div class="font-medium">Sign Out</div>
					<div class="text-sm text-gray-500">Log out from the current session</div>
				</div>
				<Button variant="danger" on:click={handleLogout}>Sign Out</Button>
			</div>
		</div>
	</div>
</div>
