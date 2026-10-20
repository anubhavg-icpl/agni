<script lang="ts">
	import { createEventDispatcher, onMount, onDestroy } from 'svelte';
	import { fade, scale } from 'svelte/transition';
	import Icon from './Icon.svelte';

	export let open = false;
	export let title = '';
	export let description = '';
	export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';
	// AlertDialog mode: an undoable-in-name-only decision never closes on
	// backdrop or Escape without the explicit confirm or cancel
	export let alert = false;

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
		if (!alert && event.target === event.currentTarget) {
			close();
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && open && !alert) {
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

	$: if (typeof document !== 'undefined') {
		document.body.style.overflow = open ? 'hidden' : '';
	}
</script>

{#if open}
	<div
		class="fixed inset-0 z-[200] bg-foreground/40"
		transition:fade={{ duration: 150 }}
		on:click={handleBackdropClick}
	></div>
	<div class="fixed inset-0 z-[300] flex items-center justify-center p-4 pointer-events-none">
		<div
			class="relative bg-popover border rounded-xl shadow-lg w-full {sizeClasses[size]} max-h-[90vh] overflow-hidden flex flex-col pointer-events-auto"
			transition:scale={{ duration: 200, start: 0.95 }}
			role="dialog"
			aria-modal="true"
			aria-labelledby={title ? 'modal-title' : undefined}
		>
			<div class="flex items-start justify-between px-6 py-4 border-b">
				<div>
					<h2 id="modal-title" class="text-xl font-semibold leading-7 text-balance">{title}</h2>
					{#if description}
						<p class="text-sm text-muted-foreground mt-1">{description}</p>
					{/if}
				</div>
				{#if !alert}
					<button
						on:click={close}
						class="text-muted-foreground hover:text-foreground transition-colors p-1 rounded-sm hover:bg-accent"
						aria-label="Close dialog"
					>
						<Icon name="ph-x" size={20} />
					</button>
				{/if}
			</div>
			<div class="px-6 py-4 overflow-y-auto flex-1">
				<slot />
			</div>
			{#if $$slots.footer}
				<div class="px-6 py-4 border-t bg-muted/50 flex justify-end gap-2">
					<slot name="footer" />
				</div>
			{/if}
		</div>
	</div>
{/if}
