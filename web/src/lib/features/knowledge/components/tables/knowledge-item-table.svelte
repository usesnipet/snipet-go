<script lang="ts">
	import FlexTable from "$lib/components/flex-table/flex-table.svelte";
	import type { ColumnDef } from "@tanstack/table-core";
	import type { KnowledgeItem } from "../../schemas";
	import { renderComponent } from "$lib/components/ui/data-table";
	import TableField from "$lib/components/flex-table/table-field.svelte";
	import TableExternalId from "./table-external-id.svelte";
	import TableTimeField from "$lib/components/flex-table/table-time-field.svelte";
	import TableBadgeField from "$lib/components/flex-table/table-badge-field.svelte";

  type Props = {
		items: KnowledgeItem[]
		isLoading?: boolean;
		class?: string;
  }

  let props: Props = $props();

  const columns: ColumnDef<KnowledgeItem>[] = [
    {
      header: "Name",
      cell: ({ row }) => renderComponent(TableField, { value: row.original.name, truncate: 30 })
    },
    {
      header: "External ID",
      cell: ({ row }) => renderComponent(TableExternalId, { value: row.original.external_id, truncate: 40 })
    },
    {
      header: "Kinds",
      cell: ({ row }) => renderComponent(TableBadgeField, { value: row.original.kind, variant: "outline" })
    },
    {
      header: "Last Modified At",
      cell: ({ row }) => renderComponent(TableTimeField, { date: row.original.last_modified })
    }
  ];
</script>

<FlexTable
  {...props}
  data={props.items}
  {columns}
  emptyMessage="No Knowledge Items Stored."
/>