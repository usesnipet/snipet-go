/** Relative expiration expression: `never`, or `+{n}{d|w|m|y}` (e.g. `+7d`, `+1y`). */
export type ExpirationExpression = string;

export type ExpirationOption = {
	label: string;
	expression: ExpirationExpression;
};

export const CUSTOM_EXPIRATION_VALUE = "__custom__";

export const DEFAULT_EXPIRATION_OPTIONS: ExpirationOption[] = [
	{ label: "Never", expression: "never" },
	{ label: "7 days", expression: "+7d" },
	{ label: "30 days", expression: "+30d" },
	{ label: "90 days", expression: "+90d" },
	{ label: "1 year", expression: "+1y" },
];

/**
 * Resolves a relative expiration expression to a Date (local), or null for never.
 * Supported: `never`, `+Nd`, `+Nw`, `+Nm`, `+Ny`.
 */
export function resolveExpirationExpression(
	expression: ExpirationExpression,
	from: Date = new Date(),
): Date | null {
	const trimmed = expression.trim().toLowerCase();
	if (!trimmed || trimmed === "never") return null;

	const match = /^(\+)(\d+)([dwmy])$/.exec(trimmed);
	if (!match) {
		throw new Error(`Invalid expiration expression: ${expression}`);
	}

	const amount = Number(match[2]);
	const unit = match[3];
	const date = new Date(from.getTime());

	switch (unit) {
		case "d":
			date.setDate(date.getDate() + amount);
			break;
		case "w":
			date.setDate(date.getDate() + amount * 7);
			break;
		case "m":
			date.setMonth(date.getMonth() + amount);
			break;
		case "y":
			date.setFullYear(date.getFullYear() + amount);
			break;
		default:
			throw new Error(`Invalid expiration expression: ${expression}`);
	}

	return date;
}

/** Serialize a Date for form storage (ISO), or empty string for never. */
export function toExpirationFormValue(date: Date | null | undefined): string {
	if (!date) return "";
	return date.toISOString();
}

/** Parse a form expiration string to Date, or null when empty/invalid. */
export function fromExpirationFormValue(value: string | undefined): Date | null {
	if (!value?.trim()) return null;
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return null;
	return date;
}
