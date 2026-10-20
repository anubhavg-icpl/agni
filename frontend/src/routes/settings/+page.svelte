<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { api, type SystemInfo } from '$lib/api/client';
	import PageHeader from '$lib/components/Common/PageHeader.svelte';
	import Button from '$lib/components/Common/Button.svelte';
	import Icon from '$lib/components/Common/Icon.svelte';
	import Skeleton from '$lib/components/Common/Skeleton.svelte';
	import StatusBanner from '$lib/components/Common/StatusBanner.svelte';

	let systemInfo: SystemInfo | null = null;
	let loading = true;

	onMount(async () => {
		try {
			systemInfo = await api.getSystemInfo();
		} catch {
			systemInfo = null;
		}
		loading = false;
	});

	async function handleLogout() {
		await auth.logout();
		goto('/login');
	}

	const spec: { label: string; value: (info: SystemInfo) => string }[] = [
		{ label: 'Version', value: (i) => i.version },
		{ label: 'Firecracker', value: (i) => i.firecracker_version },
		{ label: 'Go', value: (i) => i.go_version },
		{ label: 'Operating system', value: (i) => i.os },
		{ label: 'Architecture', value: (i) => i.arch },
		{ label: 'Host CPUs', value: (i) => String(i.num_cpu) }
	];
</script>

<svelte:head>
	<title>Settings · Agni</title>
</svelte:head>

<div class="max-w-3xl space-y-6">
	<PageHeader title="Settings" description="Account, system, and API information." />

	<section class="card p-6">
		<h2 class="text-lg font-semibold text-foreground mb-4">Account</h2>
		<dl class="space-y-3">
			<div class="flex items-center justify-between gap-4">
				<dt class="text-sm text-muted-foreground">Username</dt>
				<dd class="text-sm font-medium text-foreground">{$auth.user?.username ?? 'Unknown'}</dd>
			</div>
			<div class="flex items-center justify-between gap-4">
				<dt class="text-sm text-muted-foreground">Role</dt>
				<dd class="text-sm font-medium text-foreground">Administrator</dd>
			</div>
		</dl>
		<div class="mt-6 pt-4 border-t flex justify-end">
			<Button variant="secondary" on:click={handleLogout}>
				<Icon name="ph-sign-out" size={16} />
				Sign out
			</Button>
		</div>
	</section>

	<section class="card p-6">
		<h2 class="text-lg font-semibold text-foreground mb-4">System information</h2>
		{#if loading}
			<div class="space-y-3" aria-hidden="true">
				{#each Array(3) as _}
					<Skeleton className="h-4 w-2/3" />
				{/each}
			</div>
		{:else if systemInfo}
			<dl class="grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-4">
				{#each spec as row}
					<div>
						<dt class="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
							{row.label}
						</dt>
						<dd class="font-mono text-[13px] leading-5 text-foreground mt-1">
							{row.value(systemInfo)}
						</dd>
					</div>
				{/each}
			</dl>
		{:else}
			<StatusBanner
				variant="warning"
				message="System information is unavailable. Check that the Agni server is reachable."
			/>
		{/if}
	</section>

	<section class="card p-6">
		<h2 class="text-lg font-semibold text-foreground mb-4">API</h2>
		<dl class="space-y-3">
			<div>
				<dt class="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
					Base URL
				</dt>
				<dd class="font-mono text-[13px] leading-5 text-foreground mt-1 break-all">
					{typeof window !== 'undefined' ? window.location.origin : ''}/api
				</dd>
			</div>
		</dl>
		<p class="text-xs text-muted-foreground mt-4">
			The Agni API is served by the same process as this interface.
		</p>
	</section>

	<section class="card p-6">
		<h2 class="text-lg font-semibold text-foreground mb-4">About</h2>
		<p class="text-sm text-muted-foreground max-w-prose">
			Agni manages Firecracker microVMs: create, run, and monitor lightweight virtual machines with
			fast startup times.
		</p>
		<div class="flex gap-4 mt-4">
			<a
				href="https://github.com/anubhavg-icpl/agni"
				target="_blank"
				rel="noopener noreferrer"
				class="text-sm text-primary-accent hover:underline inline-flex items-center gap-1.5"
			>
				<Icon name="ph-github-logo" size={16} />
				GitHub
			</a>
			<a
				href="https://firecracker-microvm.github.io/"
				target="_blank"
				rel="noopener noreferrer"
				class="text-sm text-primary-accent hover:underline inline-flex items-center gap-1.5"
			>
				<Icon name="ph-book-open" size={16} />
				Firecracker documentation
			</a>
		</div>
	</section>
</div>
