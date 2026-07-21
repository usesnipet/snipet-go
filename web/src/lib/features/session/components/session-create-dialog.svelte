<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { sessionService } from "../service";
	import { agentService } from "$lib/features/agent/service";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import { createSessionSchema } from "../schemas";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormSelect from "$lib/components/form/form-select.svelte";
	import FormMetadata from "$lib/components/form/form-metadata.svelte";
	import FormError from "$lib/components/form/form-error.svelte";
	import type { Snippet } from "svelte";

	type Props = {
		clientCode: string;
		open?: boolean;
		trigger?: Snippet;
	};

	let { clientCode, open = $bindable(false), trigger }: Props = $props();

	const agentsQuery = agentService.list();
	const agentOptions = $derived(
		(agentsQuery.data ?? []).map((agent) => ({
			label: agent.name,
			value: agent.id,
		})),
	);

	const createSessionMutation = $derived(sessionService.create(clientCode));
	const form = superForm(defaults(zod4(createSessionSchema)), {
		SPA: true,
		dataType: "json",
		validators: zod4Client(createSessionSchema),
		async onUpdate({ form: formState }) {
			if (!formState.valid) return;
			createSessionMutation.mutate(formState.data, {
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
			<Dialog.Title>Create session</Dialog.Title>
			<Dialog.Description>
				Select an agent to create a new session for this client.
			</Dialog.Description>
		</Dialog.Header>

		<form use:enhance class="flex max-h-[70vh] flex-col gap-6">
			<FieldGroup>
				<FormSelect
					{form}
					field="agent_id"
					label="Agent"
					placeholder="Select an agent"
					options={agentOptions}
				/>
				{#key open}
					<FormMetadata
						{form}
						field="metadata"
						label="Metadata"
						description="Optional key-value pairs for this session."
					/>
				{/key}
				<FormError {form} />
				<Dialog.Footer>
					<Button variant="outline" type="button" onclick={handleCancel}>Cancel</Button>
					<Button disabled={createSessionMutation.isPending} type="submit">Create</Button>
				</Dialog.Footer>
			</FieldGroup>
		</form>
	</Dialog.Content>
</Dialog.Root>
