import { writable } from 'svelte/store';
import { api, type VM } from '$lib/api/client';

interface VMState {
	vms: VM[];
	loading: boolean;
	error: string | null;
}

function createVMStore() {
	const { subscribe, set, update } = writable<VMState>({
		vms: [],
		loading: false,
		error: null
	});

	return {
		subscribe,
		async fetch() {
			update((s) => ({ ...s, loading: true, error: null }));
			try {
				const vms = await api.listVMs();
				set({ vms: vms || [], loading: false, error: null });
			} catch (e) {
				let msg = (e as Error).message;
				if (msg.includes('Failed to fetch')) msg = "The zoo is closed. (Network error)";
				update((s) => ({ ...s, loading: false, error: msg }));
			}
		},
		async create(name: string, config: VM['config']) {
			update((s) => ({ ...s, loading: true, error: null }));
			try {
				const vm = await api.createVM({ name, config });
				update((s) => ({ ...s, vms: [...s.vms, vm], loading: false }));
				return vm;
			} catch (e) {
				let msg = (e as Error).message;
				if (msg.includes('already exists')) msg = "We already have a demon with that name.";
				update((s) => ({ ...s, loading: false, error: msg }));
				return null;
			}
		},
		async start(id: string) {
			try {
				await api.startVM(id);
				await this.fetch();
				return true;
			} catch (e) {
				update((s) => ({ ...s, error: "It refuses to wake up. (" + (e as Error).message + ")" }));
				return false;
			}
		},
		async stop(id: string) {
			try {
				await api.stopVM(id);
				await this.fetch();
				return true;
			} catch (e) {
				update((s) => ({ ...s, error: "It won't die. (" + (e as Error).message + ")" }));
				return false;
			}
		},
		async shutdown(id: string) {
			try {
				await api.shutdownVM(id);
				await this.fetch();
				return true;
			} catch (e) {
				update((s) => ({ ...s, error: "Clean kill failed. (" + (e as Error).message + ")" }));
				return false;
			}
		},
		async delete(id: string) {
			try {
				await api.deleteVM(id);
				update((s) => ({ ...s, vms: s.vms.filter((vm) => vm.id !== id) }));
				return true;
			} catch (e) {
				update((s) => ({ ...s, error: "Exorcism failed. (" + (e as Error).message + ")" }));
				return false;
			}
		},
		clearError() {
			update((s) => ({ ...s, error: null }));
		}
	};
}

export const vms = createVMStore();
