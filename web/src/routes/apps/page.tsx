import { Page, PageActions } from "@/components/page";
import { SecretKeyDialog } from "@/components/secret-key-dialog";
import { Button } from "@/components/ui/button";
import { AppList } from "@/features/app/components/app-list";
import { CreateAppDialog } from "@/features/app/components/create-app-dialog";
import { useTenantStore } from "@/features/tenant/store";
import { useDialog } from "@/lib/dialog/use-dialog";
import { Plus } from "lucide-react";

import type { AppWithSecret } from "@/features/app/schemas";

export function AppsPage() {
  const tenant = useTenantStore((state) => state.tenant);
  const { openDialog } = useDialog();

  const showSecret = (app: AppWithSecret) => {
    openDialog({
      component: SecretKeyDialog,
      props: {
        secret: app.key,
        title: "App key created",
        description: "Copy this key now. You will not be able to see it again.",
      },
    });
  };

  const openCreate = () => {
    if (!tenant) return;
    openDialog({
      component: CreateAppDialog,
      props: {
        tenantId: tenant.id,
        onCreated: (app) => showSecret(app),
      },
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
