<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { knowledgeService } from "../service";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import { createKnowledgeSchema, type CreateKnowledge } from "../schemas";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormInput from "$lib/components/form/form-input.svelte";
	import FormTextarea from "$lib/components/form/form-textarea.svelte";
	import FormError from "$lib/components/form/form-error.svelte";
	import FormSelect from "$lib/components/form/form-select.svelte";
	import JsonSchemaFormDialog from "$lib/sjsf/json-schema-form-dialog.svelte";
	import { ON_INPUT, ON_CHANGE } from "@sjsf/form";
	import type { FormOptions, Schema } from "@sjsf/form";
	import type { Snippet } from "svelte";

	type Props = {
		open?: boolean;
		trigger?: Snippet;
	}
	let { open = $bindable(false), trigger }: Props = $props();

	const driversQuery = knowledgeService.listDrivers();

	const createKnowledgeMutation = knowledgeService.create();
	const form = superForm(
		defaults({ configuration: {} }, zod4(createKnowledgeSchema)),
		{
			SPA: true,
			dataType: "json",
			validators: zod4Client(createKnowledgeSchema),
			onChange(event) {
				if (event.paths.includes("driver")) {
					event.set("configuration", {});
				}
			},
			async onUpdate({ form }) {
				if (!form.valid) return;
				createKnowledgeMutation.mutate(form.data);
			},
		},
	);
	const { form: formStore, enhance } = form;

	type DriverFormConfig = Pick<
		FormOptions<Record<string, unknown>>,
		"schema" | "uiSchema" | "initialValue" | "fieldsValidationMode"
	>;

	const drivers = $derived(driversQuery.data?.source_drivers ?? []);
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

	function handleConfigurationSave(value: CreateKnowledge["configuration"]) {
		formStore.update((data) => ({ ...data, configuration: value }));
		configureDialogKey++;
	}

	function handleCancel() {
		form.reset();
		open = false;
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
			<Dialog.Title>Create knowledge</Dialog.Title>
			<Dialog.Description>Fill the fields below to create a new knowledge source.</Dialog.Description>
		</Dialog.Header>

		<form use:enhance>
			<FieldGroup>
				<FormInput
				  {form}
					field="name"
					label="Name"
					placeholder="e.g. Product docs"
				/>
				<FormTextarea
				  {form}
					field="description"
					label="Description"
					placeholder="Optional"
				/>
				<FormSelect
				  {form}
					field="driver"
					label="Driver"
					options={
						(driversQuery.data?.source_drivers ?? []).map(driver => ({
							label: driver.name,
							value: driver.name,
						}))
					}
				/>
				{#if driverFormConfig}
					{#key `${$formStore.driver}-${configureDialogKey}`}
						<JsonSchemaFormDialog
							title="Configure Driver"
							description="Configure the driver to use for this knowledge source."
							formConfig={driverFormConfig}
							onSubmit={(value) =>
								handleConfigurationSave(value as CreateKnowledge["configuration"])}
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
					<Button variant="outline" type="button" onclick={handleCancel}>
						Cancel
					</Button>
					<Button disabled={createKnowledgeMutation.isPending} type="submit">
						Create
					</Button>
				</Dialog.Footer>
			</FieldGroup>
		</form>
	</Dialog.Content>
</Dialog.Root>
