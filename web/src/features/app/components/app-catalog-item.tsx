import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { useTenantStore } from "@/features/tenant/store";
import { useNavigate } from "@/hooks/use-navigate";
import { useDialog } from "@/lib/dialog";
import { ROUTES } from "@/routes";
import { BotIcon, PencilIcon, TrashIcon } from "lucide-react";

import { AppCode } from "./app-code";
import { DeleteAppDialog } from "./delete-app-dialog";
import { UpdateAppDialog } from "./update-app-dialog";

import type { App } from "../schemas";


function AuthBadges({ app }: { app: App }) {
  const methods: string[] = [];
  if (app.auth_config.oidc.enabled) methods.push("OIDC");
  if (app.auth_config.webhook.enabled) methods.push("Webhook");

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

export function AppCatalogItem({ app }: { app: App }) {
  const tenant = useTenantStore((state) => state.tenant);
  const { openDialog } = useDialog();
  const navigate = useNavigate();

  const openEdit = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation();
    event.preventDefault();
    if (!tenant) return;
    openDialog({
      component: UpdateAppDialog,
      props: { tenantId: tenant.id, app },
    });
  };

  const openDelete = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation();
    event.preventDefault();
    if (!tenant) return;
    openDialog({
      component: DeleteAppDialog,
      props: { tenantId: tenant.id, app },
    });
  };

  const goToAppPage = () => {
    navigate(ROUTES.app, { params: { appCode: app.code } });
  }

  return (
    <Card className="flex h-full flex-col cursor-pointer" onClick={goToAppPage}>
      <CardHeader className="flex flex-row items-start gap-3 space-y-0 pb-3">
        <BotIcon className="size-8 text-muted-foreground border border-border rounded-full p-1.5" />
        <div className="flex min-w-0 flex-1 items-center justify-between space-y-1">
          <div className="flex min-w-0 flex-wrap items-center gap-2">
            <h2 className="truncate text-base font-semibold leading-tight">{app.name}</h2>
            <AuthBadges app={app} />
          </div>
          <div>
            <Button
              variant="ghost"
              size="icon"
              aria-label="Edit app"
              onClick={openEdit}
            >
              <PencilIcon className="size-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              aria-label="Delete app"
              onClick={openDelete}
            >
              <TrashIcon className="size-4" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="flex flex-1 flex-col gap-3 pt-0">
        <AppCode code={app.code} />
      </CardContent>
    </Card>
  )
}
