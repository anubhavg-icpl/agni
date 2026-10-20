<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { VMConfig, Drive } from '$lib/api/client';
	import FileInput from '../Common/FileInput.svelte';

	export let config: Partial<VMConfig>;
	export let errors: Record<string, string> = {};

	const dispatch = createEventDispatcher();

	let rootDrivePath = config.root_drive?.path || '';
	let rootDriveReadOnly = config.root_drive?.read_only || false;
	let additionalDrives: Drive[] = config.additional_drives || [];

	function emitChange() {
		dispatch('change', {
			root_drive: {
				path: rootDrivePath,
				read_only: rootDriveReadOnly
			},
			additional_drives: additionalDrives.length > 0 ? additionalDrives : undefined
		});
	}

	function addDrive() {
		additionalDrives = [...additionalDrives, { path: '', read_only: false }];
	}

	function removeDrive(index: number) {
		additionalDrives = additionalDrives.filter((_, i) => i !== index);
		emitChange();
	}

	function updateDrive(index: number, field: keyof Drive, value: string | boolean) {
		additionalDrives[index] = { ...additionalDrives[index], [field]: value };
		additionalDrives = additionalDrives;
		emitChange();
	}

	function handleDrivePathInput(index: number, event: Event) {
		const target = event.target as HTMLInputElement;
		updateDrive(index, 'path', target.value);
	}

	function handleDriveReadOnlyChange(index: number, event: Event) {
		const target = event.target as HTMLInputElement;
		updateDrive(index, 'read_only', target.checked);
	}
</script>

<div class="space-y-6">
	<h3 class="text-lg font-medium">Storage Configuration</h3>

	<!-- Root Drive -->
	<div class="space-y-4">
		<h4 class="font-medium text-muted-foreground">Root Drive</h4>

		<FileInput
			label="Root Filesystem Path"
			bind:value={rootDrivePath}
			on:change={emitChange}
			placeholder="/path/to/rootfs.ext4"
			required
			error={errors.root_drive}
			hint="Path to the root filesystem image (ext4, squashfs, etc.)"
		/>

		<div class="flex items-center gap-2">
			<input
				type="checkbox"
				id="rootDriveReadOnly"
				bind:checked={rootDriveReadOnly}
				on:change={emitChange}
				class="rounded bg-secondary border-input"
			/>
			<label for="rootDriveReadOnly" class="text-sm text-muted-foreground">Mount as read-only</label>
		</div>
	</div>

	<!-- Additional Drives -->
	<div class="space-y-4">
		<div class="flex items-center justify-between">
			<h4 class="font-medium text-muted-foreground">Additional Drives</h4>
			<button type="button" on:click={addDrive} class="h-8 px-3 text-[13px] font-medium rounded-md bg-secondary text-secondary-foreground hover:bg-secondary/70 transition-colors">
				<svg class="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4v16m8-8H4"
					/>
				</svg>
				Add Drive
			</button>
		</div>

		{#if additionalDrives.length === 0}
			<div class="text-sm text-muted-foreground italic">No additional drives configured</div>
		{:else}
			{#each additionalDrives as drive, index}
				<div class="bg-muted rounded-lg p-4 space-y-3">
					<div class="flex items-center justify-between">
						<span class="text-sm font-medium text-muted-foreground">Drive {index + 1}</span>
						<button
							type="button"
							on:click={() => removeDrive(index)}
							class="text-destructive hover:text-destructive text-sm"
						>
							Remove
						</button>
					</div>

					<div>
						<label class="label text-sm" for="drive-{index}-path">Path</label>
						<input
							type="text"
							id="drive-{index}-path"
							value={drive.path}
							on:input={(e) => handleDrivePathInput(index, e)}
							class="input w-full text-sm"
							placeholder="/path/to/drive.img"
						/>
					</div>

					<div class="flex items-center gap-2">
						<input
							type="checkbox"
							id="drive-{index}-readonly"
							checked={drive.read_only}
							on:change={(e) => handleDriveReadOnlyChange(index, e)}
							class="rounded bg-secondary border-input"
						/>
						<label for="drive-{index}-readonly" class="text-sm text-muted-foreground">Read-only</label>
					</div>
				</div>
			{/each}
		{/if}
	</div>

	<!-- Storage Summary -->
	<div class="bg-muted rounded-lg p-4">
		<h4 class="text-sm font-medium text-muted-foreground mb-2">Storage Summary</h4>
		<div class="space-y-1 text-sm">
			<div class="flex items-center gap-2">
				<span class="text-muted-foreground">Root:</span>
				<span class="font-mono text-xs truncate flex-1" title={rootDrivePath}>
					{rootDrivePath || '(not set)'}
				</span>
				{#if rootDriveReadOnly}
					<span class="text-xs bg-status-warn/10 text-status-warn-text px-1 rounded">RO</span>
				{/if}
			</div>
			{#each additionalDrives as drive, index}
				<div class="flex items-center gap-2">
					<span class="text-muted-foreground">Drive {index + 1}:</span>
					<span class="font-mono text-xs truncate flex-1" title={drive.path}>
						{drive.path || '(not set)'}
					</span>
					{#if drive.read_only}
						<span class="text-xs bg-status-warn/10 text-status-warn-text px-1 rounded">RO</span>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</div>
