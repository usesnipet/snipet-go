<script lang="ts">
	import FlexTable from "$lib/components/flex-table/flex-table.svelte";
	import type { ColumnDef } from "@tanstack/table-core";
	import type { Session } from "../schemas";
	import { renderComponent } from "$lib/components/ui/data-table";
	import TableField from "$lib/components/flex-table/table-field.svelte";
	import TableActionsField from "$lib/components/flex-table/table-actions-field.svelte";
	import TrashIcon from "@lucide/svelte/icons/trash";
	import { sessionService } from "../service";
	import ConfirmDialog from "$lib/components/confirm-dialog.svelte";

	type Props = {
		sessions: Session[];
		clientCode: string;
		isLoading?: boolean;
		class?: string;
		onRowClick?: (session: Session) => void;
	};

	let props: Props = $props();

	let deleteConfirmDialogOpen = $state(false);
	let deleteConfirmDialogId = $state<string | null>(null);

	const deleteMutation = $derived(sessionService.delete(props.clientCode));

	const handleConfirmDelete = (id: string) => {
		deleteConfirmDialogOpen = true;
		deleteConfirmDialogId = id;
	};

	const handleDelete = () => {
		deleteMutation.mutate(deleteConfirmDialogId!);
		deleteConfirmDialogId = null;
	};

	const columns: ColumnDef<Session>[] = [
		{
			header: "ID",
			cell: ({ row }) =>
				renderComponent(TableField, { value: row.original.id, truncate: 8 }),
		},
		{
			header: "Name",
			cell: ({ row }) =>
				renderComponent(TableField, { value: row.original.metadata.name, truncate: 40 }),
		},
		{
			header: "Agent",
			cell: ({ row }) =>
				renderComponent(TableField, { value: row.original.agent?.name ?? "Unknown", truncate: 40 }),
		},
		{
			header: "Actions",
			cell: ({ row }) =>
				renderComponent(TableActionsField, {
					actions: [
						{
							key: "delete",
							icon: TrashIcon,
							variant: "destructive",
							onClick: () => handleConfirmDelete(row.original.id),
						},
					],
				}),
		},
	];
</script>

<FlexTable
	{...props}
	data={props.sessions}
	{columns}
	onRowClick={props.onRowClick}
	emptyMessage="No sessions found."
/>

<ConfirmDialog
	bind:open={deleteConfirmDialogOpen}
	title="Delete Session"
	danger
	description="Are you sure you want to delete this session? This action cannot be undone."
	onConfirm={handleDelete}
	onCancel={() => {
		deleteConfirmDialogId = null;
	}}
/>
