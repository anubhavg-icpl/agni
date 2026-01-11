<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import RiveAnimation from '$lib/components/RiveAnimation.svelte';

	let username = '';
	let password = '';
	let confirmPassword = '';

	async function handleLogin() {
		const success = await auth.login(username, password);
		if (success) {
			goto('/');
		}
	}

	async function handleSetup() {
		if (password !== confirmPassword) {
			return;
		}
		if (password.length < 8) {
			return;
		}
		const success = await auth.setup(username, password);
		if (success) {
			goto('/');
		}
	}

	$: passwordMismatch = confirmPassword.length > 0 && password !== confirmPassword;
	$: passwordTooShort = password.length > 0 && password.length < 8;
</script>

<svelte:head>
	<title>{$auth.setupRequired ? 'Setup' : 'Login'} - Agni</title>
</svelte:head>

<div class="min-h-screen bg-gray-900 flex items-center justify-center relative overflow-hidden">
	<!-- Background with original logo art -->
	<div class="absolute inset-0 opacity-10">
		<img src="/logo-bg.png" alt="" class="w-full h-full object-cover blur-sm" aria-hidden="true" />
	</div>
	
	<!-- Fire gradient overlay -->
	<div class="absolute inset-0 bg-fire-pattern"></div>
	
	<!-- Login/Setup Card -->
	<div class="card-fire max-w-md w-full p-8 relative z-10">
		<!-- Agni Animated Logo -->
		<div class="flex flex-col items-center mb-6 relative">
			<!-- Glow effect behind logo -->
			<div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-32 h-32 bg-orange-500/20 rounded-full blur-2xl pointer-events-none"></div>
			
			<div class="mb-2 relative z-10 w-40 h-40 flex items-center justify-center bg-black rounded-full overflow-hidden" style="mix-blend-mode: screen;">
				<div class="absolute inset-0 z-20 pointer-events-none" style="background: radial-gradient(circle, transparent 50%, #000000 90%);"></div>
				<RiveAnimation 
					src="/skull-fire-logo.riv" 
					width={160} 
					height={160}
				/>
			</div>
			<h1 class="text-4xl font-bold text-fire tracking-tight">
				Agni
			</h1>
			<p class="text-gray-400 text-sm mt-2">Firecracker MicroVM Control</p>
		</div>

		{#if $auth.setupRequired}
			<!-- Setup Form -->
			<div class="mb-6 p-4 bg-orange-500/10 border border-orange-500/30 rounded-lg">
				<h2 class="text-orange-400 font-semibold mb-1">👋 Welcome to Agni!</h2>
				<p class="text-gray-400 text-sm">Create your admin account to get started.</p>
			</div>

			<form on:submit|preventDefault={handleSetup} class="space-y-5">
				<div>
					<label class="label" for="setup-username">Username</label>
					<input 
						type="text" 
						id="setup-username" 
						bind:value={username} 
						class="input w-full" 
						placeholder="Choose a username"
						required 
					/>
				</div>
				<div>
					<label class="label" for="setup-password">Password</label>
					<input 
						type="password" 
						id="setup-password" 
						bind:value={password} 
						class="input w-full" 
						class:border-red-500={passwordTooShort}
						placeholder="Minimum 8 characters"
						required 
					/>
					{#if passwordTooShort}
						<p class="text-red-400 text-xs mt-1">Password must be at least 8 characters</p>
					{/if}
				</div>
				<div>
					<label class="label" for="setup-confirm">Confirm Password</label>
					<input 
						type="password" 
						id="setup-confirm" 
						bind:value={confirmPassword} 
						class="input w-full" 
						class:border-red-500={passwordMismatch}
						placeholder="Confirm your password"
						required 
					/>
					{#if passwordMismatch}
						<p class="text-red-400 text-xs mt-1">Passwords do not match</p>
					{/if}
				</div>
				{#if $auth.error}
					<div class="bg-red-500/10 border border-red-500/30 rounded-lg p-3">
						<p class="text-red-400 text-sm">{$auth.error}</p>
					</div>
				{/if}
				<button 
					type="submit" 
					class="btn btn-primary w-full py-3 text-lg font-semibold" 
					disabled={$auth.loading || passwordMismatch || passwordTooShort || !username || !password}
				>
					{$auth.loading ? 'Creating...' : 'Create Admin Account'}
				</button>
			</form>
		{:else}
			<!-- Login Form -->
			<form on:submit|preventDefault={handleLogin} class="space-y-5">
				<div>
					<label class="label" for="login-username">Username</label>
					<input 
						type="text" 
						id="login-username" 
						bind:value={username} 
						class="input w-full" 
						placeholder="Enter your username"
						required 
					/>
				</div>
				<div>
					<label class="label" for="login-password">Password</label>
					<input 
						type="password" 
						id="login-password" 
						bind:value={password} 
						class="input w-full" 
						placeholder="Enter your password"
						required 
					/>
				</div>
				{#if $auth.error}
					<div class="bg-red-500/10 border border-red-500/30 rounded-lg p-3">
						<p class="text-red-400 text-sm">{$auth.error}</p>
					</div>
				{/if}
				<button 
					type="submit" 
					class="btn btn-primary w-full py-3 text-lg font-semibold" 
					disabled={$auth.loading || !username || !password}
				>
					{$auth.loading ? 'Signing in...' : 'Sign In'}
				</button>
			</form>
		{/if}

		<p class="text-center text-gray-500 text-xs mt-6">
			Powered by Firecracker MicroVMs
		</p>
	</div>
</div>
