import { Page } from "@/components/page";

export const AdminPage = () => {
  return (
    <Page
      title="Admin"
      description="Overview of your Snipet workspace."
      documentTitle="Admin"
    >
      <p className="text-sm text-muted-foreground">
        Select a section from the sidebar to get started.
      </p>
    </Page>
  )
}
