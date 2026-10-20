<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { api, type VMMetrics } from '$lib/api/client';
	import CPUChart from './CPUChart.svelte';
	import MemoryChart from './MemoryChart.svelte';
	import { formatBytes, formatBytesPerSecond, formatPercent } from '$lib/utils/formatters';
	import Skeleton from '../Common/Skeleton.svelte';
	import StatusBanner from '../Common/StatusBanner.svelte';

	export let vmId: string;

	let metrics: VMMetrics | null = null;
	let metricsHistory: VMMetrics[] = [];
	let loading = true;
	let error = '';
	let pollInterval: ReturnType<typeof setInterval>;

	const MAX_HISTORY = 60;

	async function fetchMetrics() {
		try {
			metrics = await api.getVMMetrics(vmId);
			metricsHistory = [...metricsHistory, metrics].slice(-MAX_HISTORY);
			error = '';
		} catch (e) {
			error = `Failed to load metrics: ${(e as Error).message}`;
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchMetrics();
		pollInterval = setInterval(fetchMetrics, 1000);
	});

	onDestroy(() => {
		if (pollInterval) {
			clearInterval(pollInterval);
		}
	});
</script>

<div class="space-y-6">
	{#if loading && !metrics}
		<div class="space-y-3" aria-hidden="true">
			{#each Array(3) as _}
				<Skeleton className="h-5 w-1/2" />
			{/each}
		</div>
	{:else if error && !metrics}
		<StatusBanner variant="error" message={error} />
	{:else if metrics}
		{#if error}
			<StatusBanner variant="warning" message="Showing the last received metrics. {error}" />
		{/if}

		<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
			<div class="rounded-lg border bg-card shadow-sm p-5">
				<div class="flex items-center gap-2 mb-3">
					<span class="h-2 w-2 rounded-full bg-primary" aria-hidden="true"></span>
					<span class="text-xs font-medium text-muted-foreground">CPU usage</span>
				</div>
				<div class="text-2xl font-semibold tabular-nums leading-8 text-foreground">
					{formatPercent(metrics.cpu_usage)}
				</div>
			</div>
			<div class="rounded-lg border bg-card shadow-sm p-5">
				<div class="flex items-center gap-2 mb-3">
					<span class="h-2 w-2 rounded-full bg-primary" aria-hidden="true"></span>
					<span class="text-xs font-medium text-muted-foreground">Memory used</span>
				</div>
				<div class="text-2xl font-semibold tabular-nums leading-8 text-foreground">
					{formatBytes(metrics.memory_used)}
				</div>
				<p class="text-xs text-muted-foreground mt-1">of {formatBytes(metrics.memory_total)}</p>
			</div>
			<div class="rounded-lg border bg-card shadow-sm p-5">
				<div class="flex items-center gap-2 mb-3">
					<span class="h-2 w-2 rounded-full bg-primary" aria-hidden="true"></span>
					<span class="text-xs font-medium text-muted-foreground">Disk I/O</span>
				</div>
				<div class="text-lg font-semibold tabular-nums text-foreground">
					{formatBytesPerSecond(metrics.disk_read_bytes)} / {formatBytesPerSecond(
						metrics.disk_write_bytes
					)}
				</div>
				<p class="text-xs text-muted-foreground mt-1">read / write</p>
			</div>
			<div class="rounded-lg border bg-card shadow-sm p-5">
				<div class="flex items-center gap-2 mb-3">
					<span class="h-2 w-2 rounded-full bg-primary" aria-hidden="true"></span>
					<span class="text-xs font-medium text-muted-foreground">Network I/O</span>
				</div>
				<div class="text-lg font-semibold tabular-nums text-foreground">
					{formatBytesPerSecond(metrics.net_rx_bytes)} / {formatBytesPerSecond(
						metrics.net_tx_bytes
					)}
				</div>
				<p class="text-xs text-muted-foreground mt-1">rx / tx</p>
			</div>
		</div>

		<div class="grid md:grid-cols-2 gap-6">
			<div class="card p-6">
				<h3 class="text-lg font-semibold text-foreground mb-4">CPU usage</h3>
				<CPUChart data={metricsHistory} />
			</div>
			<div class="card p-6">
				<h3 class="text-lg font-semibold text-foreground mb-4">Memory usage</h3>
				<MemoryChart data={metricsHistory} />
			</div>
		</div>

		<div class="card p-6">
			<h3 class="text-lg font-semibold text-foreground mb-4">Raw metrics</h3>
			<dl class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
				<div>
					<dt class="text-xs font-medium text-muted-foreground">CPU usage</dt>
					<dd class="data text-foreground mt-1">{formatPercent(metrics.cpu_usage)}</dd>
				</div>
				<div>
					<dt class="text-xs font-medium text-muted-foreground">Memory used</dt>
					<dd class="data text-foreground mt-1">{formatBytes(metrics.memory_used)}</dd>
				</div>
				<div>
					<dt class="text-xs font-medium text-muted-foreground">Memory total</dt>
					<dd class="data text-foreground mt-1">{formatBytes(metrics.memory_total)}</dd>
				</div>
				<div>
					<dt class="text-xs font-medium text-muted-foreground">Disk read</dt>
					<dd class="data text-foreground mt-1">{formatBytes(metrics.disk_read_bytes)}</dd>
				</div>
				<div>
					<dt class="text-xs font-medium text-muted-foreground">Disk write</dt>
					<dd class="data text-foreground mt-1">{formatBytes(metrics.disk_write_bytes)}</dd>
				</div>
				<div>
					<dt class="text-xs font-medium text-muted-foreground">Net rx</dt>
					<dd class="data text-foreground mt-1">{formatBytes(metrics.net_rx_bytes)}</dd>
				</div>
				<div>
					<dt class="text-xs font-medium text-muted-foreground">Net tx</dt>
					<dd class="data text-foreground mt-1">{formatBytes(metrics.net_tx_bytes)}</dd>
				</div>
				<div>
					<dt class="text-xs font-medium text-muted-foreground">Timestamp</dt>
					<dd class="data text-foreground mt-1">
						{new Date(metrics.timestamp).toLocaleTimeString()}
					</dd>
				</div>
			</dl>
		</div>
	{/if}
</div>
