<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import type { Snippet } from 'svelte';
	import { buttonVariants } from './ui/button';
	import Button from './ui/button/button.svelte';

  type Props = {
    title: string;
    description?: string;
    danger?: boolean;
    onConfirm: () => void;
    onCancel?: () => void;
    trigger?: Snippet;
    open?: boolean;

  }

  let {
    title,
    description,
    onConfirm,
    onCancel,
    trigger,
    open = $bindable(false),
    danger
  }: Props = $props();


  function handleCancel() {
    onCancel?.();
		open = false;
	}
  function handleConfirm() {
    onConfirm();
    open = false;
  }
</script>

<Dialog.Root bind:open>
  {#if trigger}
    <Dialog.Trigger class={buttonVariants({ variant: 'outline' })}>
      {@render trigger()}
    </Dialog.Trigger>
  {/if}

	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>{title}</Dialog.Title>
			{#if description}
				<Dialog.Description>{description}</Dialog.Description>
			{/if}
		</Dialog.Header>

    <div class="flex justify-end gap-2">
      <Button onclick={handleConfirm} variant={danger ? "destructive" : "default"}>
        Confirm
      </Button>
      <Button onclick={handleCancel} variant={danger ? "default" : "outline"}>
        Cancel
      </Button>
    </div>
	</Dialog.Content>
</Dialog.Root>
