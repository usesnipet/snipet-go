<script lang="ts" generics="T extends Record<string, unknown>">
	import { Switch } from "$lib/components/ui/switch/index.js";
	import {
		Field,
		FieldLabel,
		FieldError,
		FieldDescription,
		FieldContent,
	} from "$lib/components/ui/field/index.js";
	import { formFieldProxy, type FormPathLeaves, type SuperForm } from "sveltekit-superforms";

	type Props = {
		form: SuperForm<T>;
		field: FormPathLeaves<T>;
		label?: string;
		description?: string;
	};

	let { form, field, label, description }: Props = $props();

	const { value, errors } = $derived.by(() => formFieldProxy(form, field));
	const fieldErrors = $derived.by(() => {
		const list = $errors;
		if (!Array.isArray(list)) return undefined;
		return list.map((message) => ({ message }));
	});
</script>

<Field orientation="horizontal">
	<FieldContent>
		<FieldLabel for={field}>{label}</FieldLabel>
		{#if description}
			<FieldDescription>{description}</FieldDescription>
		{/if}
		<FieldError errors={fieldErrors} />
	</FieldContent>
	<Switch id={field} bind:checked={$value as boolean} />
</Field>
