import { LoadingFallback } from "@/components/loading-fallback";
import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { ApiKeySecretDialog } from "@/features/api-key/components/api-key-secret-dialog";
import { ApiKeyTable } from "@/features/api-key/components/api-key-table";
import { CreateApiKeyDialog } from "@/features/api-key/components/create-api-key-dialog";
import { DeleteApiKeyDialog } from "@/features/api-key/components/delete-api-key-dialog";
import { RollApiKeyDialog } from "@/features/api-key/components/roll-api-key-dialog";
import { UpdateApiKeyExpirationDialog } from "@/features/api-key/components/update-api-key-expiration-dialog";
import { useListApiKey } from "@/features/api-key/hooks";
import { Plus } from "lucide-react";
import { useState } from "react";

import type { ApiKey, ApiKeyWithSecret } from "@/features/api-key/schemas";

export function AdminApiKeysPage() {
  return (
    <Page
      title="API Keys"
      description="Create and manage API keys for authenticating requests."
      documentTitle="API Keys"
    >
      <AdminApiKeysContent />
    </Page>
  )
}

function AdminApiKeysContent() {
  const { data, isLoading } = useListApiKey();

  const [selected, setSelected] = useState<ApiKey | null>(null);
  const [expirationOpen, setExpirationOpen] = useState(false);
  const [rollOpen, setRollOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [secret, setSecret] = useState<ApiKeyWithSecret | null>(null);
  const [secretTitle, setSecretTitle] = useState("API Key created");

  const openExpiration = (apiKey: ApiKey) => {
    setSelected(apiKey);
    setExpirationOpen(true);
  };

  const openRoll = (apiKey: ApiKey) => {
    setSelected(apiKey);
    setRollOpen(true);
  };

  const openDelete = (apiKey: ApiKey) => {
    setSelected(apiKey);
    setDeleteOpen(true);
  };

  const showSecret = (apiKey: ApiKeyWithSecret, title: string) => {
    setSecretTitle(title);
    setSecret(apiKey);
  };

  if (isLoading) {
    return <LoadingFallback />;
  }

  return (
    <>
      <PageActions>
        <CreateApiKeyDialog onCreated={(apiKey) => showSecret(apiKey, "API Key created")}>
          <Button>
            <Plus />
            Create API Key
          </Button>
        </CreateApiKeyDialog>
      </PageActions>

      <ApiKeyTable
        data={data?.data ?? []}
        onUpdateExpiration={openExpiration}
        onRoll={openRoll}
        onDelete={openDelete}
      />

      <UpdateApiKeyExpirationDialog
        apiKey={selected}
        open={expirationOpen}
        onOpenChange={setExpirationOpen}
      />

      <RollApiKeyDialog
        apiKey={selected}
        open={rollOpen}
        onOpenChange={setRollOpen}
        onRolled={(apiKey) => showSecret(apiKey, "API Key rolled")}
      />

      <DeleteApiKeyDialog
        apiKey={selected}
        open={deleteOpen}
        onOpenChange={setDeleteOpen}
      />

      <ApiKeySecretDialog
        open={secret !== null}
        onOpenChange={(open) => {
          if (!open) setSecret(null);
        }}
        secret={secret?.key ?? null}
        title={secretTitle}
        description="Copy this key now. You will not be able to see it again."
      />
    </>
  )
}
