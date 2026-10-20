<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
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
		if (success) {
			// The auth store flips the layout branches mid-setup, which remounts
			// this page and would drop component state. Carrying the signal in
			// the URL survives every remount: the credentials panel renders on
			// the next mount, after navigation completes.
			goto('/login?created=1');
		}
	}

	function finishSetup() {
		goto('/');
	}

	function onPasswordInput(event: Event) {
		password = (event.target as HTMLInputElement).value;
	}

	function onConfirmInput(event: Event) {
		confirmPassword = (event.target as HTMLInputElement).value;
	}

	// Reactive so a same-route navigation (?created=1) flips the view without a
	// remount; the layout branch switch also re-evaluates it.
	$: setupDone = $page.url.searchParams.get('created') === '1';
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
	<!-- Illustration panel: the Agni ember scene -->
	<div
		class="hidden lg:flex relative flex-col justify-between overflow-hidden p-10"
		style="background: linear-gradient(180deg, hsl(16 40% 5%), hsl(14 50% 9%));"
	>
		<LoginScene />

		<!-- Brand -->
		<div class="relative z-10 flex flex-col items-start gap-6 pt-8">
			<div class="flex items-center gap-4">
				<img
					src="/logo.png"
					alt=""
					class="w-12 h-12"
					aria-hidden="true"
					style="filter: drop-shadow(0 0 20px rgb(255 138 66 / 0.55));"
				/>
				<div>
					<h2
						class="text-[30px] leading-none font-semibold tracking-[0.14em] text-transparent bg-clip-text"
						style="background-image: linear-gradient(120deg, #ffe3bd, #ffa14d 55%, #ff6b35);"
					>
						AGNI
					</h2>
					<p
						class="text-[11px] font-medium uppercase tracking-[0.32em] mt-2"
						style="color: rgb(255 190 140 / 0.75);"
					>
						Firecracker microVMs
					</p>
				</div>
			</div>
			<p
				class="max-w-[40ch] text-[15px] leading-relaxed"
				style="color: rgb(255 224 197 / 0.88); text-shadow: 0 1px 14px rgb(10 4 2 / 0.6);"
			>
				Real Linux machines, lit in seconds. Boot microVMs pulled straight from OCI
				registries and drive them from a live serial console &mdash; no images to build,
				no daemons to babysit.
			</p>
		</div>

		<!-- Feature row anchored above the embers -->
		<div class="relative z-10 flex items-end justify-between gap-6 pb-1">
			<ul class="flex flex-col gap-2.5 text-[13px] font-medium">
				<li
					class="flex items-center gap-2.5"
					style="color: rgb(255 214 182 / 0.92); text-shadow: 0 1px 12px rgb(10 4 2 / 0.6);"
				>
					<span style="color: rgb(255 158 92);"><Icon name="ph-fire" size={16} /></span>
					34 distros, one click from their registries
				</li>
				<li
					class="flex items-center gap-2.5"
					style="color: rgb(255 214 182 / 0.92); text-shadow: 0 1px 12px rgb(10 4 2 / 0.6);"
				>
					<span style="color: rgb(255 158 92);"><Icon name="ph-terminal-window" size={16} /></span>
					Interactive serial console with instant echo
				</li>
				<li
					class="flex items-center gap-2.5"
					style="color: rgb(255 214 182 / 0.92); text-shadow: 0 1px 12px rgb(10 4 2 / 0.6);"
				>
					<span style="color: rgb(255 158 92);"><Icon name="ph-lightning" size={16} /></span>
					Firecracker boots in under a second
				</li>
			</ul>
			<span
				class="text-[11px] font-mono tracking-wider shrink-0"
				style="color: rgb(255 190 140 / 0.55);"
			>
				{APP_VERSION}
			</span>
		</div>
	</div>

	<!-- Form side -->
	<div class="flex flex-col items-center justify-center p-6 relative min-h-dvh">
		<div
			class="absolute inset-0 pointer-events-none"
			style="background: radial-gradient(640px 460px at 28% 18%, oklch(var(--status-warn) / 0.06), transparent 70%);"
		></div>
		{#if mounted}
			<div class="w-full max-w-[400px]" in:fade={{ duration: 200 }}>
				<div class="flex items-center justify-center gap-2.5 mb-8">
					<img src="/logo.png" alt="Agni" class="w-[30px] h-[30px]" />
					<span class="text-base font-semibold text-foreground">Agni</span>
				</div>

				<div class="card p-6">
					{#if setupDone}
						<h1 class="text-lg font-semibold text-foreground">Admin account created</h1>
						<p class="text-sm text-muted-foreground mt-1 mb-6">
							Save these credentials - they are the only way to sign in to Agni on this
							machine.
						</p>
						<div class="space-y-3 mb-6">
							<div class="rounded-md border bg-accent p-3">
								<p class="text-[11px] font-medium uppercase tracking-wide text-muted-foreground">Username</p>
								{#if username}
									<p class="font-mono text-sm text-foreground mt-1 break-all">{username}</p>
								{:else}
									<p class="text-sm text-muted-foreground mt-1">the username you just created</p>
								{/if}
							</div>
							<div class="rounded-md border bg-accent p-3">
								<div class="flex items-center justify-between">
									<p class="text-[11px] font-medium uppercase tracking-wide text-muted-foreground">Password</p>
									<button
										type="button"
										on:click={() => (showPassword = !showPassword)}
										class="text-muted-foreground hover:text-foreground inline-flex items-center gap-1 text-[11px]"
										aria-label={showPassword ? 'Hide password' : 'Show password'}
									>
										<Icon name={showPassword ? 'ph-eye-slash' : 'ph-eye'} size={14} />
										{showPassword ? 'Hide' : 'Show'}
									</button>
								</div>
								{#if password}
									<p class="font-mono text-sm text-foreground mt-1 break-all">
										{showPassword ? password : '•'.repeat(Math.max(8, password.length))}
									</p>
								{:else}
									<p class="text-sm text-muted-foreground mt-1">the password you just created</p>
								{/if}
							</div>
						</div>
						<Button variant="primary" fullWidth on:click={finishSetup}>
							Continue to dashboard
						</Button>
					{:else if $auth.setupRequired}
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
