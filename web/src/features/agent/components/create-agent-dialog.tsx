import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { useCreateAgent } from "../hooks";
import { createAgentSchema } from "../schemas";

import { AgentFormFields } from "./agent-form-fields";

import type { CreateAgent, Agent } from "../schemas";

import type { DialogInstanceProps } from "@/lib/dialog";

type CreateAgentDialogProps = DialogInstanceProps<{
  tenantId: string
  onCreated?: (agent: Agent) => void
}>;

const defaultValues: CreateAgent = {
  name: "",
  description: "",
  instructions: "",
  llm_ids: [],
};

export function CreateAgentDialog({ tenantId, onCreated, close }: CreateAgentDialogProps) {
  const form = useForm<CreateAgent>({
    resolver: zodResolver(createAgentSchema),
    defaultValues,
  });

  const { mutateAsync, isPending } = useCreateAgent();

  const onSubmit = form.handleSubmit(async (values) => {
    const result = await mutateAsync({ tenantId, data: values });
    form.reset();
    onCreated?.(result);
    close();
  });

  return (
    <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Create agent</DialogTitle>
        <DialogDescription>
          Define the agent persona and the language models it will use.
        </DialogDescription>
      </DialogHeader>
      <Form {...form}>
        <form onSubmit={onSubmit} className="flex flex-col gap-4">
          <AgentFormFields />
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
