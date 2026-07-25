import { CatalogList } from "@/components/catalog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Loading } from "@/components/ui/loading";
import { useDialog } from "@/lib/dialog";
import { BotIcon, PencilIcon, TrashIcon } from "lucide-react";

import { useListClient } from "../hooks";

import { ClientCode } from "./client-code";
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
      <Badge key={methods[0]} variant="outline" className="shrink-0 font-normal">
        {methods[0]}
      </Badge>
      {methods.length > 1 && (
        <Badge variant="outline" className="shrink-0 font-normal">
          +{methods.length - 1}
        </Badge>
      )}
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
        <Card className="flex h-full flex-col">
          <CardHeader className="flex flex-row items-start gap-3 space-y-0 pb-3">
            <BotIcon className="size-8 text-muted-foreground border border-border rounded-full p-1.5" />
            <div className="flex min-w-0 flex-1 items-center justify-between space-y-1">
              <div className="flex min-w-0 flex-wrap items-center gap-2">
                <h2 className="truncate text-base font-semibold leading-tight">{client.name}</h2>
                <AuthBadges client={client} />
              </div>
              <div>
                <Button
                  variant="ghost"
                  size="icon"
                  aria-label="Edit client"
                  onClick={() => openEdit(client)}
                >
                  <PencilIcon className="size-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  aria-label="Delete client"
                  onClick={() => openDelete(client)}
                >
                  <TrashIcon className="size-4" />
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent className="flex flex-1 flex-col gap-3 pt-0">
            <ClientCode client={client} />
          </CardContent>
        </Card>
      )}
    />
  );
}
