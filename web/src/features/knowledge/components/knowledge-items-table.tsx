import { formatUpdatedAt } from "@/components/catalog/format-updated-at";
import { DataTable } from "@/components/data-table";
import { Badge } from "@/components/ui/badge";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useParams } from "react-router";

import { useListKnowledgeItems } from "../hooks";

import { KnowledgeItemsTableNameField } from "./knowledge-items-table-name-field";

import type { DataTableColumn, DataTablePagination } from "@/components/data-table";
import type { KnowledgeItem } from "../schemas";
function useKnowledgeItemsTableQuery(pagination: DataTablePagination) {
  const { id = "", tenantSlug = "" } = useParams<{ id: string; tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  return useListKnowledgeItems(tenant?.id ?? "", id, { searchParams: pagination });
}


export function KnowledgeItemsTable() {
  const columns: DataTableColumn<KnowledgeItem>[] = [
    {
      id: "name",
      header: "Name",
      cell: (row) => <KnowledgeItemsTableNameField item={row} />,
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
