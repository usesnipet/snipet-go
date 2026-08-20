import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { AppList } from "@/features/app/components/app-list";
import { CreateAppDialog } from "@/features/app/components/create-app-dialog";
import { useTenantStore } from "@/features/tenant/store";
import { useDialog } from "@/lib/dialog/use-dialog";
import { Plus } from "lucide-react";

export function AppsPage() {
  const tenant = useTenantStore((state) => state.tenant);
  const { openDialog } = useDialog();

  const openCreate = () => {
    if (!tenant) return;
    openDialog({
      component: CreateAppDialog,
      props: { tenantId: tenant.id },
    });
  };

  return (
    <Page
      title="Apps"
      description="Create and manage apps."
      documentTitle="Apps"
    >
      <PageActions>
        <Button onClick={openCreate}>
          <Plus />
          Create App
        </Button>
      </PageActions>
      <AppList />
    </Page>
  )
}
