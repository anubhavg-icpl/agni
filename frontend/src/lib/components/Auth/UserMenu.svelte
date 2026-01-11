<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { fade, scale } from 'svelte/transition';
	import { auth } from '$lib/stores/auth';

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

	function handleSettings() {
		close();
		dispatch('settings');
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.user-menu')) {
			close();
		}
	}
</script>

<svelte:window on:click={handleClickOutside} />

<div class="relative user-menu">
	<button
		on:click={toggle}
		class="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-gray-700 transition-colors"
	>
		<div class="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center text-sm font-medium text-white">
			{username.charAt(0).toUpperCase()}
		</div>
		<span class="text-sm text-gray-300 hidden sm:block">{username}</span>
		<svg
			class="w-4 h-4 text-gray-400 transition-transform {open ? 'rotate-180' : ''}"
			fill="none"
			viewBox="0 0 24 24"
			stroke="currentColor"
		>
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
		</svg>
	</button>

	{#if open}
		<div
			transition:scale={{ duration: 150, start: 0.95 }}
			class="absolute right-0 mt-2 w-48 bg-gray-800 border border-gray-700 rounded-lg shadow-lg overflow-hidden z-50"
		>
			<div class="px-4 py-3 border-b border-gray-700">
				<p class="text-sm font-medium text-white">{username}</p>
				<p class="text-xs text-gray-400">Administrator</p>
			</div>

			<div class="py-1">
				<button
					on:click={handleSettings}
					class="w-full px-4 py-2 text-left text-sm text-gray-300 hover:bg-gray-700 flex items-center gap-2"
				>
					<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
						/>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
						/>
					</svg>
					Settings
				</button>

				<button
					on:click={handleLogout}
					class="w-full px-4 py-2 text-left text-sm text-red-400 hover:bg-gray-700 flex items-center gap-2"
				>
					<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
						/>
					</svg>
					Sign Out
				</button>
			</div>
		</div>
	{/if}
</div>
