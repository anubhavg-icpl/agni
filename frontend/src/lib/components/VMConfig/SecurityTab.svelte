<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { VMConfig, JailerConfig } from '$lib/api/client';

	export let config: Partial<VMConfig>;

	const dispatch = createEventDispatcher();

	let jailerEnabled = !!config.jailer;
	let jailerConfig: JailerConfig = config.jailer || {
		chroot_base_dir: '/srv/jailer',
		uid: 1000,
		gid: 1000,
		daemonize: true,
		exec_file: '/usr/bin/firecracker',
		numa_node: undefined
	};

	function emitChange() {
		dispatch('change', {
			jailer: jailerEnabled ? jailerConfig : undefined
		});
	}

	function toggleJailer() {
		jailerEnabled = !jailerEnabled;
		emitChange();
	}

	function updateJailerField(field: keyof JailerConfig, value: string | number | boolean) {
		jailerConfig = { ...jailerConfig, [field]: value };
		emitChange();
	}

	function handleChrootBaseInput(event: Event) {
		const target = event.target as HTMLInputElement;
		updateJailerField('chroot_base_dir', target.value);
	}

	function handleUIDInput(event: Event) {
		const target = event.target as HTMLInputElement;
		updateJailerField('uid', parseInt(target.value) || 1000);
	}

	function handleGIDInput(event: Event) {
		const target = event.target as HTMLInputElement;
		updateJailerField('gid', parseInt(target.value) || 1000);
	}

	function handleExecFileInput(event: Event) {
		const target = event.target as HTMLInputElement;
		updateJailerField('exec_file', target.value);
	}

	function handleDaemonizeChange(event: Event) {
		const target = event.target as HTMLInputElement;
		updateJailerField('daemonize', target.checked);
	}

	function handleNumaNodeInput(event: Event) {
		const target = event.target as HTMLInputElement;
		const val = target.value;
		updateJailerField('numa_node', val ? parseInt(val) : (undefined as unknown as number));
	}
</script>

<div class="space-y-6">
	<h3 class="text-lg font-medium">Security Configuration</h3>

	<!-- Jailer Toggle -->
	<div class="flex items-center gap-3 p-4 bg-muted rounded-lg">
		<input
			type="checkbox"
			id="jailerEnabled"
			checked={jailerEnabled}
			on:change={toggleJailer}
			class="rounded bg-secondary border-input text-primary-accent focus:ring-primary w-5 h-5"
		/>
		<div>
			<label for="jailerEnabled" class="font-medium">Enable Jailer</label>
			<p class="text-sm text-muted-foreground">
				Run Firecracker in a secure sandbox with restricted permissions
			</p>
		</div>
	</div>

	{#if jailerEnabled}
		<div class="space-y-4 border-l-2 border-primary pl-4">
			<!-- Chroot Base -->
			<div>
				<label for="chrootBase" class="label">Chroot Base Directory</label>
				<input
					type="text"
					id="chrootBase"
					value={jailerConfig.chroot_base_dir}
					on:input={handleChrootBaseInput}
					class="input w-full font-mono"
					placeholder="/srv/jailer"
				/>
				<p class="text-sm text-muted-foreground mt-1">
					Base directory for jailer chroot environments
				</p>
			</div>

			<!-- UID/GID -->
			<div class="grid grid-cols-2 gap-4">
				<div>
					<label for="jailerUid" class="label">User ID (UID)</label>
					<input
						type="number"
						id="jailerUid"
						value={jailerConfig.uid}
						on:input={handleUIDInput}
						class="input w-full"
						min="0"
					/>
				</div>
				<div>
					<label for="jailerGid" class="label">Group ID (GID)</label>
					<input
						type="number"
						id="jailerGid"
						value={jailerConfig.gid}
						on:input={handleGIDInput}
						class="input w-full"
						min="0"
					/>
				</div>
			</div>
			<p class="text-sm text-muted-foreground">
				The unprivileged user/group that Firecracker will run as inside the jail
			</p>

			<!-- Firecracker Executable -->
			<div>
				<label for="execFile" class="label">Firecracker Executable Path</label>
				<input
					type="text"
					id="execFile"
					value={jailerConfig.exec_file}
					on:input={handleExecFileInput}
					class="input w-full font-mono"
					placeholder="/usr/bin/firecracker"
				/>
			</div>

			<!-- Daemonize -->
			<div class="flex items-center gap-2">
				<input
					type="checkbox"
					id="jailerDaemonize"
					checked={jailerConfig.daemonize}
					on:change={handleDaemonizeChange}
					class="rounded bg-secondary border-input"
				/>
				<label for="jailerDaemonize" class="text-sm text-muted-foreground">Daemonize Process</label>
			</div>

			<!-- NUMA Node (optional) -->
			<div>
				<label for="numaNode" class="label">NUMA Node (Optional)</label>
				<input
					type="number"
					id="numaNode"
					value={jailerConfig.numa_node ?? ''}
					on:input={handleNumaNodeInput}
					class="input w-full"
					min="0"
					placeholder="Auto"
				/>
				<p class="text-sm text-muted-foreground mt-1">
					Pin the VM to a specific NUMA node for better memory locality
				</p>
			</div>
		</div>
	{/if}

	<!-- Security Info -->
	<div class="bg-primary/10 border border-primary/30 rounded-lg p-4">
		<h4 class="font-medium text-primary-accent mb-2">Security Best Practices</h4>
		<ul class="text-sm text-muted-foreground space-y-1">
			<li>- Enable jailer for production workloads</li>
			<li>- Use unprivileged user/group IDs (not root)</li>
			<li>- Consider disabling SMT for sensitive workloads</li>
			<li>- Use read-only root filesystems when possible</li>
			<li>- Limit CPU and memory resources appropriately</li>
		</ul>
	</div>

	<!-- Security Summary -->
	<div class="bg-muted rounded-lg p-4">
		<h4 class="text-sm font-medium text-muted-foreground mb-2">Security Summary</h4>
		<div class="space-y-1 text-sm">
			<div class="flex items-center gap-2">
				<span class="text-muted-foreground">Jailer:</span>
				<span class={jailerEnabled ? 'text-status-ok-text' : 'text-status-warn-text'}>
					{jailerEnabled ? 'Enabled' : 'Disabled'}
				</span>
			</div>
			{#if jailerEnabled}
				<div>
					<span class="text-muted-foreground">Chroot:</span>
					<span class="ml-2 font-mono text-xs">{jailerConfig.chroot_base_dir}</span>
				</div>
				<div>
					<span class="text-muted-foreground">User/Group:</span>
					<span class="ml-2">{jailerConfig.uid}:{jailerConfig.gid}</span>
				</div>
			{/if}
		</div>
	</div>
</div>
