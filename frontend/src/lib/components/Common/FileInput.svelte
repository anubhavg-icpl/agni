<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Icon from './Icon.svelte';

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
</script>

<div class="space-y-2">
	{#if label}
		<label for={id} class="label">
			{label}
			{#if required}
				<span class="text-destructive" aria-hidden="true">*</span>
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
				class="input pr-8 font-mono text-[13px] {error ? 'border-destructive focus:ring-destructive/30' : ''}"
			/>
			{#if value}
				<button
					type="button"
					on:click={() => {
						value = '';
						dispatch('change', '');
					}}
					class="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
					aria-label="Clear path"
				>
					<Icon name="ph-x" size={16} />
				</button>
			{/if}
		</div>
		<button type="button" class="btn-secondary h-9 px-3 rounded-md border bg-secondary text-secondary-foreground hover:bg-secondary/70 transition-colors" title="Browse files">
			<Icon name="ph-folder-open" size={16} />
		</button>
	</div>

	{#if error}
		<p class="text-[13px] text-destructive flex items-center gap-1.5">
			<Icon name="ph-warning-circle" size={14} />
			{error}
		</p>
	{:else if hint}
		<p class="text-xs text-muted-foreground">{hint}</p>
	{/if}
</div>
