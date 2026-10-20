<script lang="ts">
	import Icon from './Icon.svelte';
	import { createEventDispatcher } from 'svelte';

	export let variant: 'info' | 'warning' | 'error' = 'info';
	export let message = '';
	export let dismissible = false;

	const dispatch = createEventDispatcher();

	const meta: Record<string, { icon: string; box: string; text: string }> = {
		info: {
			icon: 'ph-info',
			box: 'bg-primary/10 border-primary/30',
			text: 'text-primary-accent'
		},
		warning: {
			icon: 'ph-warning-circle',
			box: 'bg-status-warn/10 border-status-warn/30',
			text: 'text-status-warn-text'
		},
		error: {
			icon: 'ph-x-circle',
			box: 'bg-status-error/10 border-status-error/30',
			text: 'text-status-error-text'
		}
	};

	$: meta_ = meta[variant];
</script>

<div class="flex items-start gap-3 rounded-lg border p-4 {meta_.box}" role="status">
	<span class="mt-0.5 {meta_.text}"><Icon name={meta_.icon} size={20} /></span>
	<p class="text-sm text-foreground flex-1">{message}</p>
	{#if dismissible}
		<button
			on:click={() => dispatch('dismiss')}
			class="text-muted-foreground hover:text-foreground transition-colors"
			aria-label="Dismiss"
		>
			<Icon name="ph-x" size={16} />
		</button>
	{/if}
</div>
