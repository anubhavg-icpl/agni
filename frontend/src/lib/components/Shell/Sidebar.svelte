<script lang="ts">
	import { page } from '$app/stores';
	import { createEventDispatcher } from 'svelte';
	import Icon from '../Common/Icon.svelte';
	import { APP_VERSION } from '$lib/version';

	export let collapsed = false;

	const dispatch = createEventDispatcher();

	interface NavItem {
		href: string;
		label: string;
		icon: string;
		exact?: boolean;
	}

	const groups: { section: string; items: NavItem[] }[] = [
		{
			section: 'Overview',
			items: [{ href: '/', label: 'Dashboard', icon: 'ph-squares-four', exact: true }]
		},
		{
			section: 'Manage',
			items: [
				{ href: '/vms', label: 'Virtual Machines', icon: 'ph-cubes' },
				{ href: '/configs', label: 'Templates', icon: 'ph-files' }
			]
		},
		{
			section: 'System',
			items: [{ href: '/settings', label: 'Settings', icon: 'ph-gear-six' }]
		}
	];

	function isActive(item: NavItem, pathname: string): boolean {
		return item.exact ? pathname === item.href : pathname.startsWith(item.href);
	}

	function handleNavigate() {
		dispatch('navigate');
	}
</script>

<div class="flex flex-col h-full min-h-0">
	<!-- Brand row, 58px, matches the topbar -->
	<a href="/" class="h-[58px] shrink-0 flex items-center gap-2.5 px-4" on:click={handleNavigate}>
		<img src="/logo.png" alt="Agni" class="w-[30px] h-[30px] shrink-0" />
		{#if !collapsed}
			<span class="text-base font-semibold text-foreground">Agni</span>
		{/if}
	</a>

	<nav class="flex-1 overflow-y-auto py-2" aria-label="Main">
		{#each groups as group}
			<div class={collapsed ? 'px-2' : 'px-3'}>
				{#if !collapsed}
					<p class="text-xs font-medium text-muted-foreground pt-4 pb-1.5 px-2.5">
						{group.section}
					</p>
				{:else}
					<div class="border-t my-3 mx-1"></div>
				{/if}
				<ul class="space-y-0.5">
					{#each group.items as item}
						{@const active = isActive(item, $page.url.pathname)}
						<li>
							<a
								href={item.href}
								class="h-9 flex items-center gap-2.5 rounded-md text-[13px] font-medium transition-colors duration-150
								{collapsed ? 'justify-center w-full' : 'px-2.5'}
								{active
									? 'bg-sidebar-accent text-primary-accent font-semibold'
									: 'text-muted-foreground hover:bg-accent hover:text-foreground'}"
								aria-current={active ? 'page' : undefined}
								aria-label={collapsed ? item.label : undefined}
								title={collapsed ? item.label : undefined}
								on:click={handleNavigate}
							>
								<Icon name={item.icon} size={16} />
								{#if !collapsed}
									<span class="truncate">{item.label}</span>
								{/if}
							</a>
						</li>
					{/each}
				</ul>
			</div>
		{/each}
	</nav>

	<!-- Foot pinned above a hairline -->
	<div class="shrink-0 border-t px-4 py-3">
		{#if collapsed}
			<span class="block text-center font-mono text-xs text-muted-foreground" title={APP_VERSION}>
				{APP_VERSION.slice(0, 2)}
			</span>
		{:else}
			<span class="font-mono text-xs text-muted-foreground">{APP_VERSION}</span>
		{/if}
	</div>
</div>
