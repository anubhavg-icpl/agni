/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	darkMode: 'class',
	theme: {
		extend: {
			colors: {
				// Fire-themed primary colors matching the Agni logo
				primary: {
					50: '#fff7ed',
					100: '#ffedd5',
					200: '#fed7aa',
					300: '#fdba74',
					400: '#fb923c',
					500: '#f97316',
					600: '#ea580c',
					700: '#c2410c',
					800: '#9a3412',
					900: '#7c2d12',
					950: '#431407'
				},
				// Accent flame colors
				flame: {
					orange: '#fb923c',
					red: '#ef4444',
					yellow: '#fbbf24',
					ember: '#dc2626'
				}
			},
			backgroundImage: {
				'fire-gradient': 'linear-gradient(135deg, #fb923c 0%, #ef4444 50%, #fbbf24 100%)',
				'ember-glow': 'radial-gradient(ellipse at center, rgba(251, 146, 60, 0.15) 0%, transparent 70%)'
			},
			boxShadow: {
				'fire': '0 0 20px rgba(251, 146, 60, 0.3)',
				'fire-lg': '0 0 40px rgba(251, 146, 60, 0.4)',
				'ember': '0 4px 20px rgba(234, 88, 12, 0.25)'
			},
			animation: {
				'pulse-fire': 'pulse-fire 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
				'glow': 'glow 2s ease-in-out infinite alternate'
			},
			keyframes: {
				'pulse-fire': {
					'0%, 100%': { opacity: '1' },
					'50%': { opacity: '0.7' }
				},
				'glow': {
					'from': { boxShadow: '0 0 10px rgba(251, 146, 60, 0.3)' },
					'to': { boxShadow: '0 0 25px rgba(251, 146, 60, 0.5)' }
				}
			}
		}
	},
	plugins: []
};
