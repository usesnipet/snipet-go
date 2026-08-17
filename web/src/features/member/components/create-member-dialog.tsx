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

import { useCreateMember } from "../hooks";
import { createMemberSchema } from "../schemas";

import type { CreateMember } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

const roleOptions = [
  { label: "Admin", value: "admin" },
  { label: "Member", value: "user" },
];

type CreateMemberDialogProps = DialogInstanceProps<{
  tenantId: string
}>;

export function CreateMemberDialog({ tenantId, close }: CreateMemberDialogProps) {
  const form = useForm<CreateMember>({
    resolver: zodResolver(createMemberSchema),
    defaultValues: { name: "", email: "", password: "", confirm_password: "", role: "user" },
  });

  const { mutateAsync, isPending } = useCreateMember();

  const onSubmit = form.handleSubmit(async (values) => {
    await mutateAsync({ tenantId, data: values });
    close();
  });

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Create member</DialogTitle>
        <DialogDescription>
          Create a user account directly in this tenant. They can sign in immediately with the password you set.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <FieldGroup>
            <FormInput
              name="name"
              label="Name"
              placeholder="Jane Doe"
              required
            />
            <FormInput
              name="email"
              label="Email"
              type="email"
              placeholder="teammate@company.com"
              required
            />
            <FormInput
              name="password"
              label="Password"
              type="password"
              placeholder="At least 8 characters"
              required
            />
            <FormInput
              name="confirm_password"
              label="Confirm password"
              type="password"
              placeholder="Repeat the password"
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
              Create member
            </Button>
          </DialogFooter>
        </form>
      </Form>
    </DialogContent>
  )
}
