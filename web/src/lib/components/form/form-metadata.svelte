<script lang="ts" generics="T extends Record<string, unknown>">
	import { Input } from "$lib/components/ui/input/index.js";
	import {
		Field,
		FieldLabel,
		FieldError,
		FieldDescription,
	} from "$lib/components/ui/field/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import TrashIcon from "@lucide/svelte/icons/trash";
	import { formFieldProxy, type FormPath, type FormPathLeaves, type SuperForm } from "sveltekit-superforms";

	type Entry = {
		id: string;
		key: string;
		value: string;
	};

	type Props = {
		form: SuperForm<T>;
		field: FormPath<T>;
		label?: string;
		description?: string;
	};

	let { form, field, label, description }: Props = $props();

	const { value, errors } = $derived.by(() =>
		formFieldProxy(form, field as FormPathLeaves<T>),
	);
	const fieldErrors = $derived.by(() => {
		const list = $errors;
		if (!Array.isArray(list)) return undefined;
		return list.map((message) => ({ message }));
	});

	let entries = $state<Entry[]>([]);
	let lastSerialized = $state("");

	function serializeRecord(record: Record<string, string> | null | undefined): string {
		return JSON.stringify(record ?? {});
	}

	function recordFromEntries(list: Entry[]): Record<string, string> {
		const result: Record<string, string> = {};
		for (const entry of list) {
			const key = entry.key.trim();
			if (!key) continue;
			result[key] = entry.value;
		}
		return result;
	}

	function entriesFromRecord(record: Record<string, string> | null | undefined): Entry[] {
		return Object.entries(record ?? {}).map(([key, val]) => ({
			id: crypto.randomUUID(),
			key,
			value: String(val ?? ""),
		}));
	}

	function commit(nextEntries: Entry[]) {
		const next = recordFromEntries(nextEntries);
		lastSerialized = serializeRecord(next);
		$value = next as typeof $value;
	}

	function addEntry() {
		entries = [...entries, { id: crypto.randomUUID(), key: "", value: "" }];
	}

	function removeEntry(id: string) {
		const next = entries.filter((entry) => entry.id !== id);
		entries = next;
		commit(next);
	}

	function updateEntry(id: string, patch: Partial<Pick<Entry, "key" | "value">>) {
		const next = entries.map((entry) => (entry.id === id ? { ...entry, ...patch } : entry));
		entries = next;
		commit(next);
	}

	$effect(() => {
		const current = $value as Record<string, string> | null | undefined;
		const serialized = serializeRecord(current);
		if (serialized === lastSerialized) return;
		lastSerialized = serialized;
		entries = entriesFromRecord(current);
	});
</script>

<Field>
	<div class="flex items-center justify-between gap-2">
		<FieldLabel for={field}>{label}</FieldLabel>
		<Button type="button" variant="outline" size="sm" onclick={addEntry}>
			<PlusIcon />
			Add
		</Button>
	</div>

	<div class="space-y-2">
		{#each entries as entry (entry.id)}
			<div class="flex items-center gap-2">
				<Input
					placeholder="Key"
					value={entry.key}
					oninput={(e) => updateEntry(entry.id, { key: e.currentTarget.value })}
				/>
				<Input
					placeholder="Value"
					value={entry.value}
					oninput={(e) => updateEntry(entry.id, { value: e.currentTarget.value })}
				/>
				<Button
					type="button"
					variant="ghost"
					size="icon"
					onclick={() => removeEntry(entry.id)}
					aria-label="Remove metadata entry"
				>
					<TrashIcon />
				</Button>
			</div>
		{:else}
			<p class="text-muted-foreground text-sm">No metadata entries.</p>
		{/each}
	</div>

	<FieldError errors={fieldErrors} />
	{#if description}
		<FieldDescription>{description}</FieldDescription>
	{/if}
</Field>
