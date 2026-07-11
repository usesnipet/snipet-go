<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { knowledgeService } from "../service";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import { createKnowledgeIndexSchema, type CreateKnowledgeIndex, type Knowledge } from "../schemas";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormInput from "$lib/components/form/form-input.svelte";
	import FormError from "$lib/components/form/form-error.svelte";
	import type { Snippet } from "svelte";
	import { ON_CHANGE, ON_INPUT, type FormOptions, type Schema } from "@sjsf/form";
	import FormSelect from "$lib/components/form/form-select.svelte";
	import JsonSchemaFormDialog from "$lib/sjsf/json-schema-form-dialog.svelte";

	type Props = {
		open: boolean;
		knowledge: Knowledge;
		trigger?: Snippet;
	}

	let { open = $bindable(false), knowledge, trigger }: Props = $props();
	const indexDriversQuery = knowledgeService.listIndexDrivers();

	const createKnowledgeIndexMutation = knowledgeService.createIndex();
	const form = superForm(
		defaults(zod4(createKnowledgeIndexSchema)),
		{
			SPA: true,
			validators: zod4Client(createKnowledgeIndexSchema),
			async onUpdate({ form }) {
				if (!form.valid) return;
				createKnowledgeIndexMutation.mutate({
					data: form.data,
					knowledgeId: knowledge.id,
				});
				open = false;
			},
		},
	);

	const { enhance, form: formStore } = form;

	type DriverFormConfig = Pick<
		FormOptions<Record<string, unknown>>,
		"schema" | "uiSchema" | "initialValue" | "fieldsValidationMode"
	>;

	const drivers = $derived(indexDriversQuery.data?.index_drivers ?? []);
	const selectedDriver = $derived(drivers.find((driver) => driver.name === $formStore.driver));
	const driverFormConfig = $derived.by((): DriverFormConfig | null => {
		if (!selectedDriver) return null;

		const configuration = $formStore.configuration;
		const initialValue =
			typeof configuration === "object" && configuration !== null && !Array.isArray(configuration)
				? configuration
				: {};

		return {
			schema: selectedDriver.configuration_schema as Schema,
			uiSchema: {},
			initialValue,
			fieldsValidationMode: ON_INPUT | ON_CHANGE,
		};
	});

	let configureDialogKey = $state(0);

	function handleConfigurationSave(value: CreateKnowledgeIndex["configuration"]) {
		formStore.update((data) => ({ ...data, configuration: value }));
		configureDialogKey++;
	}
</script>

<Dialog.Root bind:open>
	{#if trigger}
		<Dialog.Trigger class={buttonVariants({ variant: "outline" })}>
			{@render trigger()}
		</Dialog.Trigger>
	{/if}

	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Create index for {knowledge.name}</Dialog.Title>
			<Dialog.Description>Fill the fields below to create a new index.</Dialog.Description>
		</Dialog.Header>

		<form use:enhance>
			<FieldGroup>
				<FormInput
				  {form}
					field="name"
					label="Name"
					placeholder="e.g. Product docs"
				/>
				<FormSelect
					{form}
					field="driver"
					label="Driver"
					options={
						(indexDriversQuery.data?.index_drivers ?? []).map(driver => ({
							label: driver.name,
							value: driver.name,
						}))
					}
				/>
				{#if driverFormConfig}
				{#key `${$formStore.driver}-${configureDialogKey}`}
					<JsonSchemaFormDialog
						title="Configure Driver"
						description="Configure the driver to use for this index."
						formConfig={driverFormConfig}
						onSubmit={(value) =>
							handleConfigurationSave(value as CreateKnowledgeIndex["configuration"])}
					>
						{#snippet trigger()}
							Configure
						{/snippet}
					</JsonSchemaFormDialog>
				{/key}
			{:else}
				<Button variant="outline" disabled>Configure</Button>
			{/if}
				<FormError {form} />
				<Dialog.Footer>
					<Button variant="outline" type="button" onclick={() => form.reset()}>
						Cancel
					</Button>
					<Button disabled={createKnowledgeIndexMutation.isPending} type="submit">
						Create
					</Button>
				</Dialog.Footer>
			</FieldGroup>
		</form>
	</Dialog.Content>
</Dialog.Root>
