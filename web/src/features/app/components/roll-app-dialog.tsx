import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Spinner } from "@/components/ui/spinner";

import { useRollApp } from "../hooks";

import type { App, AppWithSecret } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type RollAppDialogProps = DialogInstanceProps<{
  app: App
  onRolled: (app: AppWithSecret) => void
}>;

export function RollAppDialog({ app, onRolled, close }: RollAppDialogProps) {
  const { mutateAsync, isPending } = useRollApp();

  const handleConfirm = async () => {
    const result = await mutateAsync({ code: app.code });
    close();
    onRolled(result);
  };

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Roll app key?</DialogTitle>
        <DialogDescription>
          This will generate a new secret for{" "}
          <span className="font-medium text-foreground">{app.name}</span>.
          The previous secret will stop working immediately.
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
          Roll key
        </Button>
      </DialogFooter>
    </DialogContent>
  )
}
