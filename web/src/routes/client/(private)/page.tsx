import { Page } from "@/components/page";

export const ClientPage = () => {
  return (
    <Page
      title="Client"
      description="Overview of your Snipet workspace."
      documentTitle="Client"
    >
      <p className="text-sm text-muted-foreground">
        Select a section from the sidebar to get started.
      </p>
    </Page>
  )
}
