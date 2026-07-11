<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { knowledgeService } from "../service";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import { updateKnowledgeIndexSchema, type Knowledge, type KnowledgeIndex } from "../schemas";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormInput from "$lib/components/form/form-input.svelte";
	import FormError from "$lib/components/form/form-error.svelte";
	import type { Snippet } from "svelte";

	type Props = {
		open: boolean;
		knowledge?: Knowledge;
		trigger?: Snippet;
		index?: KnowledgeIndex;
	}

	let { open = $bindable(false), knowledge, trigger, index }: Props = $props();

	const updateKnowledgeIndexMutation = knowledgeService.updateIndex();
	const form = superForm(
		defaults(zod4(updateKnowledgeIndexSchema)),
		{
			SPA: true,
			validators: zod4Client(updateKnowledgeIndexSchema),
			async onUpdate({ form }) {
				if (!form.valid) return;
				if (!knowledge) return;
				if (!index) return;
				updateKnowledgeIndexMutation.mutate({
					id: index.id,
					knowledgeId: knowledge.id,
					data: form.data,
				});
				open = false;
			},
		},
	);

	$effect(() => {
		if (index) {
			form.form.set({
				name: index.name,
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
			<Dialog.Title>Update index for {knowledge?.name}</Dialog.Title>
			<Dialog.Description>Fill the fields below to update the index.</Dialog.Description>
		</Dialog.Header>

		<form use:enhance>
			<FieldGroup>
				<FormInput
				  {form}
					field="name"
					label="Name"
					placeholder="e.g. Product docs"
				/>
				<FormError {form} />
				<Dialog.Footer>
					<Button variant="outline" type="button" onclick={() => form.reset()}>
						Cancel
					</Button>
					<Button disabled={updateKnowledgeIndexMutation.isPending} type="submit">
						Update
					</Button>
				</Dialog.Footer>
			</FieldGroup>
		</form>
	</Dialog.Content>
</Dialog.Root>
