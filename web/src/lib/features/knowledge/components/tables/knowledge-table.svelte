<script lang="ts">
	import { goto } from "$app/navigation";
	import { resolve } from "$app/paths";
	import FlexTable from "$lib/components/flex-table/flex-table.svelte";
	import type { ColumnDef } from "@tanstack/table-core";
	import type { Knowledge } from "../../schemas";
	import { renderComponent } from "$lib/components/ui/data-table";
	import TableTimeField from "$lib/components/flex-table/table-time-field.svelte";
	import TableBadgeField from "$lib/components/flex-table/table-badge-field.svelte";
	import TableField from "$lib/components/flex-table/table-field.svelte";
	import TableKnowledgeName from "./table-knowldege-name.svelte";
	import TableActionsField from "$lib/components/flex-table/table-actions-field.svelte";
  import RefreshCwIcon from "@lucide/svelte/icons/refresh-cw";
  import TrashIcon from "@lucide/svelte/icons/trash";
  import ArrowRightIcon from "@lucide/svelte/icons/arrow-right";
	import { knowledgeService } from "../../service";
	import ConfirmDialog from "$lib/components/confirm-dialog.svelte";

  type Props = {
		knowledges: Knowledge[]
		isLoading?: boolean;
		class?: string;
  }

  let props: Props = $props();
  let deleteConfirmDialogOpen = $state(false);
  let deleteConfirmDialogId = $state<string | null>(null);

  const syncMutation = $derived(knowledgeService.sync());
  const deleteMutation = $derived(knowledgeService.delete());

	function handleOpenKnowledge(id: string) {
		goto(resolve("/(protected)/knowledge/[id]", { id }));
	}

  const handleSync = (id: string) => {
    syncMutation.mutate({ id });
  }

  const handleConfirmDelete = (id: string) => {
    deleteConfirmDialogOpen = true;
    deleteConfirmDialogId = id;
  }

  const handleDelete = () => {
    deleteMutation.mutate(deleteConfirmDialogId!);
  }

  const columns: ColumnDef<Knowledge>[] = [
    {
      header: "Name",
      cell: ({ row }) => renderComponent(TableKnowledgeName, { knowledge: row.original })
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
      }),
    },
    {
      header: "Actions",
      cell: ({ row }) => renderComponent(TableActionsField, {
        actions: [
          {
            key: "view",
            icon: ArrowRightIcon,
            onClick: () => handleOpenKnowledge(row.original.id)
          },
          {
            key: "sync",
            icon: RefreshCwIcon,
            variant: "outline",
            onClick: () => handleSync(row.original.id)
          },
          {
            key: "delete",
            icon: TrashIcon,
            variant: "destructive",
            onClick: () => handleConfirmDelete(row.original.id)
          }
        ]
      })
    }
  ]
</script>

<FlexTable
  {...props}
  data={props.knowledges}
  {columns}
  emptyMessage="No Knowledge Found."
/>

<ConfirmDialog
  open={deleteConfirmDialogOpen}
  title="Delete Knowledge"
  danger
  description="Are you sure you want to delete this knowledge? This action cannot be undone."
  onConfirm={handleDelete}
  onCancel={() => {
    deleteConfirmDialogOpen = false;
    deleteConfirmDialogId = null;
  }}
/>