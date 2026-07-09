<script lang="ts" generics="TData, TValue">
	import {
		type ColumnDef,
		getCoreRowModel,
	} from "@tanstack/table-core";
	import { createSvelteTable, FlexRender } from "$lib/components/ui/data-table/index.js";
	import * as Table from "$lib/components/ui/table/index.js";
	import { cn } from "$lib/utils";
	import { ScrollArea } from "../ui/scroll-area";

	type DataTableProps<TData, TValue> = {
		data: TData[]
		columns: ColumnDef<TData, TValue>[];
		isLoading?: boolean;
		emptyMessage?: string;
		onRowClick?: (row: TData) => void;
		class?: string;
	};

	let {
		data,
		columns,
		isLoading = false,
		emptyMessage = "No results.",
		onRowClick,
		class: className,
	}: DataTableProps<TData, TValue> = $props();

	const table = createSvelteTable({
		get data() {
			return data;
		},
		get columns() {
			return columns;
		},
		getCoreRowModel: getCoreRowModel(),
	});
</script>

<div class={cn("flex h-full min-h-0 w-full flex-1 flex-col gap-3", className)}>
	<div class="flex min-h-0 flex-1 flex-col overflow-hidden rounded-sm border">
		<ScrollArea class="h-full min-h-0 flex-1">
			<Table.Root class="relative *:data-[slot=table-container]:overflow-visible">
				<Table.Header class="bg-accent sticky top-0 z-10">
					{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
						<Table.Row>
							{#each headerGroup.headers as header (header.id)}
								<Table.Head colspan={header.colSpan}>
									{#if !header.isPlaceholder}
										<FlexRender
											content={header.column.columnDef.header}
											context={header.getContext()}
										/>
									{/if}
								</Table.Head>
							{/each}
						</Table.Row>
					{/each}
				</Table.Header>
				<Table.Body>
					{#if isLoading}
						<Table.Row>
							<Table.Cell colspan={columns.length} class="text-muted-foreground h-24 text-center">
								Loading...
							</Table.Cell>
						</Table.Row>
					{:else}
						{#each table.getRowModel().rows as row (row.id)}
							<Table.Row
								data-state={row.getIsSelected() && "selected"}
								class={onRowClick ? "cursor-pointer" : undefined}
								onclick={() => onRowClick?.(row.original)}
							>
								{#each row.getVisibleCells() as cell (cell.id)}
									<Table.Cell>
										<FlexRender
											content={cell.column.columnDef.cell}
											context={cell.getContext()}
										/>
									</Table.Cell>
								{/each}
							</Table.Row>
						{:else}
							<Table.Row>
								<Table.Cell colspan={columns.length} class="text-muted-foreground h-24 text-center">
									{emptyMessage}
								</Table.Cell>
							</Table.Row>
						{/each}
					{/if}
				</Table.Body>
			</Table.Root>
		</ScrollArea>
	</div>
</div>
