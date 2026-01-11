<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { VMConfig } from '$lib/api/client';
	import FileInput from '../Common/FileInput.svelte';

	export let config: Partial<VMConfig>;

	const dispatch = createEventDispatcher();

	let kernelPath = config.kernel_path || './vmlinux';
	let kernelOpts =
		config.kernel_opts || 'ro console=ttyS0 noapic reboot=k panic=1 pci=off nomodules';
	let initrdPath = config.initrd_path || '';
	let metadata = config.metadata || '';
	let logLevel = config.log_level || 'Debug';

	function emitChange() {
		dispatch('change', {
			kernel_path: kernelPath,
			kernel_opts: kernelOpts,
			initrd_path: initrdPath || undefined,
			metadata: metadata || undefined,
			log_level: logLevel
		});
	}

	function applyPreset(value: string) {
		kernelOpts = value;
		emitChange();
	}

	const logLevels = ['Error', 'Warning', 'Info', 'Debug'];

	const defaultKernelOpts = [
		{ label: 'Minimal', value: 'ro console=ttyS0 noapic reboot=k panic=1 pci=off nomodules' },
		{
			label: 'Standard',
			value: 'ro console=ttyS0 noapic reboot=k panic=1 pci=off nomodules init=/sbin/init'
		},
		{ label: 'Debug', value: 'ro console=ttyS0 noapic reboot=k panic=1 pci=off nomodules debug' }
	];
</script>

<div class="space-y-6">
	<h3 class="text-lg font-medium">Kernel Configuration</h3>

	<!-- Kernel Path -->
	<FileInput
		label="Kernel Image Path"
		bind:value={kernelPath}
		on:change={emitChange}
		placeholder="./vmlinux"
		hint="Path to the uncompressed Linux kernel binary"
	/>

	<!-- Kernel Options -->
	<div>
		<label for="kernelOpts" class="label">Kernel Boot Options</label>
		<textarea
			id="kernelOpts"
			bind:value={kernelOpts}
			on:input={emitChange}
			class="input w-full h-20 font-mono text-sm"
			placeholder="ro console=ttyS0 ..."
		/>
		<div class="flex flex-wrap gap-2 mt-2">
			{#each defaultKernelOpts as preset}
				<button
					type="button"
					on:click={() => applyPreset(preset.value)}
					class="px-2 py-1 text-xs rounded bg-gray-700 text-gray-300 hover:bg-gray-600"
				>
					{preset.label}
				</button>
			{/each}
		</div>
	</div>

	<!-- Initrd Path -->
	<FileInput
		label="Initial Ramdisk (initrd) Path"
		bind:value={initrdPath}
		on:change={emitChange}
		placeholder="/path/to/initrd.img"
		hint="Optional: Path to initrd/initramfs image"
	/>

	<!-- Log Level -->
	<div>
		<label for="logLevel" class="label">Firecracker Log Level</label>
		<select id="logLevel" bind:value={logLevel} on:change={emitChange} class="input w-full">
			{#each logLevels as level}
				<option value={level}>{level}</option>
			{/each}
		</select>
	</div>

	<!-- Metadata -->
	<div>
		<label for="metadata" class="label">Instance Metadata (JSON)</label>
		<textarea
			id="metadata"
			bind:value={metadata}
			on:input={emitChange}
			class="input w-full h-24 font-mono text-sm"
			placeholder={'{"key": "value"}'}
		/>
		<p class="text-sm text-gray-500 mt-1">
			Optional JSON metadata accessible via MMDS (Microvm Metadata Service)
		</p>
	</div>

	<!-- Kernel Summary -->
	<div class="bg-gray-900 rounded-lg p-4">
		<h4 class="text-sm font-medium text-gray-400 mb-2">Kernel Summary</h4>
		<div class="space-y-2 text-sm">
			<div class="flex">
				<span class="text-gray-500 w-24">Kernel:</span>
				<span class="font-mono text-xs truncate" title={kernelPath}>{kernelPath}</span>
			</div>
			{#if initrdPath}
				<div class="flex">
					<span class="text-gray-500 w-24">Initrd:</span>
					<span class="font-mono text-xs truncate" title={initrdPath}>{initrdPath}</span>
				</div>
			{/if}
			<div class="flex">
				<span class="text-gray-500 w-24">Log Level:</span>
				<span>{logLevel}</span>
			</div>
			<div>
				<span class="text-gray-500">Boot Options:</span>
				<div class="font-mono text-xs bg-gray-800 p-2 rounded mt-1 break-all">
					{kernelOpts}
				</div>
			</div>
		</div>
	</div>
</div>
