<script lang="ts">
	import { onMount, onDestroy, tick } from 'svelte';
	import { api } from '$lib/api/client';
	import LogFilter from './LogFilter.svelte';
	import { formatTime } from '$lib/utils/formatters';

	export let vmId: string;
	export let autoScroll = true;

	interface LogEntry {
		timestamp: string;
		level: string;
		message: string;
		source?: string;
	}

	let logs: LogEntry[] = [];
	let filteredLogs: LogEntry[] = [];
	let ws: WebSocket | null = null;
	let connected = false;
	let error = '';
	let logContainer: HTMLElement;

	let levelFilter = 'all';
	let searchQuery = '';
	let paused = false;

	const levelColors: Record<string, string> = {
		debug: 'text-gray-400',
		info: 'text-blue-400',
		warn: 'text-yellow-400',
		warning: 'text-yellow-400',
		error: 'text-red-400',
		fatal: 'text-red-500'
	};

	function connect() {
		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		const host = window.location.host;
		const token = api.getToken();

		ws = new WebSocket(`${protocol}//${host}/api/vms/${vmId}/logs?token=${token}`);

		ws.onopen = () => {
			connected = true;
			error = '';
		};

		ws.onmessage = (event) => {
			if (paused) return;

			try {
				const entry: LogEntry = JSON.parse(event.data);
				logs = [...logs, entry].slice(-1000); // Keep last 1000 logs
				applyFilters();

				if (autoScroll) {
					tick().then(scrollToBottom);
				}
			} catch {
				// Handle plain text messages
				logs = [
					...logs,
					{
						timestamp: new Date().toISOString(),
						level: 'info',
						message: event.data
					}
				].slice(-1000);
				applyFilters();

				if (autoScroll) {
					tick().then(scrollToBottom);
				}
			}
		};

		ws.onerror = () => {
			error = 'WebSocket connection error';
			connected = false;
		};

		ws.onclose = () => {
			connected = false;
			// Attempt to reconnect after 3 seconds
			setTimeout(() => {
				if (!ws || ws.readyState === WebSocket.CLOSED) {
					connect();
				}
			}, 3000);
		};
	}

	function disconnect() {
		if (ws) {
			ws.close();
			ws = null;
		}
	}

	function scrollToBottom() {
		if (logContainer) {
			logContainer.scrollTop = logContainer.scrollHeight;
		}
	}

	function applyFilters() {
		filteredLogs = logs.filter((log) => {
			// Level filter
			if (levelFilter !== 'all' && log.level.toLowerCase() !== levelFilter) {
				return false;
			}

			// Search filter
			if (searchQuery && !log.message.toLowerCase().includes(searchQuery.toLowerCase())) {
				return false;
			}

			return true;
		});
	}

	function handleFilterChange(event: CustomEvent<{ level: string; search: string }>) {
		levelFilter = event.detail.level;
		searchQuery = event.detail.search;
		applyFilters();
	}

	function handlePauseToggle() {
		paused = !paused;
	}

	function handleClear() {
		logs = [];
		filteredLogs = [];
	}

	function handleExport() {
		const content = filteredLogs
			.map((log) => `[${formatTime(log.timestamp)}] [${log.level.toUpperCase()}] ${log.message}`)
			.join('\n');

		const blob = new Blob([content], { type: 'text/plain' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `vm-${vmId}-logs.txt`;
		a.click();
		URL.revokeObjectURL(url);
	}

	onMount(() => {
		connect();
	});

	onDestroy(() => {
		disconnect();
	});

	$: if (levelFilter || searchQuery !== undefined) {
		applyFilters();
	}
</script>

<div class="flex flex-col h-full">
	<!-- Header -->
	<div class="flex items-center justify-between mb-4">
		<div class="flex items-center gap-2">
			<span class="w-2 h-2 rounded-full {connected ? 'bg-green-500' : 'bg-red-500'}"></span>
			<span class="text-sm text-gray-400">
				{connected ? 'Connected' : 'Disconnected'}
			</span>
		</div>

		<div class="flex items-center gap-2">
			<button
				on:click={handlePauseToggle}
				class="btn btn-secondary text-sm px-3 py-1"
				title={paused ? 'Resume' : 'Pause'}
			>
				{#if paused}
					<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
						<path d="M8 5v14l11-7z" />
					</svg>
				{:else}
					<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
						<path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z" />
					</svg>
				{/if}
			</button>
			<button on:click={handleClear} class="btn btn-secondary text-sm px-3 py-1"> Clear </button>
			<button on:click={handleExport} class="btn btn-secondary text-sm px-3 py-1"> Export </button>
		</div>
	</div>

	<!-- Filters -->
	<LogFilter on:change={handleFilterChange} />

	<!-- Log content -->
	<div
		bind:this={logContainer}
		class="flex-1 bg-gray-900 rounded-lg p-4 font-mono text-sm overflow-y-auto"
	>
		{#if error}
			<div class="text-red-400">{error}</div>
		{:else if filteredLogs.length === 0}
			<div class="text-gray-500 text-center py-8">
				{logs.length === 0 ? 'Waiting for logs...' : 'No logs match the current filters'}
			</div>
		{:else}
			{#each filteredLogs as log}
				<div class="flex gap-2 py-0.5 hover:bg-gray-800/50">
					<span class="text-gray-500 flex-shrink-0">{formatTime(log.timestamp)}</span>
					<span
						class="flex-shrink-0 w-14 {levelColors[log.level.toLowerCase()] || 'text-gray-400'}"
					>
						[{log.level.toUpperCase()}]
					</span>
					<span class="text-gray-200 break-all">{log.message}</span>
				</div>
			{/each}
		{/if}
	</div>

	<!-- Status bar -->
	<div class="flex items-center justify-between mt-2 text-xs text-gray-500">
		<span>{filteredLogs.length} / {logs.length} entries</span>
		{#if paused}
			<span class="text-yellow-400">Paused</span>
		{/if}
		<label class="flex items-center gap-1">
			<input type="checkbox" bind:checked={autoScroll} class="rounded" />
			Auto-scroll
		</label>
	</div>
</div>
