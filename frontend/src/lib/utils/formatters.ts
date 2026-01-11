// Date and time formatting utilities

export function formatDate(date: Date | string): string {
	const d = typeof date === 'string' ? new Date(date) : date;
	return d.toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	});
}

export function formatDateTime(date: Date | string): string {
	const d = typeof date === 'string' ? new Date(date) : date;
	return d.toLocaleString('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric',
		hour: '2-digit',
		minute: '2-digit',
		second: '2-digit'
	});
}

export function formatTime(date: Date | string): string {
	const d = typeof date === 'string' ? new Date(date) : date;
	return d.toLocaleTimeString('en-US', {
		hour: '2-digit',
		minute: '2-digit',
		second: '2-digit'
	});
}

export function formatRelativeTime(date: Date | string): string {
	const d = typeof date === 'string' ? new Date(date) : date;
	const now = new Date();
	const diffMs = now.getTime() - d.getTime();
	const diffSec = Math.floor(diffMs / 1000);
	const diffMin = Math.floor(diffSec / 60);
	const diffHour = Math.floor(diffMin / 60);
	const diffDay = Math.floor(diffHour / 24);

	if (diffSec < 60) return 'just now';
	if (diffMin < 60) return `${diffMin}m ago`;
	if (diffHour < 24) return `${diffHour}h ago`;
	if (diffDay < 7) return `${diffDay}d ago`;
	return formatDate(d);
}

// Bytes formatting utilities

export function formatBytes(bytes: number, decimals = 2): string {
	if (bytes === 0) return '0 B';

	const k = 1024;
	const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
	const i = Math.floor(Math.log(bytes) / Math.log(k));

	return `${parseFloat((bytes / Math.pow(k, i)).toFixed(decimals))} ${sizes[i]}`;
}

export function formatBytesPerSecond(bytes: number): string {
	return `${formatBytes(bytes)}/s`;
}

export function parseBytes(str: string): number {
	const units: Record<string, number> = {
		b: 1,
		kb: 1024,
		mb: 1024 ** 2,
		gb: 1024 ** 3,
		tb: 1024 ** 4
	};

	const match = str.toLowerCase().match(/^(\d+(?:\.\d+)?)\s*(b|kb|mb|gb|tb)?$/);
	if (!match) return 0;

	const value = parseFloat(match[1]);
	const unit = match[2] || 'b';
	return value * (units[unit] || 1);
}

// Duration formatting utilities

export function formatDuration(seconds: number): string {
	if (seconds < 60) return `${seconds}s`;

	const hours = Math.floor(seconds / 3600);
	const minutes = Math.floor((seconds % 3600) / 60);
	const secs = seconds % 60;

	if (hours > 0) {
		return `${hours}h ${minutes}m ${secs}s`;
	}
	return `${minutes}m ${secs}s`;
}

export function formatUptime(startTime: Date | string): string {
	const start = typeof startTime === 'string' ? new Date(startTime) : startTime;
	const now = new Date();
	const diffSec = Math.floor((now.getTime() - start.getTime()) / 1000);
	return formatDuration(diffSec);
}

// Number formatting utilities

export function formatNumber(num: number): string {
	return num.toLocaleString('en-US');
}

export function formatPercent(value: number, decimals = 1): string {
	return `${value.toFixed(decimals)}%`;
}

export function formatCPU(cores: number): string {
	return cores === 1 ? '1 CPU' : `${cores} CPUs`;
}

export function formatMemory(mb: number): string {
	if (mb >= 1024) {
		return `${(mb / 1024).toFixed(1)} GB`;
	}
	return `${mb} MB`;
}

// Path formatting utilities

export function getFilename(path: string): string {
	return path.split('/').pop() || path;
}

export function getDirname(path: string): string {
	const parts = path.split('/');
	parts.pop();
	return parts.join('/') || '/';
}

export function truncatePath(path: string, maxLength = 30): string {
	if (path.length <= maxLength) return path;
	const filename = getFilename(path);
	if (filename.length >= maxLength - 3) {
		return '...' + filename.slice(-(maxLength - 3));
	}
	return '...' + path.slice(-(maxLength - 3));
}

// VM ID formatting

export function formatVMId(id: string, length = 8): string {
	return id.slice(0, length);
}
