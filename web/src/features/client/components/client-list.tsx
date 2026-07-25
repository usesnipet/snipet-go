import { CatalogCard, CatalogList } from "@/components/catalog";
import { Badge } from "@/components/ui/badge";
import { Loading } from "@/components/ui/loading";
import { useDialog } from "@/lib/dialog";
import { Pencil, Trash2 } from "lucide-react";

import { useListClient } from "../hooks";

import { DeleteClientDialog } from "./delete-client-dialog";
import { UpdateClientDialog } from "./update-client-dialog";

import type { Client } from "../schemas";

function AuthBadges({ client }: { client: Client }) {
  const methods: string[] = [];
  if (client.config.oidc.enabled) methods.push("OIDC");
  if (client.config.webhook.enabled) methods.push("Webhook");
  if (client.config.anonymous.enabled) methods.push("Anonymous");

  if (methods.length === 0) {
    return (
      <Badge variant="outline" className="shrink-0 font-normal text-muted-foreground">
        No auth
      </Badge>
    );
  }

  return (
    <>
      {methods.map((method) => (
        <Badge key={method} variant="outline" className="shrink-0 font-normal">
          {method}
        </Badge>
      ))}
    </>
  );
}

export function ClientList() {
  const { data, isLoading } = useListClient();
  const { openDialog } = useDialog();

  const openEdit = (client: Client) => {
    openDialog({
      component: UpdateClientDialog,
      props: { client },
    });
  };

  const openDelete = (client: Client) => {
    openDialog({
      component: DeleteClientDialog,
      props: { client },
    });
  };

  if (isLoading) {
    return (
      <div className="flex flex-1 items-center justify-center py-12">
        <Loading />
      </div>
    );
  }

  return (
    <CatalogList
      items={data?.data ?? []}
      emptyMessage="No clients yet."
      renderItem={(client) => (
        <CatalogCard
          title={client.name}
          badge={client.code}
          extraBadges={<AuthBadges client={client} />}
          actions={[
            {
              label: "Edit client",
              icon: <Pencil />,
              onClick: () => openEdit(client),
            },
            {
              label: "Delete client",
              icon: <Trash2 />,
              onClick: () => openDelete(client),
            },
          ]}
        />
      )}
    />
  );
}
