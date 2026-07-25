import {
  AlertDialog, AlertDialogCancel, AlertDialogContent, AlertDialogDescription,
  AlertDialogFooter, AlertDialogHeader, AlertDialogTitle
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";

import { useRollApiKey } from "../hooks";

import type { ApiKey, ApiKeyWithSecret } from "../schemas";

type RollApiKeyDialogProps = {
  apiKey: ApiKey | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onRolled: (apiKey: ApiKeyWithSecret) => void
}

export function RollApiKeyDialog({
  apiKey,
  open,
  onOpenChange,
  onRolled,
}: RollApiKeyDialogProps) {
  const { mutateAsync, isPending } = useRollApiKey();

  const handleConfirm = async () => {
    if (!apiKey) return;
    const result = await mutateAsync(apiKey.id);
    onOpenChange(false);
    onRolled(result);
  };

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Roll API key?</AlertDialogTitle>
          <AlertDialogDescription>
            This will generate a new secret for{" "}
            <span className="font-medium text-foreground">{apiKey?.name}</span>.
            The previous secret will stop working immediately.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isPending}>Cancel</AlertDialogCancel>
          <Button
            variant="destructive"
            disabled={isPending}
            onClick={handleConfirm}
          >
            {isPending && <Spinner size="sm" />}
            Roll key
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
