<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type SystemStorage } from '$lib/api/client';
	import Icon from '$lib/components/Common/Icon.svelte';

	let storage: SystemStorage | null = null;
	let failed = false;

	function fmt(bytes: number): string {
		if (bytes >= 1 << 30) return (bytes / (1 << 30)).toFixed(1) + ' GB';
		if (bytes >= 1 << 20) return (bytes / (1 << 20)).toFixed(0) + ' MB';
		return (bytes / (1 << 10)).toFixed(0) + ' KB';
	}

	onMount(() => {
		const load = () =>
			api
				.getStorageUsage()
				.then((s) => ((storage = s), (failed = false)))
				.catch(() => (failed = true));
		load();
		const t = setInterval(load, 15000);
		return () => clearInterval(t);
	});

	$: usedPct = storage && storage.disk_total_bytes
		? Math.min(100, Math.round((storage.disk_total_bytes - storage.disk_free_bytes) / storage.disk_total_bytes * 100))
		: 0;
	$: top = storage
		? [...storage.templates].filter((t) => !t.is_cache).sort((a, b) => b.bytes - a.bytes).slice(0, 4)
		: [];
</script>

<div class="rounded-lg border bg-card shadow-sm p-5">
	<div class="flex items-start justify-between gap-4">
		<div class="flex items-center gap-3">
			<span class="h-11 w-11 shrink-0 flex items-center justify-center rounded-lg bg-accent text-foreground">
				<Icon name="ph-hard-drives" size={20} />
			</span>
			<div>
				<h3 class="text-sm font-semibold text-foreground">Image storage</h3>
				<p class="text-xs text-muted-foreground mt-0.5">Kernels and root filesystems on disk</p>
			</div>
		</div>
		{#if storage}
			<div class="text-right">
				<div class="text-2xl font-semibold tabular-nums leading-8 text-foreground">{fmt(storage.total_bytes)}</div>
				<p class="text-xs text-muted-foreground">{fmt(storage.disk_free_bytes)} free on disk</p>
			</div>
		{:else if failed}
			<p class="text-xs text-muted-foreground">Unavailable</p>
		{/if}
	</div>

	{#if storage}
		<div class="mt-4">
			<div class="h-1.5 w-full rounded-full bg-accent overflow-hidden">
				<div class="h-full rounded-full bg-foreground/70" style="width: {usedPct}%"></div>
			</div>
			<div class="mt-1.5 flex justify-between text-[11px] text-muted-foreground">
				<span>Disk {usedPct}% used</span>
				<span>{storage.templates.length} cached images</span>
			</div>
		</div>

		{#if top.length}
			<ul class="mt-4 space-y-2">
				{#each top as t (t.id)}
					<li class="flex items-center justify-between gap-3 text-[13px]">
						<span class="flex items-center gap-2 min-w-0">
							<span
								class="h-1.5 w-1.5 rounded-full {t.ready ? 'bg-status-ok' : 'bg-muted-foreground/40'} shrink-0"
								title={t.ready ? 'Ready to deploy' : 'Incomplete'}
							></span>
							<span class="truncate font-mono text-muted-foreground">{t.id}</span>
						</span>
						<span class="tabular-nums text-muted-foreground">{fmt(t.bytes)}</span>
					</li>
				{/each}
			</ul>
		{/if}
	{/if}
</div>
