// Fire engine for the Agni login panel. Architecture follows the UecWeb
// living-shield engine: raw WebGL points in one interleaved buffer, a fade
// quad that leaves short trails instead of clearing, per-theme palettes,
// additive blending in dark and premultiplied normal blending in light, DPR
// capped at 2, adaptive particle counts, a still frame for reduced motion,
// and full teardown via WEBGL_lose_context.
//
// The scene is a single flame: pooled particles rise from a base line with
// turbulent sway, cool down a fire ramp (hot core to deep red), and fade as
// they taper, with a sparse twinkling ember layer above.

export type FireTheme = 'dark' | 'light';

type Rgb = readonly [number, number, number];
type Palette = {
	readonly add: boolean;
	readonly core: Rgb;
	readonly mid: Rgb;
	readonly deep: Rgb;
	readonly ember: Rgb;
	readonly bg: Rgb;
};

const PALETTES: Record<FireTheme, Palette> = {
	dark: {
		add: true,
		core: [1, 0.96, 0.85],
		mid: [1, 0.62, 0.2],
		deep: [0.86, 0.16, 0.1],
		ember: [1, 0.72, 0.3],
		bg: [1, 0.55, 0.2]
	},
	light: {
		add: false,
		core: [0.78, 0.4, 0.04],
		mid: [0.8, 0.32, 0.04],
		deep: [0.66, 0.09, 0.06],
		ember: [0.82, 0.42, 0.08],
		bg: [0.85, 0.45, 0.12]
	}
};

const TRAILS = 0.5;
const TAU = Math.PI * 2;
const FLAME_MIN = 2600;
const FLAME_MAX = 5200;
const EMBER_COUNT = 90;
const STRIDE = 7;

/** Fewer points where fill rate is scarce (dense high-DPI panels, small cores). */
export function pickParticleCount(
	dpr: number,
	viewportWidth: number,
	cores: number | undefined
): number {
	if ((dpr > 1.5 && viewportWidth < 1440) || (cores !== undefined && cores <= 4)) return FLAME_MIN;
	return FLAME_MAX;
}

const VERT = `attribute vec2 aPos; attribute vec4 aCol; attribute float aSize;
uniform vec2 uRes; uniform float uDpr; varying vec4 vCol;
void main(){ vec2 c=aPos/uRes*2.0-1.0; gl_Position=vec4(c.x,-c.y,0.0,1.0); gl_PointSize=max(1.0,aSize*uDpr); vCol=aCol; }`;
const FRAG = `precision mediump float; varying vec4 vCol;
void main(){ float d=length(gl_PointCoord-0.5); float a=pow(smoothstep(0.5,0.0,d),1.6)*vCol.a; gl_FragColor=vec4(vCol.rgb*a,a); }`;
const FADE_VERT = 'attribute vec2 aQ; void main(){ gl_Position=vec4(aQ,0.0,1.0); }';
const FADE_FRAG =
	'precision mediump float; uniform float uFade; void main(){ gl_FragColor=vec4(0.0,0.0,0.0,uFade); }';

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

export type FireOptions = {
	particles: number;
	reducedMotion: boolean;
	theme: FireTheme;
	onReady?: () => void;
};

interface Flame {
	x: number;
	y: number;
	vx: number;
	vy: number;
	life: number;
	maxLife: number;
	size: number;
	seed: number;
}

export class FireEngine {
	private readonly gl: WebGLRenderingContext;
	private readonly ptProg: WebGLProgram;
	private readonly fadeProg: WebGLProgram;
	private readonly quad: WebGLBuffer;
	private readonly vbo: WebGLBuffer;
	private readonly loc: {
		pos: number;
		col: number;
		size: number;
		q: number;
		res: WebGLUniformLocation | null;
		dpr: WebGLUniformLocation | null;
		fade: WebGLUniformLocation | null;
	};
	private readonly opts: FireOptions;
	private readonly rand: () => number = Math.random;
	private theme: FireTheme;

	private W = 0;
	private H = 0;
	private dpr = 1;
	private baseY = 0;
	private cx = 0;

	private flameCount = 0;
	private flames: Flame[] = [];
	private embers: Flame[] = [];
	private buf = new Float32Array(0);

	private t = 0;
	private raf = 0;
	private last = 0;
	private running = false;
	private readyFired = false;
	private visible = true;
	private destroyed = false;

	/** Throws when WebGL is unavailable; the caller keeps the static panel in that case. */
	constructor(
		private readonly canvas: HTMLCanvasElement,
		opts: FireOptions
	) {
		const gl = canvas.getContext('webgl', {
			premultipliedAlpha: true,
			preserveDrawingBuffer: true,
			antialias: false,
			alpha: true
		});
		if (!gl) throw new Error('WebGL unavailable');
		this.gl = gl;
		this.opts = opts;
		this.theme = opts.theme;
		this.ptProg = compile(gl, VERT, FRAG);
		this.fadeProg = compile(gl, FADE_VERT, FADE_FRAG);
		const quad = gl.createBuffer();
		const vbo = gl.createBuffer();
		if (!quad || !vbo) throw new Error('createBuffer failed');
		this.quad = quad;
		this.vbo = vbo;
		gl.bindBuffer(gl.ARRAY_BUFFER, quad);
		gl.bufferData(gl.ARRAY_BUFFER, new Float32Array([-1, -1, 1, -1, -1, 1, 1, 1]), gl.STATIC_DRAW);
		this.loc = {
			pos: gl.getAttribLocation(this.ptProg, 'aPos'),
			col: gl.getAttribLocation(this.ptProg, 'aCol'),
			size: gl.getAttribLocation(this.ptProg, 'aSize'),
			q: gl.getAttribLocation(this.fadeProg, 'aQ'),
			res: gl.getUniformLocation(this.ptProg, 'uRes'),
			dpr: gl.getUniformLocation(this.ptProg, 'uDpr'),
			fade: gl.getUniformLocation(this.fadeProg, 'uFade')
		};
	}

	resize(cssW: number, cssH: number): void {
		if (this.destroyed || cssW < 2 || cssH < 2) return;
		this.dpr = Math.min(window.devicePixelRatio || 1, 2);
		this.W = cssW;
		this.H = cssH;
		this.canvas.width = Math.round(cssW * this.dpr);
		this.canvas.height = Math.round(cssH * this.dpr);
		this.gl.viewport(0, 0, this.canvas.width, this.canvas.height);
		this.cx = cssW / 2;
		this.baseY = cssH * 0.78;
		this.build();
		if (this.opts.reducedMotion || !this.running) this.still();
	}

	/** Re-homes the palette in place; live particles carry over their positions. */
	setTheme(theme: FireTheme): void {
		if (theme === this.theme) return;
		this.theme = theme;
		this.gl.clearColor(0, 0, 0, 0);
		this.gl.clear(this.gl.COLOR_BUFFER_BIT);
		if (this.opts.reducedMotion || !this.running) this.still();
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
		const gl = this.gl;
		gl.deleteBuffer(this.vbo);
		gl.deleteBuffer(this.quad);
		gl.deleteProgram(this.ptProg);
		gl.deleteProgram(this.fadeProg);
		gl.getExtension('WEBGL_lose_context')?.loseContext();
	}

	private build(): void {
		if (this.W < 2 || this.H < 2) return;
		const count = this.opts.particles;
		if (this.flameCount !== count) {
			this.flames = [];
			for (let i = 0; i < count; i++) {
				this.flames.push(this.spawnFlame({}, true));
			}
			this.embers = [];
			for (let i = 0; i < EMBER_COUNT; i++) {
				this.embers.push(this.spawnEmber({}, true));
			}
			this.flameCount = count;
		}
		// Halo points for every 4th flame plus embers
		const capacity = (count + (count >> 2) + EMBER_COUNT + 32) * STRIDE;
		this.buf = new Float32Array(capacity);
		this.gl.bindBuffer(this.gl.ARRAY_BUFFER, this.vbo);
		this.gl.bufferData(this.gl.ARRAY_BUFFER, this.buf.byteLength, this.gl.DYNAMIC_DRAW);
		this.gl.clearColor(0, 0, 0, 0);
		this.gl.clear(this.gl.COLOR_BUFFER_BIT);
	}

	private spawnFlame(_p: Partial<Flame> = {}, warmup = false): Flame {
		const angle = this.rand() * TAU;
		const radius = Math.sqrt(this.rand()) * Math.min(this.W, this.H) * 0.035;
		const maxLife = 1.1 + this.rand() * 1.1;
		return {
			x: this.cx + Math.cos(angle) * radius,
			y: this.baseY + this.rand() * 6 + (warmup ? -this.rand() * this.H * 0.5 : 0),
			vx: (this.rand() - 0.5) * 12,
			vy: -(this.H * 0.16 + this.rand() * this.H * 0.1),
			life: warmup ? this.rand() * maxLife : 0,
			maxLife,
			size: 5 + this.rand() * 8,
			seed: this.rand() * 100
		};
	}

	private spawnEmber(_p: Partial<Flame> = {}, warmup = false): Flame {
		const maxLife = 2 + this.rand() * 1.6;
		return {
			x: this.cx + (this.rand() - 0.5) * 24,
			y: this.baseY - 10 + (warmup ? -this.rand() * this.H * 0.45 : 0),
			vx: (this.rand() - 0.5) * 46,
			vy: -(this.H * 0.28 + this.rand() * this.H * 0.16),
			life: warmup ? this.rand() * maxLife : 0,
			maxLife,
			size: 1.1 + this.rand() * 1.3,
			seed: this.rand() * 100
		};
	}

	private stepFlame(p: Flame, dt: number, i: number, out: number): number {
		p.life += dt;
		if (p.life >= p.maxLife) {
			const next = this.spawnFlame(p);
			p.x = next.x;
			p.y = next.y;
			p.vx = next.vx;
			p.vy = next.vy;
			p.life = next.life;
			p.maxLife = next.maxLife;
			p.size = next.size;
			p.seed = next.seed;
		}
		const t = p.life / p.maxLife;

		// Turbulent sway that keeps the column coherent near the base
		const sway = Math.sin(this.t * 2.1 + p.seed + p.y * 0.012) * 26;
		p.vx += (sway - p.vx) * Math.min(1, dt * 1.6);
		p.x += p.vx * dt;
		p.y += p.vy * dt;
		p.vy -= 14 * dt; // buoyancy builds as the flame rises

		const envelope = Math.sin(Math.PI * Math.pow(t, 0.62));
		const flicker = 0.82 + 0.18 * Math.sin(this.t * 21 + p.seed * 5);
		const alpha = Math.pow(1 - t, 1.35) * flicker;
		const size = p.size * envelope;

		const c = PALETTES[this.theme];
		// Ramp: core -> mid -> deep by life
		const k = t < 0.35 ? t / 0.35 : (t - 0.35) / 0.65;
		let r: number, g: number, b: number;
		if (t < 0.35) {
			r = c.core[0] + (c.mid[0] - c.core[0]) * k;
			g = c.core[1] + (c.mid[1] - c.core[1]) * k;
			b = c.core[2] + (c.mid[2] - c.core[2]) * k;
		} else {
			r = c.mid[0] + (c.deep[0] - c.mid[0]) * k;
			g = c.mid[1] + (c.deep[1] - c.mid[1]) * k;
			b = c.mid[2] + (c.deep[2] - c.mid[2]) * k;
		}

		const put = (a: number, s: number): number => {
			this.buf[out++] = p.x;
			this.buf[out++] = p.y;
			this.buf[out++] = r;
			this.buf[out++] = g;
			this.buf[out++] = b;
			this.buf[out++] = a;
			this.buf[out++] = s;
			return out;
		};

		// A wide faint halo every 4th particle builds the glow body
		if ((i & 3) === 0) out = put(alpha * 0.05, size * 6);
		out = put(alpha, size);
		return out;
	}

	private stepEmber(p: Flame, dt: number, out: number): number {
		p.life += dt;
		if (p.life >= p.maxLife) {
			const next = this.spawnEmber(p);
			p.x = next.x;
			p.y = next.y;
			p.vx = next.vx;
			p.vy = next.vy;
			p.life = next.life;
			p.maxLife = next.maxLife;
			p.size = next.size;
			p.seed = next.seed;
		}
		const t = p.life / p.maxLife;

		p.vy += 30 * dt; // sparks arc and fall back
		p.vx *= 1 - 0.22 * dt;
		p.x += p.vx * dt + Math.sin(this.t * 3 + p.seed) * 8 * dt;
		p.y += p.vy * dt;

		const twinkle = 0.6 + 0.4 * Math.sin(this.t * 17 + p.seed * 7);
		const alpha = Math.pow(1 - t, 1.1) * twinkle;
		const c = PALETTES[this.theme].ember;

		this.buf[out++] = p.x;
		this.buf[out++] = p.y;
		this.buf[out++] = c[0];
		this.buf[out++] = c[1];
		this.buf[out++] = c[2];
		this.buf[out++] = alpha;
		this.buf[out++] = p.size * (1 - t * 0.4);
		return out + 1;
	}

	private step(dt: number): void {
		this.t += dt;
		let out = 0;
		for (let i = 0; i < this.flames.length; i++) {
			out = this.stepFlame(this.flames[i], dt, i, out);
		}
		for (let i = 0; i < this.embers.length; i++) {
			out = this.stepEmber(this.embers[i], dt, out);
		}
		this.render(out / STRIDE);
	}

	private render(count: number): void {
		const gl = this.gl;
		const L = this.loc;
		gl.enable(gl.BLEND);
		// Fade the previous frame instead of clearing: short fiery trails
		gl.useProgram(this.fadeProg);
		gl.bindBuffer(gl.ARRAY_BUFFER, this.quad);
		gl.enableVertexAttribArray(L.q);
		gl.vertexAttribPointer(L.q, 2, gl.FLOAT, false, 0, 0);
		gl.blendFunc(gl.ZERO, gl.ONE_MINUS_SRC_ALPHA);
		gl.uniform1f(L.fade, this.opts.reducedMotion ? 1 : 1 - TRAILS);
		gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4);
		gl.disableVertexAttribArray(L.q);

		gl.useProgram(this.ptProg);
		gl.bindBuffer(gl.ARRAY_BUFFER, this.vbo);
		gl.bufferSubData(gl.ARRAY_BUFFER, 0, this.buf.subarray(0, count * STRIDE));
		gl.enableVertexAttribArray(L.pos);
		gl.vertexAttribPointer(L.pos, 2, gl.FLOAT, false, STRIDE * 4, 0);
		gl.enableVertexAttribArray(L.col);
		gl.vertexAttribPointer(L.col, 4, gl.FLOAT, false, STRIDE * 4, 8);
		gl.enableVertexAttribArray(L.size);
		gl.vertexAttribPointer(L.size, 1, gl.FLOAT, false, STRIDE * 4, 24);
		gl.uniform2f(L.res, this.W, this.H);
		gl.uniform1f(L.dpr, this.dpr);
		gl.blendFunc(gl.ONE, PALETTES[this.theme].add ? gl.ONE : gl.ONE_MINUS_SRC_ALPHA);
		gl.drawArrays(gl.POINTS, 0, count);
		if (!this.readyFired) {
			this.readyFired = true;
			this.opts.onReady?.();
		}
	}

	private start(): void {
		if (this.destroyed || this.running || this.opts.reducedMotion || !this.flameCount) return;
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
		const dt = Math.min(0.033, (now - this.last) / 1000);
		this.last = now;
		this.step(dt);
		this.raf = requestAnimationFrame(this.frame);
	};

	/** A settled single frame for reduced motion: a full flame, not a dead screen. */
	private still(): void {
		if (!this.flameCount) return;
		for (let i = 0; i < 45; i++) this.step(1 / 60);
		this.gl.clearColor(0, 0, 0, 0);
		this.gl.clear(this.gl.COLOR_BUFFER_BIT);
		this.step(1 / 60);
	}
}
