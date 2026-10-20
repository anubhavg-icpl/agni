<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';
	import { FireEngine, pickParticleCount } from './fireEngine';

	// Thin lifecycle wrapper around the fire engine: sizing, visibility, and
	// teardown. The scene is always the ember-night palette (the brand panel
	// keeps its cinematic look in both app themes); WebGL failure leaves the
	// warm gradient panel in place, never a broken box.
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
				theme: 'dark'
			});
		} catch {
			return;
		}

		const applySize = () => engine?.resize(container.clientWidth, container.clientHeight);
		observer = new ResizeObserver(applySize);
		observer.observe(container);
		applySize();
		// Kick the animation loop immediately; visibility events only pause it.
		engine.setVisible(!document.hidden);

		const onVisibility = () => engine?.setVisible(!document.hidden);
		document.addEventListener('visibilitychange', onVisibility);

		cleanup = () => {
			observer?.disconnect();
			observer = null;
			document.removeEventListener('visibilitychange', onVisibility);
			engine?.destroy();
			engine = null;
		};
	});

	onDestroy(() => {
		cleanup?.();
		cleanup = null;
	});
</script>

<div bind:this={container} class="absolute inset-0 pointer-events-none" aria-hidden="true">
	<canvas bind:this={canvas} class="w-full h-full block"></canvas>
</div>
