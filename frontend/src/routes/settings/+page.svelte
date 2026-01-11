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
	<title>Settings | Agni</title>
</svelte:head>

<div class="max-w-2xl mx-auto space-y-4 sm:space-y-6 px-2 sm:px-0">
	<div>
		<h1 class="text-xl sm:text-2xl font-bold">Settings</h1>
		<p class="text-gray-500 text-xs sm:text-sm">Application configuration and preferences</p>
	</div>

	<!-- Account Settings -->
	<div class="card">
		<h2 class="text-base sm:text-lg font-medium mb-3 sm:mb-4 flex items-center gap-2">
			<svg class="w-5 h-5 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
			</svg>
			<span>Account</span>
		</h2>
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
		<h2 class="text-base sm:text-lg font-medium mb-3 sm:mb-4 flex items-center gap-2">
			<svg class="w-5 h-5 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
			</svg>
			<span>System Information</span>
		</h2>
		{#if loading}
			<div class="text-gray-400 text-sm flex items-center gap-2">
				<svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
				</svg>
				Loading...
			</div>
		{:else if systemInfo}
			<div class="grid grid-cols-1 sm:grid-cols-2 gap-2 sm:gap-4 text-xs sm:text-sm">
				<div class="flex justify-between sm:block">
					<span class="text-gray-400">Version:</span>
					<span class="sm:ml-2">{systemInfo.version}</span>
				</div>
				<div class="flex justify-between sm:block">
					<span class="text-gray-400">Go Version:</span>
					<span class="sm:ml-2">{systemInfo.go_version}</span>
				</div>
				<div class="flex justify-between sm:block">
					<span class="text-gray-400">Firecracker:</span>
					<span class="sm:ml-2">{systemInfo.firecracker_version}</span>
				</div>
				<div class="flex justify-between sm:block">
					<span class="text-gray-400">OS:</span>
					<span class="sm:ml-2">{systemInfo.os}</span>
				</div>
				<div class="flex justify-between sm:block">
					<span class="text-gray-400">Architecture:</span>
					<span class="sm:ml-2">{systemInfo.arch}</span>
				</div>
				<div class="flex justify-between sm:block">
					<span class="text-gray-400">CPUs:</span>
					<span class="sm:ml-2">{systemInfo.num_cpu}</span>
				</div>
			</div>
		{:else}
			<div class="text-gray-500 text-sm">Unable to load system information</div>
		{/if}
	</div>

	<!-- API Configuration -->
	<div class="card">
		<h2 class="text-base sm:text-lg font-medium mb-3 sm:mb-4 flex items-center gap-2">
			<svg class="w-5 h-5 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
			</svg>
			<span>API Configuration</span>
		</h2>
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
		<h2 class="text-base sm:text-lg font-medium mb-3 sm:mb-4 flex items-center gap-2">
			<svg class="w-5 h-5 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
			</svg>
			<span>About Agni</span>
		</h2>
		<div class="space-y-2 text-sm text-gray-400">
			<p>Agni - Firecracker microVM management platform.</p>
			<p class="text-gray-500">
				Built on top of AWS Firecracker for lightweight, secure virtualization.
			</p>
			<div class="flex gap-4 mt-4">
				<a
					href="https://github.com/anubhavg-icpl/agni"
					target="_blank"
					rel="noopener noreferrer"
					class="text-orange-400 hover:text-orange-300 flex items-center gap-1.5"
				>
					<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
						<path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
					</svg>
					GitHub
				</a>
				<a
					href="https://firecracker-microvm.github.io/"
					target="_blank"
					rel="noopener noreferrer"
					class="text-orange-400 hover:text-orange-300 flex items-center gap-1.5"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
					</svg>
					Docs
				</a>
			</div>
		</div>
	</div>

	<!-- Danger Zone -->
	<div class="card border-red-500/30 bg-red-500/5">
		<h2 class="text-base sm:text-lg font-medium mb-3 sm:mb-4 text-red-400 flex items-center gap-2">
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
			</svg>
			<span>Session</span>
		</h2>
		<div class="space-y-4">
			<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
				<div>
					<div class="font-medium text-gray-200 text-sm sm:text-base">Sign Out</div>
					<div class="text-xs sm:text-sm text-gray-500">End your current session</div>
				</div>
				<Button variant="danger" on:click={handleLogout}>
					<svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
					</svg>
					Sign Out
				</Button>
			</div>
		</div>
	</div>
</div>
