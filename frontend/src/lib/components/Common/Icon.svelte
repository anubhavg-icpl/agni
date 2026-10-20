<script lang="ts">
	// Phosphor Bold, the one icon set (Law 11). Sizes 16 / 20 / 24 only.
	// Accepts the icon name with or without the leading "ph-" prefix.
	export let name: string;
	export let size: 12 | 14 | 16 | 20 | 24 = 16;

	$: cls = name.replace(/^ph-/, '');
</script>

<i class="ph-bold ph-{cls} leading-none" style="font-size: {size}px" aria-hidden="true"></i>
