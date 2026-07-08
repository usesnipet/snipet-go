<script lang="ts">
	import Time from "svelte-time/Time.svelte";
  import type { ConfigType, OptionType } from "dayjs";

  type Props = {
    date: ConfigType;
    format?: OptionType;
  }

  let props: Props = $props();
</script>

<Time timestamp={props.date} format={props.format} />