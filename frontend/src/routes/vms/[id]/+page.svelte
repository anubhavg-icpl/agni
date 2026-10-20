<script lang="ts">
	import { onMount } from 'svelte';
	import type { Terminal } from '@xterm/xterm';
	import type { FitAddon } from '@xterm/addon-fit';
	import '@xterm/xterm/css/xterm.css';
	import { api, type VM } from '$lib/api/client';
	import { vms as vmStore } from '$lib/stores/vms';
	import StatusBadge from '$lib/components/Common/StatusBadge.svelte';
	import Button from '$lib/components/Common/Button.svelte';
	import Icon from '$lib/components/Common/Icon.svelte';
	import { formatDate, getFilename } from '$lib/utils/formatters';
	import { goto } from '$app/navigation';
	import Modal from '$lib/components/Common/Modal.svelte';

	let destroyOpen = false;
	let destroying = false;

	async function destroyVM() {
		if (!vm) return;
		destroying = true;
		try {
			await api.deleteVM(vm.id);
			goto('/vms');
		} catch (e) {
			actionError = `Failed to destroy VM: ${(e as Error).message}`;
			destroying = false;
		}
	}

	export let params: { id: string };

	let vm: VM | null = null;
	let loadError = '';
	let actionError = '';
	let busy = false;

	let terminalEl: HTMLElement;
	let term: Terminal | undefined;
	let fit: FitAddon | undefined;
	let ws: WebSocket | undefined;
	let connected = false;
	let consoleNotice = 'Connect the console to interact with this virtual machine.';
	let retryCount = 0;
	let statusTimer: ReturnType<typeof setInterval> | undefined;
	let resizeHandler: (() => void) | undefined;

	function wsUrl(): string {
		const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		const token = api.getToken() ?? '';
		return `${proto}//${window.location.host}/api/vms/${params.id}/console?token=${encodeURIComponent(token)}`;
	}

	function connectConsole() {
		disconnectConsole();
		if (!term) return;
		if (vm && vm.status !== 'running') {
			consoleNotice = 'The virtual machine is not running. Start it to attach the serial console.';
			return;
		}
		consoleNotice = 'Connecting to serial console…';
		const sock = new WebSocket(wsUrl());
		ws = sock;
		sock.onopen = () => {
			connected = true;
			retryCount = 0;
			consoleNotice = '';
			term?.focus();
		};
		sock.onmessage = (ev) => {
			term?.write(String(ev.data));
		};
		sock.onclose = () => {
			connected = false;
			if (ws !== sock) return;
			if (vm?.status === 'running') {
				// Exponential backoff so a failing connection can never
				// rapid-loop; resets after a successful connect.
				retryCount = Math.min(retryCount + 1, 6);
				const delay = Math.min(1500 * 2 ** (retryCount - 1), 20000);
				consoleNotice = `Console disconnected — reconnecting in ${Math.round(delay / 1000)}s…`;
				setTimeout(() => {
					if (ws === sock && vm?.status === 'running') connectConsole();
				}, delay);
			} else {
				consoleNotice =
					'The virtual machine is not running. Start it to attach the serial console.';
			}
		};
	}

	function disconnectConsole() {
		if (ws) {
			const s = ws;
			ws = undefined;
			s.onclose = null;
			s.close();
		}
		connected = false;
	}

	async function refresh() {
		try {
			vm = await api.getVM(params.id);
			loadError = '';
		} catch (e) {
			loadError = (e as Error).message;
		}
	}

	async function runAction(fn: (id: string) => Promise<unknown>) {
		busy = true;
		actionError = '';
		try {
			await fn(params.id);
			await refresh();
			if (vm?.status === 'running') connectConsole();
		} catch (e) {
			actionError = (e as Error).message;
		} finally {
			busy = false;
		}
	}

	onMount(() => {
		void refresh();

		// xterm ships as CommonJS: load it at mount time so Vite's SSR import
		// analysis never sees a named CJS export.
		void (async () => {
			const [{ Terminal: TerminalCtor }, { FitAddon: FitAddonCtor }] = await Promise.all([
				import('@xterm/xterm'),
				import('@xterm/addon-fit')
			]);
			if (!terminalEl) return;

			const t = new TerminalCtor({
				fontFamily: "'JetBrains Mono', ui-monospace, SFMono-Regular, Menlo, monospace",
				fontSize: 12.5,
				lineHeight: 1.2,
				cursorBlink: true,
				convertEol: true,
				scrollback: 5000,
				theme: {
					background: '#0a0a0a',
					foreground: '#ededed',
					cursor: '#ededed',
					selectionBackground: '#525252'
				}
			});
			const f = new FitAddonCtor();
			t.loadAddon(f);
			t.open(terminalEl);
			f.fit();
			t.onData((d) => {
				if (ws && ws.readyState === WebSocket.OPEN) ws.send(d);
			});
			term = t;
			fit = f;

			resizeHandler = () => fit?.fit();
			window.addEventListener('resize', resizeHandler);

			statusTimer = setInterval(() => {
				void (async () => {
					const was = vm?.status;
					await refresh();
					if (was && was !== 'running' && vm?.status === 'running') connectConsole();
				})();
			}, 2000);

			if (vm?.status === 'running') connectConsole();
		})();

		return () => {
			clearInterval(statusTimer);
			if (resizeHandler) window.removeEventListener('resize', resizeHandler);
			disconnectConsole();
			term?.dispose();
		};
	});
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex flex-wrap items-start justify-between gap-4">
		<div class="min-w-0">
			<a
				href="/vms"
				class="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
			>
				<Icon name="ph-arrow-left" size={16} />
				Virtual machines
			</a>
			<div class="mt-2 flex items-center gap-3 flex-wrap">
				<h1 class="text-2xl font-semibold tracking-tight">{vm?.name ?? 'Virtual machine'}</h1>
				{#if vm}<StatusBadge status={vm.status} />{/if}
			</div>
			{#if vm}
				<p class="mt-1 text-sm text-muted-foreground font-mono">{vm.id}</p>
			{/if}
		</div>
		<div class="flex items-center gap-2">
			{#if vm?.status === 'stopped' || vm?.status === 'error'}
				<Button
					variant="primary"
					loading={busy}
					on:click={() => runAction((id) => vmStore.start(id))}
				>
					<span class="inline-flex items-center gap-1.5"
						><Icon name="ph-play" size={16} /> Start</span
					>
				</Button>
			{:else if vm?.status === 'running'}
				<Button
					variant="secondary"
					loading={busy}
					on:click={() => runAction((id) => vmStore.shutdown(id))}
				>
					<span class="inline-flex items-center gap-1.5"
						><Icon name="ph-power" size={16} /> Shut down</span
					>
				</Button>
				<Button
					variant="destructive"
					loading={busy}
					on:click={() => runAction((id) => vmStore.stop(id))}
				>
					<span class="inline-flex items-center gap-1.5"
						><Icon name="ph-stop" size={16} /> Force stop</span
					>
				</Button>
			{/if}
		</div>
	</div>

	{#if loadError}
		<div
			class="rounded-lg border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
		>
			{loadError}
		</div>
	{/if}
	{#if actionError}
		<div
			class="rounded-lg border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
		>
			{actionError}
		</div>
	{/if}

	<div class="grid gap-6 lg:grid-cols-[1fr_320px] items-start">
		<!-- Console (always mounted so xterm has a parent element) -->
		<div class="card overflow-hidden">
			<div class="flex items-center justify-between gap-3 border-b px-4 py-3">
				<div class="flex items-center gap-2 text-sm font-medium">
					<Icon name="ph-terminal-window" size={16} />
					Serial console
				</div>
				<div class="flex items-center gap-2">
					<span class="text-xs text-muted-foreground inline-flex items-center gap-1.5">
						{#if connected}
							<span class="relative flex h-2 w-2">
								<span
									class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"
								></span>
								<span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
							</span>
							Connected
						{:else}
							<span class="h-2 w-2 rounded-full bg-muted-foreground/50"></span>
							Disconnected
						{/if}
					</span>
					<Button variant="secondary" size="sm" on:click={connectConsole}>Reconnect</Button>
					<Button variant="destructive" size="sm" on:click={() => (destroyOpen = true)}>
						<Icon name="ph-trash" size={16} />
						Destroy
					</Button>
				</div>
			</div>
			<div class="relative">
				<div class="p-3 bg-[#0a0a0a]">
					<div id="terminal" class="terminal-el" bind:this={terminalEl}></div>
				</div>
				{#if consoleNotice}
					<div
						class="absolute inset-0 flex items-center justify-center bg-[#0a0a0a]/85 text-sm text-muted-foreground px-6 text-center"
					>
						{consoleNotice}
					</div>
				{/if}
			</div>
		</div>

		<!-- Info -->
		{#if vm}
			<div class="card p-0">
				<div class="border-b px-4 py-3 text-sm font-medium">Configuration</div>
				<dl class="divide-y text-sm">
					<div class="px-4 py-3">
						<dt class="text-muted-foreground text-xs">vCPUs</dt>
						<dd class="mt-0.5 font-medium tabular-nums">{vm.config.cpus}</dd>
					</div>
					<div class="px-4 py-3">
						<dt class="text-muted-foreground text-xs">Memory</dt>
						<dd class="mt-0.5 font-medium tabular-nums">{vm.config.memory_mb} MB</dd>
					</div>
					<div class="px-4 py-3">
						<dt class="text-muted-foreground text-xs">Kernel</dt>
						<dd class="mt-0.5 font-mono text-[13px] break-all" title={vm.config.kernel_path}>
							{getFilename(vm.config.kernel_path) || 'N/A'}
						</dd>
					</div>
					{#if vm.config.root_drive?.path}
						<div class="px-4 py-3">
							<dt class="text-muted-foreground text-xs">Root drive</dt>
							<dd class="mt-0.5 font-mono text-[13px] break-all" title={vm.config.root_drive.path}>
								{getFilename(vm.config.root_drive.path)}
							</dd>
						</div>
					{/if}
					{#if vm.config.initrd_path}
						<div class="px-4 py-3">
							<dt class="text-muted-foreground text-xs">Initrd</dt>
							<dd class="mt-0.5 font-mono text-[13px] break-all" title={vm.config.initrd_path}>
								{getFilename(vm.config.initrd_path)}
							</dd>
						</div>
					{/if}
					{#if vm.config.kernel_opts}
						<div class="px-4 py-3">
							<dt class="text-muted-foreground text-xs">Kernel arguments</dt>
							<dd class="mt-0.5 font-mono text-[13px] break-all">{vm.config.kernel_opts}</dd>
						</div>
					{/if}
					<div class="px-4 py-3">
						<dt class="text-muted-foreground text-xs">Created</dt>
						<dd class="mt-0.5 font-medium">{formatDate(vm.created_at)}</dd>
					</div>
				</dl>
			</div>
		{:else}
			<div class="card p-0">
				<div class="border-b px-4 py-3 text-sm font-medium">Configuration</div>
				<div class="px-4 py-3 text-sm text-muted-foreground">Loading configuration…</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.terminal-el :global(.xterm) {
		padding: 4px;
	}
</style>

{#if destroyOpen}
	<Modal open={destroyOpen} title="Destroy virtual machine" on:close={() => (destroyOpen = false)}>
		<p class="text-sm text-muted-foreground">
			This force-stops and permanently deletes <span class="font-medium text-foreground">{vm?.name}</span>.
			The template rootfs on disk is kept; only this VM is removed.
		</p>
		<div class="mt-5 flex justify-end gap-2">
			<Button variant="secondary" size="sm" on:click={() => (destroyOpen = false)}>Cancel</Button>
			<Button variant="destructive" size="sm" disabled={destroying} on:click={destroyVM}>
				{destroying ? 'Destroying…' : 'Destroy VM'}
			</Button>
		</div>
	</Modal>
{/if}
