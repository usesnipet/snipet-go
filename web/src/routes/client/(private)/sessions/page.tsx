import { Page } from "@/components/page";
import { ClientSessionsTable } from "@/features/client/components/client-sessions-table";

export function ClientSessionsPage() {
  return (
    <Page
      title="Sessions"
      description="Sessions for this client."
      documentTitle="Sessions"
    >
      <ClientSessionsTable />
    </Page>
  )
}
