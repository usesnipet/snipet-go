import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Spinner } from "@/components/ui/spinner";

import { useDeleteClient } from "../hooks";

import type { Client } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type DeleteClientDialogProps = DialogInstanceProps<{
  tenantId: string
  client: Client
}>;

export function DeleteClientDialog({ tenantId, client, close }: DeleteClientDialogProps) {
  const { mutateAsync, isPending } = useDeleteClient();

  const handleConfirm = async () => {
    await mutateAsync({ tenantId, code: client.code });
    close();
  };

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Delete client?</DialogTitle>
        <DialogDescription>
          This will permanently delete{" "}
          <span className="font-medium text-foreground">{client.name}</span>.
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
  )
}
