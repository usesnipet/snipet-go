import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Spinner } from "@/components/ui/spinner";

import { useDeleteAgent } from "../hooks";

import type { Agent } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type DeleteAgentDialogProps = DialogInstanceProps<{
  agent: Agent
}>;

export function DeleteAgentDialog({ agent, close }: DeleteAgentDialogProps) {
  const { mutateAsync, isPending } = useDeleteAgent();

  const handleConfirm = async () => {
    await mutateAsync({ id: agent.id });
    close();
  };

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Delete agent?</DialogTitle>
        <DialogDescription>
          This will permanently delete{" "}
          <span className="font-medium text-foreground">{agent.name}</span>.
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
