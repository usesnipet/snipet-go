import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { useUpdateClient } from "../hooks";
import { updateClientSchema } from "../schemas";

import { ClientFormFields } from "./client-form-fields";

import type { Client, UpdateClient } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type UpdateClientDialogProps = DialogInstanceProps<{
  client: Client
}>;

export function UpdateClientDialog({ client, close }: UpdateClientDialogProps) {
  const form = useForm<UpdateClient>({
    resolver: zodResolver(updateClientSchema),
    defaultValues: {
      name: client.name,
      config: {
        anonymous: {
          enabled: client.config.anonymous.enabled,
        },
        oidc: {
          enabled: client.config.oidc.enabled,
          issuer: client.config.oidc.issuer,
          audience: client.config.oidc.audience,
        },
        webhook: {
          enabled: client.config.webhook.enabled,
          url: client.config.webhook.url,
        }
      },
    },
  });

  const { mutateAsync, isPending } = useUpdateClient();

  const onSubmit = form.handleSubmit(async (values) => {
    await mutateAsync({ code: client.code, data: values });
    close();
  });

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Edit client</DialogTitle>
        <DialogDescription>
          Update the settings for{" "}
          <span className="font-medium text-foreground">{client.name}</span>.
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
              Save
            </Button>
          </DialogFooter>
        </form>
      </Form>
    </DialogContent>
  )
}
