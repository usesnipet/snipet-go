<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { knowledgeService } from "../service";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import { createKnowledgeSchema } from "../schemas";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormInput from "$lib/components/form/form-input.svelte";
	import FormTextarea from "$lib/components/form/form-textarea.svelte";
	import FormError from "$lib/components/form/form-error.svelte";

	const driversQuery = knowledgeService.listDrivers();

	const createKnowledgeMutation = knowledgeService.create();
	const form = superForm(
    defaults(zod4(createKnowledgeSchema)),
    {
      SPA: true,
      validators: zod4Client(createKnowledgeSchema),
      async onUpdate({ form }) {
        if (!form.valid) return;
        createKnowledgeMutation.mutate(form.data);
      },
    },
  );
  const { enhance } = form;
</script>

<Dialog.Root>
  <Dialog.Trigger class={buttonVariants({ variant: "outline" })}>
		<PlusIcon />
		Create knowledge
	</Dialog.Trigger>

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
				<FormInput
				  {form}
					field="driver"
					label="Driver"
					placeholder="e.g. s3, notion, github"
				/>
				<FormTextarea
				  {form}
					field="configuration"
					label="Configuration (JSON)"
					placeholder='Type your configuration here...'
				/>
				<FormError {form} />
				<Dialog.Footer>
					<Button variant="outline" type="button" onclick={() => form.reset()}>
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
