import { formatUpdatedAt } from "@/components/catalog/format-updated-at";
import { DataTable } from "@/components/data-table";
import { Badge } from "@/components/ui/badge";
import { truncate } from "@/lib/utils";
import { useParams } from "react-router";

import { useListKnowledgeItems } from "../hooks";

import type { DataTableColumn, DataTablePagination } from "@/components/data-table";
import type { KnowledgeItem } from "../schemas";
function useKnowledgeItemsTableQuery(pagination: DataTablePagination) {
  const { id = "" } = useParams();
  return useListKnowledgeItems(id, { searchParams: pagination });
}

export function KnowledgeItemsTable() {
  const columns: DataTableColumn<KnowledgeItem>[] = [
    {
      id: "name",
      header: "Name",
      cell: (row) => <span className="font-medium">{truncate(row.name, 20)}</span>,
    },
    {
      id: "kind",
      header: "Kind",
      cell: (row) => (
        <Badge variant="outline" className="font-normal capitalize">
          {row.kind}
        </Badge>
      ),
    },
    {
      id: "external_id",
      header: "External ID",
      cell: (row) => (
        <span className="font-mono text-xs text-muted-foreground">{truncate(row.external_id, 50)}</span>
      ),
    },
    {
      id: "last_modified",
      header: "Last modified",
      cell: (row) =>
        row.last_modified ? formatUpdatedAt(row.last_modified.toISOString()) : "—",
    },
  ];

  return (
    <DataTable
      columns={columns}
      useQuery={useKnowledgeItemsTableQuery}
      pageSize={200}
      getRowKey={(row) => row.id}
      emptyMessage="No items yet. Sync this knowledge source to populate items."
    />
  );
}
