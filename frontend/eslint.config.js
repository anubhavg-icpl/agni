import js from '@eslint/js';
import ts from 'typescript-eslint';
import svelte from 'eslint-plugin-svelte';
import globals from 'globals';

export default [
	js.configs.recommended,
	...ts.configs.recommended,
	...svelte.configs['flat/recommended'],
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node
			}
		}
	},
	{
		files: ['**/*.svelte'],
		languageOptions: {
			parserOptions: {
				parser: ts.parser
			}
		}
	},
	{
		// Disable overly strict rules for our use case
		rules: {
			// Allow each blocks without keys for simple cases
			'svelte/require-each-key': 'off',
			// Don't require event dispatcher types (TypeScript inference works)
			'svelte/require-event-dispatcher-types': 'off',
			// Allow navigation without resolve() - we're not using prerendering
			'svelte/no-navigation-without-resolve': 'off',
			// Allow unused vars prefixed with _
			'@typescript-eslint/no-unused-vars': [
				'error',
				{ argsIgnorePattern: '^_', varsIgnorePattern: '^_' }
			]
		}
	},
	{
		ignores: ['build/', '.svelte-kit/', 'dist/']
	}
];
