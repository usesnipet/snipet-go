import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { ClientList } from "@/features/client/components/client-list";
import { CreateClientDialog } from "@/features/client/components/create-client-dialog";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useDialog } from "@/lib/dialog/use-dialog";
import { Plus } from "lucide-react";
import { useParams } from "react-router";

export function ClientsPage() {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const { openDialog } = useDialog();

  const openCreate = () => {
    if (!tenant) return;
    openDialog({
      component: CreateClientDialog,
      props: { tenantId: tenant.id },
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
