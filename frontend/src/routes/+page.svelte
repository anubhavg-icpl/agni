<script lang="ts">
	import { onMount } from 'svelte';
	import { vms } from '$lib/stores/vms';
	import PageHeader from '$lib/components/Common/PageHeader.svelte';
	import Button from '$lib/components/Common/Button.svelte';
	import Icon from '$lib/components/Common/Icon.svelte';
	import StatCards from '$lib/components/Dashboard/StatCards.svelte';
	import VMTable from '$lib/components/Dashboard/VMTable.svelte';

	// SvelteKit passes params to all route components
	export let params: Record<string, string> = {};

	onMount(() => {
		vms.fetch();
		const interval = setInterval(() => vms.fetch(), 5000);
		return () => clearInterval(interval);
	});
</script>

<svelte:head>
	<title>Dashboard · Agni</title>
</svelte:head>

<div class="space-y-6">
	<PageHeader
		title="Dashboard"
		description="Health and activity of the Firecracker microVMs on this machine."
	>
		<div slot="actions">
			<Button variant="primary" size="sm" href="/vms/new">
				<Icon name="ph-plus" size={16} />
				Create VM
			</Button>
		</div>
	</PageHeader>

	<StatCards vms={$vms.vms} />

	<section aria-label="Virtual machines">
		<h2 class="text-lg font-semibold text-foreground">Virtual machines</h2>
		<div class="mt-4">
			<VMTable
				vms={$vms.vms}
				loading={$vms.loading}
				error={$vms.error}
				onRetry={() => vms.fetch()}
			/>
		</div>
	</section>
</div>
