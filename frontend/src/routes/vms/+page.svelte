<script lang="ts">
	import { onMount } from 'svelte';
	import { vms } from '$lib/stores/vms';
	import PageHeader from '$lib/components/Common/PageHeader.svelte';
	import Button from '$lib/components/Common/Button.svelte';
	import Icon from '$lib/components/Common/Icon.svelte';
	import VMTable from '$lib/components/Dashboard/VMTable.svelte';

	onMount(() => {
		vms.fetch();
		const interval = setInterval(() => vms.fetch(), 5000);
		return () => clearInterval(interval);
	});
</script>

<svelte:head>
	<title>Virtual Machines · Agni</title>
</svelte:head>

<div class="space-y-6">
	<PageHeader
		title="Virtual Machines"
		description="Create, start, stop, and remove Firecracker microVMs."
	>
		<div slot="actions">
			<Button variant="primary" size="sm" href="/vms/new">
				<Icon name="ph-plus" size={16} />
				Create VM
			</Button>
		</div>
	</PageHeader>

	<VMTable vms={$vms.vms} loading={$vms.loading} error={$vms.error} onRetry={() => vms.fetch()} />
</div>
