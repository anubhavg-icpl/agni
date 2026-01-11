<script lang="ts">
	import type { VM } from '$lib/api/client';

	export let vms: VM[] = [];

	$: totalVMs = vms.length;
	$: runningVMs = vms.filter((vm) => vm.status === 'running').length;
	$: stoppedVMs = vms.filter((vm) => vm.status === 'stopped').length;
	$: errorVMs = vms.filter((vm) => vm.status === 'error').length;
</script>

<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
	<div class="card">
		<div class="text-sm text-gray-400 mb-1">Total VMs</div>
		<div class="text-2xl font-bold">{totalVMs}</div>
	</div>
	<div class="card">
		<div class="text-sm text-gray-400 mb-1">Running</div>
		<div class="text-2xl font-bold text-green-400">{runningVMs}</div>
	</div>
	<div class="card">
		<div class="text-sm text-gray-400 mb-1">Stopped</div>
		<div class="text-2xl font-bold text-gray-400">{stoppedVMs}</div>
	</div>
	<div class="card">
		<div class="text-sm text-gray-400 mb-1">Errors</div>
		<div class="text-2xl font-bold text-red-400">{errorVMs}</div>
	</div>
</div>
