<script lang="ts">
	import type { VM } from '$lib/api/client';

	export let vms: VM[] = [];

	$: total = vms.length;
	$: running = vms.filter((vm) => vm.status === 'running').length;
	$: stopped = vms.filter((vm) => vm.status === 'stopped').length;
	$: errors = vms.filter((vm) => vm.status === 'error').length;
</script>

<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
	<div class="rounded-lg border bg-card shadow-sm p-5">
		<div class="flex items-center gap-2 mb-3">
			<span class="h-2 w-2 rounded-full bg-primary" aria-hidden="true"></span>
			<span class="text-xs font-medium text-muted-foreground">Total</span>
		</div>
		<div class="text-2xl font-semibold tabular-nums leading-8 text-foreground">{total}</div>
		<p class="text-xs text-muted-foreground mt-1">Virtual machines</p>
	</div>

	<div class="rounded-lg border bg-card shadow-sm p-5">
		<div class="flex items-center gap-2 mb-3">
			<span class="h-2 w-2 rounded-full bg-status-ok" aria-hidden="true"></span>
			<span class="text-xs font-medium text-muted-foreground">Running</span>
		</div>
		<div class="text-2xl font-semibold tabular-nums leading-8 text-foreground">{running}</div>
		<p class="text-xs text-muted-foreground mt-1">Active microVMs</p>
	</div>

	<div class="rounded-lg border bg-card shadow-sm p-5">
		<div class="flex items-center gap-2 mb-3">
			<span class="h-2 w-2 rounded-full bg-status-idle" aria-hidden="true"></span>
			<span class="text-xs font-medium text-muted-foreground">Stopped</span>
		</div>
		<div class="text-2xl font-semibold tabular-nums leading-8 text-foreground">{stopped}</div>
		<p class="text-xs text-muted-foreground mt-1">Available to start</p>
	</div>

	<div class="rounded-lg border bg-card shadow-sm p-5">
		<div class="flex items-center gap-2 mb-3">
			<span
				class="h-2 w-2 rounded-full {errors > 0 ? 'bg-status-error pulse-dot' : 'bg-status-idle'}"
				aria-hidden="true"
			></span>
			<span class="text-xs font-medium text-muted-foreground">Errors</span>
		</div>
		<div class="text-2xl font-semibold tabular-nums leading-8 text-foreground">{errors}</div>
		<p class="text-xs text-muted-foreground mt-1">{errors > 0 ? 'Needs attention' : 'No failures'}</p>
	</div>
</div>
