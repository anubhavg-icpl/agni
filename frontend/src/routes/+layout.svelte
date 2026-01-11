<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	onMount(async () => {
		await auth.init();
	});

	$: if (
		!$auth.loading &&
		!$isAuthenticated &&
		!$auth.setupRequired &&
		!$page.url.pathname.includes('/login')
	) {
		goto('/login');
	}

	async function handleLogout() {
		await auth.logout();
		goto('/login');
	}
</script>

{#if $auth.loading}
	<div class="min-h-screen bg-gray-900 flex items-center justify-center">
		<div class="text-gray-400">Loading...</div>
	</div>
{:else if $auth.setupRequired}
	<slot />
{:else if $isAuthenticated}
	<div class="min-h-screen bg-gray-900">
		<!-- Navigation -->
		<nav class="bg-gray-800 border-b border-gray-700">
			<div class="max-w-7xl mx-auto px-4">
				<div class="flex items-center justify-between h-16">
					<div class="flex items-center gap-8">
						<a href="/" class="text-xl font-bold text-primary-400">Firectl</a>
						<div class="flex gap-4">
							<a
								href="/"
								class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium"
							>
								Dashboard
							</a>
							<a
								href="/vms/new"
								class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium"
							>
								New VM
							</a>
							<a
								href="/configs"
								class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium"
							>
								Templates
							</a>
						</div>
					</div>
					<div class="flex items-center gap-4">
						<span class="text-gray-400 text-sm">{$auth.user?.username}</span>
						<button
							on:click={handleLogout}
							class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium"
						>
							Logout
						</button>
					</div>
				</div>
			</div>
		</nav>

		<!-- Main content -->
		<main class="max-w-7xl mx-auto px-4 py-8">
			<slot />
		</main>
	</div>
{:else}
	<slot />
{/if}
