<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { clientService } from "../service";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import { createClientSchema } from "../schemas";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormInput from "$lib/components/form/form-input.svelte";
	import FormError from "$lib/components/form/form-error.svelte";
	import ClientConfigFields from "./client-config-fields.svelte";
	import type { Snippet } from "svelte";

	type Props = {
		open?: boolean;
		trigger?: Snippet;
	};

	let { open = $bindable(false), trigger }: Props = $props();

	const createClientMutation = clientService.create();
	const form = superForm(defaults(zod4(createClientSchema)), {
		SPA: true,
		dataType: "json",
		validators: zod4Client(createClientSchema),
		async onUpdate({ form: formState }) {
			if (!formState.valid) return;
			createClientMutation.mutate(formState.data, {
				onSuccess: () => {
					form.reset();
					open = false;
				},
			});
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

	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>Create client</Dialog.Title>
			<Dialog.Description>Fill the fields below to create a new client.</Dialog.Description>
		</Dialog.Header>

		<form use:enhance class="flex max-h-[70vh] flex-col gap-6">
			<FieldGroup>
				<FormInput {form} field="name" label="Name" placeholder="e.g. Acme Corp" />
				<ClientConfigFields {form} />
				<FormError {form} />
				<Dialog.Footer>
					<Button variant="outline" type="button" onclick={handleCancel}>Cancel</Button>
					<Button disabled={createClientMutation.isPending} type="submit">Create</Button>
				</Dialog.Footer>
			</FieldGroup>
		</form>
	</Dialog.Content>
</Dialog.Root>
