import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Spinner } from "@/components/ui/spinner";

import { useRemoveInvitation } from "../hooks";

import type { Invitation } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type CancelInvitationDialogProps = DialogInstanceProps<{
  tenantId: string
  invitation: Invitation
}>;

export function CancelInvitationDialog({ tenantId, invitation, close }: CancelInvitationDialogProps) {
  const { mutateAsync, isPending } = useRemoveInvitation();

  const handleConfirm = async () => {
    await mutateAsync({ tenantId, id: invitation.id });
    close();
  };

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Cancel invitation?</DialogTitle>
        <DialogDescription>
          This will cancel the pending invitation sent to{" "}
          <span className="font-medium text-foreground">{invitation.email}</span>.
          They will no longer be able to accept it.
        </DialogDescription>
      </DialogHeader>
      <DialogFooter>
        <DialogClose asChild>
          <Button type="button" variant="outline" disabled={isPending}>
            Keep invitation
          </Button>
        </DialogClose>
        <Button
          variant="destructive"
          disabled={isPending}
          onClick={handleConfirm}
        >
          {isPending && <Spinner size="sm" />}
          Cancel invitation
        </Button>
      </DialogFooter>
    </DialogContent>
  )
}
