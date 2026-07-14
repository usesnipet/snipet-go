<script lang="ts">
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import * as Dialog from "$lib/components/ui/dialog";
	import { FieldGroup } from "$lib/components/ui/field";
	import FormError from "$lib/components/form/form-error.svelte";
	import FormInput from "$lib/components/form/form-input.svelte";
	import FormExpiration from "$lib/components/form/form-expiration.svelte";
	import { DEFAULT_EXPIRATION_OPTIONS, fromExpirationFormValue } from "$lib/components/form/expiration";
	import { apiKeyService } from "../service";
	import {
		createApiKeyFormSchema,
		type APIKeyWithSecret,
	} from "../schemas";
	import { defaults, superForm } from "sveltekit-superforms";
	import { zod4, zod4Client } from "sveltekit-superforms/adapters";
	import type { Snippet } from "svelte";
	import CopyIcon from "@lucide/svelte/icons/copy";
	import CheckIcon from "@lucide/svelte/icons/check";
	import { toast } from "svelte-sonner";

	type Props = {
		open?: boolean;
		trigger?: Snippet;
	};

	let { open = $bindable(false), trigger }: Props = $props();

	let createdKey = $state<APIKeyWithSecret | null>(null);
	let copied = $state(false);

	const createApiKeyMutation = apiKeyService.create();
	const form = superForm(defaults(zod4(createApiKeyFormSchema)), {
		SPA: true,
		dataType: "json",
		validators: zod4Client(createApiKeyFormSchema),
		async onUpdate({ form: formState }) {
			if (!formState.valid) return;
			createApiKeyMutation.mutate(
				{
					name: formState.data.name,
					expires_at: fromExpirationFormValue(formState.data.expires_at),
				},
				{
					onSuccess: (data) => {
						createdKey = data;
						form.reset();
					},
				},
			);
		},
	});
	const { enhance } = form;

	function handleClose() {
		form.reset();
		createdKey = null;
		copied = false;
		open = false;
	}

	async function handleCopy() {
		if (!createdKey?.key) return;
		try {
			await navigator.clipboard.writeText(createdKey.key);
			copied = true;
			toast.success("API key copied.");
			setTimeout(() => {
				copied = false;
			}, 2000);
		} catch {
			toast.error("Failed to copy API key.");
		}
	}
</script>

<Dialog.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (!isOpen) {
			form.reset();
			createdKey = null;
			copied = false;
		}
	}}
>
	{#if trigger}
		<Dialog.Trigger class={buttonVariants({ variant: "outline" })}>
			{@render trigger()}
		</Dialog.Trigger>
	{/if}

	<Dialog.Content>
		{#if createdKey}
			<Dialog.Header>
				<Dialog.Title>API key created</Dialog.Title>
				<Dialog.Description>
					Copy your API key now. You will not be able to see it again.
				</Dialog.Description>
			</Dialog.Header>

			<div class="flex flex-col gap-3">
				<div class="bg-muted flex items-center gap-2 rounded-xl p-3">
					<code class="min-w-0 flex-1 truncate text-sm">{createdKey.key}</code>
					<Button variant="outline" size="icon-sm" onclick={handleCopy} aria-label="Copy API key">
						{#if copied}
							<CheckIcon />
						{:else}
							<CopyIcon />
						{/if}
					</Button>
				</div>
			</div>

			<Dialog.Footer>
				<Button type="button" onclick={handleClose}>Done</Button>
			</Dialog.Footer>
		{:else}
			<Dialog.Header>
				<Dialog.Title>Create API key</Dialog.Title>
				<Dialog.Description>Fill the fields below to create a new API key.</Dialog.Description>
			</Dialog.Header>

			<form use:enhance>
				<FieldGroup>
					<FormInput {form} field="name" label="Name" placeholder="e.g. Production" />
					<FormExpiration
						{form}
						field="expires_at"
						label="Expires at"
						options={DEFAULT_EXPIRATION_OPTIONS}
						allowCustom
						description="Choose a preset or pick a custom date."
					/>
					<FormError {form} />
					<Dialog.Footer>
						<Button variant="outline" type="button" onclick={handleClose}>Cancel</Button>
						<Button disabled={createApiKeyMutation.isPending} type="submit">Create</Button>
					</Dialog.Footer>
				</FieldGroup>
			</form>
		{/if}
	</Dialog.Content>
</Dialog.Root>
