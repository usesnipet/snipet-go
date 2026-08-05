import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { truncate } from "@/lib/utils";

import type { KnowledgeItem } from "../schemas";
export function KnowledgeItemsTableNameField({ item }: { item: KnowledgeItem }) {

  return (
    <>
      <span className="font-medium">
        {truncate(item.name, 40)}
      </span>
      <div className="flex items-center gap-2">
        <Tooltip>
          <TooltipTrigger>
            <span className="font-mono text-xs text-muted-foreground">
              {truncate(item.external_id, 50)}
            </span>
          </TooltipTrigger>
          <TooltipContent>
            {item.external_id}
          </TooltipContent>
        </Tooltip>
      </div>
    </>
  )

}