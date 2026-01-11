<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Button from '../Common/Button.svelte';
	import { validate, username, minLength } from '$lib/utils/validation';

	export let loading = false;
	export let error = '';

	let usernameValue = '';
	let passwordValue = '';
	let usernameError = '';
	let passwordError = '';

	const dispatch = createEventDispatcher();

	function validateForm(): boolean {
		let valid = true;

		const usernameResult = validate(usernameValue, username());
		if (!usernameResult.valid) {
			usernameError = usernameResult.error || '';
			valid = false;
		} else {
			usernameError = '';
		}

		const passwordResult = validate(passwordValue, minLength(1, 'Password'));
		if (!passwordResult.valid) {
			passwordError = passwordResult.error || '';
			valid = false;
		} else {
			passwordError = '';
		}

		return valid;
	}

	function handleSubmit() {
		if (!validateForm()) return;

		dispatch('submit', {
			username: usernameValue,
			password: passwordValue
		});
	}
</script>

<form on:submit|preventDefault={handleSubmit} class="space-y-6">
	<div class="space-y-4">
		<div>
			<label for="username" class="label">Username</label>
			<input
				type="text"
				id="username"
				bind:value={usernameValue}
				class="input w-full {usernameError ? 'border-red-500' : ''}"
				placeholder="admin"
				autocomplete="username"
				disabled={loading}
			/>
			{#if usernameError}
				<p class="text-sm text-red-400 mt-1">{usernameError}</p>
			{/if}
		</div>

		<div>
			<label for="password" class="label">Password</label>
			<input
				type="password"
				id="password"
				bind:value={passwordValue}
				class="input w-full {passwordError ? 'border-red-500' : ''}"
				placeholder="Enter your password"
				autocomplete="current-password"
				disabled={loading}
			/>
			{#if passwordError}
				<p class="text-sm text-red-400 mt-1">{passwordError}</p>
			{/if}
		</div>
	</div>

	{#if error}
		<div class="bg-red-500/20 text-red-400 border border-red-500/50 rounded-lg p-3 text-sm">
			{error}
		</div>
	{/if}

	<Button type="submit" variant="primary" fullWidth {loading}>
		{loading ? 'Signing in...' : 'Sign In'}
	</Button>
</form>
