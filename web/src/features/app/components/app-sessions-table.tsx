import { DataTable } from "@/components/data-table";
import { useListSession } from "@/features/session/hooks";
import { useParams } from "react-router";

import type { DataTableColumn, DataTablePagination } from "@/components/data-table";
import type { Session } from "@/features/session/schemas";

function useAppSessionsTableQuery(pagination: DataTablePagination) {
  const { appCode = "" } = useParams();
  return useListSession(appCode, {
    searchParams: { ...pagination, include: "agent" },
  })
}

export function AppSessionsTable() {
  const columns: DataTableColumn<Session>[] = [
    {
      id: "name",
      header: "Name",
      cell: (row) => (
        <span className="font-medium">
          {typeof row.metadata.name === "string" && row.metadata.name
            ? row.metadata.name
            : "Untitled"}
        </span>
      ),
    },
    {
      id: "agent",
      header: "Agent",
      cell: (row) => (
        <span className="text-muted-foreground">{row.agent?.name ?? "—"}</span>
      ),
    },
    {
      id: "id",
      header: "ID",
      cell: (row) => (
        <span className="font-mono text-xs text-muted-foreground">{row.id}</span>
      ),
    },
  ];

  return (
    <DataTable
      columns={columns}
      useQuery={useAppSessionsTableQuery}
      getRowKey={(row) => row.id}
      emptyMessage="No sessions yet."
    />
  )
}
