<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { scale } from 'svelte/transition';
	import Icon from '../Common/Icon.svelte';

	export let username = '';

	let open = false;
	const dispatch = createEventDispatcher();

	function toggle() {
		open = !open;
	}

	function close() {
		open = false;
	}

	function handleLogout() {
		close();
		dispatch('logout');
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') close();
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.user-menu')) {
			close();
		}
	}
</script>

<svelte:window on:click={handleClickOutside} on:keydown={handleKeydown} />

<div class="relative user-menu">
	<button
		on:click={toggle}
		class="flex items-center gap-2 px-2 py-1.5 rounded-md hover:bg-accent transition-colors"
		aria-haspopup="menu"
		aria-expanded={open}
	>
		<span
			class="w-7 h-7 rounded-full bg-primary text-primary-foreground flex items-center justify-center text-[13px] font-medium"
			aria-hidden="true"
		>
			{username.charAt(0).toUpperCase()}
		</span>
		<span class="text-[13px] font-medium text-foreground hidden sm:block">{username}</span>
		<span
			class="text-muted-foreground transition-transform duration-150 {open ? 'rotate-180' : ''}"
			aria-hidden="true"
		>
			<Icon name="ph-caret-down" size={14} />
		</span>
	</button>

	{#if open}
		<div
			transition:scale={{ duration: 150, start: 0.95 }}
			class="absolute right-0 mt-2 w-52 bg-popover border rounded-lg shadow-lg overflow-hidden z-[100]"
			role="menu"
		>
			<div class="px-4 py-3 border-b">
				<p class="text-sm font-medium text-foreground">{username}</p>
				<p class="text-xs text-muted-foreground mt-0.5">Administrator</p>
			</div>
			<div class="py-1">
				<a
					href="/settings"
					class="w-full px-4 py-2 text-left text-sm text-foreground hover:bg-accent flex items-center gap-2"
					role="menuitem"
					on:click={close}
				>
					<span class="text-muted-foreground"><Icon name="ph-gear-six" size={16} /></span>
					Settings
				</a>
				<button
					on:click={handleLogout}
					class="w-full px-4 py-2 text-left text-sm text-destructive hover:bg-accent flex items-center gap-2"
					role="menuitem"
				>
					<Icon name="ph-sign-out" size={16} />
					Sign out
				</button>
			</div>
		</div>
	{/if}
</div>
