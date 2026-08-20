import { DataTable } from "@/components/data-table";
import { useListAppUser } from "@/features/app-user/hooks";
import { useParams } from "react-router";

import type { DataTableColumn, DataTablePagination } from "@/components/data-table";
import type { User } from "@/features/app-user/schemas";

function useAppUsersTableQuery(pagination: DataTablePagination) {
  const { appCode = "" } = useParams();
  return useListAppUser(appCode, { searchParams: pagination })
}

export function AppUsersTable() {
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
      useQuery={useAppUsersTableQuery}
      getRowKey={(row) => row.id}
      emptyMessage="No users yet."
    />
  )
}
