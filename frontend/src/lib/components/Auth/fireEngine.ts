// Fire engine for the Agni login panel: one full-panel fragment shader in
// the spirit of skeeto's webgl-fire and the cinematic ember references -
// the canvas paints the entire panel (near-black night, saturated
// hand-tuned fire palette, rising fbm flame tongues, drifting embers,
// breathing ground glow) so the flame blends into the UI. No filmic
// haze: colors stay punchy, only highlights above white get a soft knee.
// One draw call per frame; reduced motion renders one settled frame.

export type FireTheme = 'dark' | 'light';

export type FireOptions = {
	/** Legacy particle budget; only the quality tier (high/low) is honored. */
	particles: number;
	reducedMotion: boolean;
	theme: FireTheme;
	onReady?: () => void;
};

/** Lower octaves and resolution where fill rate or cores are scarce. */
export function pickParticleCount(
	dpr: number,
	viewportWidth: number,
	cores: number | undefined
): number {
	if ((dpr > 1.5 && viewportWidth < 1440) || (cores !== undefined && cores <= 4)) return 2600;
	return 5200;
}

const VERT = `attribute vec2 aQ;
void main(){ gl_Position = vec4(aQ, 0.0, 1.0); }`;

const FRAG = `precision highp float;
uniform vec2 uRes;
uniform float uTime;
uniform float uDark;
uniform float uOct;

float hash(vec2 p){ p = fract(p*vec2(233.34, 851.73)); p += dot(p, p+23.45); return fract(p.x*p.y); }

float noise(vec2 p){
	vec2 i = floor(p), f = fract(p);
	vec2 u = f*f*(3.0-2.0*f);
	return mix(mix(hash(i), hash(i+vec2(1,0)), u.x), mix(hash(i+vec2(0,1)), hash(i+vec2(1,1)), u.x), u.y);
}

float fbm(vec2 p){
	float a = 0.5, s = 0.0;
	for (int i = 0; i < 5; i++){
		if (float(i) >= uOct) break;
		s += a * noise(p);
		p = p*2.03 + vec2(11.7, 5.3);
		a *= 0.52;
	}
	return s;
}

// Hand-tuned saturated fire palette: ember red -> hot orange -> amber ->
// near-white core. Hard stops keep the color punchy instead of pastel.
vec3 firePal(float t){
	t = clamp(t, 0.0, 1.0);
	vec3 c = mix(vec3(0.45, 0.015, 0.0), vec3(0.98, 0.25, 0.008), smoothstep(0.0, 0.30, t));
	c = mix(c, vec3(1.0, 0.54, 0.03), smoothstep(0.30, 0.60, t));
	c = mix(c, vec3(1.0, 0.82, 0.28), smoothstep(0.60, 0.82, t));
	c = mix(c, vec3(1.0, 0.97, 0.85), smoothstep(0.82, 0.97, t));
	return c;
}

void main(){
	vec2 uv = gl_FragCoord.xy / uRes;          // y up: 0 bottom, 1 top
	float aspect = uRes.x / max(uRes.y, 1.0);
	vec2 p = vec2(uv.x * aspect, uv.y);
	float t = uTime;

	// --- Flame field -------------------------------------------------------
	vec2 q = vec2(p.x * 2.6, p.y * 2.2 - t * 1.15);
	float w = fbm(q * 1.4 + vec2(0.0, -t * 0.22));
	float n1 = fbm(q * 2.1 + w * 1.35);
	float n2 = fbm(q * 5.6 - w * 0.8 + 7.3);

	// Full-width fire bed, tallest behind the center (the wordmark sits above).
	float center = mix(0.55, 1.0, exp(-pow((uv.x - 0.5) * 2.5, 2.0) * 2.2));
	float bedWarp = sin(uv.x * 9.0 + t * 0.7) * 0.035 + (n2 - 0.5) * 0.16;
	float tongue = (0.38 + 0.52 * n1) * center + 0.09;
	float h = uv.y - bedWarp;
	float body = smoothstep(tongue, 0.04, h);
	float core = smoothstep(tongue * 0.38, 0.0, h) * (0.72 + 0.28 * n2);
	float heat = clamp(body * (0.45 + 0.55 * n2) + core * 0.65, 0.0, 1.0);
	heat *= 0.95 + 0.05 * sin(t * 10.0 + n1 * 6.2831);
	float flameA = clamp(heat * 1.35 + core * 0.35, 0.0, 1.0);
	flameA *= smoothstep(1.05, 0.55, uv.y + tongue * 0.3);

	// Breathing ground glow, saturated orange.
	float breathe = 0.9 + 0.1 * sin(t * 0.9);
	float glow = exp(-h * 5.5) * (0.4 + 0.6 * n1) * (0.55 + 0.45 * center) * breathe;

	// Embers: rising cells with a hot flicker toward amber-white.
	vec2 esc = vec2(p.x * 24.0, uv.y * 9.0 - t * 1.4);
	float sn = noise(esc);
	float sn2 = noise(esc * 1.7 + 4.7);
	float ember = smoothstep(0.82, 0.995, sn)
		* smoothstep(0.95, 0.3, uv.y + 0.08 * sin(t * 0.7 + sn2 * 6.2831))
		* (0.55 + 0.45 * sin(t * 9.0 + sn2 * 40.0));

	// --- Compose -----------------------------------------------------------
	vec3 col;
	if (uDark > 0.5) {
		// Night: clean near-black, saturated flame, tight glow. No haze.
		vec3 bg = mix(vec3(0.010, 0.005, 0.004), vec3(0.030, 0.008, 0.003), pow(1.0 - uv.y, 1.8));
		vec3 fire = firePal(heat);
		col = bg + fire * flameA;
		col += vec3(1.0, 0.34, 0.08) * glow * 0.55;
		vec3 emberCol = mix(vec3(1.0, 0.42, 0.06), vec3(1.0, 0.85, 0.45), sn2);
		col += emberCol * ember * 1.6;
		// Soft knee only above white; saturated colors pass through untouched.
		col = col / (1.0 + max(col - 1.0, vec3(0.0)) * 0.5);
	} else {
		// Day: warm cream paper with a saturated ochre flame.
		vec3 bg = mix(vec3(0.99, 0.965, 0.93), vec3(0.955, 0.86, 0.74), pow(1.0 - uv.y, 1.5));
		vec3 flame = firePal(heat) * vec3(0.92, 0.62, 0.38) + bg * 0.05;
		col = mix(bg, flame, clamp(flameA * 0.92, 0.0, 1.0));
		col += vec3(1.0, 0.62, 0.25) * ember * 0.6;
	}

	gl_FragColor = vec4(clamp(col, 0.0, 1.0), 1.0);
}
`;

function compile(gl: WebGLRenderingContext, vs: string, fs: string): WebGLProgram {
	const p = gl.createProgram();
	if (!p) throw new Error('createProgram failed');
	for (const [type, src] of [
		[gl.VERTEX_SHADER, vs],
		[gl.FRAGMENT_SHADER, fs]
	] as const) {
		const sh = gl.createShader(type);
		if (!sh) throw new Error('createShader failed');
		gl.shaderSource(sh, src);
		gl.compileShader(sh);
		if (!gl.getShaderParameter(sh, gl.COMPILE_STATUS))
			throw new Error(gl.getShaderInfoLog(sh) ?? 'shader compile failed');
		gl.attachShader(p, sh);
	}
	gl.linkProgram(p);
	if (!gl.getProgramParameter(p, gl.LINK_STATUS))
		throw new Error(gl.getProgramInfoLog(p) ?? 'program link failed');
	return p;
}

export class FireEngine {
	private readonly gl: WebGLRenderingContext;
	private readonly prog: WebGLProgram;
	private readonly quad: WebGLBuffer;
	private readonly loc: {
		q: number;
		res: WebGLUniformLocation | null;
		time: WebGLUniformLocation | null;
		dark: WebGLUniformLocation | null;
		oct: WebGLUniformLocation | null;
	};
	private readonly highQuality: boolean;
	private readonly reducedMotion: boolean;
	private onReady?: () => void;
	private theme: FireTheme;

	private W = 0;
	private H = 0;
	private t = Math.random() * 20;
	private raf = 0;
	private last = 0;
	private running = false;
	private readyFired = false;
	private visible = true;
	private destroyed = false;

	/** Throws when WebGL is unavailable; the caller keeps the static panel. */
	constructor(
		private readonly canvas: HTMLCanvasElement,
		opts: FireOptions
	) {
		const gl = canvas.getContext('webgl', {
			alpha: false,
			antialias: false,
			premultipliedAlpha: false
		});
		if (!gl) throw new Error('WebGL unavailable');
		this.gl = gl;
		this.theme = opts.theme;
		this.reducedMotion = opts.reducedMotion;
		this.highQuality = opts.particles >= 4000;
		this.onReady = opts.onReady;
		this.prog = compile(gl, VERT, FRAG);
		const quad = gl.createBuffer();
		if (!quad) throw new Error('createBuffer failed');
		this.quad = quad;
		gl.bindBuffer(gl.ARRAY_BUFFER, quad);
		gl.bufferData(gl.ARRAY_BUFFER, new Float32Array([-1, -1, 1, -1, -1, 1, 1, 1]), gl.STATIC_DRAW);
		this.loc = {
			q: gl.getAttribLocation(this.prog, 'aQ'),
			res: gl.getUniformLocation(this.prog, 'uRes'),
			time: gl.getUniformLocation(this.prog, 'uTime'),
			dark: gl.getUniformLocation(this.prog, 'uDark'),
			oct: gl.getUniformLocation(this.prog, 'uOct')
		};
	}

	resize(cssW: number, cssH: number): void {
		if (this.destroyed || cssW < 2 || cssH < 2) return;
		const dpr = Math.min(window.devicePixelRatio || 1, this.highQuality ? 1.75 : 1.25);
		let scale = dpr;
		if (cssW * cssH * scale * scale > 2.6e6) scale *= Math.sqrt(2.6e6 / (cssW * cssH * scale * scale));
		this.W = cssW;
		this.H = cssH;
		this.canvas.width = Math.max(2, Math.round(cssW * scale));
		this.canvas.height = Math.max(2, Math.round(cssH * scale));
		this.gl.viewport(0, 0, this.canvas.width, this.canvas.height);
		if (this.reducedMotion || !this.running) this.still();
	}

	setTheme(theme: FireTheme): void {
		if (theme === this.theme) return;
		this.theme = theme;
		if (this.reducedMotion || !this.running) this.still();
	}

	setVisible(visible: boolean): void {
		this.visible = visible;
		if (visible) this.start();
		else this.stop();
	}

	destroy(): void {
		if (this.destroyed) return;
		this.stop();
		this.destroyed = true;
		this.gl.deleteBuffer(this.quad);
		this.gl.deleteProgram(this.prog);
		this.gl.getExtension('WEBGL_lose_context')?.loseContext();
	}

	private render(): void {
		const gl = this.gl;
		gl.useProgram(this.prog);
		gl.bindBuffer(gl.ARRAY_BUFFER, this.quad);
		gl.enableVertexAttribArray(this.loc.q);
		gl.vertexAttribPointer(this.loc.q, 2, gl.FLOAT, false, 0, 0);
		gl.uniform2f(this.loc.res, this.canvas.width, this.canvas.height);
		gl.uniform1f(this.loc.time, this.t);
		gl.uniform1f(this.loc.dark, this.theme === 'dark' ? 1 : 0);
		gl.uniform1f(this.loc.oct, this.highQuality ? 5 : 4);
		gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4);
		if (!this.readyFired) {
			this.readyFired = true;
			Promise.resolve().then(() => this.onReady?.());
		}
	}

	private start(): void {
		if (this.destroyed || this.running || this.reducedMotion || this.W < 2 || !this.visible)
			return;
		this.running = true;
		this.last = performance.now();
		this.raf = requestAnimationFrame(this.frame);
	}

	private stop(): void {
		this.running = false;
		if (this.raf) cancelAnimationFrame(this.raf);
		this.raf = 0;
	}

	private readonly frame = (now: number): void => {
		if (!this.running) return;
		const dt = Math.min(0.05, (now - this.last) / 1000);
		this.last = now;
		this.t += dt;
		this.render();
		this.raf = requestAnimationFrame(this.frame);
	};

	/** One settled frame for reduced motion. */
	private still(): void {
		if (this.W < 2) return;
		this.t += 3.2;
		this.render();
	}
}
