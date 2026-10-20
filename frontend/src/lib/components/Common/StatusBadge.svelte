<script lang="ts">
	import { vmStatus } from '$lib/utils/status';

	export let status: string;
	export let showLabel = true;

	$: info = vmStatus(status);
</script>

<span class="inline-flex items-center gap-1.5 whitespace-nowrap">
	{#if info.spinner}
		<svg class="animate-spin h-3 w-3 text-status-{info.tone}" viewBox="0 0 24 24" aria-hidden="true">
			<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none" />
			<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
		</svg>
	{:else}
		<span
			class="h-2 w-2 rounded-full bg-status-{info.tone} {info.pulse ? 'pulse-dot' : ''}"
			aria-hidden="true"
		></span>
	{/if}
	{#if showLabel}
		<span class="text-xs font-medium text-status-{info.tone}-text">{info.label}</span>
	{/if}
</span>
