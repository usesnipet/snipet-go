<script lang="ts">
	import FlexTable from "$lib/components/flex-table/flex-table.svelte";
	import type { ColumnDef } from "@tanstack/table-core";
	import type { APIKey } from "../schemas";
	import { renderComponent } from "$lib/components/ui/data-table";
	import TableField from "$lib/components/flex-table/table-field.svelte";
	import TableTimeField from "$lib/components/flex-table/table-time-field.svelte";
	import TableBadgeField from "$lib/components/flex-table/table-badge-field.svelte";
	import TableActionsField from "$lib/components/flex-table/table-actions-field.svelte";
	import PencilIcon from "@lucide/svelte/icons/pencil";
	import ApiKeyUpdateDialog from "./api-key-update-dialog.svelte";

	type Props = {
		apiKeys: APIKey[];
		isLoading?: boolean;
		class?: string;
	};

	let props: Props = $props();

	let updateDialogOpen = $state(false);
	let updateDialogApiKey: APIKey | null = $state(null);

	function handleOpenUpdateDialog(apiKey: APIKey) {
		updateDialogApiKey = apiKey;
		updateDialogOpen = true;
	}

	const columns: ColumnDef<APIKey>[] = [
		{
			header: "Name",
			cell: ({ row }) => renderComponent(TableField, { value: row.original.name }),
		},
		{
			header: "Prefix",
			cell: ({ row }) =>
				renderComponent(TableField, { value: row.original.key_id, truncate: 9 }),
		},
		{
			header: "Expires at",
			cell: ({ row }) =>
				renderComponent(TableTimeField, { date: row.original.expires_at }),
		},
		{
			header: "Status",
			cell: ({ row }) =>
				renderComponent(TableBadgeField, {
					value: row.original.active ? "Active" : "Inactive",
				}),
		},
		{
			header: "Actions",
			cell: ({ row }) =>
				renderComponent(TableActionsField, {
					actions: [
						{
							key: "update",
							icon: PencilIcon,
							onClick: () => handleOpenUpdateDialog(row.original),
						},
					],
				}),
		},
	];
</script>

<FlexTable
	{...props}
	data={props.apiKeys}
	{columns}
	emptyMessage="No Api Found."
/>

<ApiKeyUpdateDialog bind:open={updateDialogOpen} apiKey={updateDialogApiKey ?? undefined} />
