<script lang="ts">
	import * as defaults from '$lib/sjsf/defaults';
	import * as post from '$lib/post';
	import { BasicForm, createForm } from '@sjsf/form';

	const form = createForm<post.Model>({ ...defaults, ...post, onSubmit: console.log });
</script>

<BasicForm {form} novalidate style="padding: 2rem;" />
