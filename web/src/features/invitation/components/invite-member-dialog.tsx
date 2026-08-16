import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { FormInput } from "@/components/form/input";
import { FormSelect } from "@/components/form/select";
import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { FieldGroup } from "@/components/ui/field";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";

import { useCreateInvitation } from "../hooks";
import { createInvitationSchema } from "../schemas";

import type { CreateInvitation } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

const roleOptions = [
  { label: "Admin", value: "admin" },
  { label: "Member", value: "user" },
];

type InviteMemberDialogProps = DialogInstanceProps<{
  tenantId: string
}>;

export function InviteMemberDialog({ tenantId, close }: InviteMemberDialogProps) {
  const form = useForm<CreateInvitation>({
    resolver: zodResolver(createInvitationSchema),
    defaultValues: { email: "", role: "user" },
  });

  const { mutateAsync, isPending } = useCreateInvitation();

  const onSubmit = form.handleSubmit(async (values) => {
    await mutateAsync({ tenantId, data: values });
    close();
  });

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Invite member</DialogTitle>
        <DialogDescription>
          Send an invitation to join this tenant. They will receive an email with a link to accept.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <FieldGroup>
            <FormInput
              name="email"
              label="Email"
              type="email"
              placeholder="teammate@company.com"
              required
            />
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
              Send invitation
            </Button>
          </DialogFooter>
        </form>
      </Form>
    </DialogContent>
  )
}
