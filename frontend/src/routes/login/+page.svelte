<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	let username = '';
	let password = '';

	async function handleSubmit() {
		const success = await auth.login(username, password);
		if (success) {
			goto('/');
		}
	}
</script>

<svelte:head>
	<title>Login - Firectl</title>
</svelte:head>

<div class="min-h-screen bg-gray-900 flex items-center justify-center">
	<div class="card max-w-md w-full p-8">
		<h1 class="text-2xl font-bold text-center mb-6">Login to Firectl</h1>

		<form on:submit|preventDefault={handleSubmit} class="space-y-4">
			<div>
				<label class="label" for="username">Username</label>
				<input
					type="text"
					id="username"
					bind:value={username}
					class="input w-full"
					required
				/>
			</div>
			<div>
				<label class="label" for="password">Password</label>
				<input
					type="password"
					id="password"
					bind:value={password}
					class="input w-full"
					required
				/>
			</div>
			{#if $auth.error}
				<p class="text-red-400 text-sm">{$auth.error}</p>
			{/if}
			<button type="submit" class="btn btn-primary w-full" disabled={$auth.loading}>
				{$auth.loading ? 'Logging in...' : 'Login'}
			</button>
		</form>
	</div>
</div>
