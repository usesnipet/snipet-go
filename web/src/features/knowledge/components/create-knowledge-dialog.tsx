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

import { useCreateKnowledge } from "../hooks";
import { createKnowledgeSchema } from "../schemas";

import { KnowledgeFormFields } from "./knowledge-form-fields";

import type { CreateKnowledge, CreateKnowledgeResponse } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type CreateKnowledgeDialogProps = DialogInstanceProps<{
  onCreated?: (result: CreateKnowledgeResponse) => void;
}>;

const defaultValues: CreateKnowledge = {
  name: "",
  description: "",
  driver: "",
  configuration: {},
};

export function CreateKnowledgeDialog({ onCreated, close }: CreateKnowledgeDialogProps) {
  const form = useForm<CreateKnowledge>({
    resolver: zodResolver(createKnowledgeSchema),
    defaultValues,
  });

  const { mutateAsync, isPending } = useCreateKnowledge();

  const onSubmit = form.handleSubmit(async (values) => {
    const result = await mutateAsync({ data: values });
    form.reset();
    onCreated?.(result);
    close();
  });

  return (
    <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Create Knowledge</DialogTitle>
        <DialogDescription>
          Add a knowledge source and configure its driver.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <KnowledgeFormFields />
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
