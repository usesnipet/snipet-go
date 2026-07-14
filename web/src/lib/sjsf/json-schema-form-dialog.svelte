<script lang="ts" generics="T extends Record<string, unknown> = Record<string, unknown>">
	import * as defaults from './defaults';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import {
		Content,
		Form,
		createForm,
		setFormContext,
		type FormOptions
	} from '@sjsf/form';
	import type { Snippet } from 'svelte';
	import ScrollArea from '$lib/components/ui/scroll-area/scroll-area.svelte';

	type JsonSchemaFormConfig = Pick<
		FormOptions<T>,
		'schema' | 'uiSchema' | 'initialValue' | 'fieldsValidationMode'
	>;
	let {
		title,
		description,
		formConfig,
		onSubmit,
		trigger,
		open = $bindable(false)
	}: {
		title: string;
		description?: string;
		formConfig: JsonSchemaFormConfig;
		onSubmit?: (value: T) => void;
		trigger?: Snippet;
		open?: boolean;
	} = $props();

	const form = $derived.by(() => createForm<T>({
		...defaults,
		...formConfig,
		onSubmit(value) {
			onSubmit?.(value);
			open = false;
			form.reset();
		}
	}));

	$effect(() => {
		setFormContext(form);
	});

	function handleCancel() {
		form.reset();
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}>
		{#if trigger}
			{@render trigger()}
		{:else}
			Open form
		{/if}
	</Dialog.Trigger>

	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>{title}</Dialog.Title>
			{#if description}
				<Dialog.Description>{description}</Dialog.Description>
			{/if}
		</Dialog.Header>

		<Form attributes={{ novalidate: true }}>
		  <ScrollArea class="max-h-[80vh] pr-4">
				<Content />
			</ScrollArea>



			<Dialog.Footer>
				<Button variant="outline" type="button" onclick={handleCancel}>Cancel</Button>
				<Button type="submit">Save</Button>
			</Dialog.Footer>
		</Form>
	</Dialog.Content>
</Dialog.Root>
