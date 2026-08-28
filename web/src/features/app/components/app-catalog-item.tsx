import { SecretKeyDialog } from "@/components/secret-key-dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { DateFormat } from "@/components/ui/date";
import {
  DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger
} from "@/components/ui/dropdown-menu";
import { useNavigate } from "@/hooks/use-navigate";
import { useDialog } from "@/lib/dialog";
import { ROUTES } from "@/routes";
import {
  BotIcon, MoreHorizontal, PencilIcon, PowerIcon, PowerOffIcon, RefreshCwIcon, TrashIcon
} from "lucide-react";

import { useSetActiveApp } from "../hooks";

import { AppCode } from "./app-code";
import { DeleteAppDialog } from "./delete-app-dialog";
import { RollAppDialog } from "./roll-app-dialog";
import { UpdateAppDialog } from "./update-app-dialog";

import type { App, AppWithSecret } from "../schemas";

export function AppCatalogItem({ app }: { app: App }) {
  const { openDialog } = useDialog();
  const navigate = useNavigate();
  const { mutate: setActive } = useSetActiveApp();

  const showSecret = (app: AppWithSecret) => {
    openDialog({
      component: SecretKeyDialog,
      props: {
        secret: app.key,
        title: "App key rolled",
        description: "Copy this key now. You will not be able to see it again.",
      },
    });
  };

  const openRoll = () => {
    openDialog({
      component: RollAppDialog,
      props: { app, onRolled: (rolled) => showSecret(rolled) },
    });
  };

  const toggleActive = () => {
    setActive({ code: app.code, active: app.status !== "active" });
  };

  const openEdit = () => {
    openDialog({
      component: UpdateAppDialog,
      props: { app },
    });
  };

  const openDelete = () => {
    openDialog({
      component: DeleteAppDialog,
      props: { app },
    });
  };

  const goToAppPage = () => {
    navigate(ROUTES.app, { params: { appCode: app.code } });
  }

  const isActive = app.status === "active";

  return (
    <Card className="flex h-full flex-col cursor-pointer" onClick={goToAppPage}>
      <CardHeader className="flex flex-row items-start gap-3 space-y-0 pb-3">
        <BotIcon className="size-8 text-muted-foreground border border-border rounded-full p-1.5" />
        <div className="flex min-w-0 flex-1 items-center justify-between space-y-1">
          <div className="flex min-w-0 flex-wrap items-center gap-2">
            <h2 className="truncate text-base font-semibold leading-tight">{app.name}</h2>
            <Badge variant={app.status === "active" ? "default" : "secondary"}>
              {app.status}
            </Badge>
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                aria-label="Open menu"
                onClick={(event) => event.stopPropagation()}
              >
                <MoreHorizontal className="size-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" onClick={(event) => event.stopPropagation()}>
              <DropdownMenuItem onClick={openRoll}>
                <RefreshCwIcon />
                Roll key
              </DropdownMenuItem>
              {
                (app.status !== "pending" || app.public) && (
                  <DropdownMenuItem onClick={toggleActive}>
                    {isActive ? <PowerOffIcon /> : <PowerIcon />}
                    {isActive ? "Deactivate" : "Activate"}
                  </DropdownMenuItem>
                )
              }
              <DropdownMenuItem onClick={openEdit}>
                <PencilIcon />
                Edit
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                className="text-destructive focus:text-destructive"
                onClick={openDelete}
              >
                <TrashIcon />
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardHeader>
      <CardContent className="flex flex-1 flex-col gap-3 pt-0">
        <AppCode code={app.code} />
        <p className="mt-auto text-xs text-muted-foreground">
          Last verified: <DateFormat date={app.last_verified_at} emptyValue="Never" />
        </p>
      </CardContent>
    </Card>
  )
}
