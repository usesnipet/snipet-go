import { Page } from "@/components/page";

export const AppPage = () => {
  return (
    <Page
      title="App"
      description="Overview of your Snipet workspace."
      documentTitle="App"
    >
      <p className="text-sm text-muted-foreground">
        Select a section from the sidebar to get started.
      </p>
    </Page>
  )
}
