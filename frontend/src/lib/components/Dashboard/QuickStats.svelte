<script lang="ts">
	import type { VM } from '$lib/api/client';

	export let vms: VM[] = [];

	$: totalVMs = vms.length;
	$: runningVMs = vms.filter((vm) => vm.status === 'running').length;
	$: stoppedVMs = vms.filter((vm) => vm.status === 'stopped').length;
	$: errorVMs = vms.filter((vm) => vm.status === 'error').length;

	const stats = [
		{ label: 'Total Minions', icon: '👹', getValue: () => totalVMs, color: 'text-orange-400' },
		{ label: 'Alive & Kicking', icon: '🔥', getValue: () => runningVMs, color: 'text-green-400' },
		{ label: 'Taking a Nap', icon: '💀', getValue: () => stoppedVMs, color: 'text-gray-400' },
		{ label: 'Having Issues', icon: '💥', getValue: () => errorVMs, color: 'text-red-400' }
	];
</script>

<div class="grid grid-cols-2 lg:grid-cols-4 gap-2 sm:gap-4">
	{#each stats as stat, i}
		<div class="card hover:border-orange-500/30 transition-all duration-200 group p-3 sm:p-4">
			<div class="flex items-center gap-1.5 sm:gap-2 mb-1 sm:mb-2">
				<span class="text-sm sm:text-lg group-hover:scale-110 transition-transform">{stat.icon}</span>
				<span class="text-xs sm:text-sm text-gray-400 truncate">{stat.label}</span>
			</div>
			<div class="text-xl sm:text-3xl font-bold {stat.color}">
				{stat.getValue()}
			</div>
		</div>
	{/each}
</div>
