import { Button } from "@/components/ui/button";
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Spinner } from "@/components/ui/spinner";

import { useDeleteKnowledgeIndex } from "../hooks";

import type { KnowledgeIndex } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type DeleteKnowledgeIndexDialogProps = DialogInstanceProps<{
  knowledgeID: string;
  index: KnowledgeIndex;
}>;

export function DeleteKnowledgeIndexDialog({
  knowledgeID,
  index,
  close,
}: DeleteKnowledgeIndexDialogProps) {
  const { mutateAsync, isPending } = useDeleteKnowledgeIndex();

  const handleConfirm = async () => {
    await mutateAsync({ knowledgeID, id: index.id });
    close();
  };

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Delete index?</DialogTitle>
        <DialogDescription>
          This will permanently delete{" "}
          <span className="font-medium text-foreground">{index.name}</span>.
          This action cannot be undone.
        </DialogDescription>
      </DialogHeader>
      <DialogFooter>
        <DialogClose asChild>
          <Button type="button" variant="outline" disabled={isPending}>
            Cancel
          </Button>
        </DialogClose>
        <Button variant="destructive" disabled={isPending} onClick={handleConfirm}>
          {isPending && <Spinner size="sm" />}
          Delete
        </Button>
      </DialogFooter>
    </DialogContent>
  );
}
