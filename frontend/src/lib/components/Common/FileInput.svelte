<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let value = '';
	export let label = '';
	export let placeholder = '/path/to/file';
	export let required = false;
	export let error = '';
	export let hint = '';
	export let id = `file-input-${Math.random().toString(36).slice(2)}`;

	const dispatch = createEventDispatcher();

	function handleInput(event: Event) {
		const target = event.target as HTMLInputElement;
		value = target.value;
		dispatch('change', value);
	}

	function handleBrowse() {
		// In a Wails app, this would trigger a native file dialog
		// For now, we dispatch an event that can be handled by parent components
		dispatch('browse');
	}
</script>

<div class="space-y-1">
	{#if label}
		<label for={id} class="label">
			{label}
			{#if required}
				<span class="text-red-400">*</span>
			{/if}
		</label>
	{/if}

	<div class="flex gap-2">
		<div class="relative flex-1">
			<input
				type="text"
				{id}
				{value}
				{placeholder}
				{required}
				on:input={handleInput}
				on:blur
				class="input w-full pr-10 font-mono text-sm {error
					? 'border-red-500 focus:border-red-500'
					: ''}"
			/>
			{#if value}
				<button
					type="button"
					on:click={() => {
						value = '';
						dispatch('change', '');
					}}
					class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-white"
					aria-label="Clear"
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
			{/if}
		</div>
		<button
			type="button"
			on:click={handleBrowse}
			class="btn btn-secondary px-3"
			title="Browse files"
		>
			<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
				/>
			</svg>
		</button>
	</div>

	{#if error}
		<p class="text-sm text-red-400">{error}</p>
	{:else if hint}
		<p class="text-sm text-gray-500">{hint}</p>
	{/if}
</div>
