<script lang="ts" generics="T extends Record<string, unknown>">
  import { Input } from "$lib/components/ui/input/index.js";
  import {
    Field,
    FieldLabel,
    FieldError,
	FieldDescription,
  } from "$lib/components/ui/field/index.js";
	import type { HTMLAttributes, HTMLInputTypeAttribute } from "svelte/elements";
	import { formFieldProxy, type FormPathLeaves, type SuperForm } from "sveltekit-superforms";

  type Props = {
    form: SuperForm<T>;
    field: FormPathLeaves<T>;
    label?: string;
    description?: string;
    type?: Exclude<HTMLInputTypeAttribute, "file">;
  } & Omit<HTMLAttributes<HTMLInputElement>, "type">;

  let { form, field, label, description, type, ...rest }: Props = $props();
  const { value, constraints, errors } = $derived.by(() => formFieldProxy(form, field));
  const fieldErrors = $derived.by(() => {
    const list = $errors;
    if (!Array.isArray(list)) return undefined;
    return list.map((message) => ({ message: message }));
  });
</script>

<Field>
  <div class="flex items-center">
    <FieldLabel for={field}>{label}</FieldLabel>
  </div>
  <Input
    id={field}
    {type}
    bind:value={$value}
    {...$constraints}
    {...rest}
  />
  <FieldError errors={fieldErrors} />
  {#if description}
    <FieldDescription>{description}</FieldDescription>
  {/if}
</Field>