import { DataTable } from "@/components/data-table";
import { useListUser } from "@/features/user/hooks";
import { useParams } from "react-router";

import type { DataTableColumn } from "@/components/data-table";
import type { User } from "@/features/user/schemas";

export function ClientUsersTable() {
  const { clientCode = "" } = useParams();
  const { data, isLoading } = useListUser(clientCode);

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
      loading={isLoading}
      columns={columns}
      data={data?.data ?? []}
      getRowKey={(row) => row.id}
      emptyMessage="No users yet."
    />
  )
}
