import { writable, get } from 'svelte/store';
import { api, type DeployTemplate, type ImageJob, type TemplatesResponse } from '$lib/api/client';

interface TemplateState {
	templates: DeployTemplate[];
	runtime: TemplatesResponse['runtime'];
	kvmAvailable: boolean;
	loading: boolean;
	error: string | null;
}

function createTemplateStore() {
	const { subscribe, set, update } = writable<TemplateState>({
		templates: [],
		runtime: { installed: false, version: 'v1.12.1' },
		kvmAvailable: true,
		loading: false,
		error: null
	});

	async function fetch() {
		update((s) => ({ ...s, loading: true, error: null }));
		try {
			const res = await api.listTemplates();
			set({
				templates: res.templates || [],
				runtime: res.runtime,
				kvmAvailable: res.kvm_available,
				loading: false,
				error: null
			});
		} catch (e) {
			update((s) => ({
				...s,
				loading: false,
				error: `Failed to load templates. ${(e as Error).message}`
			}));
		}
	}

	/** Starts (or joins) the download job for a template. */
	async function prepare(templateId: string): Promise<ImageJob> {
		const res = await api.prepareTemplate(templateId);
		return res.job;
	}

	/** Polls a download job until it finishes; returns final state. */
	async function waitForJob(jobId: string, onUpdate?: (job: ImageJob) => void): Promise<ImageJob> {
		for (;;) {
			const res = await api.getImageJob(jobId);
			onUpdate?.(res.job);
			if (res.job.state !== 'running') return res.job;
			await new Promise((r) => setTimeout(r, 800));
		}
	}

	/**
	 * Single-click deploy: downloads missing artifacts (with progress),
	 * then creates and optionally starts the VM.
	 */
	async function deploy(
		template: DeployTemplate,
		options: { name: string; cpus: number; memoryMb: number; autoStart: boolean },
		onProgress?: (job: ImageJob | null) => void
	): Promise<{ ok: boolean; error?: string; startError?: string }> {
		try {
			if (!template.ready) {
				const job = await prepare(template.id);
				const final = await waitForJob(job.id, onProgress);
				if (final.state === 'error') {
					return { ok: false, error: final.error || 'Download failed.' };
				}
			}
			const res = await api.deployFromTemplate({
				template_id: template.id,
				name: options.name,
				cpus: options.cpus,
				memory_mb: options.memoryMb,
				auto_start: options.autoStart
			});
			await fetch();
			return res.started || !options.autoStart
				? { ok: true }
				: { ok: true, startError: res.start_error };
		} catch (e) {
			return { ok: false, error: (e as Error).message };
		}
	}

	/** Installs the firecracker runtime, tracking progress. */
	async function installRuntime(
		onProgress?: (job: ImageJob | null) => void
	): Promise<{ ok: boolean; error?: string }> {
		try {
			const job = await prepare('firecracker-runtime');
			const final = await waitForJob(job.id, onProgress);
			await fetch();
			return final.state === 'error'
				? { ok: false, error: final.error || 'Installation failed.' }
				: { ok: true };
		} catch (e) {
			return { ok: false, error: (e as Error).message };
		}
	}

	function templateById(id: string): DeployTemplate | undefined {
		return get({ subscribe }).templates.find((t) => t.id === id);
	}

	return {
		subscribe,
		fetch,
		prepare,
		waitForJob,
		deploy,
		installRuntime,
		templateById,
		clearError: () => update((s) => ({ ...s, error: null }))
	};
}

export const templates = createTemplateStore();
