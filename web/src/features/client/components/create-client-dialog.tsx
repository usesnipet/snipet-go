import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { useCreateClient } from "../hooks";
import { createClientSchema } from "../schemas";

import { ClientFormFields } from "./client-form-fields";

import type { Client, CreateClient } from "../schemas";

import type { DialogInstanceProps } from "@/lib/dialog";

type CreateClientDialogProps = DialogInstanceProps<{
  onCreated?: (client: Client) => void
}>;

export function CreateClientDialog({ onCreated, close }: CreateClientDialogProps) {
  const form = useForm<CreateClient>({
    resolver: zodResolver(createClientSchema),
    defaultValues: {
      name: "",
      config: {
        anonymous: {
          enabled: false,
        },
        oidc: {
          enabled: false,
          issuer: "",
          audience: "",
        },
        webhook: {
          enabled: false,
          url: "",
        },
      }
    },
  });

  const { mutateAsync, isPending } = useCreateClient();

  const onSubmit = form.handleSubmit(async (values) => {
    const result = await mutateAsync(values);
    form.reset();
    onCreated?.(result);
    close();
  });

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Create client</DialogTitle>
        <DialogDescription>
          Create a new client and configure how requests are authenticated.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <ClientFormFields />
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
  )
}
