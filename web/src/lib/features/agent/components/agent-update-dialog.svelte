<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { agentService } from "../service";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import { updateAgentSchema, type Agent } from "../schemas";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormInput from "$lib/components/form/form-input.svelte";
	import FormTextarea from "$lib/components/form/form-textarea.svelte";
	import FormError from "$lib/components/form/form-error.svelte";
	import type { Snippet } from "svelte";

	type Props = {
		open: boolean;
		agent?: Agent;
		trigger?: Snippet;
	};

	let { open = $bindable(false), agent, trigger }: Props = $props();

	const updateAgentMutation = agentService.update();
	const form = superForm(defaults(zod4(updateAgentSchema)), {
		SPA: true,
		validators: zod4Client(updateAgentSchema),
		async onUpdate({ form }) {
			if (!form.valid) return;
			if (!agent) return;
			updateAgentMutation.mutate(
				{
					data: form.data,
					id: agent.id,
				},
				{
					onSuccess: () => {
						open = false;
					},
				},
			);
		},
	});

	$effect(() => {
		if (agent) {
			form.form.set({
				name: agent.name,
				description: agent.description,
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
			<Dialog.Title>Update agent</Dialog.Title>
			<Dialog.Description>Fill the fields below to update the agent.</Dialog.Description>
		</Dialog.Header>

		<form use:enhance>
			<FieldGroup>
				<FormInput {form} field="name" label="Name" placeholder="e.g. Support agent" />
				<FormTextarea {form} field="description" label="Description" placeholder="Optional" />
				<FormError {form} />
				<Dialog.Footer>
					<Button variant="outline" type="button" onclick={() => form.reset()}>
						Cancel
					</Button>
					<Button disabled={updateAgentMutation.isPending} type="submit">Update</Button>
				</Dialog.Footer>
			</FieldGroup>
		</form>
	</Dialog.Content>
</Dialog.Root>
