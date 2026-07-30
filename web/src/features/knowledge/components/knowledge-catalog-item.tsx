import { CatalogCard } from "@/components/catalog";
import { Badge } from "@/components/ui/badge";
import { BookOpenIcon } from "lucide-react";

import type { Knowledge, SyncStatus } from "../schemas";

const syncStatusLabel: Record<SyncStatus, string> = {
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
  return (
    <CatalogCard
      icon={<BookOpenIcon className="size-8 shrink-0 text-muted-foreground border border-border rounded-full p-1.5" />}
      title={knowledge.name}
      badge={knowledge.driver}
      description={knowledge.description || undefined}
      updatedAt={knowledge.last_synced_at?.toISOString()}
      extraBadges={<SyncStatusBadge status={knowledge.sync_status} />}
    />
  );
}
