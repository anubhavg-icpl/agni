<script lang="ts">
	import { goto } from '$app/navigation';
	import { vms } from '$lib/stores/vms';
	import type { VMConfig } from '$lib/api/client';

	let name = '';
	let kernelPath = './vmlinux';
	let kernelOpts = 'ro console=ttyS0 noapic reboot=k panic=1 pci=off nomodules';
	let rootDrivePath = '';
	let rootDriveReadOnly = false;
	let cpus = 1;
	let memoryMB = 512;
	let cpuTemplate = '';
	let disableSMT = false;
	let logLevel = 'Debug';

	let loading = false;
	let error = '';

	async function handleSubmit() {
		if (!name || !rootDrivePath) {
			error = 'Name and root drive path are required';
			return;
		}

		loading = true;
		error = '';

		const config: VMConfig = {
			name,
			kernel_path: kernelPath,
			kernel_opts: kernelOpts,
			root_drive: {
				path: rootDrivePath,
				read_only: rootDriveReadOnly
			},
			cpus,
			memory_mb: memoryMB,
			cpu_template: cpuTemplate || undefined,
			disable_smt: disableSMT,
			log_level: logLevel
		};

		const vm = await vms.create(name, config);
		loading = false;

		if (vm) {
			goto('/');
		} else {
			error = $vms.error || 'Failed to summon demon. You lack mana.';
		}
	}
</script>

<svelte:head>
	<title>New VM - Firectl</title>
</svelte:head>

<div class="max-w-2xl mx-auto">
	<h1 class="text-2xl font-bold mb-6">Summon New Demon</h1>

	<form on:submit|preventDefault={handleSubmit} class="space-y-6">
		<!-- Basic Info -->
		<div class="card">
			<h2 class="text-lg font-medium mb-4">The Basics (Don't mess this up)</h2>
			<div class="space-y-4">
				<div>
					<label class="label" for="name">Demon Name</label>
					<input
						type="text"
						id="name"
						bind:value={name}
						class="input w-full"
						placeholder="my-little-demon"
						required
					/>
				</div>
			</div>
		</div>

		<!-- Kernel Configuration -->
		<div class="card">
			<h2 class="text-lg font-medium mb-4">Brain (Kernel)</h2>
			<div class="space-y-4">
				<div>
					<label class="label" for="kernelPath">Kernel Image Path</label>
					<input
						type="text"
						id="kernelPath"
						bind:value={kernelPath}
						class="input w-full"
						placeholder="/path/to/vmlinux"
					/>
				</div>
				<div>
					<label class="label" for="kernelOpts">Kernel Options</label>
					<input type="text" id="kernelOpts" bind:value={kernelOpts} class="input w-full" />
				</div>
			</div>
		</div>

		<!-- Storage -->
		<div class="card">
			<h2 class="text-lg font-medium mb-4">Soul Container (Storage)</h2>
			<div class="space-y-4">
				<div>
					<label class="label" for="rootDrivePath">Root Drive Path</label>
					<input
						type="text"
						id="rootDrivePath"
						bind:value={rootDrivePath}
						class="input w-full"
						placeholder="/path/to/rootfs.ext4"
						required
					/>
				</div>
				<div class="flex items-center gap-2">
					<input
						type="checkbox"
						id="rootDriveReadOnly"
						bind:checked={rootDriveReadOnly}
						class="rounded bg-gray-700 border-gray-600"
					/>
					<label for="rootDriveReadOnly" class="text-sm text-gray-300">Look but don't touch (Read Only)</label>
				</div>
			</div>
		</div>

		<!-- Resources -->
		<div class="card">
			<h2 class="text-lg font-medium mb-4">Sacrifices (Resources)</h2>
			<div class="grid grid-cols-2 gap-4">
				<div>
					<label class="label" for="cpus">Brain Cells (CPUs)</label>
					<input type="number" id="cpus" bind:value={cpus} class="input w-full" min="1" max="32" />
				</div>
				<div>
					<label class="label" for="memoryMB">Short-term Memory (MB)</label>
					<input
						type="number"
						id="memoryMB"
						bind:value={memoryMB}
						class="input w-full"
						min="128"
						max="32768"
						step="128"
					/>
				</div>
			</div>
			<div class="mt-4 grid grid-cols-2 gap-4">
				<div>
					<label class="label" for="cpuTemplate">CPU Template</label>
					<select id="cpuTemplate" bind:value={cpuTemplate} class="input w-full">
						<option value="">None</option>
						<option value="C3">C3</option>
						<option value="T2">T2</option>
					</select>
				</div>
				<div class="flex items-end">
					<div class="flex items-center gap-2">
						<input
							type="checkbox"
							id="disableSMT"
							bind:checked={disableSMT}
							class="rounded bg-gray-700 border-gray-600"
						/>
						<label for="disableSMT" class="text-sm text-gray-300">Disable SMT</label>
					</div>
				</div>
			</div>
		</div>

		<!-- Error -->
		{#if error}
			<div class="bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg p-4">
				{error}
			</div>
		{/if}

		<!-- Actions -->
		<div class="flex gap-4">
			<button type="submit" class="btn btn-primary flex-1" disabled={loading}>
				{loading ? 'Summoning...' : 'Summon'}
			</button>
			<a href="/" class="btn btn-secondary">Abort Ritual</a>
		</div>
	</form>
</div>
