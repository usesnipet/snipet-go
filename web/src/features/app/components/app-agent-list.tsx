import { Button } from "@/components/ui/button";
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Link } from "@/components/ui/link";
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue
} from "@/components/ui/select";
import { useListAgent } from "@/features/agent/hooks";
import { ROUTES } from "@/routes";
import { XIcon } from "lucide-react";
import { useState } from "react";
import { useFormContext } from "react-hook-form";

type AppAgentListProps = {
  name?: string
  label?: string
}

export function AppAgentList({
  name = "agent_ids",
  label = "Agents",
}: AppAgentListProps) {
  const form = useFormContext();
  const { data, isLoading } = useListAgent();
  const agents = data?.data ?? [];
  const agentById = new Map(agents.map((agent) => [agent.id, agent]));
  const [selectKey, setSelectKey] = useState(0);

  return (
    <FormField
      control={form.control}
      name={name}
      render={({ field }) => {
        const selected: string[] = field.value ?? [];
        const available = agents.filter((agent) => !selected.includes(agent.id));
        const disabled = form.formState.isSubmitting || isLoading;

        const addAgent = (id: string) => {
          if (!id || selected.includes(id)) return;
          field.onChange([...selected, id]);
          setSelectKey((key) => key + 1);
        };

        const removeAgent = (id: string) => {
          field.onChange(selected.filter((value) => value !== id));
        };

        return (
          <FormItem>
            <FormLabel>{label}</FormLabel>
            <FormControl>
              <div className="flex flex-col gap-3">
                {isLoading ? (
                  <p className="text-sm text-muted-foreground">Loading agents...</p>
                ) : agents.length === 0 ? (
                  <p className="text-sm text-muted-foreground">
                    No agents created yet.{" "}
                    <Link href={ROUTES.agent} className="underline underline-offset-2">
                      Create one in Agents
                    </Link>
                    .
                  </p>
                ) : (
                  <>
                    <Select
                      key={selectKey}
                      disabled={disabled || available.length === 0}
                      onValueChange={addAgent}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue
                          placeholder={
                            available.length === 0
                              ? "All agents selected"
                              : "Select an agent to add"
                          }
                        />
                      </SelectTrigger>
                      <SelectContent>
                        {available.map((agent) => (
                          <SelectItem key={agent.id} value={agent.id}>
                            {agent.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>

                    {selected.length > 0 ? (
                      <ul className="flex flex-col gap-2">
                        {selected.map((id) => {
                          const agent = agentById.get(id);
                          return (
                            <li
                              key={id}
                              className="flex items-center gap-2 rounded-md border border-border px-3 py-2"
                            >
                              <span className="flex min-w-0 flex-1 flex-col gap-0.5">
                                <span className="truncate text-sm font-medium leading-none">
                                  {agent?.name ?? id}
                                </span>
                                {agent?.description ? (
                                  <span className="truncate text-xs text-muted-foreground">
                                    {agent.description}
                                  </span>
                                ) : null}
                              </span>
                              <Button
                                type="button"
                                variant="ghost"
                                size="icon"
                                aria-label="Remove agent"
                                disabled={disabled}
                                onClick={() => removeAgent(id)}
                              >
                                <XIcon className="size-4" />
                              </Button>
                            </li>
                          );
                        })}
                      </ul>
                    ) : (
                      <p className="text-sm text-muted-foreground">
                        This app is not linked to any agent yet.
                      </p>
                    )}
                  </>
                )}
              </div>
            </FormControl>
            <FormMessage />
          </FormItem>
        );
      }}
    />
  );
}
