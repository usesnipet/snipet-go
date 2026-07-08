<script lang="ts">
	import FlexTable from "$lib/components/flex-table/flex-table.svelte";
	import type { ColumnDef } from "@tanstack/table-core";
	import type { APIKey } from "../schemas";

  type Props = {
		apiKeys: APIKey[]
		isLoading?: boolean;
		class?: string;
  }

  let props: Props = $props();

  const columns: ColumnDef<APIKey>[] = [
    {
      header: "Name",
      cell: ({ row }) => row.original.name
    },
    {
      header: "Prefix",
      cell: ({ row }) => row.original.key_id
    },
    {
      header: "Expires at",
      cell: ({ row }) => row.original.expires_at
    },
    {
      header: "Status",
      cell: ({ row }) => row.original.active
    },
  ]
</script>

<FlexTable
  {...props}
  data={props.apiKeys}
  {columns}
  emptyMessage="No Api Found."
/>