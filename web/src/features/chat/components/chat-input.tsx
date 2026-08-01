"use client";

import { Button } from "@/components/ui/button";
import {
  PromptInput, PromptInputAction, PromptInputActions, PromptInputTextarea
} from "@/components/ui/prompt-input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useListClientAgents } from "@/features/client/hooks";
import { cn } from "@/lib/utils";
import { ArrowUp, Square } from "lucide-react";
import * as React from "react";

export type ChatInputSubmit = {
  message: string;
  agentId: string;
};

export type Props = {
  clientCode: string;
  containerClassName?: string;
  placeholder?: string;
  disabled?: boolean;
  isLoading?: boolean;
  onSubmit?: (value: ChatInputSubmit) => void;
  onStop?: () => void;
  className?: string;
};

function ChatInput({
  clientCode,
  containerClassName,
  placeholder = "Ask me anything...",
  disabled,
  isLoading = false,
  onSubmit,
  onStop,
  className,
}: Props) {
  const [value, setValue] = React.useState("");
  const { data: agentsPage, isLoading: isLoadingAgents } =
    useListClientAgents(clientCode);
  const agents = agentsPage?.data ?? [];
  const [agentId, setAgentId] = React.useState("");
  const selectedAgentId = agentId || agents[0]?.id || "";

  const canSubmit =
    !disabled && !isLoading && !!selectedAgentId && !!value.trim();
  const canStop = isLoading && !!onStop;

  const handleSubmit = () => {
    if (!canSubmit) return;

    const message = value.trim();
    onSubmit?.({ message, agentId: selectedAgentId });
    setValue("");
  };

  const handleAction = () => {
    if (canStop) {
      onStop?.();
      return;
    }
    handleSubmit();
  };

  return (
    <PromptInput
      value={value}
      onValueChange={setValue}
      onSubmit={handleSubmit}
      isLoading={isLoading}
      disabled={disabled}
      className={cn(
        "w-full",
        containerClassName,
        className,
      )}
    >
      <PromptInputTextarea placeholder={placeholder} disabled={disabled} />
      <PromptInputActions className="justify-between pt-2">
        <Select
          value={selectedAgentId || undefined}
          onValueChange={setAgentId}
          disabled={isLoadingAgents || agents.length === 0 || disabled}
        >
          <SelectTrigger className="h-8 w-auto min-w-24 border-0 bg-transparent px-2 shadow-none">
            <SelectValue
              placeholder={
                isLoadingAgents
                  ? "Loading agents..."
                  : agents.length === 0
                    ? "No agents available"
                    : "Select an agent"
              }
            />
          </SelectTrigger>
          <SelectContent>
            {agents.map((agent) => (
              <SelectItem key={agent.id} value={agent.id}>
                {agent.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <PromptInputAction
          tooltip={isLoading ? (canStop ? "Stop" : "Running...") : "Send"}
        >
          <Button
            type="button"
            size="icon"
            className="size-8 shrink-0 rounded-full"
            disabled={!canSubmit && !canStop}
            onClick={handleAction}
          >
            {isLoading ? (
              <Square className="size-3.5 fill-current" />
            ) : (
              <ArrowUp className="size-4" />
            )}
          </Button>
        </PromptInputAction>
      </PromptInputActions>
    </PromptInput>
  );
}

export { ChatInput };
