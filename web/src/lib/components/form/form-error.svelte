<script lang="ts" generics="T extends Record<string, unknown>">
  import type { SuperForm } from "sveltekit-superforms";
	import { FieldError } from "../ui/field";

  type Props = {
    form: SuperForm<T>;
  }

  let { form }: Props = $props();
  const errors = $derived(form.errors);
</script>

<FieldError errors={$errors._errors?.map((message) => ({ message }))} />
