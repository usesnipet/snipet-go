import { Page } from "@/components/page";
import { ClientUsersTable } from "@/features/client/components/client-users-table";

export function ClientUsersPage() {
  return (
    <Page
      title="Users"
      description="Users associated with this client."
      documentTitle="Users"
    >
      <ClientUsersTable />
    </Page>
  )
}
