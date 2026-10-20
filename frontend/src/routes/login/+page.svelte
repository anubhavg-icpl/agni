<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { fade } from 'svelte/transition';
	import { auth } from '$lib/stores/auth';
	import Icon from '$lib/components/Common/Icon.svelte';
	import Button from '$lib/components/Common/Button.svelte';
	import LoginScene from '$lib/components/Auth/LoginScene.svelte';
	import { APP_VERSION } from '$lib/version';

	let username = '';
	let password = '';
	let confirmPassword = '';
	let showPassword = false;
	let mounted = false;

	onMount(() => {
		mounted = true;
	});

	async function handleLogin() {
		const success = await auth.login(username, password);
		if (success) goto('/');
	}

	async function handleSetup() {
		if (password !== confirmPassword || password.length < 8) return;
		const success = await auth.setup(username, password);
		if (success) goto('/');
	}

	function onPasswordInput(event: Event) {
		password = (event.target as HTMLInputElement).value;
	}

	function onConfirmInput(event: Event) {
		confirmPassword = (event.target as HTMLInputElement).value;
	}

	$: passwordMismatch = confirmPassword.length > 0 && password !== confirmPassword;
	$: passwordTooShort = password.length > 0 && password.length < 8;
	$: setupInvalid =
		$auth.setupRequired &&
		(passwordMismatch || passwordTooShort || !username || !password || !confirmPassword);
</script>

<svelte:head>
	<title>{$auth.setupRequired ? 'Create admin account' : 'Sign in'} · Agni</title>
</svelte:head>

<div class="min-h-dvh grid lg:grid-cols-2">
	<!-- Illustration panel: the Agni fire over a warm ground -->
	<div
		class="hidden lg:flex relative items-center justify-center overflow-hidden"
		style="background: oklch(var(--status-warn) / 0.06);"
	>
		<LoginScene />
		<div class="relative z-10 flex flex-col items-center gap-4">
			<img src="/logo.png" alt="" class="w-12 h-12" aria-hidden="true" />
			<div class="text-center">
				<p class="text-2xl font-semibold text-foreground">Agni</p>
				<p class="text-sm text-muted-foreground mt-1">Firecracker microVM manager</p>
			</div>
		</div>
	</div>

	<!-- Form side -->
	<div class="flex flex-col items-center justify-center p-6 relative min-h-dvh">
		{#if mounted}
			<div class="w-full max-w-[400px]" in:fade={{ duration: 200 }}>
				<div class="flex items-center justify-center gap-2.5 mb-8">
					<img src="/logo.png" alt="Agni" class="w-[30px] h-[30px]" />
					<span class="text-base font-semibold text-foreground">Agni</span>
				</div>

				<div class="card p-6">
					{#if $auth.setupRequired}
						<h1 class="text-lg font-semibold text-foreground">Create the admin account</h1>
						<p class="text-sm text-muted-foreground mt-1 mb-6">
							First run of Agni on this machine. Create an administrator to continue.
						</p>

						{#if $auth.error}
							<div
								class="mb-4 rounded-md border border-destructive/30 bg-destructive/10 p-3 flex items-start gap-2"
								role="alert"
							>
								<span class="text-destructive mt-0.5"><Icon name="ph-x-circle" size={16} /></span>
								<p class="text-sm text-destructive">{$auth.error}</p>
							</div>
						{/if}

						<form on:submit|preventDefault={handleSetup} class="space-y-4">
							<div>
								<label class="label" for="setup-username">Username</label>
								<input
									id="setup-username"
									type="text"
									bind:value={username}
									class="input h-11"
									autocomplete="username"
									required
								/>
							</div>
							<div>
								<label class="label" for="setup-password">Password</label>
								<div class="relative">
									<input
										id="setup-password"
										type={showPassword ? 'text' : 'password'}
										value={password}
										on:input={onPasswordInput}
										class="input h-11 pr-10 {passwordTooShort
											? 'border-destructive focus:ring-destructive/30'
											: ''}"
										autocomplete="new-password"
										minlength="8"
										required
									/>
									<button
										type="button"
										on:click={() => (showPassword = !showPassword)}
										class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
										aria-label={showPassword ? 'Hide password' : 'Show password'}
									>
										<Icon name={showPassword ? 'ph-eye-slash' : 'ph-eye'} size={16} />
									</button>
								</div>
								{#if passwordTooShort}
									<p class="text-[13px] text-destructive mt-1.5">
										Password must be at least 8 characters.
									</p>
								{/if}
							</div>
							<div>
								<label class="label" for="setup-confirm">Confirm password</label>
								<input
									id="setup-confirm"
									type={showPassword ? 'text' : 'password'}
									value={confirmPassword}
									on:input={onConfirmInput}
									class="input h-11 {passwordMismatch
										? 'border-destructive focus:ring-destructive/30'
										: ''}"
									autocomplete="new-password"
									required
								/>
								{#if passwordMismatch}
									<p class="text-[13px] text-destructive mt-1.5">Passwords do not match.</p>
								{/if}
							</div>
							<Button
								type="submit"
								variant="primary"
								fullWidth
								loading={$auth.loading}
								disabled={setupInvalid}
							>
								Create admin account
							</Button>
						</form>
					{:else}
						<h1 class="text-lg font-semibold text-foreground">Sign in</h1>
						<p class="text-sm text-muted-foreground mt-1 mb-6">
							Use your Agni administrator account.
						</p>

						{#if $auth.error}
							<div
								class="mb-4 rounded-md border border-destructive/30 bg-destructive/10 p-3 flex items-start gap-2"
								role="alert"
							>
								<span class="text-destructive mt-0.5"><Icon name="ph-x-circle" size={16} /></span>
								<p class="text-sm text-destructive">{$auth.error}</p>
							</div>
						{/if}

						<form on:submit|preventDefault={handleLogin} class="space-y-4">
							<div>
								<label class="label" for="login-username">Username</label>
								<input
									id="login-username"
									type="text"
									bind:value={username}
									class="input h-11"
									autocomplete="username"
									required
								/>
							</div>
							<div>
								<label class="label" for="login-password">Password</label>
								<div class="relative">
									<input
										id="login-password"
										type={showPassword ? 'text' : 'password'}
										value={password}
										on:input={onPasswordInput}
										class="input h-11 pr-10"
										autocomplete="current-password"
										required
									/>
									<button
										type="button"
										on:click={() => (showPassword = !showPassword)}
										class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
										aria-label={showPassword ? 'Hide password' : 'Show password'}
									>
										<Icon name={showPassword ? 'ph-eye-slash' : 'ph-eye'} size={16} />
									</button>
								</div>
							</div>
							<Button
								type="submit"
								variant="primary"
								fullWidth
								loading={$auth.loading}
								disabled={!username || !password}
							>
								Sign in
							</Button>
						</form>
					{/if}
				</div>
			</div>
		{/if}

		<p class="absolute bottom-4 left-4 font-mono text-xs text-muted-foreground">{APP_VERSION}</p>
	</div>
</div>
