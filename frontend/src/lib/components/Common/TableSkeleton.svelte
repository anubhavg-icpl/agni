<script lang="ts">
	export let rows = 6;
</script>

<div class="p-4 space-y-4" aria-hidden="true">
	{#each Array(rows) as _}
		<div class="flex items-center gap-4">
			<div class="skeleton h-4 w-40"></div>
			<div class="skeleton h-4 w-20"></div>
			<div class="skeleton h-4 w-14 ml-auto"></div>
			<div class="skeleton h-4 w-14"></div>
			<div class="skeleton h-4 w-24"></div>
		</div>
	{/each}
</div>
