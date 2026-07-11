<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { knowledgeService } from "../service";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import { updateKnowledgeSchema, type Knowledge } from "../schemas";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormInput from "$lib/components/form/form-input.svelte";
	import FormTextarea from "$lib/components/form/form-textarea.svelte";
	import FormError from "$lib/components/form/form-error.svelte";
	import type { Snippet } from "svelte";

	type Props = {
		open: boolean;
		knowledge?: Knowledge;
		trigger?: Snippet;
	}

	let { open = $bindable(false), knowledge, trigger }: Props = $props();

	const updateKnowledgeMutation = knowledgeService.update();
	const form = superForm(
		defaults(zod4(updateKnowledgeSchema)),
		{
			SPA: true,
			validators: zod4Client(updateKnowledgeSchema),
			async onUpdate({ form }) {
				if (!form.valid) return;
				if (!knowledge) return;
				updateKnowledgeMutation.mutate({
					data: form.data,
					id: knowledge.id,
				});
				open = false;
			},
		},
	);

	$effect(() => {
		if (knowledge) {
			form.form.set({
				name: knowledge.name,
				description: knowledge.description
			});
		}
	});
	const { enhance } = form;
</script>

<Dialog.Root bind:open>
	{#if trigger}
		<Dialog.Trigger class={buttonVariants({ variant: "outline" })}>
			{@render trigger()}
		</Dialog.Trigger>
	{/if}

	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Update knowledge</Dialog.Title>
			<Dialog.Description>Fill the fields below to update the knowledge source.</Dialog.Description>
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
				<FormError {form} />
				<Dialog.Footer>
					<Button variant="outline" type="button" onclick={() => form.reset()}>
						Cancel
					</Button>
					<Button disabled={updateKnowledgeMutation.isPending} type="submit">
						Update
					</Button>
				</Dialog.Footer>
			</FieldGroup>
		</form>
	</Dialog.Content>
</Dialog.Root>
