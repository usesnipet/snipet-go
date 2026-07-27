"use client"
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
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

export type Props = Omit<React.ComponentProps<"textarea">, "onSubmit"> & {
  containerclassname?: string;
  onSubmit?: (value: string) => void;
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
      ...props
    },
    ref,
  ) => {
    const innerRef = React.useRef<HTMLTextAreaElement>(null);
    const scrollAreaRef = React.useRef<HTMLDivElement>(null);
    const [areaHeight, setAreaHeight] = React.useState<number>();

    React.useImperativeHandle(ref, () => innerRef.current!);

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

    const handleSubmit = () => {
      const el = innerRef.current;
      if (!el) return;

      const value = el.value.trim();
      if (!value) return;

      onSubmit?.(value);

      if (props.value === undefined) {
        el.value = "";
        adjustHeight(el);
      }
    };

    return (
      <div
        className={cn(
          "relative flex items-end rounded-md border border-input bg-background",
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
        <Button
          type="button"
          variant="outline"
          size="icon"
          className="absolute right-2 bottom-1"
          disabled={props.disabled}
          onClick={handleSubmit}
        >
          <Send className="w-4 h-4" />
        </Button>
      </div>
    );
  },
);
ChatInput.displayName = "ChatInput";

export { ChatInput };
