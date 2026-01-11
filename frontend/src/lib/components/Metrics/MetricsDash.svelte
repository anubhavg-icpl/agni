<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { api, type VMMetrics } from '$lib/api/client';
	import CPUChart from './CPUChart.svelte';
	import MemoryChart from './MemoryChart.svelte';
	import { formatBytes, formatBytesPerSecond, formatPercent } from '$lib/utils/formatters';

	export let vmId: string;

	let metrics: VMMetrics | null = null;
	let metricsHistory: VMMetrics[] = [];
	let loading = true;
	let error = '';
	let pollInterval: ReturnType<typeof setInterval>;

	const MAX_HISTORY = 60; // Keep 60 data points (1 minute at 1s interval)

	async function fetchMetrics() {
		try {
			metrics = await api.getVMMetrics(vmId);
			metricsHistory = [...metricsHistory, metrics].slice(-MAX_HISTORY);
			error = '';
		} catch (e) {
			error = (e as Error).message;
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
		<div class="text-center py-8 text-gray-400">Loading metrics...</div>
	{:else if error}
		<div class="bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg p-4">
			{error}
		</div>
	{:else if metrics}
		<!-- Summary cards -->
		<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
			<div class="card">
				<div class="text-sm text-gray-400 mb-1">CPU Usage</div>
				<div class="text-2xl font-bold text-blue-400">
					{formatPercent(metrics.cpu_usage)}
				</div>
			</div>
			<div class="card">
				<div class="text-sm text-gray-400 mb-1">Memory Used</div>
				<div class="text-2xl font-bold text-green-400">
					{formatBytes(metrics.memory_used)}
				</div>
				<div class="text-xs text-gray-500">
					of {formatBytes(metrics.memory_total)}
				</div>
			</div>
			<div class="card">
				<div class="text-sm text-gray-400 mb-1">Disk I/O</div>
				<div class="text-lg font-bold">
					<span class="text-green-400">{formatBytesPerSecond(metrics.disk_read_bytes)}</span>
					<span class="text-gray-500 mx-1">/</span>
					<span class="text-orange-400">{formatBytesPerSecond(metrics.disk_write_bytes)}</span>
				</div>
				<div class="text-xs text-gray-500">read / write</div>
			</div>
			<div class="card">
				<div class="text-sm text-gray-400 mb-1">Network I/O</div>
				<div class="text-lg font-bold">
					<span class="text-green-400">{formatBytesPerSecond(metrics.net_rx_bytes)}</span>
					<span class="text-gray-500 mx-1">/</span>
					<span class="text-orange-400">{formatBytesPerSecond(metrics.net_tx_bytes)}</span>
				</div>
				<div class="text-xs text-gray-500">rx / tx</div>
			</div>
		</div>

		<!-- Charts -->
		<div class="grid md:grid-cols-2 gap-6">
			<div class="card">
				<h3 class="text-lg font-medium mb-4">CPU Usage</h3>
				<CPUChart data={metricsHistory} />
			</div>
			<div class="card">
				<h3 class="text-lg font-medium mb-4">Memory Usage</h3>
				<MemoryChart data={metricsHistory} />
			</div>
		</div>

		<!-- Raw metrics -->
		<div class="card">
			<h3 class="text-lg font-medium mb-4">Raw Metrics</h3>
			<div class="grid grid-cols-2 md:grid-cols-3 gap-4 text-sm">
				<div>
					<span class="text-gray-400">CPU Usage:</span>
					<span class="ml-2">{formatPercent(metrics.cpu_usage)}</span>
				</div>
				<div>
					<span class="text-gray-400">Memory Used:</span>
					<span class="ml-2">{formatBytes(metrics.memory_used)}</span>
				</div>
				<div>
					<span class="text-gray-400">Memory Total:</span>
					<span class="ml-2">{formatBytes(metrics.memory_total)}</span>
				</div>
				<div>
					<span class="text-gray-400">Disk Read:</span>
					<span class="ml-2">{formatBytes(metrics.disk_read_bytes)}</span>
				</div>
				<div>
					<span class="text-gray-400">Disk Write:</span>
					<span class="ml-2">{formatBytes(metrics.disk_write_bytes)}</span>
				</div>
				<div>
					<span class="text-gray-400">Net RX:</span>
					<span class="ml-2">{formatBytes(metrics.net_rx_bytes)}</span>
				</div>
				<div>
					<span class="text-gray-400">Net TX:</span>
					<span class="ml-2">{formatBytes(metrics.net_tx_bytes)}</span>
				</div>
				<div class="col-span-2">
					<span class="text-gray-400">Timestamp:</span>
					<span class="ml-2">{new Date(metrics.timestamp).toLocaleTimeString()}</span>
				</div>
			</div>
		</div>
	{/if}
</div>
