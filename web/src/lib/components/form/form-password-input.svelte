<script lang="ts" generics="T extends Record<string, unknown>">
  import { Input } from "$lib/components/ui/input/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import {
    Field,
    FieldLabel,
    FieldError,
    FieldDescription,
  } from "$lib/components/ui/field/index.js";
  import EyeIcon from "@lucide/svelte/icons/eye";
  import EyeOffIcon from "@lucide/svelte/icons/eye-off";
  import type { HTMLAttributes } from "svelte/elements";
  import { formFieldProxy, type FormPathLeaves, type SuperForm } from "sveltekit-superforms";

  type Props = {
    form: SuperForm<T>;
    field: FormPathLeaves<T>;
    label?: string;
    description?: string;
    disabled?: boolean;
  } & Omit<HTMLAttributes<HTMLInputElement>, "type">;

  let { form, field, label, description, disabled, ...rest }: Props = $props();
  const { value, constraints, errors } = $derived.by(() => formFieldProxy(form, field));
  const fieldErrors = $derived.by(() => {
    const list = $errors;
    if (!Array.isArray(list)) return undefined;
    return list.map((message) => ({ message: message }));
  });

  let visible = $state(false);
  const inputType = $derived(visible ? "text" : "password");
</script>

<Field>
  <div class="flex items-center">
    <FieldLabel for={field}>{label}</FieldLabel>
  </div>
  <div class="relative">
    <Input
      id={field}
      type={inputType}
      class="pr-10"
      bind:value={$value}
      {disabled}
      {...$constraints}
      {...rest}
    />
    <Button
      type="button"
      variant="ghost"
      size="icon-sm"
      class="absolute top-1/2 right-1 -translate-y-1/2"
      {disabled}
      onclick={() => (visible = !visible)}
      aria-label={visible ? "Hide password" : "Show password"}
    >
      {#if visible}
        <EyeOffIcon />
      {:else}
        <EyeIcon />
      {/if}
    </Button>
  </div>
  <FieldError errors={fieldErrors} />
  {#if description}
    <FieldDescription>{description}</FieldDescription>
  {/if}
</Field>
