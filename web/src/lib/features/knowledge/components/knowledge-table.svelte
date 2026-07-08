<script lang="ts">
	import FlexTable from "$lib/components/flex-table/flex-table.svelte";
	import type { ColumnDef } from "@tanstack/table-core";
	import type { Knowledge } from "../schemas";
	import { renderComponent } from "$lib/components/ui/data-table";
	import TableTimeField from "$lib/components/flex-table/table-time-field.svelte";
	import TableBadgeField from "$lib/components/flex-table/table-badge-field.svelte";
	import TableField from "$lib/components/flex-table/table-field.svelte";

  type Props = {
		knowledges: Knowledge[]
		isLoading?: boolean;
		class?: string;
  }

  let props: Props = $props();

  const columns: ColumnDef<Knowledge>[] = [
    {
      header: "Name",
      cell: ({ row }) => renderComponent(TableField, { value: row.original.name })
    },
    {
      header: "Description",
      cell: ({ row }) => renderComponent(TableField, { value: row.original.description, truncate: 40 })
    },
    {
      header: "Last Synced At",
      cell: ({ row }) => renderComponent(TableTimeField, { date: row.original.last_synced_at })
    },
    {
      header: "Sync Status",
      cell: ({ row }) => renderComponent(TableBadgeField, {
        value: row.original.sync_status,
        variant: row.original.sync_status === "success" ? "default" : "destructive"
      })
    },
  ]
</script>

<FlexTable
  {...props}
  data={props.knowledges}
  {columns}
  emptyMessage="No Knowledge Found."
/>