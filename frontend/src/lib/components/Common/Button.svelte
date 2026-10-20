<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let variant: 'primary' | 'secondary' | 'ghost' | 'destructive' = 'primary';
	export let size: 'sm' | 'md' | 'lg' = 'md';
	export let disabled = false;
	export let loading = false;
	export let type: 'button' | 'submit' | 'reset' = 'button';
	export let fullWidth = false;
	// When href is set the control renders a link: navigation is links (Law 12)
	export let href: string | undefined = undefined;
	export let title: string | undefined = undefined;

	const dispatch = createEventDispatcher();

	const variantClasses: Record<string, string> = {
		primary: 'bg-primary text-primary-foreground hover:bg-primary/85',
		secondary: 'bg-secondary text-secondary-foreground hover:bg-secondary/70',
		ghost: 'bg-transparent text-foreground hover:bg-accent hover:text-accent-foreground',
		destructive: 'bg-destructive text-destructive-foreground hover:bg-destructive/85'
	};

	const sizeClasses: Record<string, string> = {
		sm: 'h-8 px-3 text-[13px] gap-1.5',
		md: 'h-9 px-4 text-sm gap-2',
		lg: 'h-10 px-4 text-sm gap-2'
	};

	function handleClick(event: MouseEvent) {
		if (!disabled && !loading) {
			dispatch('click', event);
		}
	}
</script>

{#if href}
	<a
		{href}
		{title}
		class="inline-flex items-center justify-center rounded-md font-medium transition-colors duration-150 active:scale-[0.98] {variantClasses[
			variant
		]}
		{sizeClasses[size]} {fullWidth ? 'w-full' : ''}"
	>
		<slot />
	</a>
{:else}
	<button
		{type}
		{title}
		disabled={disabled || loading}
		on:click={handleClick}
		class="inline-flex items-center justify-center rounded-md font-medium transition-colors duration-150 active:scale-[0.98] disabled:opacity-50 disabled:pointer-events-none
		{variantClasses[variant]} {sizeClasses[size]} {fullWidth ? 'w-full' : ''}"
	>
		{#if loading}
			<svg class="animate-spin h-4 w-4" viewBox="0 0 24 24" aria-hidden="true">
				<circle
					class="opacity-25"
					cx="12"
					cy="12"
					r="10"
					stroke="currentColor"
					stroke-width="4"
					fill="none"
				/>
				<path
					class="opacity-75"
					fill="currentColor"
					d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
				/>
			</svg>
		{/if}
		<slot />
	</button>
{/if}
