<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { VMConfig, NIC, Vsock } from '$lib/api/client';

	export let config: Partial<VMConfig>;

	const dispatch = createEventDispatcher();

	let networkInterfaces: NIC[] = config.network_interfaces || [];
	let vsockDevices: Vsock[] = config.vsock_devices || [];

	function emitChange() {
		dispatch('change', {
			network_interfaces: networkInterfaces.length > 0 ? networkInterfaces : undefined,
			vsock_devices: vsockDevices.length > 0 ? vsockDevices : undefined
		});
	}

	// Network Interfaces
	function addNIC() {
		networkInterfaces = [
			...networkInterfaces,
			{
				device: '',
				mac_address: generateMAC()
			}
		];
	}

	function removeNIC(index: number) {
		networkInterfaces = networkInterfaces.filter((_, i) => i !== index);
		emitChange();
	}

	function updateNIC(index: number, field: keyof NIC, value: string) {
		networkInterfaces[index] = { ...networkInterfaces[index], [field]: value };
		networkInterfaces = networkInterfaces;
		emitChange();
	}

	function handleNICInput(index: number, field: keyof NIC, event: Event) {
		const target = event.target as HTMLInputElement;
		updateNIC(index, field, target.value);
	}

	// Vsock Devices
	function addVsock() {
		vsockDevices = [
			...vsockDevices,
			{
				cid: 3,
				path: ''
			}
		];
	}

	function removeVsock(index: number) {
		vsockDevices = vsockDevices.filter((_, i) => i !== index);
		emitChange();
	}

	function updateVsock(index: number, field: keyof Vsock, value: string | number) {
		vsockDevices[index] = { ...vsockDevices[index], [field]: value };
		vsockDevices = vsockDevices;
		emitChange();
	}

	function handleVsockCIDInput(index: number, event: Event) {
		const target = event.target as HTMLInputElement;
		updateVsock(index, 'cid', parseInt(target.value) || 3);
	}

	function handleVsockPathInput(index: number, event: Event) {
		const target = event.target as HTMLInputElement;
		updateVsock(index, 'path', target.value);
	}

	// Generate random MAC address
	function generateMAC(): string {
		const hex = () =>
			Math.floor(Math.random() * 256)
				.toString(16)
				.padStart(2, '0');
		// Use locally administered, unicast MAC address
		return `02:${hex()}:${hex()}:${hex()}:${hex()}:${hex()}`;
	}
</script>

<div class="space-y-6">
	<h3 class="text-lg font-medium">Network Configuration</h3>

	<!-- Network Interfaces -->
	<div class="space-y-4">
		<div class="flex items-center justify-between">
			<h4 class="font-medium text-gray-300">Network Interfaces (TAP)</h4>
			<button type="button" on:click={addNIC} class="btn btn-secondary text-sm px-3 py-1">
				<svg class="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4v16m8-8H4"
					/>
				</svg>
				Add NIC
			</button>
		</div>

		{#if networkInterfaces.length === 0}
			<div class="text-sm text-gray-500 italic">No network interfaces configured</div>
		{:else}
			{#each networkInterfaces as nic, index}
				<div class="bg-gray-900 rounded-lg p-4 space-y-3">
					<div class="flex items-center justify-between">
						<span class="text-sm font-medium text-gray-400">Interface {index + 1}</span>
						<button
							type="button"
							on:click={() => removeNIC(index)}
							class="text-red-400 hover:text-red-300 text-sm"
						>
							Remove
						</button>
					</div>

					<div class="grid md:grid-cols-2 gap-4">
						<div>
							<label class="label text-sm" for="nic-{index}-device">TAP Device</label>
							<input
								type="text"
								id="nic-{index}-device"
								value={nic.device}
								on:input={(e) => handleNICInput(index, 'device', e)}
								class="input w-full text-sm"
								placeholder="tap0"
							/>
						</div>
						<div>
							<label class="label text-sm" for="nic-{index}-mac">Guest MAC Address</label>
							<div class="flex gap-2">
								<input
									type="text"
									id="nic-{index}-mac"
									value={nic.mac_address}
									on:input={(e) => handleNICInput(index, 'mac_address', e)}
									class="input flex-1 text-sm font-mono"
									placeholder="02:XX:XX:XX:XX:XX"
								/>
								<button
									type="button"
									on:click={() => updateNIC(index, 'mac_address', generateMAC())}
									class="btn btn-secondary text-xs px-2"
									title="Generate random MAC"
								>
									<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
										/>
									</svg>
								</button>
							</div>
						</div>
					</div>
				</div>
			{/each}
		{/if}
	</div>

	<!-- Vsock Devices -->
	<div class="space-y-4">
		<div class="flex items-center justify-between">
			<h4 class="font-medium text-gray-300">Vsock Devices</h4>
			<button type="button" on:click={addVsock} class="btn btn-secondary text-sm px-3 py-1">
				<svg class="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4v16m8-8H4"
					/>
				</svg>
				Add Vsock
			</button>
		</div>

		{#if vsockDevices.length === 0}
			<div class="text-sm text-gray-500 italic">No vsock devices configured</div>
		{:else}
			{#each vsockDevices as vsock, index}
				<div class="bg-gray-900 rounded-lg p-4 space-y-3">
					<div class="flex items-center justify-between">
						<span class="text-sm font-medium text-gray-400">Vsock {index + 1}</span>
						<button
							type="button"
							on:click={() => removeVsock(index)}
							class="text-red-400 hover:text-red-300 text-sm"
						>
							Remove
						</button>
					</div>

					<div class="grid md:grid-cols-2 gap-4">
						<div>
							<label class="label text-sm" for="vsock-{index}-cid">Guest CID</label>
							<input
								type="number"
								id="vsock-{index}-cid"
								value={vsock.cid}
								on:input={(e) => handleVsockCIDInput(index, e)}
								class="input w-full text-sm"
								min="3"
							/>
							<p class="text-xs text-gray-500 mt-1">Context ID (must be >= 3)</p>
						</div>
						<div>
							<label class="label text-sm" for="vsock-{index}-path">Unix Socket Path</label>
							<input
								type="text"
								id="vsock-{index}-path"
								value={vsock.path}
								on:input={(e) => handleVsockPathInput(index, e)}
								class="input w-full text-sm font-mono"
								placeholder="/tmp/vsock.sock"
							/>
						</div>
					</div>
				</div>
			{/each}
		{/if}
	</div>

	<!-- Network Summary -->
	<div class="bg-gray-900 rounded-lg p-4">
		<h4 class="text-sm font-medium text-gray-400 mb-2">Network Summary</h4>
		<div class="space-y-1 text-sm">
			<div>
				<span class="text-gray-500">NICs:</span>
				<span class="ml-2">{networkInterfaces.length}</span>
			</div>
			<div>
				<span class="text-gray-500">Vsock Devices:</span>
				<span class="ml-2">{vsockDevices.length}</span>
			</div>
		</div>
	</div>
</div>
