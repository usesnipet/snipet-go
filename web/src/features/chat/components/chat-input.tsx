"use client"
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useListAgent } from "@/features/agent/hooks";
import { cn } from "@/lib/utils";
import { Send } from "lucide-react";
import * as React from "react";

const MAX_ROWS = 10;

function getMaxHeight(el: HTMLTextAreaElement) {
  const style = getComputedStyle(el);
  const lineHeight = Number.parseFloat(style.lineHeight) || 20;
  const paddingY =
    Number.parseFloat(style.paddingTop) + Number.parseFloat(style.paddingBottom);
  return lineHeight * MAX_ROWS + paddingY;
}

function getViewport(root: HTMLElement | null) {
  return root?.querySelector<HTMLElement>("[data-radix-scroll-area-viewport]") ?? null;
}

export type ChatInputSubmit = {
  message: string;
  agentId: string;
};

export type Props = Omit<React.ComponentProps<"textarea">, "onSubmit"> & {
  containerclassname?: string;
  onSubmit?: (value: ChatInputSubmit) => void;
};

const ChatInput = React.forwardRef<HTMLTextAreaElement, Props>(
  (
    {
      className,
      containerclassname,
      onSubmit,
      onKeyDown,
      onChange,
      rows = 1,
      disabled,
      ...props
    },
    ref,
  ) => {
    const innerRef = React.useRef<HTMLTextAreaElement>(null);
    const scrollAreaRef = React.useRef<HTMLDivElement>(null);
    const [areaHeight, setAreaHeight] = React.useState<number>();
    const { data: agentsPage, isLoading: isLoadingAgents } = useListAgent();
    const agents = agentsPage?.data ?? [];
    const [agentId, setAgentId] = React.useState("");

    React.useImperativeHandle(ref, () => innerRef.current!);

    React.useEffect(() => {
      if (agentId || !agentsPage?.data?.length) return;
      setAgentId(agentsPage.data[0].id);
    }, [agentId, agentsPage?.data]);

    const adjustHeight = (el: HTMLTextAreaElement | null) => {
      if (!el) return;

      el.style.height = "auto";
      const scrollHeight = el.scrollHeight;
      el.style.height = `${scrollHeight}px`;
      el.style.overflow = "hidden";

      setAreaHeight(Math.min(scrollHeight, getMaxHeight(el)));
    };

    React.useLayoutEffect(() => {
      adjustHeight(innerRef.current);
    }, [props.value]);

    const canSubmit = !disabled && !!agentId;

    const handleSubmit = () => {
      const el = innerRef.current;
      if (!el || !canSubmit) return;

      const message = el.value.trim();
      if (!message) return;

      onSubmit?.({ message, agentId });

      if (props.value === undefined) {
        el.value = "";
        adjustHeight(el);
      }
    };

    return (
      <div
        className={cn(
          "relative flex flex-col rounded-md border border-input bg-background",
          containerclassname,
        )}
      >
        <ScrollArea
          ref={scrollAreaRef}
          className="w-full rounded-md"
          style={areaHeight != null ? { height: areaHeight } : undefined}
        >
          <textarea
            {...props}
            ref={innerRef}
            rows={rows}
            disabled={disabled}
            className={cn(
              "flex min-h-10 w-full resize-none border-0 bg-transparent px-3 py-2 pr-12 text-base placeholder:text-muted-foreground outline-none focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
              className,
            )}
            onChange={(event) => {
              onChange?.(event);
              adjustHeight(event.currentTarget);

              const atEnd =
                event.currentTarget.selectionStart === event.currentTarget.value.length;
              if (!atEnd) return;

              requestAnimationFrame(() => {
                const viewport = getViewport(scrollAreaRef.current);
                if (viewport) viewport.scrollTop = viewport.scrollHeight;
              });
            }}
            onKeyDown={(event) => {
              onKeyDown?.(event);
              if (event.defaultPrevented) return;
              if (event.key === "Enter" && !event.shiftKey) {
                event.preventDefault();
                handleSubmit();
              }
            }}
          />
        </ScrollArea>
        <div className="flex items-center justify-between gap-2 px-2 pb-2">
          <Select
            value={agentId || undefined}
            onValueChange={setAgentId}
            disabled={isLoadingAgents || agents.length === 0 || disabled}
          >
            <SelectTrigger className="h-8 w-auto min-w-24 border-0 bg-transparent shadow-none px-2">
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
          <Button
            type="button"
            variant="outline"
            size="icon"
            className="size-8 shrink-0"
            disabled={!canSubmit}
            onClick={handleSubmit}
          >
            <Send className="w-4 h-4" />
          </Button>
        </div>
      </div>
    );
  },
);
ChatInput.displayName = "ChatInput";

export { ChatInput };
