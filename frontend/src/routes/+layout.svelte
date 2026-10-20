<script lang="ts">
	import '@phosphor-icons/web/bold';
	import '../app.css';
	import { onMount } from 'svelte';
	import { fly, fade } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import { theme, toggleTheme } from '$lib/stores/theme';
	import Sidebar from '$lib/components/Shell/Sidebar.svelte';
	import UserMenu from '$lib/components/Auth/UserMenu.svelte';
	import Icon from '$lib/components/Common/Icon.svelte';
	import Toast from '$lib/components/Common/Toast.svelte';

	let initialized = false;
	let collapsed = false;
	let mobileOpen = false;

	const SIDEBAR_KEY = 'agni-sidebar-collapsed';

	onMount(async () => {
		collapsed = window.localStorage.getItem(SIDEBAR_KEY) === '1';
		if (!initialized) {
			initialized = true;
			await auth.init();
		}
	});

	function toggleCollapsed() {
		collapsed = !collapsed;
		window.localStorage.setItem(SIDEBAR_KEY, collapsed ? '1' : '0');
	}

	function closeMobile() {
		mobileOpen = false;
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') closeMobile();
	}

	function handleLogout() {
		auth.logout();
		goto('/login');
	}

	// Route unauthenticated visitors to the login screen, which also owns
	// first-run setup
	$: if (
		initialized &&
		!$auth.loading &&
		(!$isAuthenticated || $auth.setupRequired) &&
		!$page.url.pathname.includes('/login')
	) {
		goto('/login');
	}

	// Close the mobile sheet on navigation
	$: if ($page.url.pathname) {
		mobileOpen = false;
	}

	$: sidebarWidth = collapsed ? '76px' : '262px';
</script>

<svelte:window on:keydown={handleKeydown} />

<Toast />

{#if !initialized || $auth.loading}
	<div class="min-h-dvh bg-background flex flex-col items-center justify-center gap-4">
		<img src="/logo.png" alt="Agni" class="w-10 h-10" />
		<div class="skeleton h-4 w-28"></div>
	</div>
{:else if $auth.setupRequired || $page.url.pathname.includes('/login')}
	<slot />
{:else if $isAuthenticated}
	<div class="min-h-dvh">
		<!-- Mobile backdrop -->
		{#if mobileOpen}
			<div
				class="fixed inset-0 z-[200] bg-foreground/40 md:hidden"
				transition:fade={{ duration: 150 }}
				on:click={closeMobile}
				aria-hidden="true"
			></div>
		{/if}

		<!-- Desktop sidebar -->
		<aside
			class="hidden md:flex fixed inset-y-0 left-0 z-40 bg-sidebar border-r transition-[width] duration-200 ease-out"
			style="width: {sidebarWidth}"
		>
			<Sidebar {collapsed} />
		</aside>

		<!-- Collapse circle on the rail edge -->
		<button
			class="hidden md:flex fixed top-[72px] z-[45] h-7 w-7 items-center justify-center rounded-full bg-card border shadow-sm text-muted-foreground hover:text-foreground transition-colors"
			style="left: calc({sidebarWidth} - 14px)"
			on:click={toggleCollapsed}
			aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
			aria-expanded={!collapsed}
		>
			<Icon name={collapsed ? 'ph-caret-right' : 'ph-caret-left'} size={14} />
		</button>

		<!-- Mobile sheet -->
		{#if mobileOpen}
			<div
				class="fixed inset-y-0 left-0 z-[201] w-[262px] bg-sidebar border-r md:hidden"
				transition:fly={{ x: -262, duration: 200 }}
			>
				<Sidebar collapsed={false} on:navigate={closeMobile} />
			</div>
		{/if}

		<div
			class="flex flex-col min-h-dvh transition-[margin] duration-200 ease-out {collapsed
				? 'md:ml-[76px]'
				: 'md:ml-[262px]'}"
		>
			<!-- Topbar: left side empty, controls right, user chip last -->
			<header
				class="sticky top-0 z-50 h-[58px] shrink-0 border-b bg-background flex items-center justify-between px-4 gap-2"
			>
				<button
					class="md:hidden h-9 w-9 inline-flex items-center justify-center rounded-md text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
					on:click={() => (mobileOpen = true)}
					aria-label="Open navigation"
					aria-expanded={mobileOpen}
				>
					<Icon name="ph-list" size={20} />
				</button>
				<div class="flex items-center gap-1.5 ml-auto">
					<button
						class="h-9 w-9 inline-flex items-center justify-center rounded-md text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
						on:click={toggleTheme}
						aria-label={$theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}
						title={$theme === 'dark' ? 'Light theme' : 'Dark theme'}
					>
						<Icon name={$theme === 'dark' ? 'ph-sun' : 'ph-moon'} size={20} />
					</button>
					<UserMenu username={$auth.user?.username ?? ''} on:logout={handleLogout} />
				</div>
			</header>

			<main class="flex-1 px-4 py-6 md:px-[30px] md:py-[26px]">
				<div class="space-y-6">
					<slot />
				</div>
			</main>
		</div>
	</div>
{:else}
	<slot />
{/if}
