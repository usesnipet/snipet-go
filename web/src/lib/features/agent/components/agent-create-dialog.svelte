<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { agentService } from "../service";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import { createAgentSchema } from "../schemas";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormInput from "$lib/components/form/form-input.svelte";
	import FormTextarea from "$lib/components/form/form-textarea.svelte";
	import FormError from "$lib/components/form/form-error.svelte";
	import type { Snippet } from "svelte";

	type Props = {
		open?: boolean;
		trigger?: Snippet;
	};

	let { open = $bindable(false), trigger }: Props = $props();

	const createAgentMutation = agentService.create();
	const form = superForm(defaults(zod4(createAgentSchema)), {
		SPA: true,
		dataType: "json",
		validators: zod4Client(createAgentSchema),
		async onUpdate({ form: formState }) {
			if (!formState.valid) return;
			createAgentMutation.mutate(
				{
					name: formState.data.name,
					description: formState.data.description ?? "",
					configuration: { llms: [] },
				},
				{
					onSuccess: () => {
						form.reset();
						open = false;
					},
				},
			);
		},
	});
	const { enhance } = form;

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
			<Dialog.Title>Create agent</Dialog.Title>
			<Dialog.Description>Fill the fields below to create a new agent.</Dialog.Description>
		</Dialog.Header>

		<form use:enhance>
			<FieldGroup>
				<FormInput {form} field="name" label="Name" placeholder="e.g. Support agent" />
				<FormTextarea {form} field="description" label="Description" placeholder="Optional" />
				<FormError {form} />
				<Dialog.Footer>
					<Button variant="outline" type="button" onclick={handleCancel}>Cancel</Button>
					<Button disabled={createAgentMutation.isPending} type="submit">Create</Button>
				</Dialog.Footer>
			</FieldGroup>
		</form>
	</Dialog.Content>
</Dialog.Root>
