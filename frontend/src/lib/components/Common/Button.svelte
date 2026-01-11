<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let variant: 'primary' | 'secondary' | 'success' | 'danger' | 'ghost' = 'primary';
	export let size: 'sm' | 'md' | 'lg' = 'md';
	export let disabled = false;
	export let loading = false;
	export let type: 'button' | 'submit' | 'reset' = 'button';
	export let fullWidth = false;

	const dispatch = createEventDispatcher();

	const variantClasses: Record<string, string> = {
		primary: 'bg-blue-600 hover:bg-blue-700 text-white border-blue-600',
		secondary: 'bg-gray-700 hover:bg-gray-600 text-gray-200 border-gray-600',
		success: 'bg-green-600 hover:bg-green-700 text-white border-green-600',
		danger: 'bg-red-600 hover:bg-red-700 text-white border-red-600',
		ghost: 'bg-transparent hover:bg-gray-700 text-gray-300 border-transparent'
	};

	const sizeClasses: Record<string, string> = {
		sm: 'px-2 py-1 text-sm',
		md: 'px-4 py-2',
		lg: 'px-6 py-3 text-lg'
	};

	function handleClick(event: MouseEvent) {
		if (!disabled && !loading) {
			dispatch('click', event);
		}
	}
</script>

<button
	{type}
	disabled={disabled || loading}
	on:click={handleClick}
	class="
		inline-flex items-center justify-center gap-2
		rounded-lg border font-medium
		transition-colors duration-200
		focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-900
		disabled:opacity-50 disabled:cursor-not-allowed
		{variantClasses[variant]}
		{sizeClasses[size]}
		{fullWidth ? 'w-full' : ''}
	"
>
	{#if loading}
		<svg class="animate-spin h-4 w-4" viewBox="0 0 24 24">
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
