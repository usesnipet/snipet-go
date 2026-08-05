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

import { useCreateLlm } from "../hooks";
import { createLlmSchema } from "../schemas";

import { LlmFormFields } from "./llm-form-fields";

import type { CreateLlm, Llm } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type CreateLlmDialogProps = DialogInstanceProps<{
  onCreated?: (llm: Llm) => void;
}>;

const defaultValues: CreateLlm = {
  name: "",
  provider: "",
  configuration: {},
};

export function CreateLlmDialog({ onCreated, close }: CreateLlmDialogProps) {
  const form = useForm<CreateLlm>({
    resolver: zodResolver(createLlmSchema),
    defaultValues,
  });

  const { mutateAsync, isPending } = useCreateLlm();

  const onSubmit = form.handleSubmit(async (values) => {
    const result = await mutateAsync(values);
    form.reset();
    onCreated?.(result);
    close();
  });

  return (
    <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Create LLM</DialogTitle>
        <DialogDescription>
          Add a named language model provider configuration.
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
              Create
            </Button>
          </DialogFooter>
        </form>
      </Form>
    </DialogContent>
  );
}
