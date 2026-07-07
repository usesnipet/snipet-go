<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Card from "$lib/components/ui/card/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import {
    FieldGroup,
    Field,
    FieldLabel,
    FieldError,
  } from "$lib/components/ui/field/index.js";
  import { defaults, setError, superForm } from "sveltekit-superforms";
  import { zod4, zod4Client } from "sveltekit-superforms/adapters";
  import { api } from "$lib/api";
  import { ApiError } from "$lib/http";
	import { login } from "$lib/store/auth.svelte";
	import { goto } from "$app/navigation";
	import { resolve } from "$app/paths";

  const { form, constraints, enhance, errors } = superForm(
    defaults(zod4(api.apiKey.me.schema)),
    {
      SPA: true,
      validators: zod4Client(api.apiKey.me.schema),
      async onUpdate({ form }) {
        if (!form.valid) {
          return;
        }

        try {
          await api.apiKey.me.run(form.data.apiKey);
          login(form.data.apiKey);
          return goto(resolve("/"));
        } catch (error) {
          if (ApiError.is(error)) {
            setError(form, "", error.message);
            return;
          }

          setError(form, "", "Something went wrong. Please try again.");
        }
      },
    },
  );
</script>

<Card.Root class="mx-auto w-full max-w-sm">
  <Card.Header>
    <Card.Title class="text-2xl">API Key</Card.Title>
    <Card.Description>Enter your API key below to get started</Card.Description>
  </Card.Header>
  <Card.Content>
    <form use:enhance>
      <FieldGroup>
        <Field>
          <div class="flex items-center">
            <FieldLabel for="api-key">API Key</FieldLabel>
          </div>
          <Input id="api-key" type="password" bind:value={$form.apiKey} {...$constraints.apiKey} />
          <FieldError errors={$errors.apiKey?.map((message) => ({ message }))} />
        </Field>
        <FieldError errors={$errors._errors?.map((message) => ({ message }))} />
        <Field>
          <Button type="submit" class="w-full">Submit</Button>
        </Field>
      </FieldGroup>
    </form>
  </Card.Content>
</Card.Root>