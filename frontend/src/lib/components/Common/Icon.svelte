<script lang="ts">
	// Phosphor Bold, the one icon set (Law 11). Sizes 16 / 20 / 24 only.
	export let name: string;
	export let size: 16 | 20 | 24 = 16;
</script>

<i
	class="ph-bold ph-{name} leading-none"
	style="font-size: {size}px"
	aria-hidden="true"
></i>
