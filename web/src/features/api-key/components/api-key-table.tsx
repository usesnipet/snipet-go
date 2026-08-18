import { DataTable } from "@/components/data-table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { DateFormat } from "@/components/ui/date";
import {
  DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger
} from "@/components/ui/dropdown-menu";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useDialog } from "@/lib/dialog";
import { CalendarClock, MoreHorizontal, RefreshCw, Trash2 } from "lucide-react";
import { useParams } from "react-router";

import { useListApiKey } from "../hooks";

import { ApiKeySecretDialog } from "./api-key-secret-dialog";
import { DeleteApiKeyDialog } from "./delete-api-key-dialog";
import { RollApiKeyDialog } from "./roll-api-key-dialog";
import { UpdateApiKeyExpirationDialog } from "./update-api-key-expiration-dialog";

import type { DataTableColumn, DataTablePagination } from "@/components/data-table";
import type { ApiKey, ApiKeyWithSecret } from "../schemas";

function useApiKeyTableQuery(pagination: DataTablePagination) {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  return useListApiKey(tenant?.id ?? "", { searchParams: pagination })
}

export function ApiKeyTable() {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const { openDialog } = useDialog();

  const showSecret = (apiKey: ApiKeyWithSecret) => {
    openDialog({
      component: ApiKeySecretDialog,
      props: {
        secret: apiKey.key,
        title: "API Key rolled",
        description: "Copy this key now. You will not be able to see it again.",
      },
    });
  };

  const openExpiration = (apiKey: ApiKey) => {
    if (!tenant) return;
    openDialog({
      component: UpdateApiKeyExpirationDialog,
      props: { tenantId: tenant.id, apiKey },
    });
  };

  const openRoll = (apiKey: ApiKey) => {
    if (!tenant) return;
    openDialog({
      component: RollApiKeyDialog,
      props: { tenantId: tenant.id, apiKey, onRolled: (rolled) => showSecret(rolled) },
    });
  };

  const openDelete = (apiKey: ApiKey) => {
    if (!tenant) return;
    openDialog({
      component: DeleteApiKeyDialog,
      props: { tenantId: tenant.id, apiKey },
    });
  };

  const columns: DataTableColumn<ApiKey>[] = [
    {
      id: "name",
      header: "Name",
      cell: (row) => (
        <div className="flex flex-col gap-0.5">
          <span className="font-medium">{row.name}</span>
          <span className="font-mono text-xs text-muted-foreground">{row.key_id}</span>
        </div>
      ),
    },
    {
      id: "status",
      header: "Status",
      cell: (row) => (
        <Badge variant={row.active ? "default" : "secondary"}>
          {row.active ? "Active" : "Disabled"}
        </Badge>
      ),
    },
    {
      id: "expires_at",
      header: "Expires",
      cell: (row) => (
        <DateFormat className="text-muted-foreground" emptyValue="Never" date={row.expires_at} />
      ),
    },
    {
      id: "created_at",
      header: "Created",
      cell: (row) => (
        <DateFormat className="text-muted-foreground" date={row.created_at} />
      ),
    },
    {
      id: "actions",
      header: <span className="sr-only">Actions</span>,
      headerClassName: "w-12",
      className: "text-right",
      cell: (row) => (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon-sm">
              <MoreHorizontal />
              <span className="sr-only">Open menu</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => openExpiration(row)}>
              <CalendarClock />
              Update expiration
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => openRoll(row)}>
              <RefreshCw />
              Roll key
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              onClick={() => openDelete(row)}
            >
              <Trash2 />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ),
    },
  ];

  return (
    <DataTable
      columns={columns}
      useQuery={useApiKeyTableQuery}
      getRowKey={(row) => row.id}
      emptyMessage="No API keys yet."
    />
  )
}
