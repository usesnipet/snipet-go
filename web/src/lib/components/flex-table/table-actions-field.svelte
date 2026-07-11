<script lang="ts">
	import { Button, type ButtonProps } from "../ui/button";
  import type { LucideIcon } from '@lucide/svelte';

  type Props = {
    actions: {
      key: string;
      icon?: LucideIcon;
      label?: string;
      variant?: ButtonProps["variant"];
      onClick: () => void;
    }[];
  }

  let { actions }: Props = $props();
</script>

<div class="flex gap-1">
  {#each actions as action (action.key)}
    {@const Icon = action.icon}
    <Button
      onclick={action.onClick}
      variant={action.variant}
      size={action.label ? "default" : "icon"}
    >
      <Icon /> {action.label}
    </Button>
  {/each}
</div>