export type JsonValue =
	| string
	| number
	| boolean
	| null
	| JsonValue[]
	| { [key: string]: JsonValue };

export function isPlainObject(value: unknown): value is Record<string, JsonValue> {
	return typeof value === "object" && value !== null && !Array.isArray(value);
}

export function isSimplePrimitive(value: unknown): value is string | number | boolean {
	const type = typeof value;
	return type === "string" || type === "number" || type === "boolean";
}

export function formatKey(key: string): string {
	return key
		.replace(/_/g, " ")
		.replace(/([a-z])([A-Z])/g, "$1 $2")
		.replace(/^./, (char) => char.toUpperCase());
}

export function summarizeValue(value: JsonValue): string {
	if (value === null) return "Empty";
	if (typeof value === "boolean") return value ? "Yes" : "No";
	if (typeof value === "number") return String(value);
	if (typeof value === "string") return value.length > 48 ? `${value.slice(0, 48)}…` : value;
	if (Array.isArray(value)) return `${value.length} ${value.length === 1 ? "item" : "items"}`;
	if (isPlainObject(value)) {
		const count = Object.keys(value).length;
		return `${count} ${count === 1 ? "property" : "properties"}`;
	}
	return "";
}

export function shouldDefaultOpen(value: JsonValue, depth: number): boolean {
	if (depth === 0) {
		if (isPlainObject(value)) return Object.keys(value).length <= 6;
		if (Array.isArray(value)) {
			return value.length <= 4 && value.every(isSimplePrimitive);
		}
		return true;
	}

	if (isPlainObject(value)) return Object.keys(value).length <= 3;
	if (Array.isArray(value)) {
		return value.length <= 2 && value.every(isSimplePrimitive);
	}

	return true;
}

export function isComplexArray(value: JsonValue[]): boolean {
	if (value.length > 4) return true;
	return value.some((item) => isPlainObject(item) || Array.isArray(item));
}
