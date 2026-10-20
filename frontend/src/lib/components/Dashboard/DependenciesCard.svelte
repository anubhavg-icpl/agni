<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Dependency } from '$lib/api/client';
	import Icon from '$lib/components/Common/Icon.svelte';

	let deps: Dependency[] = [];
	let failed = false;

	onMount(() => {
		const load = () =>
			api
				.getDependencies()
				.then((d) => ((deps = d.dependencies), (failed = false)))
				.catch(() => (failed = true));
		load();
		const t = setInterval(load, 30000);
		return () => clearInterval(t);
	});

	$: missingRequired = deps.filter((d) => d.required && !d.available);
</script>

<div class="rounded-lg border bg-card shadow-sm p-5">
	<div class="flex items-center justify-between gap-4">
		<div class="flex items-center gap-3">
			<span class="h-11 w-11 shrink-0 flex items-center justify-center rounded-lg bg-accent text-foreground">
				<Icon name="ph-seal-check" size={20} />
			</span>
			<div>
				<h3 class="text-sm font-semibold text-foreground">Environment</h3>
				<p class="text-xs text-muted-foreground mt-0.5">Host tools checked before you deploy</p>
			</div>
		</div>
		{#if !failed && deps.length}
			<span
				class="inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-medium
				{missingRequired.length ? 'text-status-err border-status-err/30 bg-status-err/10' : 'text-status-ok border-status-ok/30 bg-status-ok/10'}"
			>
				{missingRequired.length ? `${missingRequired.length} missing` : 'All required present'}
			</span>
		{:else if failed}
			<p class="text-xs text-muted-foreground">Unavailable</p>
		{/if}
	</div>

	<ul class="mt-4 space-y-2.5">
		{#each deps as dep (dep.id)}
			{@const tone = dep.available
				? 'text-status-ok'
				: dep.required
					? 'text-status-err'
					: 'text-status-warn'}
			<li class="flex items-start gap-2.5">
				<span class={tone} style="line-height:0"><Icon name={dep.available ? 'ph-check-circle' : 'ph-warning-circle'} size={16} /></span>
				<span class="min-w-0">
					<span class="block text-[13px] text-foreground">
						{dep.label}
						{#if !dep.required}
							<span class="text-muted-foreground">· optional</span>
						{/if}
					</span>
					{#if dep.available && dep.detail}
						<span class="block truncate font-mono text-[11px] text-muted-foreground">{dep.detail}</span>
					{:else if dep.hint}
						<span class="block text-[11px] text-muted-foreground">{dep.hint}</span>
					{/if}
				</span>
			</li>
		{/each}
	</ul>
</div>
