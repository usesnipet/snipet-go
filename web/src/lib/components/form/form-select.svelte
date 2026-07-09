<script lang="ts" generics="T extends Record<string, unknown>">
	import * as Select from "$lib/components/ui/select/index.js";
	import {
		Field,
		FieldLabel,
		FieldError,
		FieldDescription,
	} from "$lib/components/ui/field/index.js";
	import { formFieldProxy, type FormPathLeaves, type SuperForm } from "sveltekit-superforms";

	type Option = {
		label: string;
		value: string;
		disabled?: boolean;
	};

	type Props = {
		form: SuperForm<T>;
		field: FormPathLeaves<T>;
		label?: string;
		description?: string;
		placeholder?: string;
		options: Option[];
	};

	let {
		form,
		field,
		label,
		description,
		placeholder = "Select...",
		options,
	}: Props = $props();

	const { value, constraints, errors } = $derived.by(() => formFieldProxy(form, field));
	const fieldErrors = $derived.by(() => {
		const list = $errors;
		if (!Array.isArray(list)) return undefined;
		return list.map((message) => ({ message }));
	});
	const triggerContent = $derived(
		options.find((option) => option.value === $value)?.label ?? placeholder,
	);
</script>

<Field>
	<div class="flex items-center">
		<FieldLabel for={field}>{label}</FieldLabel>
	</div>
	<Select.Root type="single" name={field} bind:value={$value as string}>
		<Select.Trigger class="w-full" id={field} {...$constraints}>
			{triggerContent}
		</Select.Trigger>
		<Select.Content>
			<Select.Group>
				{#each options as option (option.value)}
					<Select.Item
						value={option.value}
						label={option.label}
						disabled={option.disabled}
					>
						{option.label}
					</Select.Item>
				{/each}
			</Select.Group>
		</Select.Content>
	</Select.Root>
	<FieldError errors={fieldErrors} />
	{#if description}
		<FieldDescription>{description}</FieldDescription>
	{/if}
</Field>
