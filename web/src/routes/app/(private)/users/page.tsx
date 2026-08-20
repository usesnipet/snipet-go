import { Page } from "@/components/page";
import { AppUsersTable } from "@/features/app/components/app-users-table";

export function AppUsersPage() {
  return (
    <Page
      title="Users"
      description="Users associated with this app."
      documentTitle="Users"
    >
      <AppUsersTable />
    </Page>
  )
}
