import { FormInput } from "@/components/form/input";
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

import { useUpdateKnowledgeIndex } from "../hooks";
import { updateKnowledgeIndexSchema } from "../schemas";

import type { KnowledgeIndex, UpdateKnowledgeIndex } from "../schemas";
import type { DialogInstanceProps } from "@/lib/dialog";

type UpdateKnowledgeIndexDialogProps = DialogInstanceProps<{
  tenantId: string;
  knowledgeID: string;
  index: KnowledgeIndex;
}>;

export function UpdateKnowledgeIndexDialog({
  tenantId,
  knowledgeID,
  index,
  close,
}: UpdateKnowledgeIndexDialogProps) {
  const form = useForm<UpdateKnowledgeIndex>({
    resolver: zodResolver(updateKnowledgeIndexSchema),
    defaultValues: {
      name: index.name,
    },
  });

  const { mutateAsync, isPending } = useUpdateKnowledgeIndex();

  const onSubmit = form.handleSubmit(async (data) => {
    await mutateAsync({ tenantId, knowledgeID, id: index.id, data });
    close();
  });

  return (
    <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Edit Index</DialogTitle>
        <DialogDescription>
          Update name for{" "}
          <span className="font-medium text-foreground">{index.name}</span>.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <FieldGroup>
            <FormInput name="name" label="Name" required />
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
