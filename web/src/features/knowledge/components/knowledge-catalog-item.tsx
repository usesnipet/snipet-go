import { CatalogCard } from "@/components/catalog";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useNavigate } from "@/hooks/use-navigate";
import { useDialog } from "@/lib/dialog";
import { ROUTES } from "@/routes";
import {
  BookOpenIcon,
  MoreHorizontal,
  PencilIcon,
  RefreshCw,
  RotateCcw,
  Trash2,
} from "lucide-react";
import { useParams } from "react-router";

import { useSyncKnowledge } from "../hooks";

import { DeleteKnowledgeDialog } from "./delete-knowledge-dialog";
import { SyncStatusBadge } from "./sync-status-badge";
import { UpdateKnowledgeDialog } from "./update-knowledge-dialog";

import type { Knowledge } from "../schemas";

export function KnowledgeCatalogItem({ knowledge }: { knowledge: Knowledge }) {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const tenantId = tenant?.id ?? "";
  const { openDialog } = useDialog();
  const navigate = useNavigate();
  const { mutate: sync, isPending: isSyncing } = useSyncKnowledge();

  const isSyncBusy =
    knowledge.sync_status === "pending" || knowledge.sync_status === "in_progress";

  const openEdit = () => {
    openDialog({
      component: UpdateKnowledgeDialog,
      props: { tenantId, knowledge },
    });
  };

  const openDelete = () => {
    openDialog({
      component: DeleteKnowledgeDialog,
      props: { tenantId, knowledge },
    });
  };

  const goToDetail = () => {
    navigate(ROUTES.knowledgeDetail, { params: { id: knowledge.id } });
  };

  return (
    <CatalogCard
      icon={
        <BookOpenIcon className="size-8 shrink-0 rounded-full border border-border p-1.5 text-muted-foreground" />
      }
      title={knowledge.name}
      badge={knowledge.driver}
      description={knowledge.description || undefined}
      updatedAt={knowledge.last_synced_at?.toISOString()}
      extraBadges={<SyncStatusBadge status={knowledge.sync_status} />}
      onClick={goToDetail}
      headerActions={
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              size="icon"
              aria-label="Knowledge actions"
              onClick={(event) => event.stopPropagation()}
            >
              <MoreHorizontal className="size-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" onClick={(event: React.MouseEvent<HTMLDivElement>) => event.stopPropagation()}>
            <DropdownMenuItem
              disabled={isSyncBusy || isSyncing}
              onClick={() => sync({ tenantId, id: knowledge.id })}
            >
              <RefreshCw />
              Sync
            </DropdownMenuItem>
            <DropdownMenuItem
              disabled={isSyncBusy || isSyncing}
              onClick={() => sync({ tenantId, id: knowledge.id, force: true })}
            >
              <RotateCcw />
              Full resync
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={openEdit}>
              <PencilIcon />
              Edit
            </DropdownMenuItem>
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              onClick={openDelete}
            >
              <Trash2 />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      }
    />
  );
}
