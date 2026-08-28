import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Spinner } from "@/components/ui/spinner";

import { useDeleteApp } from "../hooks";

import type { App } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type DeleteAppDialogProps = DialogInstanceProps<{
  app: App
}>;

export function DeleteAppDialog({ app, close }: DeleteAppDialogProps) {
  const { mutateAsync, isPending } = useDeleteApp();

  const handleConfirm = async () => {
    await mutateAsync({ code: app.code });
    close();
  };

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Delete app?</DialogTitle>
        <DialogDescription>
          This will permanently delete{" "}
          <span className="font-medium text-foreground">{app.name}</span>.
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
