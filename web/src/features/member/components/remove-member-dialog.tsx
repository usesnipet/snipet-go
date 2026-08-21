import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Spinner } from "@/components/ui/spinner";

import { useRemoveMember } from "../hooks";

import type { DialogInstanceProps } from "@/lib/dialog";
import type { Member } from "@/models/member";

type RemoveMemberDialogProps = DialogInstanceProps<{
  tenantId: string
  member: Member
}>;

export function RemoveMemberDialog({ tenantId, member, close }: RemoveMemberDialogProps) {
  const { mutateAsync, isPending } = useRemoveMember();

  const handleConfirm = async () => {
    await mutateAsync({ tenantId, id: member.id });
    close();
  };

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Remove member?</DialogTitle>
        <DialogDescription>
          This will remove{" "}
          <span className="font-medium text-foreground">{member.user?.name}</span>{" "}
          from this tenant. They will lose access immediately. This action cannot be undone.
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
          Remove
        </Button>
      </DialogFooter>
    </DialogContent>
  )
}
