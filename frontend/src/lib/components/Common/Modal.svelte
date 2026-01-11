<script lang="ts">
	import { createEventDispatcher, onMount, onDestroy } from 'svelte';
	import { fade, scale } from 'svelte/transition';

	export let open = false;
	export let title = '';
	export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';
	export let closeOnBackdrop = true;
	export let closeOnEscape = true;

	const dispatch = createEventDispatcher();

	const sizeClasses: Record<string, string> = {
		sm: 'max-w-sm',
		md: 'max-w-md',
		lg: 'max-w-lg',
		xl: 'max-w-xl'
	};

	function close() {
		dispatch('close');
	}

	function handleBackdropClick(event: MouseEvent) {
		if (closeOnBackdrop && event.target === event.currentTarget) {
			close();
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (closeOnEscape && event.key === 'Escape' && open) {
			close();
		}
	}

	onMount(() => {
		if (typeof window !== 'undefined') {
			window.addEventListener('keydown', handleKeydown);
		}
	});

	onDestroy(() => {
		if (typeof window !== 'undefined') {
			window.removeEventListener('keydown', handleKeydown);
		}
	});

	$: if (open && typeof document !== 'undefined') {
		document.body.style.overflow = 'hidden';
	} else if (typeof document !== 'undefined') {
		document.body.style.overflow = '';
	}
</script>

{#if open}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center p-4"
		transition:fade={{ duration: 150 }}
		on:click={handleBackdropClick}
		role="dialog"
		aria-modal="true"
		aria-labelledby={title ? 'modal-title' : undefined}
	>
		<!-- Backdrop -->
		<div class="absolute inset-0 bg-black/70" />

		<!-- Modal content -->
		<div
			class="relative bg-gray-800 rounded-lg shadow-xl w-full {sizeClasses[
				size
			]} max-h-[90vh] overflow-hidden flex flex-col"
			transition:scale={{ duration: 150, start: 0.95 }}
		>
			<!-- Header -->
			{#if title || $$slots.header}
				<div class="flex items-center justify-between px-6 py-4 border-b border-gray-700">
					{#if $$slots.header}
						<slot name="header" />
					{:else}
						<h2 id="modal-title" class="text-lg font-semibold text-white">{title}</h2>
					{/if}
					<button
						on:click={close}
						class="text-gray-400 hover:text-white transition-colors p-1 rounded hover:bg-gray-700"
						aria-label="Close modal"
					>
						<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M6 18L18 6M6 6l12 12"
							/>
						</svg>
					</button>
				</div>
			{/if}

			<!-- Body -->
			<div class="px-6 py-4 overflow-y-auto flex-1">
				<slot />
			</div>

			<!-- Footer -->
			{#if $$slots.footer}
				<div class="px-6 py-4 border-t border-gray-700 bg-gray-800/50">
					<slot name="footer" />
				</div>
			{/if}
		</div>
	</div>
{/if}
