<script lang="ts">
	import { resolve } from '$app/paths';
	import JsonSchemaFormDialog from '$lib/sjsf/json-schema-form-dialog.svelte';
	import * as post from '$lib/post';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { toast } from 'svelte-sonner';
</script>

<div class="flex flex-col gap-4 p-8">
	<a class="text-primary underline" href={resolve('/demo/sjsf')}>sjsf (full page)</a>

	<JsonSchemaFormDialog
		title="Create post"
		description="Fill in the fields below to create a new post."
		formConfig={post}
		onSubmit={(value) => {
			console.log(value);
			toast.success('Post saved');
		}}
	>
		{#snippet trigger()}
			<PlusIcon />
			Create post
		{/snippet}
	</JsonSchemaFormDialog>
</div>
