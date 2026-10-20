// Canvas charts cannot resolve CSS custom properties, so the live token
// values are read per draw and converted once, here. The only legal
// conversion site for chart colour work.
export interface ChartTheme {
	foreground: string;
	muted: string;
	accent: string;
	gridline: string;
	font: string;
}

function oklchToRgb(L: number, C: number, H: number): [number, number, number] {
	const hRad = (H * Math.PI) / 180;
	const a = C * Math.cos(hRad);
	const b = C * Math.sin(hRad);

	const l_ = L + 0.3963377774 * a + 0.2158037573 * b;
	const m_ = L - 0.1055613458 * a - 0.0638541728 * b;
	const s_ = L - 0.0894841775 * a - 1.291485548 * b;

	const l = l_ * l_ * l_;
	const m = m_ * m_ * m_;
	const s = s_ * s_ * s_;

	const r = 4.0767416621 * l - 3.3077115913 * m + 0.2309699292 * s;
	const g = -1.2684380046 * l + 2.6097574011 * m - 0.3413193965 * s;
	const bb = -0.0041960863 * l - 0.7034186147 * m + 1.707614701 * s;

	const gamma = (c: number) => {
		const v = c <= 0.0031308 ? 12.92 * c : 1.055 * Math.pow(c, 1 / 2.4) - 0.055;
		return Math.round(Math.min(1, Math.max(0, v)) * 255);
	};

	return [gamma(r), gamma(g), gamma(bb)];
}

function channelColor(name: string, alpha = 1): string {
	const raw = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
	const [channels, alphaPart] = raw.split('/');
	const [L, C = 0, H = 0] = channels.trim().split(/\s+/).map(Number);
	const a = alphaPart ? parseFloat(alphaPart) * alpha : alpha;
	const [r, g, b] = oklchToRgb(L, C, H);
	return `rgba(${r}, ${g}, ${b}, ${a})`;
}

export function chartTheme(): ChartTheme {
	return {
		foreground: channelColor('--foreground'),
		muted: channelColor('--muted-foreground'),
		accent: channelColor('--primary'),
		gridline: channelColor('--foreground', 0.07),
		font: '12px "Geist Mono", ui-monospace, monospace'
	};
}
