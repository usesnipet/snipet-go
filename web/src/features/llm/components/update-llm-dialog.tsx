import { Button } from "@/components/ui/button";
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { useUpdateLlm } from "../hooks";
import { createLlmSchema } from "../schemas";

import { LlmFormFields } from "./llm-form-fields";

import type { CreateLlm, Llm } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type UpdateLlmDialogProps = DialogInstanceProps<{
  llm: Llm;
}>;

export function UpdateLlmDialog({ llm, close }: UpdateLlmDialogProps) {
  const form = useForm<CreateLlm>({
    resolver: zodResolver(createLlmSchema),
    defaultValues: {
      name: llm.name,
      provider: llm.provider,
      configuration: llm.configuration,
    },
  });

  const { mutateAsync, isPending } = useUpdateLlm();

  const onSubmit = form.handleSubmit(async (data) => {
    await mutateAsync({ tenantId, id: llm.id, data });
    close();
  });

  return (
    <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Edit LLM</DialogTitle>
        <DialogDescription>
          Update settings for{" "}
          <span className="font-medium text-foreground">{llm.name}</span>.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <LlmFormFields />
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
  );
}
