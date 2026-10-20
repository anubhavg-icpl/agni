<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import PageHeader from '$lib/components/Common/PageHeader.svelte';
	import Button from '$lib/components/Common/Button.svelte';
	import Icon from '$lib/components/Common/Icon.svelte';
	import Modal from '$lib/components/Common/Modal.svelte';
	import StatusBanner from '$lib/components/Common/StatusBanner.svelte';
	import { toasts } from '$lib/components/Common/Toast.svelte';
	import { templates } from '$lib/stores/templates';
	import { formatBytes } from '$lib/utils/formatters';
	import type { DeployTemplate, ImageJob } from '$lib/api/client';

	let modalOpen = false;
	let deploying: DeployTemplate | null = null;
	let vmName = '';
	let cpus = 1;
	let memoryMb = 256;
	let autoStart = true;
	let busy = false;
	let installing = false;
	let activeJobs: Record<string, ImageJob> = {};

	$: osTemplates = $templates.templates.filter((t) => t.category === 'os');
	$: runtimeInstalled = $templates.runtime.installed;
	$: runtimeTemplate = $templates.templates.find((t) => t.id === 'firecracker-runtime');
	$: runtimeJob =
		activeJobs['firecracker-runtime'] ??
		(runtimeTemplate?.latest_job?.state === 'running' ? runtimeTemplate.latest_job : undefined);
	$: modalJob = deploying ? activeJobs[deploying.id] : undefined;

	onMount(() => {
		templates.fetch();
	});

	function sizeHint(t: DeployTemplate): string {
		const total = t.artifacts.reduce((acc, a) => acc + a.size_hint_mb, 0);
		return total >= 1024 ? `${(total / 1024).toFixed(1)} GB` : `${total} MB`;
	}

	function jobFor(t: DeployTemplate): ImageJob | undefined {
		return activeJobs[t.id] ?? (t.latest_job?.state === 'running' ? t.latest_job : undefined);
	}

	function jobPercent(j: ImageJob | undefined): number {
		if (!j) return 0;
		if (j.state === 'done') return 100;
		let done = 0;
		let total = 0;
		for (const a of j.artifacts) {
			done += a.bytes_done;
			total += a.bytes_total;
		}
		if (total <= 0) return 10;
		return Math.min(99, Math.round((done / total) * 100));
	}

	function openDeploy(t: DeployTemplate) {
		deploying = t;
		vmName = `${t.id}-${Date.now().toString(36).slice(-4)}`;
		cpus = t.default_cpus;
		memoryMb = t.default_memory_mb;
		autoStart = true;
		modalOpen = true;
	}

	async function confirmDeploy() {
		if (!deploying || busy) return;
		busy = true;
		const target = deploying;
		const res = await templates.deploy(
			target,
			{ name: vmName, cpus, memoryMb, autoStart },
			(job) => {
				if (job) activeJobs = { ...activeJobs, [target.id]: job };
			}
		);
		busy = false;
		if (!res.ok) {
			toasts.error(res.error ?? 'Deployment failed.');
			await templates.fetch();
			return;
		}
		modalOpen = false;
		if (res.startError) {
			toasts.warning(`VM created, but the start failed: ${res.startError}`);
		} else {
			toasts.success(autoStart ? 'VM deployed and started.' : 'VM created.');
		}
		goto('/vms');
	}

	async function installRuntime() {
		if (installing) return;
		installing = true;
		const res = await templates.installRuntime((job) => {
			if (job) activeJobs = { ...activeJobs, 'firecracker-runtime': job };
		});
		installing = false;
		if (res.ok) {
			toasts.success('Firecracker runtime installed.');
		} else {
			toasts.error(res.error ?? 'Installation failed.');
		}
	}
</script>

<div class="space-y-6">
	<PageHeader
		title="Deploy"
		description="One-click microVMs from prebuilt guest images."
		crumbs={[{ label: 'Deploy' }]}
	/>

	{#if !$templates.kvmAvailable}
		<StatusBanner
			variant="error"
			message="KVM is unavailable on this host. VMs cannot start until /dev/kvm is accessible."
		/>
	{/if}

	{#if !$templates.loading && !runtimeInstalled}
		<div
			class="flex items-start gap-3 rounded-lg border p-4 bg-status-warn/10 border-status-warn/30"
			role="status"
		>
			<span class="mt-0.5 text-status-warn-text"><Icon name="ph-warning-circle" size={20} /></span>
			<div class="flex-1">
				<p class="text-sm text-foreground">
					The Firecracker runtime is not installed. VMs cannot start without it.
				</p>
				<div class="mt-3">
					<Button size="sm" on:click={installRuntime} loading={installing}>
						<Icon name="ph-download-simple" size={14} />
						Install Firecracker {$templates.runtime.version}
					</Button>
				</div>
				{#if runtimeJob && runtimeJob.state === 'running'}
					<div class="mt-3">
						<div class="h-1.5 rounded-full bg-muted overflow-hidden">
							<div
								class="h-full bg-foreground/70 transition-all duration-300"
								style="width: {jobPercent(runtimeJob)}%"
							></div>
						</div>
						<p class="mt-1.5 font-mono text-xs text-muted-foreground">
							{runtimeJob.phase} · {jobPercent(runtimeJob)}%
						</p>
					</div>
				{/if}
			</div>
		</div>
	{/if}

	{#if runtimeInstalled}
		<div class="flex items-center gap-2 text-xs text-muted-foreground">
			<Icon name="ph-check-circle" size={14} />
			Firecracker runtime <span class="font-mono">{$templates.runtime.version}</span> installed at
			<span class="font-mono truncate">{$templates.runtime.path}</span>
		</div>
	{/if}

	{#if $templates.error}
		<StatusBanner variant="error" message={$templates.error} />
	{/if}

	{#if $templates.loading}
		<div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
			{#each Array(3) as _}
				<div class="card p-5 space-y-3">
					<div class="skeleton h-5 w-2/3"></div>
					<div class="skeleton h-3 w-full"></div>
					<div class="skeleton h-3 w-4/5"></div>
				</div>
			{/each}
		</div>
	{:else}
		<div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
			{#each osTemplates as t (t.id)}
				{@const job = jobFor(t)}
				{@const failed = !job && t.latest_job?.state === 'error'}
				<div class="card flex flex-col p-5 gap-3">
					<div class="flex items-start justify-between gap-3">
						<div class="flex items-center gap-2.5 min-w-0">
							<span
								class="flex h-9 w-9 shrink-0 items-center justify-center rounded-md bg-muted"
								aria-hidden="true"
							>
								<img
									src={t.logo === 'firecracker' ? '/logo.png' : `/brand/${t.logo}.svg`}
									alt=""
									class="h-5 w-5"
									loading="lazy"
								/>
							</span>
							<div class="min-w-0">
								<h3 class="text-sm font-semibold truncate">{t.name}</h3>
								<p class="font-mono text-xs text-muted-foreground">
									{t.os} · {t.version} · {t.arch}
								</p>
							</div>
						</div>
						{#if t.ready}
							<span
								class="inline-flex items-center gap-1.5 text-xs font-medium text-status-ok-text shrink-0"
							>
								<span class="pulse-dot" style="background: oklch(var(--status-ok))"></span>
								Ready
							</span>
						{:else if job}
							<span class="text-xs font-medium text-status-warn-text shrink-0">Downloading</span>
						{:else if failed}
							<span
								class="inline-flex items-center gap-1.5 text-xs font-medium text-status-error-text shrink-0"
							>
								<span class="pulse-dot" style="background: oklch(var(--status-error))"></span>
								Failed
							</span>
						{:else}
							<span
								class="inline-flex items-center gap-1.5 text-xs font-medium text-muted-foreground shrink-0"
							>
								<span class="h-1.5 w-1.5 rounded-full bg-muted-foreground/50"></span>
								Not downloaded
							</span>
						{/if}
					</div>

					<p class="text-[13px] leading-relaxed text-muted-foreground">{t.description}</p>

					<div class="space-y-1.5 text-xs">
						{#each t.artifacts as a (a.filename)}
							<div class="flex items-center justify-between gap-2">
								<span class="text-muted-foreground truncate">{a.name}</span>
								<span class="font-mono text-muted-foreground shrink-0">
									{a.ready && a.size_bytes ? formatBytes(a.size_bytes) : `~${a.size_hint_mb} MB`}
								</span>
							</div>
						{/each}
						<div class="flex items-center justify-between gap-2 pt-1 border-t">
							<span class="text-muted-foreground">Download size</span>
							<span class="font-mono text-muted-foreground">{sizeHint(t)}</span>
						</div>
					</div>

					{#if job && job.state === 'running'}
						<div>
							<div class="h-1.5 rounded-full bg-muted overflow-hidden">
								<div
									class="h-full bg-foreground/70 transition-all duration-300"
									style="width: {jobPercent(job)}%"
								></div>
							</div>
							<p class="mt-1.5 font-mono text-xs text-muted-foreground">
								{job.phase} · {jobPercent(job)}%
							</p>
						</div>
					{/if}
					{#if failed}
						<p class="text-xs text-status-error-text">{t.latest_job?.error}</p>
					{/if}

					<div class="mt-auto pt-2 flex items-center justify-between gap-2">
						<a
							href={t.docs_url}
							target="_blank"
							rel="noreferrer"
							class="text-xs text-muted-foreground hover:text-foreground inline-flex items-center gap-1"
						>
							<Icon name="ph-arrow-square-out" size={12} />
							Source
						</a>
						<Button size="sm" on:click={() => openDeploy(t)} disabled={job?.state === 'running'}>
							<Icon name="ph-rocket-launch" size={14} />
							{t.ready ? 'Deploy' : 'Download and deploy'}
						</Button>
					</div>
				</div>
			{/each}
		</div>

		<p class="text-xs text-muted-foreground">
			Images are cached in <span class="font-mono">~/.local/share/agni/images</span> and pinned from
			the official Firecracker CI asset bucket. Guest consoles stream on
			<span class="font-mono">ttyS0</span> over the VM logs view.
		</p>
	{/if}

	<Modal
		open={modalOpen}
		on:close={() => (modalOpen = false)}
		title={deploying ? `Deploy ${deploying.name}` : ''}
		description={deploying && !deploying.ready
			? 'The guest images are downloaded first, then the VM is created and started.'
			: 'The VM is created from the cached images and started immediately.'}
	>
		{#if deploying}
			<div class="space-y-4">
				{#if !$templates.kvmAvailable}
					<StatusBanner
						variant="error"
						message="KVM is unavailable on this host. The VM will be created but cannot start."
					/>
				{/if}
				{#if !runtimeInstalled}
					<StatusBanner
						variant="warning"
						message="The Firecracker runtime is not installed. Install it above first, or the start will fail."
					/>
				{/if}
				<div>
					<label class="label" for="deploy-name">VM name</label>
					<input
						id="deploy-name"
						class="input mt-1.5"
						type="text"
						bind:value={vmName}
						spellcheck="false"
					/>
				</div>
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label class="label" for="deploy-cpus">vCPUs</label>
						<input
							id="deploy-cpus"
							class="input mt-1.5"
							type="number"
							min="1"
							max="16"
							bind:value={cpus}
						/>
					</div>
					<div>
						<label class="label" for="deploy-mem">Memory (MB)</label>
						<input
							id="deploy-mem"
							class="input mt-1.5"
							type="number"
							min="128"
							max="65536"
							step="128"
							bind:value={memoryMb}
						/>
					</div>
				</div>
				<label class="flex items-center gap-2.5 text-sm cursor-pointer select-none">
					<input type="checkbox" class="h-4 w-4 accent-foreground" bind:checked={autoStart} />
					Start the VM immediately after creation
				</label>
				<div
					class="rounded-md bg-muted/60 border border-border p-3 space-y-1.5 font-mono text-xs text-muted-foreground"
				>
					<div class="flex justify-between gap-3">
						<span>kernel</span>
						<span class="truncate"
							>{deploying.artifacts.find((a) => a.kind === 'kernel')?.filename}</span
						>
					</div>
					{#if deploying.artifacts.find((a) => a.kind === 'rootfs')}
						<div class="flex justify-between gap-3">
							<span>rootfs</span>
							<span class="truncate"
								>{deploying.artifacts.find((a) => a.kind === 'rootfs')?.filename}</span
							>
						</div>
					{/if}
					{#if deploying.artifacts.find((a) => a.kind === 'initrd')}
						<div class="flex justify-between gap-3">
							<span>initrd</span>
							<span class="truncate"
								>{deploying.artifacts.find((a) => a.kind === 'initrd')?.filename}</span
							>
						</div>
					{/if}
					<div class="flex justify-between gap-3">
						<span>boot</span>
						<span class="truncate">{deploying.kernel_opts}</span>
					</div>
				</div>
				{#if modalJob && modalJob.state === 'running'}
					<div>
						<div class="h-1.5 rounded-full bg-muted overflow-hidden">
							<div
								class="h-full bg-foreground/70 transition-all duration-300"
								style="width: {jobPercent(modalJob)}%"
							></div>
						</div>
						<p class="mt-1.5 font-mono text-xs text-muted-foreground">
							{modalJob.phase} · {jobPercent(modalJob)}%
						</p>
					</div>
				{/if}
			</div>
		{/if}
		<svelte:fragment slot="footer">
			<Button variant="ghost" on:click={() => (modalOpen = false)} disabled={busy}>Cancel</Button>
			<Button on:click={confirmDeploy} loading={busy}>
				<Icon name="ph-rocket-launch" size={14} />
				{deploying?.ready ? 'Deploy and launch' : 'Download, deploy and launch'}
			</Button>
		</svelte:fragment>
	</Modal>
</div>
