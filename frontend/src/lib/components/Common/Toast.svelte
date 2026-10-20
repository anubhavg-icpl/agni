<script lang="ts" context="module">
	import { writable } from 'svelte/store';

	export interface ToastMessage {
		id: string;
		type: 'success' | 'error' | 'warning' | 'info';
		message: string;
	}

	// Severity decides the duration, never the caller (Law 16)
	const DURATIONS: Record<ToastMessage['type'], number> = {
		success: 4000,
		info: 4000,
		warning: 7000,
		error: 0
	};

	// Three visible, the fourth queues
	const MAX_VISIBLE = 3;

	function createToastStore() {
		const { subscribe, update } = writable<ToastMessage[]>([]);
		const timers = new Map<string, ReturnType<typeof setTimeout>>();

		function remove(id: string) {
			const timer = timers.get(id);
			if (timer) {
				clearTimeout(timer);
				timers.delete(id);
			}
			update((toasts) => toasts.filter((t) => t.id !== id));
		}

		function add(type: ToastMessage['type'], message: string) {
			const id = Math.random().toString(36).slice(2);
			update((toasts) => {
				const next = [...toasts, { id, type, message }];
				return next.slice(Math.max(0, next.length - MAX_VISIBLE));
			});
			const duration = DURATIONS[type];
			if (duration > 0) {
				timers.set(
					id,
					setTimeout(() => remove(id), duration)
				);
			}
			return id;
		}

		return {
			subscribe,
			success: (message: string) => add('success', message),
			error: (message: string) => add('error', message),
			warning: (message: string) => add('warning', message),
			info: (message: string) => add('info', message),
			remove
		};
	}

	export const toasts = createToastStore();
</script>

<script lang="ts">
	import { fly, fade } from 'svelte/transition';
	import { flip } from 'svelte/animate';
	import Icon from './Icon.svelte';

	const typeMeta: Record<string, { icon: string; text: string }> = {
		success: { icon: 'ph-check-circle', text: 'text-status-ok-text' },
		error: { icon: 'ph-x-circle', text: 'text-status-error-text' },
		warning: { icon: 'ph-warning-circle', text: 'text-status-warn-text' },
		info: { icon: 'ph-info', text: 'text-primary-accent' }
	};
</script>

<div
	class="fixed z-[400] flex flex-col gap-2 w-[calc(100vw-32px)] max-w-sm pointer-events-none bottom-4 right-4 max-[640px]:top-4 max-[640px]:bottom-auto max-[640px]:left-1/2 max-[640px]:-translate-x-1/2 max-[640px]:right-auto"
>
	{#each $toasts as toast (toast.id)}
		<div
			animate:flip={{ duration: 200 }}
			in:fly={{ y: 8, duration: 200 }}
			out:fade={{ duration: 160 }}
			class="pointer-events-auto bg-popover border rounded-lg shadow-lg p-4"
			role="status"
		>
			<div class="flex items-start gap-3">
				<span class="mt-0.5 {typeMeta[toast.type].text}"
					><Icon name={typeMeta[toast.type].icon} size={20} /></span
				>
				<p class="text-sm text-foreground flex-1 break-words">{toast.message}</p>
				<button
					on:click={() => toasts.remove(toast.id)}
					class="text-muted-foreground hover:text-foreground transition-colors"
					aria-label="Dismiss notification"
				>
					<Icon name="ph-x" size={16} />
				</button>
			</div>
		</div>
	{/each}
</div>
