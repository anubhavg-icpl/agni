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
			<div class="max-w-7xl mx-auto px-3 sm:px-4">
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
								Mission Control
							</a>
							<a
								href="/vms"
								class="nav-link px-3 py-2 rounded-md text-sm font-medium"
								class:nav-link-active={$page.url.pathname.startsWith('/vms')}
							>
								The Zoo
							</a>
							<a
								href="/configs"
								class="nav-link px-3 py-2 rounded-md text-sm font-medium"
								class:nav-link-active={$page.url.pathname.startsWith('/configs')}
							>
								Cookie Cutters
							</a>
						</div>
					</div>
					<div class="flex items-center gap-2 sm:gap-4">
						<a href="/settings" class="hidden sm:block text-gray-400 hover:text-orange-400 text-sm transition-colors">
							{$auth.user?.username}
						</a>
						<button
							on:click={handleLogout}
							class="hidden sm:block text-gray-300 hover:text-orange-400 px-3 py-2 rounded-md text-sm font-medium transition-colors"
						>
							Rage Quit
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
							class="block nav-link px-3 py-2.5 rounded-md text-sm font-medium"
							class:nav-link-active={$page.url.pathname === '/'}
						>
							🎮 Mission Control
						</a>
						<a
							href="/vms"
							class="block nav-link px-3 py-2.5 rounded-md text-sm font-medium"
							class:nav-link-active={$page.url.pathname.startsWith('/vms')}
						>
							🦁 The Zoo
						</a>
						<a
							href="/configs"
							class="block nav-link px-3 py-2.5 rounded-md text-sm font-medium"
							class:nav-link-active={$page.url.pathname.startsWith('/configs')}
						>
							🍪 Cookie Cutters
						</a>
						<a
							href="/settings"
							class="block nav-link px-3 py-2.5 rounded-md text-sm font-medium"
							class:nav-link-active={$page.url.pathname.startsWith('/settings')}
						>
							⚙️ Settings
						</a>
						<div class="border-t border-gray-700/50 pt-2 mt-2">
							<div class="px-3 py-1 text-xs text-gray-500">Signed in as {$auth.user?.username}</div>
							<button
								on:click={handleLogout}
								class="w-full text-left px-3 py-2.5 rounded-md text-sm font-medium text-red-400 hover:bg-red-500/10 transition-colors"
							>
								🚪 Rage Quit
							</button>
						</div>
					</div>
				</div>
			{/if}
		</nav>

		<!-- Main content -->
		<main class="max-w-7xl mx-auto px-3 sm:px-4 py-4 sm:py-8">
			<slot />
		</main>
	</div>
{:else}
	<slot />
{/if}
