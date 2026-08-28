import { Button } from "@/components/ui/button";
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Link } from "@/components/ui/link";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useListLlm } from "@/features/llm/hooks";
import { ROUTES } from "@/routes";
import { ArrowDownIcon, ArrowUpIcon, XIcon } from "lucide-react";
import { useState } from "react";
import { useFormContext } from "react-hook-form";

import type { Llm } from "@/features/llm/schemas";
type AgentLlmListProps = {
  name?: string
  label?: string
}

function moveItem(ids: string[], index: number, direction: -1 | 1): string[] {
  const nextIndex = index + direction;
  if (nextIndex < 0 || nextIndex >= ids.length) return ids;
  const next = [...ids];
  const [item] = next.splice(index, 1);
  next.splice(nextIndex, 0, item);
  return next;
}

export function AgentLlmList({
  name = "llm_ids",
  label = "LLMs",
}: AgentLlmListProps) {
  const form = useFormContext();
  const { data, isLoading } = useListLlm();
  const llms = data?.data ?? [];
  const llmById = new Map(llms.map((llm) => [llm.id, llm]));
  const [selectKey, setSelectKey] = useState(0);
  const llmsHref = ROUTES.llms;

  return (
    <FormField
      control={form.control}
      name={name}
      render={({ field }) => {
        const selected: string[] = field.value ?? [];
        const available = llms.filter((llm) => !selected.includes(llm.id));
        const disabled = form.formState.isSubmitting || isLoading;

        const addLlm = (id: string) => {
          if (!id || selected.includes(id)) return;
          field.onChange([...selected, id]);
          setSelectKey((key) => key + 1);
        };

        const removeLlm = (id: string) => {
          field.onChange(selected.filter((value) => value !== id));
        };

        return (
          <FormItem>
            <FormLabel>{label}</FormLabel>
            <FormControl>
              <div className="flex flex-col gap-3">
                {isLoading ? (
                  <p className="text-sm text-muted-foreground">Loading LLMs...</p>
                ) : llms.length === 0 ? (
                  <p className="text-sm text-muted-foreground">
                    No LLMs configured yet.{" "}
                    <Link href={llmsHref} className="underline underline-offset-2">
                      Create one in LLMs
                    </Link>
                    .
                  </p>
                ) : (
                  <>
                    <Select
                      key={selectKey}
                      disabled={disabled || available.length === 0}
                      onValueChange={addLlm}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue
                          placeholder={
                            available.length === 0
                              ? "All LLMs selected"
                              : "Select an LLM to add"
                          }
                        />
                      </SelectTrigger>
                      <SelectContent>
                        {available.map((llm) => (
                          <SelectItem key={llm.id} value={llm.id}>
                            {llm.name}
                            {" · "}
                            {llm.provider}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>

                    {selected.length > 0 ? (
                      <ul className="flex flex-col gap-2">
                        {selected.map((id, index) => {
                          const llm = llmById.get(id);
                          return (
                            <LLMListItem
                              key={id}
                              llm={llm}
                              fallbackId={id}
                              index={index}
                              total={selected.length}
                              disabled={disabled}
                              onMoveUp={() => field.onChange(moveItem(selected, index, -1))}
                              onMoveDown={() => field.onChange(moveItem(selected, index, 1))}
                              onRemove={() => removeLlm(id)}
                            />
                          );
                        })}
                      </ul>
                    ) : (
                      <p className="text-sm text-muted-foreground">
                        Add at least one LLM. Order sets priority (first = highest).
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

type LLMListItemProps = {
  llm?: Llm
  fallbackId: string
  index: number
  total: number
  disabled: boolean
  onMoveUp: () => void
  onMoveDown: () => void
  onRemove: () => void
}

function LLMListItem({
  llm,
  fallbackId,
  index,
  total,
  disabled,
  onMoveUp,
  onMoveDown,
  onRemove,
}: LLMListItemProps) {
  return (
    <li className="flex items-center gap-2 rounded-md border border-border px-3 py-2">
      <span className="flex size-6 shrink-0 items-center justify-center rounded-full bg-muted text-xs font-medium text-muted-foreground">
        {index + 1}
      </span>
      <span className="flex min-w-0 flex-1 flex-col gap-0.5">
        <span className="truncate text-sm font-medium leading-none">
          {llm?.name ?? fallbackId}
        </span>
        {llm ? (
          <span className="text-xs text-muted-foreground">{llm.provider}</span>
        ) : null}
      </span>
      <div className="flex shrink-0 items-center">
        <Button
          type="button"
          variant="ghost"
          size="icon"
          aria-label="Move up"
          disabled={disabled || index === 0}
          onClick={onMoveUp}
        >
          <ArrowUpIcon className="size-4" />
        </Button>
        <Button
          type="button"
          variant="ghost"
          size="icon"
          aria-label="Move down"
          disabled={disabled || index === total - 1}
          onClick={onMoveDown}
        >
          <ArrowDownIcon className="size-4" />
        </Button>
        <Button
          type="button"
          variant="ghost"
          size="icon"
          aria-label="Remove LLM"
          disabled={disabled}
          onClick={onRemove}
        >
          <XIcon className="size-4" />
        </Button>
      </div>
    </li>
  );
}
