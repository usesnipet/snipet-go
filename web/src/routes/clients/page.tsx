import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { ClientList } from "@/features/client/components/client-list";
import { CreateClientDialog } from "@/features/client/components/create-client-dialog";
import { useDialog } from "@/lib/dialog/use-dialog";
import { Plus } from "lucide-react";

export function ClientsPage() {
  const { openDialog } = useDialog();

  const openCreate = () => {
    openDialog({
      component: CreateClientDialog,
      props: {},
    });
  };

  return (
    <Page
      title="Clients"
      description="Create and manage clients."
      documentTitle="Clients"
    >
      <PageActions>
        <Button onClick={openCreate}>
          <Plus />
          Create Client
        </Button>
      </PageActions>
      <ClientList />
    </Page>
  )
}
