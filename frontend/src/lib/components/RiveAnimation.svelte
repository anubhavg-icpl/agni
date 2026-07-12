<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';

	export let src: string;
	export let width: number = 120;
	export let height: number = 120;
	export let stateMachines: string | string[] | undefined = undefined;
	export let animations: string | string[] | undefined = undefined;
	export let loop: boolean = true;

	let canvas: HTMLCanvasElement;
	let riveInstance: any = null;
	let handleResize: (() => void) | null = null;

	onMount(() => {
		if (!browser) return;

		// Async load without returning a cleanup from the promise (Svelte onMount
		// cannot return a destructor from an async callback).
		void (async () => {
			const rive = await import('@rive-app/canvas');
			const Rive = rive.Rive || rive.default?.Rive;
			const Layout = rive.Layout || rive.default?.Layout;
			const Fit = rive.Fit || rive.default?.Fit;
			const Alignment = rive.Alignment || rive.default?.Alignment;

			if (!Rive) {
				console.error('Rive not found');
				return;
			}

			const config: any = {
				src,
				canvas,
				autoplay: true,
				useDevicePixelRatio: true,
				onLoad: () => {
					riveInstance?.resizeDrawingSurfaceToCanvas();
				},
				onLoop: () => {
					// Animation will loop automatically
				},
				onStop: () => {
					// Restart animation if loop is enabled
					if (loop && riveInstance) {
						riveInstance.play();
					}
				}
			};

			if (Layout && Fit && Alignment) {
				config.layout = new Layout({
					fit: Fit.Contain,
					alignment: Alignment.Center
				});
			}

			if (stateMachines) config.stateMachines = stateMachines;
			if (animations) config.animations = animations;

			riveInstance = new Rive(config);

			handleResize = () => {
				riveInstance?.resizeDrawingSurfaceToCanvas();
			};
			window.addEventListener('resize', handleResize);
		})();
	});

	onDestroy(() => {
		if (handleResize) {
			window.removeEventListener('resize', handleResize);
			handleResize = null;
		}
		if (riveInstance) {
			riveInstance.cleanup();
			riveInstance = null;
		}
	});
</script>

{#if browser}
	<canvas
		bind:this={canvas}
		width={width * 2}
		height={height * 2}
		style="width: {width}px; height: {height}px;"
	/>
{:else}
	<div style="width: {width}px; height: {height}px;"></div>
{/if}
