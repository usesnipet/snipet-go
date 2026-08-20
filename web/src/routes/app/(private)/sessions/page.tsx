import { Page } from "@/components/page";
import { AppSessionsTable } from "@/features/app/components/app-sessions-table";

export function AppSessionsPage() {
  return (
    <Page
      title="Sessions"
      description="Sessions for this app."
      documentTitle="Sessions"
    >
      <AppSessionsTable />
    </Page>
  )
}
