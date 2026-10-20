<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';
	import { theme } from '$lib/stores/theme';
	import { FireEngine, pickParticleCount } from './fireEngine';

	// Thin lifecycle wrapper around the fire engine: sizing, visibility,
	// theme swaps, and teardown. WebGL failure leaves the tinted panel and
	// lockup in place, never a broken box.
	let container: HTMLElement;
	let canvas: HTMLCanvasElement;

	let engine: FireEngine | null = null;
	let observer: ResizeObserver | null = null;
	let cleanup: (() => void) | null = null;

	onMount(() => {
		if (!browser || !canvas) return;

		const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
		const particles = pickParticleCount(
			window.devicePixelRatio,
			window.innerWidth,
			navigator.hardwareConcurrency
		);

		try {
			engine = new FireEngine(canvas, {
				particles,
				reducedMotion: reduced,
				theme: document.documentElement.classList.contains('dark') ? 'dark' : 'light'
			});
		} catch {
			return;
		}

		const applySize = () => engine?.resize(container.clientWidth, container.clientHeight);
		observer = new ResizeObserver(applySize);
		observer.observe(container);
		applySize();

		const onVisibility = () => engine?.setVisible(!document.hidden);
		document.addEventListener('visibilitychange', onVisibility);

		const unsubscribeTheme = theme.subscribe((t) => engine?.setTheme(t));

		cleanup = () => {
			observer?.disconnect();
			observer = null;
			document.removeEventListener('visibilitychange', onVisibility);
			unsubscribeTheme();
			engine?.destroy();
			engine = null;
		};
	});

	onDestroy(() => {
		cleanup?.();
		cleanup = null;
	});
</script>

<div bind:this={container} class="absolute inset-0" aria-hidden="true">
	<canvas bind:this={canvas} class="w-full h-full block"></canvas>
</div>
