<script lang="ts">
	import { goto } from '$app/navigation';
	import { vms } from '$lib/stores/vms';
	import type { VMConfig } from '$lib/api/client';
	import { validate, vmName, filePath, cpuCount, memoryMB } from '$lib/utils/validation';
	import PageHeader from '$lib/components/Common/PageHeader.svelte';
	import Button from '$lib/components/Common/Button.svelte';
	import Icon from '$lib/components/Common/Icon.svelte';
	import StatusBanner from '$lib/components/Common/StatusBanner.svelte';

	// SvelteKit passes params to all route components
	export let params: Record<string, string> = {};

	let name = '';
	let kernelPath = './vmlinux';
	let kernelOpts = 'ro console=ttyS0 noapic reboot=k panic=1 pci=off nomodules';
	let rootDrivePath = '';
	let rootDriveReadOnly = false;
	let cpus = 1;
	let memory = 512;
	let cpuTemplate = '';
	let disableSMT = false;
	let logLevel = 'Debug';

	let loading = false;
	let submitError = '';
	let touched: Record<string, boolean> = {};

	const rules: Record<string, { check: (value: unknown) => string }> = {
		name: { check: (v) => errorOf(v, vmName()) },
		kernelPath: { check: (v) => errorOf(v, filePath()) },
		rootDrivePath: { check: (v) => errorOf(v, filePath()) },
		cpus: { check: (v) => errorOf(v, cpuCount()) },
		memory: { check: (v) => errorOf(v, memoryMB()) }
	};

	function errorOf(value: unknown, rule: Parameters<typeof validate>[1]): string {
		const result = validate(value, rule);
		return result.valid ? '' : (result.error ?? '');
	}

	// Errors only appear after the field is blurred (validate on blur)
	$: errors = Object.fromEntries(
		Object.entries(rules).map(([key, rule]) => [key, touched[key] ? rule.check(getValue(key)) : ''])
	);

	function getValue(key: string): unknown {
		switch (key) {
			case 'name':
				return name;
			case 'kernelPath':
				return kernelPath;
			case 'rootDrivePath':
				return rootDrivePath;
			case 'cpus':
				return cpus;
			case 'memory':
				return memory;
			default:
				return '';
		}
	}

	function markTouched(key: string) {
		touched = { ...touched, [key]: true };
	}

	async function handleSubmit() {
		touched = Object.fromEntries(Object.keys(rules).map((key) => [key, true]));
		if (Object.values(errors).some(Boolean)) return;

		loading = true;
		submitError = '';

		const config: VMConfig = {
			name,
			kernel_path: kernelPath,
			kernel_opts: kernelOpts,
			root_drive: {
				path: rootDrivePath,
				read_only: rootDriveReadOnly
			},
			cpus,
			memory_mb: memory,
			cpu_template: cpuTemplate || undefined,
			disable_smt: disableSMT,
			log_level: logLevel
		};

		const vm = await vms.create(name, config);
		loading = false;

		if (vm) {
			goto('/');
		} else {
			submitError = $vms.error ?? 'Failed to create the VM.';
		}
	}
</script>

<svelte:head>
	<title>New VM · Agni</title>
</svelte:head>

<div class="max-w-2xl space-y-6">
	<PageHeader
		title="New VM"
		description="Configure a Firecracker microVM before launch."
		crumbs={[{ label: 'Virtual Machines', href: '/vms' }, { label: 'New VM' }]}
	/>

	<form on:submit|preventDefault={handleSubmit} class="space-y-6" novalidate>
		<section class="card p-6">
			<h2 class="text-lg font-semibold text-foreground mb-4">Basics</h2>
			<div class="space-y-4">
				<div>
					<label class="label" for="name">Name</label>
					<input
						id="name"
						type="text"
						bind:value={name}
						on:blur={() => markTouched('name')}
						class="input {errors.name ? 'border-destructive focus:ring-destructive/30' : ''}"
						placeholder="web-server-1"
					/>
					{#if errors.name}
						<p class="text-[13px] text-destructive mt-1.5 flex items-center gap-1.5">
							<Icon name="ph-warning-circle" size={14} />
							{errors.name}
						</p>
					{:else}
						<p class="text-xs text-muted-foreground mt-1.5">Letters, numbers, underscores, and hyphens.</p>
					{/if}
				</div>
			</div>
		</section>

		<section class="card p-6">
			<h2 class="text-lg font-semibold text-foreground mb-4">Kernel</h2>
			<div class="space-y-4">
				<div>
					<label class="label" for="kernelPath">Kernel image</label>
					<input
						id="kernelPath"
						type="text"
						bind:value={kernelPath}
						on:blur={() => markTouched('kernelPath')}
						class="input font-mono text-[13px] {errors.kernelPath ? 'border-destructive focus:ring-destructive/30' : ''}"
						placeholder="/path/to/vmlinux"
					/>
					{#if errors.kernelPath}
						<p class="text-[13px] text-destructive mt-1.5 flex items-center gap-1.5">
							<Icon name="ph-warning-circle" size={14} />
							{errors.kernelPath}
						</p>
					{:else}
						<p class="text-xs text-muted-foreground mt-1.5">Path to an uncompressed kernel image on this machine.</p>
					{/if}
				</div>
				<div>
					<label class="label" for="kernelOpts">Boot arguments</label>
					<input id="kernelOpts" type="text" bind:value={kernelOpts} class="input font-mono text-[13px]" />
				</div>
			</div>
		</section>

		<section class="card p-6">
			<h2 class="text-lg font-semibold text-foreground mb-4">Storage</h2>
			<div class="space-y-4">
				<div>
					<label class="label" for="rootDrivePath">Root drive</label>
					<input
						id="rootDrivePath"
						type="text"
						bind:value={rootDrivePath}
						on:blur={() => markTouched('rootDrivePath')}
						class="input font-mono text-[13px] {errors.rootDrivePath ? 'border-destructive focus:ring-destructive/30' : ''}"
						placeholder="/path/to/rootfs.ext4"
					/>
					{#if errors.rootDrivePath}
						<p class="text-[13px] text-destructive mt-1.5 flex items-center gap-1.5">
							<Icon name="ph-warning-circle" size={14} />
							{errors.rootDrivePath}
						</p>
					{:else}
						<p class="text-xs text-muted-foreground mt-1.5">Path to the root filesystem image.</p>
					{/if}
				</div>
				<div class="flex items-center gap-2">
					<input id="rootDriveReadOnly" type="checkbox" bind:checked={rootDriveReadOnly} class="h-4 w-4 rounded-sm accent-primary" />
					<label for="rootDriveReadOnly" class="text-sm text-foreground">Mount root drive read-only</label>
				</div>
			</div>
		</section>

		<section class="card p-6">
			<h2 class="text-lg font-semibold text-foreground mb-4">Resources</h2>
			<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
				<div>
					<label class="label" for="cpus">vCPUs</label>
					<input
						id="cpus"
						type="number"
						min="1"
						max="32"
						bind:value={cpus}
						on:blur={() => markTouched('cpus')}
						class="input tabular-nums {errors.cpus ? 'border-destructive focus:ring-destructive/30' : ''}"
					/>
					{#if errors.cpus}
						<p class="text-[13px] text-destructive mt-1.5 flex items-center gap-1.5">
							<Icon name="ph-warning-circle" size={14} />
							{errors.cpus}
						</p>
					{/if}
				</div>
				<div>
					<label class="label" for="memory">Memory (MB)</label>
					<input
						id="memory"
						type="number"
						min="128"
						max="32768"
						step="128"
						bind:value={memory}
						on:blur={() => markTouched('memory')}
						class="input tabular-nums {errors.memory ? 'border-destructive focus:ring-destructive/30' : ''}"
					/>
					{#if errors.memory}
						<p class="text-[13px] text-destructive mt-1.5 flex items-center gap-1.5">
							<Icon name="ph-warning-circle" size={14} />
							{errors.memory}
						</p>
					{/if}
				</div>
				<div>
					<label class="label" for="cpuTemplate">CPU template</label>
					<select id="cpuTemplate" bind:value={cpuTemplate} class="input">
						<option value="">None</option>
						<option value="C3">C3</option>
						<option value="T2">T2</option>
					</select>
				</div>
				<div>
					<label class="label" for="logLevel">Log level</label>
					<select id="logLevel" bind:value={logLevel} class="input">
						<option value="Debug">Debug</option>
						<option value="Info">Info</option>
						<option value="Warn">Warn</option>
						<option value="Error">Error</option>
					</select>
				</div>
				<div class="flex items-center gap-2 sm:col-span-2">
					<input id="disableSMT" type="checkbox" bind:checked={disableSMT} class="h-4 w-4 rounded-sm accent-primary" />
					<label for="disableSMT" class="text-sm text-foreground">Disable CPU simultaneous multithreading</label>
				</div>
			</div>
		</section>

		{#if submitError}
			<StatusBanner variant="error" message={submitError} />
		{/if}

		<div class="flex justify-end gap-2">
			<Button variant="ghost" href="/vms">Cancel</Button>
			<Button variant="primary" type="submit" loading={loading}>Create VM</Button>
		</div>
	</form>
</div>
