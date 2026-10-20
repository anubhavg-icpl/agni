<script lang="ts">
	import { onMount, onDestroy, tick } from 'svelte';
	import { api } from '$lib/api/client';
	import LogFilter from './LogFilter.svelte';
	import Button from '../Common/Button.svelte';
	import Icon from '../Common/Icon.svelte';
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
	let destroyed = false;
	let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

	let levelFilter = 'all';
	let searchQuery = '';
	let paused = false;

	// Severity colours come from the status vocabulary
	const levelColors: Record<string, string> = {
		debug: 'text-muted-foreground',
		info: 'text-primary-accent',
		warn: 'text-status-warn-text',
		warning: 'text-status-warn-text',
		error: 'text-status-error-text',
		fatal: 'text-status-error-text'
	};

	function appendLog(entry: LogEntry) {
		logs = [...logs, entry].slice(-1000);
		applyFilters();
		if (autoScroll) {
			tick().then(scrollToBottom);
		}
	}

	function parseLogMessage(raw: string): LogEntry | null {
		try {
			const data = JSON.parse(raw);
			if (data && typeof data === 'object' && data.type === 'log' && data.payload) {
				return data.payload as LogEntry;
			}
			if (data && typeof data === 'object' && (data.message != null || data.level != null)) {
				return data as LogEntry;
			}
			return {
				timestamp: new Date().toISOString(),
				level: 'info',
				message: typeof data === 'string' ? data : raw
			};
		} catch {
			return {
				timestamp: new Date().toISOString(),
				level: 'info',
				message: raw
			};
		}
	}

	function connect() {
		if (destroyed) return;

		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		const host = window.location.host;
		const token = api.getToken();

		ws = new WebSocket(`${protocol}//${host}/api/vms/${vmId}/logs?token=${token}`);

		ws.onopen = () => {
			connected = true;
			error = '';
		};

		ws.onmessage = (event) => {
			if (paused || destroyed) return;
			const entry = parseLogMessage(event.data);
			if (entry) appendLog(entry);
		};

		ws.onerror = () => {
			error = 'WebSocket connection error.';
			connected = false;
		};

		ws.onclose = () => {
			connected = false;
			if (destroyed) return;
			reconnectTimer = setTimeout(() => {
				if (!destroyed && (!ws || ws.readyState === WebSocket.CLOSED)) {
					connect();
				}
			}, 3000);
		};
	}

	function disconnect() {
		destroyed = true;
		if (reconnectTimer) {
			clearTimeout(reconnectTimer);
			reconnectTimer = null;
		}
		if (ws) {
			ws.onclose = null;
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
			if (levelFilter !== 'all' && log.level.toLowerCase() !== levelFilter) {
				return false;
			}
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
	<div class="flex items-center justify-between mb-4">
		<div class="flex items-center gap-2">
			<span
				class="w-2 h-2 rounded-full {connected ? 'bg-status-ok' : 'bg-status-error'}"
				aria-hidden="true"
			></span>
			<span class="text-sm text-muted-foreground">{connected ? 'Connected' : 'Disconnected'}</span>
		</div>

		<div class="flex items-center gap-2">
			<Button
				variant="secondary"
				size="sm"
				on:click={handlePauseToggle}
				title={paused ? 'Resume' : 'Pause'}
			>
				<Icon name={paused ? 'ph-play' : 'ph-pause'} size={16} />
				{paused ? 'Resume' : 'Pause'}
			</Button>
			<Button variant="secondary" size="sm" on:click={handleClear}>Clear</Button>
			<Button variant="secondary" size="sm" on:click={handleExport}>Export</Button>
		</div>
	</div>

	<LogFilter on:change={handleFilterChange} />

	<div
		bind:this={logContainer}
		class="flex-1 bg-muted rounded-lg p-4 font-mono text-[13px] overflow-y-auto"
	>
		{#if error}
			<p class="text-status-error-text">{error}</p>
		{:else if filteredLogs.length === 0}
			<p class="text-muted-foreground text-center py-8">
				{logs.length === 0 ? 'Waiting for logs' : 'No logs match the current filters'}
			</p>
		{:else}
			{#each filteredLogs as log}
				<div class="flex gap-2 py-0.5 hover:bg-accent/50 rounded-sm">
					<span class="text-muted-foreground shrink-0">{formatTime(log.timestamp)}</span>
					<span
						class="shrink-0 w-16 {levelColors[log.level.toLowerCase()] || 'text-muted-foreground'}"
					>
						[{log.level.toUpperCase()}]
					</span>
					<span class="text-foreground break-all">{log.message}</span>
				</div>
			{/each}
		{/if}
	</div>

	<div class="flex items-center justify-between mt-2 text-xs text-muted-foreground">
		<span class="tabular-nums">{filteredLogs.length} / {logs.length} entries</span>
		{#if paused}
			<span class="text-status-warn-text">Paused</span>
		{/if}
		<label class="flex items-center gap-1.5">
			<input type="checkbox" bind:checked={autoScroll} class="rounded-sm accent-primary" />
			Auto-scroll
		</label>
	</div>
</div>
