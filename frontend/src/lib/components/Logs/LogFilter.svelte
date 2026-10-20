<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Icon from '../Common/Icon.svelte';

	let level = 'all';
	let search = '';
	let searchTimeout: ReturnType<typeof setTimeout>;

	const dispatch = createEventDispatcher();

	const levels = [
		{ value: 'all', label: 'All levels' },
		{ value: 'debug', label: 'Debug' },
		{ value: 'info', label: 'Info' },
		{ value: 'warn', label: 'Warning' },
		{ value: 'error', label: 'Error' }
	];

	function emitChange() {
		dispatch('change', { level, search });
	}

	function handleSearchInput() {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(emitChange, 300);
	}

	function clearSearch() {
		search = '';
		emitChange();
	}
</script>

<div class="flex flex-wrap items-center gap-3 mb-4">
	<div class="flex items-center gap-2">
		<label for="log-level" class="text-sm text-muted-foreground">Level</label>
		<select
			id="log-level"
			bind:value={level}
			on:change={emitChange}
			class="input h-8 w-auto py-0 px-2 text-[13px]"
		>
			{#each levels as levelOption}
				<option value={levelOption.value}>{levelOption.label}</option>
			{/each}
		</select>
	</div>

	<div class="flex-1 min-w-[200px] relative">
		<input
			type="text"
			bind:value={search}
			on:input={handleSearchInput}
			placeholder="Search logs"
			class="input h-8 py-0 pr-8 text-[13px]"
		/>
		{#if search}
			<button
				on:click={clearSearch}
				class="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
				aria-label="Clear search"
			>
				<Icon name="ph-x" size={16} />
			</button>
		{/if}
	</div>

	<div class="flex items-center gap-1.5">
		<button
			type="button"
			on:click={() => {
				level = 'error';
				emitChange();
			}}
			class="h-8 px-2.5 text-xs font-medium rounded-sm transition-colors {level === 'error'
				? 'bg-primary text-primary-foreground'
				: 'bg-secondary text-secondary-foreground hover:bg-secondary/70'}"
		>
			Errors only
		</button>
		<button
			type="button"
			on:click={() => {
				level = 'warn';
				emitChange();
			}}
			class="h-8 px-2.5 text-xs font-medium rounded-sm transition-colors {level === 'warn'
				? 'bg-primary text-primary-foreground'
				: 'bg-secondary text-secondary-foreground hover:bg-secondary/70'}"
		>
			Warnings
		</button>
	</div>
</div>
