<script lang="ts" context="module">
	import { writable } from 'svelte/store';

	export interface ToastMessage {
		id: string;
		type: 'success' | 'error' | 'warning' | 'info';
		message: string;
		duration?: number;
	}

	function createToastStore() {
		const { subscribe, update } = writable<ToastMessage[]>([]);

		function add(type: ToastMessage['type'], message: string, duration = 5000) {
			const id = Math.random().toString(36).slice(2);
			const toast: ToastMessage = { id, type, message, duration };

			update((toasts) => [...toasts, toast]);

			if (duration > 0) {
				setTimeout(() => {
					remove(id);
				}, duration);
			}

			return id;
		}

		function remove(id: string) {
			update((toasts) => toasts.filter((t) => t.id !== id));
		}

		return {
			subscribe,
			success: (message: string, duration?: number) => add('success', message, duration),
			error: (message: string, duration?: number) => add('error', message, duration),
			warning: (message: string, duration?: number) => add('warning', message, duration),
			info: (message: string, duration?: number) => add('info', message, duration),
			remove
		};
	}

	export const toasts = createToastStore();
</script>

<script lang="ts">
	import { fly, fade } from 'svelte/transition';
	import { flip } from 'svelte/animate';

	const typeStyles: Record<string, { bg: string; border: string; icon: string }> = {
		success: {
			bg: 'bg-green-500/20',
			border: 'border-green-500/50',
			icon: 'M5 13l4 4L19 7'
		},
		error: {
			bg: 'bg-red-500/20',
			border: 'border-red-500/50',
			icon: 'M6 18L18 6M6 6l12 12'
		},
		warning: {
			bg: 'bg-yellow-500/20',
			border: 'border-yellow-500/50',
			icon: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z'
		},
		info: {
			bg: 'bg-blue-500/20',
			border: 'border-blue-500/50',
			icon: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
		}
	};

	const typeTextColors: Record<string, string> = {
		success: 'text-green-400',
		error: 'text-red-400',
		warning: 'text-yellow-400',
		info: 'text-blue-400'
	};
</script>

<div class="fixed top-4 right-4 z-50 flex flex-col gap-2 max-w-sm w-full pointer-events-none">
	{#each $toasts as toast (toast.id)}
		<div
			animate:flip={{ duration: 200 }}
			in:fly={{ x: 100, duration: 200 }}
			out:fade={{ duration: 150 }}
			class="pointer-events-auto {typeStyles[toast.type].bg} border {typeStyles[toast.type]
				.border} rounded-lg p-4 shadow-lg"
		>
			<div class="flex items-start gap-3">
				<svg
					class="w-5 h-5 flex-shrink-0 {typeTextColors[toast.type]}"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d={typeStyles[toast.type].icon}
					/>
				</svg>
				<p class="text-sm text-gray-200 flex-1">{toast.message}</p>
				<button
					on:click={() => toasts.remove(toast.id)}
					class="text-gray-400 hover:text-white transition-colors"
					aria-label="Dismiss"
				>
					<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						/>
					</svg>
				</button>
			</div>
		</div>
	{/each}
</div>
