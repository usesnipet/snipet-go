import { FormSelect } from "@/components/form/select";
import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { FieldGroup } from "@/components/ui/field";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { useUpdateRoleMember } from "../hooks";
import { updateMemberRoleSchema } from "../schemas";

import type { Member, UpdateMemberRole } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

const roleOptions = [
  { label: "Admin", value: "admin" },
  { label: "Member", value: "user" },
];

type UpdateMemberRoleDialogProps = DialogInstanceProps<{
  tenantId: string
  member: Member
}>;

export function UpdateMemberRoleDialog({ tenantId, member, close }: UpdateMemberRoleDialogProps) {
  const form = useForm<UpdateMemberRole>({
    resolver: zodResolver(updateMemberRoleSchema),
    defaultValues: { role: member.role },
  });

  const { mutateAsync, isPending } = useUpdateRoleMember();

  const onSubmit = form.handleSubmit(async (values) => {
    await mutateAsync({ tenantId, id: member.id, data: values });
    close();
  });

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Update member role</DialogTitle>
        <DialogDescription>
          Change the role of{" "}
          <span className="font-medium text-foreground">{member.user?.name}</span>{" "}
          within this tenant.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <FieldGroup>
            <FormSelect name="role" label="Role" options={roleOptions} />
          </FieldGroup>
          <DialogFooter>
            <DialogClose asChild>
              <Button type="button" variant="outline" disabled={isPending}>
                Cancel
              </Button>
            </DialogClose>
            <Button type="submit" disabled={isPending}>
              {isPending && <Spinner size="sm" />}
              Save
            </Button>
          </DialogFooter>
        </form>
      </Form>
    </DialogContent>
  )
}
