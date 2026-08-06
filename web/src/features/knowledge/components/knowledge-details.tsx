import { formatUpdatedAt } from "@/components/catalog/format-updated-at";
import { JsonViewer } from "@/components/json-viewer";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";

import { KnowledgeIndexesSection } from "./knowledge-indexes-section";
import { SyncStatusBadge } from "./sync-status-badge";

import type { Knowledge } from "../schemas";

export function KnowledgeDetails({ knowledge }: { knowledge: Knowledge }) {
  return (
    <ScrollArea className="h-full">
      <div className="flex flex-col gap-4 pr-3">
        <Card>
          <CardHeader className="pb-3">
            <div className="flex flex-wrap items-center gap-2">
              <CardTitle className="text-base">{knowledge.name}</CardTitle>
              <Badge variant="secondary" className="font-normal">
                {knowledge.driver}
              </Badge>
              <SyncStatusBadge status={knowledge.sync_status} />
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            {knowledge.description ? (
              <p className="text-sm text-muted-foreground">{knowledge.description}</p>
            ) : (
              <p className="text-sm italic text-muted-foreground">No description</p>
            )}
            <dl className="space-y-3 text-sm">
              <div className="space-y-1">
                <dt className="text-xs font-medium tracking-wide text-muted-foreground">
                  Last synced
                </dt>
                <dd>
                  {knowledge.last_synced_at
                    ? formatUpdatedAt(knowledge.last_synced_at.toISOString())
                    : "Never"}
                </dd>
              </div>
              {knowledge.sync_error ? (
                <div className="space-y-1">
                  <dt className="text-xs font-medium tracking-wide text-muted-foreground">
                    Sync error
                  </dt>
                  <dd className="wrap-break-word text-destructive">{knowledge.sync_error}</dd>
                </div>
              ) : null}
            </dl>
          </CardContent>
        </Card>

        <JsonViewer title={`Configuration (${knowledge.driver})`} value={knowledge.configuration} />

        <KnowledgeIndexesSection knowledgeID={knowledge.id} />
      </div>
    </ScrollArea>
  );
}
