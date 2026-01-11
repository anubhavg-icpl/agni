const API_BASE = '/api';

interface ApiError {
	error: string;
}

class ApiClient {
	private token: string | null = null;

	setToken(token: string | null) {
		this.token = token;
		if (token) {
			localStorage.setItem('auth_token', token);
		} else {
			localStorage.removeItem('auth_token');
		}
	}

	getToken(): string | null {
		if (!this.token) {
			this.token = localStorage.getItem('auth_token');
		}
		return this.token;
	}

	private async request<T>(method: string, path: string, body?: unknown): Promise<T> {
		const headers: HeadersInit = {
			'Content-Type': 'application/json'
		};

		const token = this.getToken();
		if (token) {
			headers['Authorization'] = `Bearer ${token}`;
		}

		const response = await fetch(`${API_BASE}${path}`, {
			method,
			headers,
			body: body ? JSON.stringify(body) : undefined
		});

		if (!response.ok) {
			const error: ApiError = await response.json().catch(() => ({
				error: `HTTP ${response.status}`
			}));
			throw new Error(error.error);
		}

		return response.json();
	}

	// Auth
	async getAuthStatus(): Promise<{ setup_required: boolean }> {
		return this.request('GET', '/auth/status');
	}

	async setup(username: string, password: string): Promise<{ success: boolean; user: User }> {
		return this.request('POST', '/auth/setup', { username, password });
	}

	async login(username: string, password: string): Promise<LoginResponse> {
		const response = await this.request<LoginResponse>('POST', '/auth/login', {
			username,
			password
		});
		this.setToken(response.token);
		return response;
	}

	async logout(): Promise<void> {
		await this.request('POST', '/auth/logout');
		this.setToken(null);
	}

	async getMe(): Promise<User> {
		return this.request('GET', '/auth/me');
	}

	// VMs
	async listVMs(): Promise<VM[]> {
		return this.request('GET', '/vms');
	}

	async getVM(id: string): Promise<VM> {
		return this.request('GET', `/vms/${id}`);
	}

	async createVM(data: CreateVMRequest): Promise<VM> {
		return this.request('POST', '/vms', data);
	}

	async deleteVM(id: string): Promise<void> {
		await this.request('DELETE', `/vms/${id}`);
	}

	async startVM(id: string): Promise<VMActionResponse> {
		return this.request('POST', `/vms/${id}/start`);
	}

	async stopVM(id: string): Promise<VMActionResponse> {
		return this.request('POST', `/vms/${id}/stop`);
	}

	async shutdownVM(id: string): Promise<VMActionResponse> {
		return this.request('POST', `/vms/${id}/shutdown`);
	}

	async getVMMetrics(id: string): Promise<VMMetrics> {
		return this.request('GET', `/vms/${id}/metrics`);
	}

	// Configs
	async listConfigs(): Promise<ConfigTemplate[]> {
		return this.request('GET', '/configs');
	}

	async getConfig(id: string): Promise<ConfigTemplate> {
		return this.request('GET', `/configs/${id}`);
	}

	async createConfig(data: CreateConfigRequest): Promise<ConfigTemplate> {
		return this.request('POST', '/configs', data);
	}

	async deleteConfig(id: string): Promise<void> {
		await this.request('DELETE', `/configs/${id}`);
	}

	// Health
	async getHealth(): Promise<HealthStatus> {
		return this.request('GET', '/health');
	}

	async getSystemInfo(): Promise<SystemInfo> {
		return this.request('GET', '/system/info');
	}
}

export const api = new ApiClient();

// Types
export interface User {
	id: string;
	username: string;
	role: string;
	created_at: string;
	last_login_at?: string;
}

export interface LoginResponse {
	token: string;
	expires_at: string;
	user: User;
}

export interface VM {
	id: string;
	name: string;
	status: 'stopped' | 'starting' | 'running' | 'stopping' | 'error';
	config: VMConfig;
	metrics?: VMMetrics;
	error?: string;
	created_at: string;
	started_at?: string;
	stopped_at?: string;
}

export interface VMConfig {
	name: string;
	kernel_path: string;
	kernel_opts: string;
	initrd_path?: string;
	root_drive: Drive;
	additional_drives?: Drive[];
	cpus: number;
	memory_mb: number;
	cpu_template?: string;
	disable_smt: boolean;
	network_interfaces?: NIC[];
	vsock_devices?: Vsock[];
	metadata?: string;
	log_level: string;
	jailer?: JailerConfig;
}

export interface JailerConfig {
	chroot_base: string;
	uid: number;
	gid: number;
	daemonize: boolean;
	exec_file: string;
	numa_node?: number;
}

export interface Drive {
	path: string;
	read_only: boolean;
	part_uuid?: string;
}

export interface NIC {
	device: string;
	mac_address: string;
}

export interface Vsock {
	path: string;
	cid: number;
}

export interface VMMetrics {
	cpu_usage: number;
	memory_used: number;
	memory_total: number;
	disk_read_bytes: number;
	disk_write_bytes: number;
	net_rx_bytes: number;
	net_tx_bytes: number;
	timestamp: string;
}

export interface VMActionResponse {
	success: boolean;
	message: string;
	vm_id: string;
}

export interface CreateVMRequest {
	name: string;
	config: VMConfig;
}

export interface ConfigTemplate {
	id: string;
	name: string;
	description?: string;
	config: VMConfig;
	created_at: string;
	updated_at: string;
}

export interface CreateConfigRequest {
	name: string;
	description?: string;
	config: VMConfig;
}

export interface HealthStatus {
	status: string;
	version: string;
	uptime: string;
	timestamp: string;
	components: Record<string, { status: string; message?: string }>;
}

export interface SystemInfo {
	version: string;
	firecracker_version: string;
	go_version: string;
	os: string;
	arch: string;
	num_cpu: number;
	total_memory: number;
}
