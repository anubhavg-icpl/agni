<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let initialized = false;

	onMount(async () => {
		if (!initialized) {
			initialized = true;
			await auth.init();
		}
	});

	// Only redirect to login if not authenticated and not on login page
	$: if (
		initialized &&
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

{#if !initialized || $auth.loading}
	<div class="min-h-screen bg-gray-900 flex items-center justify-center">
		<div class="flex flex-col items-center gap-4">
			<img src="/logo.png" alt="Agni" class="w-16 h-16 logo-glow" />
			<div class="text-gray-400">Loading...</div>
		</div>
	</div>
{:else if $auth.setupRequired || $page.url.pathname.includes('/login')}
	<slot />
{:else if $isAuthenticated}
	<div class="min-h-screen bg-gray-900">
		<!-- Navigation -->
		<nav class="bg-gray-800/80 border-b border-gray-700/50 backdrop-blur-sm sticky top-0 z-50">
			<div class="max-w-7xl mx-auto px-4">
				<div class="flex items-center justify-between h-16">
					<div class="flex items-center gap-8">
						<a href="/" class="flex items-center gap-2">
							<img src="/logo.png" alt="Agni" class="w-8 h-8" />
							<span class="text-xl font-bold text-fire">Agni</span>
						</a>
						<div class="flex gap-1">
							<a
								href="/"
								class="nav-link px-3 py-2 rounded-md text-sm font-medium"
								class:nav-link-active={$page.url.pathname === '/'}
							>
								Dashboard
							</a>
							<a
								href="/vms"
								class="nav-link px-3 py-2 rounded-md text-sm font-medium"
								class:nav-link-active={$page.url.pathname.startsWith('/vms')}
							>
								VMs
							</a>
							<a
								href="/configs"
								class="nav-link px-3 py-2 rounded-md text-sm font-medium"
								class:nav-link-active={$page.url.pathname.startsWith('/configs')}
							>
								Templates
							</a>
						</div>
					</div>
					<div class="flex items-center gap-4">
						<span class="text-gray-400 text-sm">{$auth.user?.username}</span>
						<button
							on:click={handleLogout}
							class="text-gray-300 hover:text-orange-400 px-3 py-2 rounded-md text-sm font-medium transition-colors"
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
