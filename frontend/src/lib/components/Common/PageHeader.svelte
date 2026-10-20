<script lang="ts">
	import Icon from './Icon.svelte';

	export let title = '';
	export let description = '';
	// Trail renders only from two levels deep: pass two or more crumbs
	export let crumbs: { label: string; href?: string }[] = [];
</script>

<div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
	<div class="min-w-0">
		{#if crumbs.length >= 2}
			<nav aria-label="Breadcrumb" class="mb-1">
				<ol class="flex items-center gap-1 text-xs text-muted-foreground">
					{#each crumbs as crumb, i}
						{#if i > 0}
							<li class="text-muted-foreground/60" aria-hidden="true"><Icon name="ph-caret-right" size={12} /></li>
						{/if}
						<li>
							{#if crumb.href && i < crumbs.length - 1}
								<a href={crumb.href} class="hover:text-foreground transition-colors">{crumb.label}</a>
							{:else}
								<span aria-current="page" class="text-foreground">{crumb.label}</span>
							{/if}
						</li>
					{/each}
				</ol>
			</nav>
		{/if}
		<h1 class="text-base font-medium leading-6 text-foreground text-balance">{title}</h1>
		{#if description}
			<p class="text-sm text-muted-foreground mt-1">{description}</p>
		{/if}
	</div>
	{#if $$slots.actions}
		<div class="flex items-center gap-2 shrink-0">
			<slot name="actions" />
		</div>
	{/if}
</div>
