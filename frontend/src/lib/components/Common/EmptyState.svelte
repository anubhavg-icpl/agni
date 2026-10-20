<script lang="ts">
	import Icon from './Icon.svelte';
	import Button from './Button.svelte';

	export let icon = 'ph-cubes';
	export let title = '';
	export let description = '';
	export let actionLabel = '';
	export let actionHref: string | undefined = undefined;
</script>

<div class="flex flex-col items-center justify-center text-center py-16 px-6">
	<div class="text-muted-foreground mb-4">
		<Icon name={icon} size={24} />
	</div>
	<h3 class="text-xl font-semibold leading-7 text-foreground text-balance">{title}</h3>
	{#if description}
		<p class="text-sm text-muted-foreground mt-2 max-w-sm text-pretty">{description}</p>
	{/if}
	{#if actionLabel}
		<div class="mt-6">
			<Button variant="primary" size="md" href={actionHref}>{actionLabel}</Button>
		</div>
	{/if}
</div>
