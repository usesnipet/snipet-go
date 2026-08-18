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

import { useCreateKnowledgeIndex } from "../hooks";
import { createKnowledgeIndexSchema } from "../schemas";

import { KnowledgeIndexFormFields } from "./knowledge-index-form-fields";

import type { CreateKnowledgeIndex, KnowledgeIndex } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type CreateKnowledgeIndexDialogProps = DialogInstanceProps<{
  tenantId: string;
  knowledgeID: string;
  onCreated?: (result: KnowledgeIndex) => void;
}>;

const defaultValues: CreateKnowledgeIndex = {
  name: "",
  driver: "",
  configuration: {},
};

export function CreateKnowledgeIndexDialog({
  tenantId,
  knowledgeID,
  onCreated,
  close,
}: CreateKnowledgeIndexDialogProps) {
  const form = useForm<CreateKnowledgeIndex>({
    resolver: zodResolver(createKnowledgeIndexSchema),
    defaultValues,
  });

  const { mutateAsync, isPending } = useCreateKnowledgeIndex();

  const onSubmit = form.handleSubmit(async (values) => {
    const result = await mutateAsync({ tenantId, knowledgeID, data: values });
    form.reset();
    onCreated?.(result);
    close();
  });

  return (
    <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Create Index</DialogTitle>
        <DialogDescription>
          Add an index to this knowledge source and configure its driver.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <KnowledgeIndexFormFields />
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
