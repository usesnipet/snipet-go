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

import { useDeleteLlm } from "../hooks";

import type { Llm } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type DeleteLlmDialogProps = DialogInstanceProps<{
  llm: Llm;
}>;

export function DeleteLlmDialog({ llm, close }: DeleteLlmDialogProps) {
  const { mutateAsync, isPending } = useDeleteLlm();

  const handleConfirm = async () => {
    await mutateAsync({ id: llm.id });
    close();
  };

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Delete LLM?</DialogTitle>
        <DialogDescription>
          This will permanently delete{" "}
          <span className="font-medium text-foreground">{llm.name}</span>.
          This action cannot be undone.
        </DialogDescription>
      </DialogHeader>
      <DialogFooter>
        <DialogClose asChild>
          <Button type="button" variant="outline" disabled={isPending}>
            Cancel
          </Button>
        </DialogClose>
        <Button
          variant="destructive"
          disabled={isPending}
          onClick={handleConfirm}
        >
          {isPending && <Spinner size="sm" />}
          Delete
        </Button>
      </DialogFooter>
    </DialogContent>
  );
}
