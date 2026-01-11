<script lang="ts">
	import { onMount } from 'svelte';
	import type { VMMetrics } from '$lib/api/client';

	export let data: VMMetrics[] = [];

	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D | null;

	const CHART_HEIGHT = 200;
	const PADDING = 40;

	function draw() {
		if (!ctx || !canvas) return;
		const context = ctx; // TypeScript narrowing for callbacks

		const width = canvas.width;
		const height = CHART_HEIGHT;

		// Clear canvas
		ctx.clearRect(0, 0, width, height);

		// Draw background
		ctx.fillStyle = '#1f2937';
		ctx.fillRect(0, 0, width, height);

		// Draw grid lines
		ctx.strokeStyle = '#374151';
		ctx.lineWidth = 1;

		// Horizontal grid lines (0%, 25%, 50%, 75%, 100%)
		for (let i = 0; i <= 4; i++) {
			const y = PADDING + ((height - PADDING * 2) * i) / 4;
			ctx.beginPath();
			ctx.moveTo(PADDING, y);
			ctx.lineTo(width - PADDING, y);
			ctx.stroke();

			// Labels
			ctx.fillStyle = '#9ca3af';
			ctx.font = '10px monospace';
			ctx.textAlign = 'right';
			ctx.fillText(`${100 - i * 25}%`, PADDING - 5, y + 4);
		}

		if (data.length < 2) {
			ctx.fillStyle = '#6b7280';
			ctx.font = '14px sans-serif';
			ctx.textAlign = 'center';
			ctx.fillText('Collecting data...', width / 2, height / 2);
			return;
		}

		// Draw CPU line
		const chartWidth = width - PADDING * 2;
		const chartHeight = height - PADDING * 2;
		const pointSpacing = chartWidth / (data.length - 1);

		// Create gradient
		const gradient = ctx.createLinearGradient(0, PADDING, 0, height - PADDING);
		gradient.addColorStop(0, 'rgba(59, 130, 246, 0.3)');
		gradient.addColorStop(1, 'rgba(59, 130, 246, 0)');

		// Draw filled area
		ctx.beginPath();
		ctx.moveTo(PADDING, height - PADDING);

		data.forEach((point, i) => {
			const x = PADDING + i * pointSpacing;
			const y = PADDING + chartHeight * (1 - point.cpu_usage / 100);
			context.lineTo(x, y);
		});

		ctx.lineTo(PADDING + (data.length - 1) * pointSpacing, height - PADDING);
		ctx.closePath();
		ctx.fillStyle = gradient;
		ctx.fill();

		// Draw line
		ctx.beginPath();
		ctx.strokeStyle = '#3b82f6';
		ctx.lineWidth = 2;

		data.forEach((point, i) => {
			const x = PADDING + i * pointSpacing;
			const y = PADDING + chartHeight * (1 - point.cpu_usage / 100);

			if (i === 0) {
				context.moveTo(x, y);
			} else {
				context.lineTo(x, y);
			}
		});

		ctx.stroke();

		// Draw current value
		if (data.length > 0) {
			const latest = data[data.length - 1];
			ctx.fillStyle = '#3b82f6';
			ctx.font = 'bold 16px sans-serif';
			ctx.textAlign = 'right';
			ctx.fillText(`${latest.cpu_usage.toFixed(1)}%`, width - PADDING, PADDING - 10);
		}
	}

	function handleResize() {
		if (canvas) {
			canvas.width = canvas.offsetWidth;
			draw();
		}
	}

	onMount(() => {
		ctx = canvas.getContext('2d');
		handleResize();

		window.addEventListener('resize', handleResize);
		return () => window.removeEventListener('resize', handleResize);
	});

	$: if (ctx && data) {
		draw();
	}
</script>

<div class="w-full" style="height: {CHART_HEIGHT}px;">
	<canvas bind:this={canvas} class="w-full h-full rounded" style="height: {CHART_HEIGHT}px;" />
</div>
