<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { VMConfig } from '$lib/api/client';

	export let config: Partial<VMConfig>;

	const dispatch = createEventDispatcher();

	let cpus = config.cpus || 1;
	let memoryMB = config.memory_mb || 512;
	let cpuTemplate = config.cpu_template || '';
	let disableSMT = config.disable_smt || false;

	function emitChange() {
		dispatch('change', {
			cpus,
			memory_mb: memoryMB,
			cpu_template: cpuTemplate || undefined,
			disable_smt: disableSMT
		});
	}

	const cpuTemplates = [
		{ value: '', label: 'None (Default)' },
		{ value: 'C3', label: 'C3 (AWS C3 instance compatible)' },
		{ value: 'T2', label: 'T2 (AWS T2 instance compatible)' }
	];

	const memoryPresets = [256, 512, 1024, 2048, 4096, 8192, 16384];
</script>

<div class="space-y-6">
	<h3 class="text-lg font-medium">Resource Allocation</h3>

	<!-- CPUs -->
	<div>
		<label for="cpus" class="label">CPU Cores</label>
		<div class="flex items-center gap-4">
			<input
				type="range"
				id="cpus"
				bind:value={cpus}
				on:change={emitChange}
				min="1"
				max="32"
				class="flex-1 h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer"
			/>
			<input
				type="number"
				bind:value={cpus}
				on:change={emitChange}
				min="1"
				max="32"
				class="input w-20 text-center"
			/>
		</div>
		<p class="text-sm text-gray-500 mt-1">Number of virtual CPU cores (1-32)</p>
	</div>

	<!-- Memory -->
	<div>
		<label for="memory" class="label">Memory (MB)</label>
		<div class="flex items-center gap-4">
			<input
				type="range"
				id="memory"
				bind:value={memoryMB}
				on:change={emitChange}
				min="128"
				max="32768"
				step="128"
				class="flex-1 h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer"
			/>
			<input
				type="number"
				bind:value={memoryMB}
				on:change={emitChange}
				min="128"
				max="32768"
				step="128"
				class="input w-24 text-center"
			/>
		</div>
		<div class="flex gap-2 mt-2">
			{#each memoryPresets as preset}
				<button
					type="button"
					on:click={() => {
						memoryMB = preset;
						emitChange();
					}}
					class="px-2 py-1 text-xs rounded {memoryMB === preset
						? 'bg-blue-600 text-white'
						: 'bg-gray-700 text-gray-300 hover:bg-gray-600'}"
				>
					{preset >= 1024 ? `${preset / 1024}GB` : `${preset}MB`}
				</button>
			{/each}
		</div>
		<p class="text-sm text-gray-500 mt-1">RAM allocation (128MB - 32GB)</p>
	</div>

	<!-- CPU Template -->
	<div>
		<label for="cpuTemplate" class="label">CPU Template</label>
		<select id="cpuTemplate" bind:value={cpuTemplate} on:change={emitChange} class="input w-full">
			{#each cpuTemplates as template}
				<option value={template.value}>{template.label}</option>
			{/each}
		</select>
		<p class="text-sm text-gray-500 mt-1">
			Mask CPUID flags to emulate specific CPU types for compatibility
		</p>
	</div>

	<!-- SMT -->
	<div class="flex items-center gap-3">
		<input
			type="checkbox"
			id="disableSMT"
			bind:checked={disableSMT}
			on:change={emitChange}
			class="rounded bg-gray-700 border-gray-600 text-blue-600 focus:ring-blue-500"
		/>
		<div>
			<label for="disableSMT" class="font-medium">Disable SMT (Hyper-Threading)</label>
			<p class="text-sm text-gray-500">Disable simultaneous multithreading for enhanced security</p>
		</div>
	</div>

	<!-- Summary -->
	<div class="bg-gray-900 rounded-lg p-4">
		<h4 class="text-sm font-medium text-gray-400 mb-2">Resource Summary</h4>
		<div class="grid grid-cols-2 gap-2 text-sm">
			<div>
				<span class="text-gray-500">CPUs:</span>
				<span class="ml-2">{cpus} {cpus === 1 ? 'core' : 'cores'}</span>
			</div>
			<div>
				<span class="text-gray-500">Memory:</span>
				<span class="ml-2"
					>{memoryMB >= 1024 ? `${(memoryMB / 1024).toFixed(1)} GB` : `${memoryMB} MB`}</span
				>
			</div>
			<div>
				<span class="text-gray-500">Template:</span>
				<span class="ml-2">{cpuTemplate || 'Default'}</span>
			</div>
			<div>
				<span class="text-gray-500">SMT:</span>
				<span class="ml-2">{disableSMT ? 'Disabled' : 'Enabled'}</span>
			</div>
		</div>
	</div>
</div>
