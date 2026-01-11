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
				update((s) => ({ ...s, loading: false, error: (e as Error).message }));
			}
		},
		async create(name: string, config: VM['config']) {
			update((s) => ({ ...s, loading: true, error: null }));
			try {
				const vm = await api.createVM({ name, config });
				update((s) => ({ ...s, vms: [...s.vms, vm], loading: false }));
				return vm;
			} catch (e) {
				update((s) => ({ ...s, loading: false, error: (e as Error).message }));
				return null;
			}
		},
		async start(id: string) {
			try {
				await api.startVM(id);
				await this.fetch();
				return true;
			} catch (e) {
				update((s) => ({ ...s, error: (e as Error).message }));
				return false;
			}
		},
		async stop(id: string) {
			try {
				await api.stopVM(id);
				await this.fetch();
				return true;
			} catch (e) {
				update((s) => ({ ...s, error: (e as Error).message }));
				return false;
			}
		},
		async shutdown(id: string) {
			try {
				await api.shutdownVM(id);
				await this.fetch();
				return true;
			} catch (e) {
				update((s) => ({ ...s, error: (e as Error).message }));
				return false;
			}
		},
		async delete(id: string) {
			try {
				await api.deleteVM(id);
				update((s) => ({ ...s, vms: s.vms.filter((vm) => vm.id !== id) }));
				return true;
			} catch (e) {
				update((s) => ({ ...s, error: (e as Error).message }));
				return false;
			}
		},
		clearError() {
			update((s) => ({ ...s, error: null }));
		}
	};
}

export const vms = createVMStore();
