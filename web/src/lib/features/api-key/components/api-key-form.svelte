<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Card from "$lib/components/ui/card/index.js";
  import { FieldGroup } from "$lib/components/ui/field/index.js";
  import { defaults, setError, superForm } from "sveltekit-superforms";
  import { zod4, zod4Client } from "sveltekit-superforms/adapters";
  import { ApiError } from "$lib/http";
	import { login } from "$lib/features/api-key/stores/api-key-auth.svelte";
	import { goto } from "$app/navigation";
	import { resolve } from "$app/paths";
  import { checkApiKeySchema } from "../schemas";
	import { apiKeyService } from "../service";
	import FormError from "$lib/components/form/form-error.svelte";
	import FormPasswordInput from "$lib/components/form/form-password-input.svelte";

  const check = apiKeyService.check();
  const form = superForm(
    defaults(zod4(checkApiKeySchema)),
    {
      SPA: true,
      validators: zod4Client(checkApiKeySchema),
      async onUpdate({ form }) {
        if (!form.valid) return;

        check.mutate(form.data.apiKey, {
          onSuccess: () => {
            login(form.data.apiKey);
            return goto(resolve("/"));
          },
          onError: (error) => {
            if (ApiError.is(error)) {
              setError(form, "", error.message);
              return;
            }

            setError(form, "", "Something went wrong. Please try again.");
          }
        });
      },
    },
  );
  const { enhance } = form;
</script>

<Card.Root class="mx-auto w-full max-w-sm">
  <Card.Header>
    <Card.Title class="text-2xl">API Key</Card.Title>
    <Card.Description>Enter your API key below to get started</Card.Description>
  </Card.Header>
  <Card.Content>
    <form use:enhance>
      <FieldGroup>
        <FormPasswordInput
          {form}
          field="apiKey"
          label="API Key"
          description="Enter your API key below to get started"
        />
        <FormError {form} />
        <Button type="submit" class="w-full">Login</Button>
      </FieldGroup>
    </form>
  </Card.Content>
</Card.Root>