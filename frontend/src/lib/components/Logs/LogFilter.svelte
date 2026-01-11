<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	let level = 'all';
	let search = '';
	let searchTimeout: ReturnType<typeof setTimeout>;

	const dispatch = createEventDispatcher();

	const levels = [
		{ value: 'all', label: 'All Levels' },
		{ value: 'debug', label: 'Debug' },
		{ value: 'info', label: 'Info' },
		{ value: 'warn', label: 'Warning' },
		{ value: 'error', label: 'Error' }
	];

	function emitChange() {
		dispatch('change', { level, search });
	}

	function handleLevelChange() {
		emitChange();
	}

	function handleSearchInput() {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			emitChange();
		}, 300);
	}

	function clearSearch() {
		search = '';
		emitChange();
	}
</script>

<div class="flex flex-wrap items-center gap-3 mb-4">
	<!-- Level filter -->
	<div class="flex items-center gap-2">
		<label for="log-level" class="text-sm text-gray-400">Level:</label>
		<select
			id="log-level"
			bind:value={level}
			on:change={handleLevelChange}
			class="input py-1 px-2 text-sm"
		>
			{#each levels as levelOption}
				<option value={levelOption.value}>{levelOption.label}</option>
			{/each}
		</select>
	</div>

	<!-- Search -->
	<div class="flex-1 min-w-[200px] relative">
		<input
			type="text"
			bind:value={search}
			on:input={handleSearchInput}
			placeholder="Search logs..."
			class="input w-full py-1 px-3 pr-8 text-sm"
		/>
		{#if search}
			<button
				on:click={clearSearch}
				class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-white"
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

	<!-- Quick filters -->
	<div class="flex items-center gap-1">
		<button
			on:click={() => {
				level = 'error';
				emitChange();
			}}
			class="px-2 py-1 text-xs rounded {level === 'error'
				? 'bg-red-600 text-white'
				: 'bg-gray-700 text-gray-300 hover:bg-gray-600'}"
		>
			Errors Only
		</button>
		<button
			on:click={() => {
				level = 'warn';
				emitChange();
			}}
			class="px-2 py-1 text-xs rounded {level === 'warn'
				? 'bg-yellow-600 text-white'
				: 'bg-gray-700 text-gray-300 hover:bg-gray-600'}"
		>
			Warnings
		</button>
	</div>
</div>
