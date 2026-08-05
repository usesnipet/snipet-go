import { Badge } from "@/components/ui/badge";

import type { SyncStatus } from "../schemas";

const syncStatusLabel: Record<SyncStatus, string> = {
  pending: "Pending",
  in_progress: "Syncing",
  failed: "Failed",
  success: "Synced",
};

export function SyncStatusBadge({ status }: { status: SyncStatus | null }) {
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
