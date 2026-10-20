<script lang="ts">
	import { page } from '$app/stores';
	import Button from '$lib/components/Common/Button.svelte';

	// SvelteKit passes params to all route components
	export let params: Record<string, string> = {};

	$: status = $page.status ?? 404;
	$: heading = status === 404 ? 'Page not found' : 'Something went wrong';
	$: body =
		status === 404
			? 'The page you are looking for does not exist or has moved.'
			: 'An unexpected error occurred. Return to the dashboard and try again.';
</script>

<svelte:head>
	<title>{heading} · Agni</title>
</svelte:head>

<div class="min-h-dvh bg-background flex items-center justify-center p-6">
	<div class="text-center max-w-sm">
		<div class="flex items-center justify-center gap-2.5 mb-6">
			<img src="/logo.png" alt="Agni" class="w-[30px] h-[30px]" />
			<span class="text-base font-semibold text-foreground">Agni</span>
		</div>
		<p class="font-mono text-sm text-muted-foreground mb-2">{status}</p>
		<h1 class="text-xl font-semibold leading-7 text-foreground text-balance">{heading}</h1>
		<p class="text-sm text-muted-foreground mt-2 text-pretty">{body}</p>
		<div class="mt-6 flex justify-center">
			<Button variant="primary" href="/">Go to dashboard</Button>
		</div>
	</div>
</div>
