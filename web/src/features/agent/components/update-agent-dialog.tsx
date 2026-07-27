import { Button } from "@/components/ui/button";
import {
  DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Form } from "@/components/ui/form";
import { Spinner } from "@/components/ui/spinner";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { useUpdateAgent } from "../hooks";
import { updateAgentSchema } from "../schemas";

import { AgentFormFields } from "./agent-form-fields";

import type { DialogInstanceProps } from "@/lib/dialog";
import type { Agent, UpdateAgent } from "../schemas";

type UpdateAgentDialogProps = DialogInstanceProps<{
  agent: Agent
}>;

export function UpdateAgentDialog({ agent, close }: UpdateAgentDialogProps) {
  const form = useForm<UpdateAgent>({
    resolver: zodResolver(updateAgentSchema),
    defaultValues: {
      name: agent.name,
      description: agent.description,
      instructions: agent.instructions,
      llms: agent.configuration.llms,
      tools: agent.configuration.tools,
    },
  });

  const { mutateAsync, isPending } = useUpdateAgent();

  const onSubmit = form.handleSubmit(async (data) => {
    await mutateAsync({ id: agent.id, data });
    close();
  });

  return (
    <DialogContent className="sm:max-w-lg max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Edit agent</DialogTitle>
        <DialogDescription>
          Update settings for{" "}
          <span className="font-medium text-foreground">{agent.name}</span>.
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
              Save
            </Button>
          </DialogFooter>
        </form>
      </Form>
    </DialogContent>
  );
}
