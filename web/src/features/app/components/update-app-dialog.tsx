import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useUpdateApp, useUpdateAppAuthConfig } from "../hooks";
import { updateAppAuthConfigSchema, updateAppSchema } from "../schemas";

import { AppAuthConfigFields } from "./app-auth-config-fields";
import { AppFormFields } from "./app-form-fields";

import type { App } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

const formSchema = updateAppSchema.and(updateAppAuthConfigSchema);
type FormValues = z.infer<typeof formSchema>;

type UpdateAppDialogProps = DialogInstanceProps<{
  app: App
}>;

export function UpdateAppDialog({ app, close }: UpdateAppDialogProps) {
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: app.name,
      description: app.description,
      public: app.public,
      auth_config: app.auth_config,
    },
  });

  const { mutateAsync: updateApp, isPending: isUpdatingApp } = useUpdateApp();
  const { mutateAsync: updateAuthConfig, isPending: isUpdatingAuthConfig } = useUpdateAppAuthConfig();
  const isPending = isUpdatingApp || isUpdatingAuthConfig;

  const onSubmit = form.handleSubmit(async ({ auth_config: app_config, ...data }) => {
    await Promise.all([
      updateApp({ code: app.code, data }),
      updateAuthConfig({ code: app.code, data: { auth_config: app_config } }),
    ]);
    close();
  });

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Edit app</DialogTitle>
        <DialogDescription>
          Update the settings for{" "}
          <span className="font-medium text-foreground">{app.name}</span>.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <AppFormFields />
          <AppAuthConfigFields />
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
