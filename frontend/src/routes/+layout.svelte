<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let initialized = false;
	let mobileMenuOpen = false;

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

	// Close mobile menu on navigation
	$: if ($page.url.pathname) {
		mobileMenuOpen = false;
	}

	async function handleLogout() {
		await auth.logout();
		goto('/login');
	}

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}
</script>

{#if !initialized || $auth.loading}
	<div class="min-h-screen bg-gray-900 flex items-center justify-center">
		<div class="flex flex-col items-center gap-4">
			<img src="/logo.png" alt="Agni" class="w-16 h-16 logo-glow" />
			<div class="text-gray-400">Warming up the hamster wheel...</div>
		</div>
	</div>
{:else if $auth.setupRequired || $page.url.pathname.includes('/login')}
	<slot />
{:else if $isAuthenticated}
	<div class="min-h-screen bg-gray-900">
		<!-- Navigation -->
		<nav class="bg-gray-800/80 border-b border-gray-700/50 backdrop-blur-sm sticky top-0 z-50">
			<div class="px-4 sm:px-6 lg:px-8">
				<div class="flex items-center justify-between h-14 sm:h-16">
					<div class="flex items-center gap-4 sm:gap-8">
						<a href="/" class="flex items-center gap-2">
							<img src="/logo.png" alt="Agni" class="w-7 h-7 sm:w-8 sm:h-8" />
							<span class="text-lg sm:text-xl font-bold text-fire">Agni</span>
						</a>
						<!-- Desktop nav -->
						<div class="hidden sm:flex gap-1">
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
					<div class="flex items-center gap-2 sm:gap-4">
						<a href="/settings" class="hidden sm:block text-gray-400 hover:text-orange-400 text-sm transition-colors">
							{$auth.user?.username}
						</a>
						<button
							on:click={handleLogout}
							class="hidden sm:block text-gray-300 hover:text-red-400 px-3 py-2 rounded-md text-sm font-medium transition-colors"
						>
							Sign Out
						</button>
						<!-- Mobile menu button -->
						<button
							on:click={toggleMobileMenu}
							class="sm:hidden p-2 rounded-md text-gray-400 hover:text-orange-400 hover:bg-gray-700/50 transition-colors"
							aria-label="Toggle menu"
						>
							{#if mobileMenuOpen}
								<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
								</svg>
							{:else}
								<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
								</svg>
							{/if}
						</button>
					</div>
				</div>
			</div>

			<!-- Mobile menu -->
			{#if mobileMenuOpen}
				<div class="sm:hidden border-t border-gray-700/50 bg-gray-800/95 backdrop-blur-sm">
					<div class="px-3 py-3 space-y-1">
						<a
							href="/"
							class="flex items-center gap-2 nav-link px-3 py-2.5 rounded-md text-sm font-medium"
							class:nav-link-active={$page.url.pathname === '/'}
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z" />
							</svg>
							Dashboard
						</a>
						<a
							href="/vms"
							class="flex items-center gap-2 nav-link px-3 py-2.5 rounded-md text-sm font-medium"
							class:nav-link-active={$page.url.pathname.startsWith('/vms')}
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2" />
							</svg>
							Virtual Machines
						</a>
						<a
							href="/configs"
							class="flex items-center gap-2 nav-link px-3 py-2.5 rounded-md text-sm font-medium"
							class:nav-link-active={$page.url.pathname.startsWith('/configs')}
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
							</svg>
							Templates
						</a>
						<a
							href="/settings"
							class="flex items-center gap-2 nav-link px-3 py-2.5 rounded-md text-sm font-medium"
							class:nav-link-active={$page.url.pathname.startsWith('/settings')}
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
							</svg>
							Settings
						</a>
						<div class="border-t border-gray-700/50 pt-2 mt-2">
							<div class="px-3 py-1 text-xs text-gray-500">Signed in as {$auth.user?.username}</div>
							<button
								on:click={handleLogout}
								class="w-full text-left px-3 py-2.5 rounded-md text-sm font-medium text-red-400 hover:bg-red-500/10 transition-colors flex items-center gap-2"
							>
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
								</svg>
								Sign Out
							</button>
						</div>
					</div>
				</div>
			{/if}
		</nav>

		<!-- Main content -->
		<main class="px-4 sm:px-6 lg:px-8 py-4 sm:py-6">
			<slot />
		</main>
	</div>
{:else}
	<slot />
{/if}
