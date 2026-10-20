/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	darkMode: 'class',
	theme: {
		colors: {
			transparent: 'transparent',
			current: 'current',
			inherit: 'inherit',
			background: 'oklch(var(--background) / <alpha-value>)',
			foreground: 'oklch(var(--foreground) / <alpha-value>)',
			card: {
				DEFAULT: 'oklch(var(--card) / <alpha-value>)',
				foreground: 'oklch(var(--foreground) / <alpha-value>)'
			},
			popover: {
				DEFAULT: 'oklch(var(--popover) / <alpha-value>)',
				foreground: 'oklch(var(--foreground) / <alpha-value>)'
			},
			primary: {
				DEFAULT: 'oklch(var(--primary) / <alpha-value>)',
				foreground: 'oklch(var(--primary-foreground) / <alpha-value>)',
				accent: 'oklch(var(--primary-accent) / <alpha-value>)'
			},
			secondary: {
				DEFAULT: 'oklch(var(--secondary) / <alpha-value>)',
				foreground: 'oklch(var(--secondary-foreground) / <alpha-value>)'
			},
			muted: {
				DEFAULT: 'oklch(var(--muted) / <alpha-value>)',
				foreground: 'oklch(var(--muted-foreground) / <alpha-value>)'
			},
			accent: {
				DEFAULT: 'oklch(var(--accent) / <alpha-value>)',
				foreground: 'oklch(var(--accent-foreground) / <alpha-value>)'
			},
			destructive: {
				DEFAULT: 'oklch(var(--destructive) / <alpha-value>)',
				foreground: 'oklch(var(--destructive-foreground) / <alpha-value>)'
			},
			border: 'oklch(var(--border) / var(--border-alpha, 1))',
			input: 'oklch(var(--input) / var(--input-alpha, 1))',
			ring: 'oklch(var(--ring) / <alpha-value>)',
			sidebar: {
				DEFAULT: 'oklch(var(--sidebar) / <alpha-value>)',
				accent: 'oklch(var(--primary) / 0.15)'
			},
			status: {
				ok: 'oklch(var(--status-ok) / <alpha-value>)',
				warn: 'oklch(var(--status-warn) / <alpha-value>)',
				error: 'oklch(var(--status-error) / <alpha-value>)',
				idle: 'oklch(var(--status-idle) / <alpha-value>)',
				'ok-text': 'oklch(var(--status-ok-text) / <alpha-value>)',
				'warn-text': 'oklch(var(--status-warn-text) / <alpha-value>)',
				'error-text': 'oklch(var(--status-error-text) / <alpha-value>)',
				'idle-text': 'oklch(var(--status-idle-text) / <alpha-value>)'
			}
		},
		fontFamily: {
			sans: ['Geist', 'Geist Fallback', 'ui-sans-serif', 'system-ui', 'sans-serif'],
			mono: ['Geist Mono', 'ui-monospace', 'Menlo', 'monospace']
		},
		/* Radius anchor 10px; sm 6 / md 8 / lg 10 / xl 14, nothing past xl */
		borderRadius: {
			none: '0px',
			sm: '6px',
			md: '8px',
			DEFAULT: '8px',
			lg: '10px',
			xl: '14px',
			full: '9999px'
		},
		keyframes: {
			shimmer: {
				'0%, 100%': { opacity: '1' },
				'50%': { opacity: '0.55' }
			},
			'status-pulse': {
				'0%, 100%': { opacity: '1' },
				'50%': { opacity: '0.45' }
			}
		},
		animation: {
			shimmer: 'shimmer 1.6s ease-in-out infinite',
			'status-pulse': 'status-pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite'
		}
	},
	plugins: []
};
