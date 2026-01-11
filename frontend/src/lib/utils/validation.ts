// Form validation utilities

export interface ValidationResult {
	valid: boolean;
	error?: string;
}

export interface ValidationRule {
	validate: (value: unknown) => ValidationResult;
}

// Basic validators

export function required(fieldName = 'This field'): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			if (value === null || value === undefined || value === '') {
				return { valid: false, error: `${fieldName} is required` };
			}
			return { valid: true };
		}
	};
}

export function minLength(min: number, fieldName = 'This field'): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			if (typeof value !== 'string') {
				return { valid: false, error: `${fieldName} must be a string` };
			}
			if (value.length < min) {
				return { valid: false, error: `${fieldName} must be at least ${min} characters` };
			}
			return { valid: true };
		}
	};
}

export function maxLength(max: number, fieldName = 'This field'): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			if (typeof value !== 'string') {
				return { valid: false, error: `${fieldName} must be a string` };
			}
			if (value.length > max) {
				return { valid: false, error: `${fieldName} must be at most ${max} characters` };
			}
			return { valid: true };
		}
	};
}

export function minValue(min: number, fieldName = 'This field'): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			const num = typeof value === 'string' ? parseFloat(value) : value;
			if (typeof num !== 'number' || isNaN(num)) {
				return { valid: false, error: `${fieldName} must be a number` };
			}
			if (num < min) {
				return { valid: false, error: `${fieldName} must be at least ${min}` };
			}
			return { valid: true };
		}
	};
}

export function maxValue(max: number, fieldName = 'This field'): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			const num = typeof value === 'string' ? parseFloat(value) : value;
			if (typeof num !== 'number' || isNaN(num)) {
				return { valid: false, error: `${fieldName} must be a number` };
			}
			if (num > max) {
				return { valid: false, error: `${fieldName} must be at most ${max}` };
			}
			return { valid: true };
		}
	};
}

export function pattern(regex: RegExp, message: string): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			if (typeof value !== 'string') {
				return { valid: false, error: 'Value must be a string' };
			}
			if (!regex.test(value)) {
				return { valid: false, error: message };
			}
			return { valid: true };
		}
	};
}

// Domain-specific validators

export function vmName(): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			if (typeof value !== 'string') {
				return { valid: false, error: 'VM name must be a string' };
			}
			if (value.length < 1) {
				return { valid: false, error: 'VM name is required' };
			}
			if (value.length > 64) {
				return { valid: false, error: 'VM name must be at most 64 characters' };
			}
			if (!/^[a-zA-Z0-9][a-zA-Z0-9_-]*$/.test(value)) {
				return {
					valid: false,
					error:
						'VM name must start with alphanumeric and contain only letters, numbers, underscores, and hyphens'
				};
			}
			return { valid: true };
		}
	};
}

export function filePath(): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			if (typeof value !== 'string') {
				return { valid: false, error: 'Path must be a string' };
			}
			if (value.length === 0) {
				return { valid: false, error: 'Path is required' };
			}
			if (!value.startsWith('/') && !value.startsWith('./') && !value.startsWith('../')) {
				return {
					valid: false,
					error: 'Path must be absolute or relative (start with /, ./, or ../)'
				};
			}
			return { valid: true };
		}
	};
}

export function cpuCount(): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			const num = typeof value === 'string' ? parseInt(value, 10) : value;
			if (typeof num !== 'number' || isNaN(num)) {
				return { valid: false, error: 'CPU count must be a number' };
			}
			if (num < 1) {
				return { valid: false, error: 'CPU count must be at least 1' };
			}
			if (num > 32) {
				return { valid: false, error: 'CPU count must be at most 32' };
			}
			return { valid: true };
		}
	};
}

export function memoryMB(): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			const num = typeof value === 'string' ? parseInt(value, 10) : value;
			if (typeof num !== 'number' || isNaN(num)) {
				return { valid: false, error: 'Memory must be a number' };
			}
			if (num < 128) {
				return { valid: false, error: 'Memory must be at least 128 MB' };
			}
			if (num > 32768) {
				return { valid: false, error: 'Memory must be at most 32768 MB (32 GB)' };
			}
			return { valid: true };
		}
	};
}

export function username(): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			if (typeof value !== 'string') {
				return { valid: false, error: 'Username must be a string' };
			}
			if (value.length < 3) {
				return { valid: false, error: 'Username must be at least 3 characters' };
			}
			if (value.length > 32) {
				return { valid: false, error: 'Username must be at most 32 characters' };
			}
			if (!/^[a-zA-Z][a-zA-Z0-9_-]*$/.test(value)) {
				return {
					valid: false,
					error:
						'Username must start with a letter and contain only letters, numbers, underscores, and hyphens'
				};
			}
			return { valid: true };
		}
	};
}

export function password(): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			if (typeof value !== 'string') {
				return { valid: false, error: 'Password must be a string' };
			}
			if (value.length < 8) {
				return { valid: false, error: 'Password must be at least 8 characters' };
			}
			if (!/[a-z]/.test(value)) {
				return { valid: false, error: 'Password must contain at least one lowercase letter' };
			}
			if (!/[A-Z]/.test(value)) {
				return { valid: false, error: 'Password must contain at least one uppercase letter' };
			}
			if (!/[0-9]/.test(value)) {
				return { valid: false, error: 'Password must contain at least one number' };
			}
			return { valid: true };
		}
	};
}

export function macAddress(): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			if (typeof value !== 'string') {
				return { valid: false, error: 'MAC address must be a string' };
			}
			if (value === '') {
				return { valid: true }; // Optional field
			}
			if (!/^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$/.test(value)) {
				return { valid: false, error: 'Invalid MAC address format (expected XX:XX:XX:XX:XX:XX)' };
			}
			return { valid: true };
		}
	};
}

export function ipAddress(): ValidationRule {
	return {
		validate: (value: unknown): ValidationResult => {
			if (typeof value !== 'string') {
				return { valid: false, error: 'IP address must be a string' };
			}
			if (value === '') {
				return { valid: true }; // Optional field
			}
			const ipv4Regex = /^(\d{1,3}\.){3}\d{1,3}$/;
			if (!ipv4Regex.test(value)) {
				return { valid: false, error: 'Invalid IP address format' };
			}
			const parts = value.split('.').map(Number);
			if (parts.some((p) => p < 0 || p > 255)) {
				return { valid: false, error: 'IP address octets must be between 0 and 255' };
			}
			return { valid: true };
		}
	};
}

// Composite validator

export function validate(value: unknown, ...rules: ValidationRule[]): ValidationResult {
	for (const rule of rules) {
		const result = rule.validate(value);
		if (!result.valid) {
			return result;
		}
	}
	return { valid: true };
}

// Form validation helper

export interface FormErrors {
	[key: string]: string | undefined;
}

export function validateForm<T extends Record<string, unknown>>(
	values: T,
	rules: { [K in keyof T]?: ValidationRule[] }
): { valid: boolean; errors: FormErrors } {
	const errors: FormErrors = {};
	let valid = true;

	for (const [field, fieldRules] of Object.entries(rules)) {
		if (fieldRules) {
			const result = validate(values[field], ...fieldRules);
			if (!result.valid) {
				errors[field] = result.error;
				valid = false;
			}
		}
	}

	return { valid, errors };
}
