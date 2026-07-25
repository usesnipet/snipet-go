import { FormInput } from "@/components/form/input";
import { Button } from "@/components/ui/button";
import {
  Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { FieldGroup } from "@/components/ui/field";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { useUpdateExpirationApiKey } from "../hooks";

import type { ApiKey } from "../schemas";

const formSchema = z.object({
  expires_at: z.string().optional(),
});

type FormValues = z.infer<typeof formSchema>;

function toDatetimeLocalValue(date?: Date | null) {
  if (!date) return "";
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`;
}

type UpdateApiKeyExpirationDialogProps = {
  apiKey: ApiKey | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function UpdateApiKeyExpirationDialog({
  apiKey,
  open,
  onOpenChange,
}: UpdateApiKeyExpirationDialogProps) {
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: { expires_at: "" },
  });

  const { mutateAsync, isPending } = useUpdateExpirationApiKey();

  useEffect(() => {
    if (open && apiKey) {
      form.reset({ expires_at: toDatetimeLocalValue(apiKey.expires_at) });
    }
  }, [apiKey, form, open]);

  const onSubmit = form.handleSubmit(async (values) => {
    if (!apiKey) return;
    await mutateAsync({
      id: apiKey.id,
      data: {
        expires_at: values.expires_at ? new Date(values.expires_at) : undefined,
      },
    });
    onOpenChange(false);
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Update expiration</DialogTitle>
          <DialogDescription>
            Set a new expiration date for{" "}
            <span className="font-medium text-foreground">{apiKey?.name}</span>.
            Leave empty to remove expiration.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <FieldGroup>
              <FormInput
                name="expires_at"
                label="Expiration"
                type="datetime-local"
              />
            </FieldGroup>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isPending}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending && <Spinner size="sm" />}
                Save
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
