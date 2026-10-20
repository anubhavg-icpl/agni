import type { VM } from '$lib/api/client';

export type StatusTone = 'ok' | 'warn' | 'error' | 'idle';

export interface VMStatusInfo {
	label: string;
	tone: StatusTone;
	spinner: boolean;
	pulse: boolean;
}

/**
 * One resolver for the VM status vocabulary (Law 17). In-flight states take
 * the tone of where they are going; only the spinner says "not yet".
 */
export function vmStatus(status: VM['status'] | string): VMStatusInfo {
	switch (status) {
		case 'running':
			return { label: 'Running', tone: 'ok', spinner: false, pulse: false };
		case 'starting':
			return { label: 'Starting', tone: 'ok', spinner: true, pulse: false };
		case 'stopping':
			return { label: 'Stopping', tone: 'idle', spinner: true, pulse: false };
		case 'stopped':
			return { label: 'Stopped', tone: 'idle', spinner: false, pulse: false };
		case 'error':
			return { label: 'Error', tone: 'error', spinner: false, pulse: true };
		default:
			return { label: 'Unknown', tone: 'idle', spinner: false, pulse: false };
	}
}
