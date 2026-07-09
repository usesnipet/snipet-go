<script lang="ts" generics="T extends Record<string, unknown>">
  import { Textarea } from "$lib/components/ui/textarea/index.js";
  import {
    Field,
    FieldLabel,
    FieldError,
	FieldDescription,
  } from "$lib/components/ui/field/index.js";
	import type { HTMLAttributes } from "svelte/elements";
	import { formFieldProxy, type FormPathLeaves, type SuperForm } from "sveltekit-superforms";

  type Props = {
    form: SuperForm<T>;
    field: FormPathLeaves<T>;
    label?: string;
    description?: string;
  } & HTMLAttributes<HTMLTextAreaElement>;

  let { form, field, label, description, ...rest }: Props = $props();
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
  <Textarea
    id={field}
    bind:value={$value as string}
    {...$constraints}
    {...rest}
  />
  <FieldError errors={fieldErrors} />
  {#if description}
    <FieldDescription>{description}</FieldDescription>
  {/if}
</Field>