import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Spinner } from "@/components/ui/spinner";
import { useDialog } from "@/lib/dialog";
import { MoreHorizontal, PencilIcon, PlusIcon, Trash2 } from "lucide-react";

import { useListKnowledgeIndexes } from "../hooks";

import { CreateKnowledgeIndexDialog } from "./create-knowledge-index-dialog";
import { DeleteKnowledgeIndexDialog } from "./delete-knowledge-index-dialog";
import { UpdateKnowledgeIndexDialog } from "./update-knowledge-index-dialog";

import type { KnowledgeIndex } from "../schemas";

export function KnowledgeIndexesSection({
  knowledgeID,
}: {
  knowledgeID: string;
}) {
  const { openDialog } = useDialog();
  const { data, isLoading } = useListKnowledgeIndexes(knowledgeID);
  const indexes = data?.data ?? [];

  const openCreate = () => {
    openDialog({
      component: CreateKnowledgeIndexDialog,
      props: { knowledgeID },
    });
  };

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between gap-2">
          <CardTitle className="text-base">Indexes</CardTitle>
          <Button size="sm" onClick={openCreate}>
            <PlusIcon />
            Add index
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="flex justify-center py-6">
            <Spinner size="sm" />
          </div>
        ) : indexes.length === 0 ? (
          <p className="text-sm italic text-muted-foreground">
            No indexes yet. Add one to start indexing this knowledge source.
          </p>
        ) : (
          <ul className="flex flex-col gap-2">
            {indexes.map((index) => (
              <KnowledgeIndexRow
                key={index.id}
                knowledgeID={knowledgeID}
                index={index}
              />
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
}

function KnowledgeIndexRow({
  knowledgeID,
  index,
}: {
  knowledgeID: string;
  index: KnowledgeIndex;
}) {
  const { openDialog } = useDialog();

  const openEdit = () => {
    openDialog({
      component: UpdateKnowledgeIndexDialog,
      props: { knowledgeID, index },
    });
  };

  const openDelete = () => {
    openDialog({
      component: DeleteKnowledgeIndexDialog,
      props: { knowledgeID, index },
    });
  };

  return (
    <li className="flex items-center justify-between gap-2 rounded-md border border-border px-3 py-2">
      <div className="flex items-center gap-2">
        <span className="text-sm font-medium">{index.name}</span>
        <Badge variant="secondary" className="font-normal">
          {index.driver}
        </Badge>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon" aria-label="Index actions">
            <MoreHorizontal className="size-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
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
    </li>
  );
}
