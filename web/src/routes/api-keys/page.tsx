import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { ApiKeySecretDialog } from "@/features/api-key/components/api-key-secret-dialog";
import { ApiKeyTable } from "@/features/api-key/components/api-key-table";
import { CreateApiKeyDialog } from "@/features/api-key/components/create-api-key-dialog";
import { useTenantStore } from "@/features/tenant/store";
import { useDialog } from "@/lib/dialog/use-dialog";
import { Plus } from "lucide-react";

import type { ApiKeyWithSecret } from "@/features/api-key/schemas";
export function ApiKeysPage() {
  const tenant = useTenantStore((state) => state.tenant);
  const { openDialog } = useDialog();

  const showSecret = (apiKey: ApiKeyWithSecret) => {
    openDialog({
      component: ApiKeySecretDialog,
      props: {
        secret: apiKey.key,
        title: "API Key created",
        description: "Copy this key now. You will not be able to see it again.",
      },
    });
  };

  const openCreate = () => {
    if (!tenant) return;
    openDialog({
      component: CreateApiKeyDialog,
      props: {
        tenantId: tenant.id,
        onCreated: (apiKey) => showSecret(apiKey),
      },
    });
  };

  return (
    <Page
      title="API Keys"
      description="Create and manage API keys for authenticating requests."
      documentTitle="API Keys"
    >
      <PageActions>
        <Button onClick={openCreate}>
          <Plus />
          Create API Key
        </Button>
      </PageActions>
      <ApiKeyTable />
    </Page>
  )
}