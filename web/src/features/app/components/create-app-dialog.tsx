import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { useCreateApp } from "../hooks";
import { createAppSchema } from "../schemas";

import { AppFormFields } from "./app-form-fields";

import type { AppWithSecret, CreateApp } from "../schemas";

import type { DialogInstanceProps } from "@/lib/dialog";

type CreateAppDialogProps = DialogInstanceProps<{
  tenantId: string
  onCreated?: (app: AppWithSecret) => void
}>;

export function CreateAppDialog({ tenantId, onCreated, close }: CreateAppDialogProps) {
  const form = useForm<CreateApp>({
    resolver: zodResolver(createAppSchema),
    defaultValues: {
      name: "",
      description: "",
      public: false,
    },
  });

  const { mutateAsync, isPending } = useCreateApp();

  const onSubmit = form.handleSubmit(async (values) => {
    const result = await mutateAsync({ tenantId, data: values });
    form.reset();
    onCreated?.(result);
    close();
  });

  return (
    <DialogContent className="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Create app</DialogTitle>
        <DialogDescription>
          Create a new app. You can configure how it authenticates its users afterwards.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <AppFormFields />
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
