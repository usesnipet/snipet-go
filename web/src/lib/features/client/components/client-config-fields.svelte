<script lang="ts">
	import {
		FieldGroup,
		FieldSet,
		FieldSeparator,
	} from "$lib/components/ui/field";
	import FormInput from "$lib/components/form/form-input.svelte";
	import FormSwitch from "$lib/components/form/form-switch.svelte";
	import type { SuperForm } from "sveltekit-superforms";
	import type { ClientConfig } from "../schemas";

	type ClientFormData = {
		name: string;
		config: ClientConfig;
	};

	type Props = {
		form: SuperForm<ClientFormData>;
	};

	let { form }: Props = $props();
	const formData = $derived(form.form);
</script>

<FieldSet>
	<FieldGroup>
		<FormSwitch
			{form}
			field="config.oidc.enabled"
			label="Enabled OIDC"
			description="Enable OpenID Connect authentication."
		/>
		{#if $formData.config.oidc.enabled}
			<FormInput
				{form}
				field="config.oidc.issuer"
				label="Issuer"
				placeholder="https://auth.example.com"
			/>
			<FormInput
				{form}
				field="config.oidc.audience"
				label="Audience"
				placeholder="https://api.example.com"
			/>
		{/if}
	</FieldGroup>
</FieldSet>

<FieldSeparator />

<FieldSet>
	<FieldGroup>
		<FormSwitch
			{form}
			field="config.webhook.enabled"
			label="Enabled Webhook"
			description="Enable webhook notifications."
		/>
		{#if $formData.config.webhook.enabled}
			<FormInput
				{form}
				field="config.webhook.url"
				label="URL"
				type="url"
				placeholder="https://hooks.example.com/snipet"
			/>
		{/if}
	</FieldGroup>
</FieldSet>
