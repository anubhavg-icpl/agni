<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import RiveAnimation from '$lib/components/RiveAnimation.svelte';
	import { onMount } from 'svelte';
	import { fly, fade, scale } from 'svelte/transition';
	import { elasticOut, quintOut } from 'svelte/easing';

	// SvelteKit passes params to all route components
	export let params: Record<string, string> = {};

	let username = '';
	let password = '';
	let confirmPassword = '';
	let mounted = false;

	onMount(() => {
		mounted = true;
	});

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

	$: passwordInputClass = `w-full pl-12 pr-4 py-3 bg-gray-900/50 border rounded-xl text-white placeholder-gray-500 focus:ring-2 transition-all duration-200 ${
		passwordTooShort
			? 'border-red-500 focus:border-red-500 focus:ring-red-500/20'
			: 'border-gray-700/50 focus:border-orange-500/50 focus:ring-orange-500/20'
	}`;

	$: confirmInputClass = `w-full pl-12 pr-4 py-3 bg-gray-900/50 border rounded-xl text-white placeholder-gray-500 focus:ring-2 transition-all duration-200 ${
		passwordMismatch
			? 'border-red-500 focus:border-red-500 focus:ring-red-500/20'
			: 'border-gray-700/50 focus:border-orange-500/50 focus:ring-orange-500/20'
	}`;
</script>

<svelte:head>
	<title>{$auth.setupRequired ? 'Fresh Meat Detected' : 'Prove Yourself'} | Agni</title>
</svelte:head>

<div class="min-h-screen bg-gray-950 flex">
	<!-- Left Side - Rive Animation (50%) -->
	<div class="hidden lg:flex w-1/2 relative overflow-hidden bg-gray-950 items-center justify-center">
		<!-- Animated gradient background -->
		<div class="absolute inset-0 bg-gradient-to-br from-orange-900/20 via-gray-950 to-red-900/20"></div>

		<!-- Floating particles effect -->
		<div class="absolute inset-0 overflow-hidden">
			<div class="particle particle-1"></div>
			<div class="particle particle-2"></div>
			<div class="particle particle-3"></div>
			<div class="particle particle-4"></div>
			<div class="particle particle-5"></div>
		</div>

		<!-- Rive Animation Container -->
		<div class="relative z-10 flex flex-col items-center">
			<!-- Glow effect behind animation -->
			<div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-orange-500/30 rounded-full blur-[100px] animate-pulse-slow"></div>

			{#if mounted}
				<div
					class="relative w-72 h-72 flex items-center justify-center rounded-full overflow-hidden rive-container"
					in:scale={{ duration: 800, delay: 200, easing: elasticOut }}
				>
					<div class="absolute inset-0 bg-gray-950 rounded-full"></div>
					<div class="relative z-10">
						<RiveAnimation
							src="/skull-fire-logo.riv"
							width={280}
							height={280}
						/>
					</div>
				</div>

				<div
					class="mt-8 text-center"
					in:fly={{ y: 30, duration: 600, delay: 500, easing: quintOut }}
				>
					<h1 class="text-4xl font-bold bg-gradient-to-r from-orange-400 via-red-500 to-orange-400 bg-clip-text text-transparent animate-gradient">
						AGNI
					</h1>
					<p class="text-gray-500 mt-3 text-sm tracking-wider uppercase">
						Firecracker VM Manager
					</p>
					<p class="text-gray-600 mt-2 text-xs italic">
						"Your VMs won't manage themselves. Unfortunately."
					</p>
				</div>
			{/if}
		</div>

		<!-- Bottom decorative line -->
		<div class="absolute bottom-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-orange-500/50 to-transparent"></div>
	</div>

	<!-- Right Side - Login Form (50%) -->
	<div class="w-full lg:w-1/2 flex items-center justify-center p-8 relative">
		<!-- Mobile logo (shown only on small screens) -->
		<div class="absolute top-8 left-1/2 -translate-x-1/2 lg:hidden">
			<div class="w-20 h-20 flex items-center justify-center rounded-full overflow-hidden bg-gray-950">
				<RiveAnimation src="/skull-fire-logo.riv" width={76} height={76} />
			</div>
		</div>

		<div class="w-full max-w-md">
			{#if mounted}
				<div in:fly={{ x: 30, duration: 600, delay: 100, easing: quintOut }}>
					{#if $auth.setupRequired}
						<!-- Setup Header -->
						<div class="mb-8 lg:mt-0 mt-16">
							<h2 class="text-3xl font-bold text-white mb-2">Fresh Meat Detected</h2>
							<p class="text-gray-400">Someone has to run this circus. Might as well be you.</p>
						</div>

						<!-- Setup Notice -->
						<div
							class="mb-6 p-4 bg-gradient-to-r from-orange-500/10 to-red-500/10 border border-orange-500/20 rounded-xl backdrop-blur-sm"
							in:scale={{ duration: 400, delay: 300, easing: quintOut }}
						>
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 rounded-full bg-orange-500/20 flex items-center justify-center">
									<svg class="w-5 h-5 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
									</svg>
								</div>
								<div>
									<h3 class="text-orange-400 font-semibold text-sm">First Time Here?</h3>
									<p class="text-gray-500 text-xs">Create your admin account. No pressure.</p>
								</div>
							</div>
						</div>

						<!-- Setup Form -->
						<form on:submit|preventDefault={handleSetup} class="space-y-5">
							<div class="space-y-1" in:fly={{ y: 20, duration: 400, delay: 200 }}>
								<label class="block text-sm font-medium text-gray-300" for="setup-username">Username</label>
								<div class="relative">
									<div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
										<svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
										</svg>
									</div>
									<input
										type="text"
										id="setup-username"
										bind:value={username}
										class="w-full pl-12 pr-4 py-3 bg-gray-900/50 border border-gray-700/50 rounded-xl text-white placeholder-gray-500 focus:border-orange-500/50 focus:ring-2 focus:ring-orange-500/20 transition-all duration-200"
										placeholder="Something you won't forget (hopefully)"
										required
									/>
								</div>
							</div>

							<div class="space-y-1" in:fly={{ y: 20, duration: 400, delay: 300 }}>
								<label class="block text-sm font-medium text-gray-300" for="setup-password">Password</label>
								<div class="relative">
									<div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
										<svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
										</svg>
									</div>
									<input
										type="password"
										id="setup-password"
										bind:value={password}
										class={passwordInputClass}
										placeholder="Not 'password123' please"
										required
									/>
								</div>
								{#if passwordTooShort}
									<p class="text-red-400 text-xs mt-1 flex items-center gap-1" in:fade={{ duration: 200 }}>
										<svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
											<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
										</svg>
										8 characters minimum. Your cat's name doesn't count.
									</p>
								{/if}
							</div>

							<div class="space-y-1" in:fly={{ y: 20, duration: 400, delay: 400 }}>
								<label class="block text-sm font-medium text-gray-300" for="setup-confirm">Confirm Password</label>
								<div class="relative">
									<div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
										<svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
										</svg>
									</div>
									<input
										type="password"
										id="setup-confirm"
										bind:value={confirmPassword}
										class={confirmInputClass}
										placeholder="Prove you can type it twice"
										required
									/>
								</div>
								{#if passwordMismatch}
									<p class="text-red-400 text-xs mt-1 flex items-center gap-1" in:fade={{ duration: 200 }}>
										<svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
											<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
										</svg>
										Passwords don't match. Shocking, I know.
									</p>
								{/if}
							</div>

							{#if $auth.error}
								<div class="bg-red-500/10 border border-red-500/30 rounded-xl p-4" in:scale={{ duration: 200 }}>
									<p class="text-red-400 text-sm flex items-center gap-2">
										<svg class="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
										</svg>
										{$auth.error}
									</p>
								</div>
							{/if}

							<button
								type="submit"
								class="w-full py-3.5 bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white font-semibold rounded-xl transition-all duration-300 transform hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none shadow-lg shadow-orange-500/25"
								disabled={$auth.loading || passwordMismatch || passwordTooShort || !username || !password}
								in:fly={{ y: 20, duration: 400, delay: 500 }}
							>
								{#if $auth.loading}
									<span class="flex items-center justify-center gap-2">
										<svg class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
										</svg>
										Summoning powers...
									</span>
								{:else}
									Crown Me Admin
								{/if}
							</button>
						</form>
					{:else}
						<!-- Login Header -->
						<div class="mb-8 lg:mt-0 mt-16">
							<h2 class="text-3xl font-bold text-white mb-2">Back Already?</h2>
							<p class="text-gray-400">Prove you're worthy of entering the command center.</p>
						</div>

						<!-- Login Form -->
						<form on:submit|preventDefault={handleLogin} class="space-y-5">
							<div class="space-y-1" in:fly={{ y: 20, duration: 400, delay: 200 }}>
								<label class="block text-sm font-medium text-gray-300" for="login-username">Username</label>
								<div class="relative">
									<div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
										<svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
										</svg>
									</div>
									<input
										type="text"
										id="login-username"
										bind:value={username}
										class="w-full pl-12 pr-4 py-3 bg-gray-900/50 border border-gray-700/50 rounded-xl text-white placeholder-gray-500 focus:border-orange-500/50 focus:ring-2 focus:ring-orange-500/20 transition-all duration-200"
										placeholder="State your identity, mortal"
										required
									/>
								</div>
							</div>

							<div class="space-y-1" in:fly={{ y: 20, duration: 400, delay: 300 }}>
								<label class="block text-sm font-medium text-gray-300" for="login-password">Password</label>
								<div class="relative">
									<div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
										<svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
										</svg>
									</div>
									<input
										type="password"
										id="login-password"
										bind:value={password}
										class="w-full pl-12 pr-4 py-3 bg-gray-900/50 border border-gray-700/50 rounded-xl text-white placeholder-gray-500 focus:border-orange-500/50 focus:ring-2 focus:ring-orange-500/20 transition-all duration-200"
										placeholder="Whisper the secret phrase"
										required
									/>
								</div>
							</div>

							{#if $auth.error}
								<div class="bg-red-500/10 border border-red-500/30 rounded-xl p-4" in:scale={{ duration: 200 }}>
									<p class="text-red-400 text-sm flex items-center gap-2">
										<svg class="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
										</svg>
										{$auth.error}
									</p>
								</div>
							{/if}

							<button
								type="submit"
								class="w-full py-3.5 bg-gradient-to-r from-orange-500 to-red-500 hover:from-orange-600 hover:to-red-600 text-white font-semibold rounded-xl transition-all duration-300 transform hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none shadow-lg shadow-orange-500/25"
								disabled={$auth.loading || !username || !password}
								in:fly={{ y: 20, duration: 400, delay: 400 }}
							>
								{#if $auth.loading}
									<span class="flex items-center justify-center gap-2">
										<svg class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
										</svg>
										Verifying your existence...
									</span>
								{:else}
									Grant Me Access
								{/if}
							</button>
						</form>
					{/if}

					<p class="text-center text-gray-600 text-xs mt-8" in:fade={{ duration: 400, delay: 600 }}>
						Powered by mass quantities of caffeine
					</p>
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	@keyframes gradient {
		0%, 100% { background-position: 0% 50%; }
		50% { background-position: 100% 50%; }
	}

	@keyframes float {
		0%, 100% { transform: translateY(0) rotate(0deg); opacity: 0; }
		10% { opacity: 1; }
		90% { opacity: 1; }
		100% { transform: translateY(-100vh) rotate(720deg); opacity: 0; }
	}

	@keyframes pulse-slow {
		0%, 100% { opacity: 0.3; transform: translate(-50%, -50%) scale(1); }
		50% { opacity: 0.5; transform: translate(-50%, -50%) scale(1.1); }
	}

	.animate-gradient {
		background-size: 200% auto;
		animation: gradient 3s ease infinite;
	}

	.animate-pulse-slow {
		animation: pulse-slow 4s ease-in-out infinite;
	}

	.particle {
		position: absolute;
		width: 4px;
		height: 4px;
		background: linear-gradient(to top, #f97316, #ef4444);
		border-radius: 50%;
		animation: float 8s ease-in-out infinite;
	}

	.particle-1 { left: 10%; animation-delay: 0s; animation-duration: 8s; }
	.particle-2 { left: 30%; animation-delay: 1.5s; animation-duration: 10s; }
	.particle-3 { left: 50%; animation-delay: 3s; animation-duration: 7s; }
	.particle-4 { left: 70%; animation-delay: 4.5s; animation-duration: 9s; }
	.particle-5 { left: 90%; animation-delay: 6s; animation-duration: 11s; }

	.rive-container :global(canvas) {
		border-radius: 50%;
	}
</style>
