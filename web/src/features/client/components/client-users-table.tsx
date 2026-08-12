import { DataTable } from "@/components/data-table";
import { useListClientUser } from "@/features/client-user/hooks";
import { useParams } from "react-router";

import type { DataTableColumn, DataTablePagination } from "@/components/data-table";
import type { User } from "@/features/client-user/schemas";

function useClientUsersTableQuery(pagination: DataTablePagination) {
  const { clientCode = "" } = useParams();
  return useListClientUser(clientCode, { searchParams: pagination })
}

export function ClientUsersTable() {
  const columns: DataTableColumn<User>[] = [
    {
      id: "name",
      header: "Name",
      cell: (row) => <span className="font-medium">{row.name}</span>,
    },
    {
      id: "email",
      header: "Email",
      cell: (row) => (
        <span className="text-muted-foreground">{row.email ?? "—"}</span>
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
      useQuery={useClientUsersTableQuery}
      getRowKey={(row) => row.id}
      emptyMessage="No users yet."
    />
  )
}
