import { CatalogCard } from "@/components/catalog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useDialog } from "@/lib/dialog";
import {
  BookOpenIcon,
  MoreHorizontal,
  PencilIcon,
  RefreshCw,
  RotateCcw,
  Trash2,
} from "lucide-react";

import { useSyncKnowledge } from "../hooks";

import { DeleteKnowledgeDialog } from "./delete-knowledge-dialog";
import { UpdateKnowledgeDialog } from "./update-knowledge-dialog";

import type { Knowledge, SyncStatus } from "../schemas";

const syncStatusLabel: Record<SyncStatus, string> = {
  pending: "Pending",
  in_progress: "Syncing",
  failed: "Failed",
  success: "Synced",
};

function SyncStatusBadge({ status }: { status: SyncStatus | null }) {
  if (!status) {
    return (
      <Badge variant="outline" className="shrink-0 font-normal text-muted-foreground">
        Never synced
      </Badge>
    );
  }

  return (
    <Badge
      variant={status === "failed" ? "destructive" : "outline"}
      className="shrink-0 font-normal"
    >
      {syncStatusLabel[status]}
    </Badge>
  );
}

export function KnowledgeCatalogItem({ knowledge }: { knowledge: Knowledge }) {
  const { openDialog } = useDialog();
  const { mutate: sync, isPending: isSyncing } = useSyncKnowledge();

  const isSyncBusy =
    knowledge.sync_status === "pending" || knowledge.sync_status === "in_progress";

  const openEdit = () => {
    openDialog({
      component: UpdateKnowledgeDialog,
      props: { knowledge },
    });
  };

  const openDelete = () => {
    openDialog({
      component: DeleteKnowledgeDialog,
      props: { knowledge },
    });
  };

  return (
    <CatalogCard
      icon={<BookOpenIcon className="size-8 shrink-0 text-muted-foreground border border-border rounded-full p-1.5" />}
      title={knowledge.name}
      badge={knowledge.driver}
      description={knowledge.description || undefined}
      updatedAt={knowledge.last_synced_at?.toISOString()}
      extraBadges={<SyncStatusBadge status={knowledge.sync_status} />}
      headerActions={
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" aria-label="Knowledge actions">
              <MoreHorizontal className="size-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem
              disabled={isSyncBusy || isSyncing}
              onClick={() => sync({ id: knowledge.id })}
            >
              <RefreshCw />
              Sync
            </DropdownMenuItem>
            <DropdownMenuItem
              disabled={isSyncBusy || isSyncing}
              onClick={() => sync({ id: knowledge.id, force: true })}
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
