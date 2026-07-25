import { FormInput } from "@/components/form/input";
import { Button } from "@/components/ui/button";
import {
  Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle,
  DialogTrigger
} from "@/components/ui/dialog";
import { FieldGroup } from "@/components/ui/field";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useCreateApiKey } from "../hooks";

import type { ApiKeyWithSecret } from "../schemas";

const formSchema = z.object({
  name: z.string().min(1).max(255),
  expires_at: z.string().optional(),
});

type FormValues = z.infer<typeof formSchema>;

type CreateApiKeyDialogProps = {
  onCreated: (apiKey: ApiKeyWithSecret) => void
  children: React.ReactNode
}

export function CreateApiKeyDialog({
  onCreated,
  children,
}: CreateApiKeyDialogProps) {
  const [open, setOpen] = useState(false);
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: { name: "", expires_at: "" },
  });

  const { mutateAsync, isPending } = useCreateApiKey();

  const onSubmit = form.handleSubmit(async (values) => {
    const result = await mutateAsync({
      name: values.name,
      expires_at: values.expires_at ? new Date(values.expires_at) : undefined,
    });
    form.reset();
    setOpen(false);
    onCreated(result);
  });

  return (
    <Dialog
      open={open}
      onOpenChange={(next) => {
        if (!next) form.reset();
        setOpen(next);
      }}
    >
      <DialogTrigger asChild>
        {children}
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Create API Key</DialogTitle>
          <DialogDescription>
            Create a new API key. The secret will be shown only once.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <FieldGroup>
              <FormInput
                name="name"
                label="Name"
                placeholder="Production"
                required
              />
              <FormInput
                name="expires_at"
                label="Expiration"
                type="datetime-local"
                description="Optional. Leave empty for no expiration."
              />
            </FieldGroup>
            <DialogFooter>
              <DialogClose asChild>
                <Button type="button" variant="outline" disabled={isPending}>
                  Cancel
                </Button>
              </DialogClose>
              <Button type="submit" disabled={isPending}>
                {isPending && <Spinner size="sm" />}
                Create
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
