import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";

export const AdminPage = () => {
  return (
    <Page
      title="Admin"
      description="Admin"
      documentTitle="Admin"
    >
      <PageActions>
        <Button>
          <Plus className="size-4" /> Create Api Key
        </Button>
      </PageActions>
      <div>a</div>
    </Page>
  )
}