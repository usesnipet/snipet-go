import { FormInput } from "@/components/form/input";
import { FormTextarea } from "@/components/form/textarea";
import { Button } from "@/components/ui/button";
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { FieldGroup } from "@/components/ui/field";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { useUpdateKnowledge } from "../hooks";
import { updateKnowledgeSchema } from "../schemas";

import type { Knowledge, UpdateKnowledge } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type UpdateKnowledgeDialogProps = DialogInstanceProps<{
  knowledge: Knowledge;
}>;

export function UpdateKnowledgeDialog({ knowledge, close }: UpdateKnowledgeDialogProps) {
  const form = useForm<UpdateKnowledge>({
    resolver: zodResolver(updateKnowledgeSchema),
    defaultValues: {
      name: knowledge.name,
      description: knowledge.description,
    },
  });

  const { mutateAsync, isPending } = useUpdateKnowledge();

  const onSubmit = form.handleSubmit(async (data) => {
    await mutateAsync({ id: knowledge.id, data });
    close();
  });

  return (
    <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Edit Knowledge</DialogTitle>
        <DialogDescription>
          Update name and description for{" "}
          <span className="font-medium text-foreground">{knowledge.name}</span>.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <FieldGroup>
            <FormInput name="name" label="Name" required />
            <FormTextarea
              name="description"
              label="Description"
              placeholder="What this knowledge source contains"
              rows={2}
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
              Save
            </Button>
          </DialogFooter>
        </form>
      </Form>
    </DialogContent>
  );
}
