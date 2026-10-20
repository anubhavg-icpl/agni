<script lang="ts">
	import { onMount } from 'svelte';
	import type { VMMetrics } from '$lib/api/client';
	import { formatBytes } from '$lib/utils/formatters';
	import { chartTheme } from '$lib/utils/chartTheme';

	export let data: VMMetrics[] = [];

	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D | null;

	const CHART_HEIGHT = 200;
	const PADDING = 40;

	function draw() {
		if (!ctx || !canvas) return;
		const context = ctx;
		const theme = chartTheme();

		const width = canvas.width;
		const height = CHART_HEIGHT;

		context.clearRect(0, 0, width, height);

		context.strokeStyle = theme.gridline;
		context.lineWidth = 1;

		for (let i = 0; i <= 4; i++) {
			const y = PADDING + ((height - PADDING * 2) * i) / 4;
			context.beginPath();
			context.moveTo(PADDING, y);
			context.lineTo(width - PADDING, y);
			context.stroke();

			context.fillStyle = theme.muted;
			context.font = theme.font;
			context.textAlign = 'right';
			context.fillText(`${100 - i * 25}%`, PADDING - 5, y + 4);
		}

		if (data.length < 2) {
			context.fillStyle = theme.muted;
			context.font = theme.font;
			context.textAlign = 'center';
			context.fillText('Collecting data...', width / 2, height / 2);
			return;
		}

		const chartWidth = width - PADDING * 2;
		const chartHeight = height - PADDING * 2;
		const pointSpacing = chartWidth / (data.length - 1);

		const gradient = context.createLinearGradient(0, PADDING, 0, height - PADDING);
		gradient.addColorStop(0, theme.accent.replace(')', ', 0.18)').replace('rgba', 'rgba'));
		gradient.addColorStop(1, theme.accent.replace(')', ', 0)').replace('rgba', 'rgba'));

		context.beginPath();
		context.moveTo(PADDING, height - PADDING);
		data.forEach((point, i) => {
			const x = PADDING + i * pointSpacing;
			const memPercent = point.memory_total > 0 ? point.memory_used / point.memory_total : 0;
			const y = PADDING + chartHeight * (1 - memPercent);
			context.lineTo(x, y);
		});
		context.lineTo(PADDING + (data.length - 1) * pointSpacing, height - PADDING);
		context.closePath();
		context.fillStyle = gradient;
		context.fill();

		context.beginPath();
		context.strokeStyle = theme.accent;
		context.lineWidth = 2;
		data.forEach((point, i) => {
			const x = PADDING + i * pointSpacing;
			const memPercent = point.memory_total > 0 ? point.memory_used / point.memory_total : 0;
			const y = PADDING + chartHeight * (1 - memPercent);
			if (i === 0) context.moveTo(x, y);
			else context.lineTo(x, y);
		});
		context.stroke();

		const latest = data[data.length - 1];
		const memPercent =
			latest.memory_total > 0
				? ((latest.memory_used / latest.memory_total) * 100).toFixed(1)
				: '0.0';
		context.fillStyle = theme.accent;
		context.font = '600 13px "Geist Mono", ui-monospace, monospace';
		context.textAlign = 'right';
		context.fillText(
			`${formatBytes(latest.memory_used)} / ${formatBytes(latest.memory_total)} (${memPercent}%)`,
			width - PADDING,
			PADDING - 10
		);
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
	<canvas
		bind:this={canvas}
		class="w-full h-full"
		style="height: {CHART_HEIGHT}px;"
		aria-label="Memory usage over time"
	></canvas>
</div>
