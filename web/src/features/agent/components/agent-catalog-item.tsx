import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { useDialog } from "@/lib/dialog";
import { BotIcon, PencilIcon, TrashIcon } from "lucide-react";

import { DeleteAgentDialog } from "./delete-agent-dialog";
import { UpdateAgentDialog } from "./update-agent-dialog";

import type { Agent } from "../schemas";

export function AgentCatalogItem({ agent }: { agent: Agent }) {
  const { openDialog } = useDialog();

  const openEdit = () => {
    openDialog({
      component: UpdateAgentDialog,
      props: { agent },
    });
  };

  const openDelete = () => {
    openDialog({
      component: DeleteAgentDialog,
      props: { agent },
    });
  };

  const llmLabels = [...agent.llms]
    .sort((a, b) => a.priority - b.priority)
    .map((rel) => rel.llm.name || rel.llm.provider);

  return (
    <Card className="flex h-full flex-col">
      <CardHeader className="flex flex-row items-start gap-3 space-y-0 pb-3">
        <BotIcon className="size-8 shrink-0 text-muted-foreground border border-border rounded-full p-1.5" />
        <div className="flex min-w-0 flex-1 items-start justify-between gap-2">
          <div className="flex min-w-0 flex-col gap-1.5">
            <div className="flex min-w-0 flex-wrap items-center gap-2">
              <h2 className="truncate text-base font-semibold leading-tight">{agent.name}</h2>
              {llmLabels.length === 0 ? (
                <Badge variant="outline" className="shrink-0 font-normal text-muted-foreground">
                  No LLM
                </Badge>
              ) : (
                llmLabels.map((label, index) => (
                  <Badge key={`${label}-${index}`} variant="outline" className="shrink-0 font-normal">
                    {label}
                  </Badge>
                ))
              )}
            </div>
            {agent.description ? (
              <p className="line-clamp-2 text-sm text-muted-foreground">{agent.description}</p>
            ) : null}
          </div>
          <div className="flex shrink-0">
            <Button
              variant="ghost"
              size="icon"
              aria-label="Edit agent"
              onClick={openEdit}
            >
              <PencilIcon className="size-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              aria-label="Delete agent"
              onClick={openDelete}
            >
              <TrashIcon className="size-4" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="flex flex-1 flex-col gap-3 pt-0">
        {agent.instructions ? (
          <p className="line-clamp-3 text-xs text-muted-foreground">{agent.instructions}</p>
        ) : null}
      </CardContent>
    </Card>
  );
}
