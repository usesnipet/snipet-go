<script lang="ts">
	import PageHead from "$lib/components/page-head.svelte";
  import ApiKeyForm from "$lib/features/api-key/components/api-key-form.svelte";
</script>
<PageHead title="API Key Login" description="Login using an API key to access the API" />
<div class="flex h-screen w-full items-center justify-center px-4">
  <ApiKeyForm />
</div>