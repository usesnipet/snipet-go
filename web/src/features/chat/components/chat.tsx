import { ChatContainerContent, ChatContainerRoot } from "@/components/ui/chat-container";
import { Loader } from "@/components/ui/loader";
import { Message, MessageContent } from "@/components/ui/message";
import { ScrollButton } from "@/components/ui/scroll-button";
import { useNavigate } from "@/hooks/use-navigate";
import { cn } from "@/lib/utils";
import { ROUTES } from "@/routes";
import { useCallback, useEffect, useRef } from "react";
import { useLocation, useParams } from "react-router";

import { useSessionChat } from "../hooks";

import { ChatInput } from "./chat-input";

import type { ToolCallResult } from "../hooks";
import type { Message as ChatMessage, ToolCall } from "@/features/session/schemas";
import type { ChatInputSubmit } from "./chat-input";

type ChatLocationState = {
  initialMessage?: string;
};

function isVisibleMessage(message: ChatMessage) {
  return message.role === "user" || message.role === "assistant";
}

function formatToolArguments(args: Record<string, unknown>) {
  const keys = Object.keys(args);
  if (keys.length === 0) return null;
  try {
    return JSON.stringify(args, null, 2);
  } catch {
    return null;
  }
}

function ToolCallItem({
  call,
  result,
}: {
  call: ToolCall;
  result?: ToolCallResult;
}) {
  const args = formatToolArguments(call.arguments);
  const isPending = !result;
  const failed = Boolean(result?.error);

  return (
    <div className="rounded-md border border-border/60 bg-muted/40 px-3 py-2 text-xs text-muted-foreground">
      <div className="flex items-center gap-2 font-medium text-foreground">
        <span>{call.tool}</span>
        {isPending ? (
          <Loader variant="text-shimmer" text="Running" className="text-xs" />
        ) : failed ? (
          <span className="text-destructive">Failed</span>
        ) : (
          <span className="text-muted-foreground">Done</span>
        )}
      </div>
      {args && (
        <pre className="mt-1 max-h-32 overflow-auto whitespace-pre-wrap break-all font-mono text-[11px] opacity-80">
          {args}
        </pre>
      )}
      {result?.error && (
        <p className="mt-1 text-destructive">{result.error}</p>
      )}
      {result?.result && (
        <pre className="mt-1 max-h-40 overflow-auto whitespace-pre-wrap break-all font-mono text-[11px] opacity-80">
          {result.result}
        </pre>
      )}
    </div>
  );
}

function ChatMessageItem({
  message,
  toolResults,
}: {
  message: ChatMessage;
  toolResults: Record<string, ToolCallResult>;
}) {
  const isUser = message.role === "user";
  const toolCalls = message.tool_calls ?? [];

  return (
    <Message
      className={cn(
        "mx-auto w-full max-w-3xl flex-col gap-1 px-2",
        isUser ? "items-end" : "items-start",
      )}
    >
      {message.content ? (
        <MessageContent
          markdown={!isUser}
          className={cn("max-w-[90%]", !isUser && "bg-transparent")}
        >
          {message.content}
        </MessageContent>
      ) : null}
      {!isUser && toolCalls.length > 0 && (
        <div className="flex w-full max-w-[90%] flex-col gap-2">
          {toolCalls.map((call) => (
            <ToolCallItem
              key={call.id}
              call={call}
              result={toolResults[call.id]}
            />
          ))}
        </div>
      )}
    </Message>
  );
}

export function Chat() {
  const { clientCode: clientCodeParam, sessionId: sessionIdParam } = useParams<{
    clientCode: string;
    sessionId: string;
  }>();
  const clientCode = clientCodeParam ?? "";
  const sessionId = sessionIdParam ?? "";
  const location = useLocation();
  const navigate = useNavigate();
  const initialSentRef = useRef(false);

  const { messages, toolResults, isLoading, isRunning, error, sendMessage, stop } =
    useSessionChat(clientCode, sessionId);

  const visibleMessages = messages.filter(isVisibleMessage);
  const lastVisible = visibleMessages[visibleMessages.length - 1];
  const toolsFinishedOnLast =
    lastVisible?.role === "assistant" &&
    (lastVisible.tool_calls?.length ?? 0) > 0 &&
    lastVisible.tool_calls!.every((call) => toolResults[call.id]);
  const waitingForFirstChunk =
    isRunning &&
    (!lastVisible ||
      lastVisible.role === "user" ||
      (lastVisible.role === "assistant" &&
        !lastVisible.content &&
        !(lastVisible.tool_calls?.length)) ||
      toolsFinishedOnLast);

  const handleSubmit = useCallback(
    ({ message }: ChatInputSubmit) => {
      void sendMessage(message);
    },
    [sendMessage],
  );

  useEffect(() => {
    initialSentRef.current = false;
  }, [sessionId]);

  useEffect(() => {
    const state = (location.state ?? null) as ChatLocationState | null;
    const initialMessage = state?.initialMessage?.trim();
    if (!initialMessage || !sessionId || initialSentRef.current || isLoading) return;

    initialSentRef.current = true;
    navigate(ROUTES.clientChatSession, {
      params: { clientCode, sessionId },
      replace: true,
      state: {},
    });
    void sendMessage(initialMessage);
  }, [location.state, sessionId, clientCode, isLoading, navigate, sendMessage]);

  return (
    <div className="flex h-full flex-col items-center gap-4 p-4">
      <div className="relative min-h-0 w-full flex-1">
        <ChatContainerRoot className="h-full">
          <ChatContainerContent className="gap-4 px-2 pb-4">
            {isLoading && visibleMessages.length === 0 ? (
              <p className="py-8 text-center text-sm text-muted-foreground">
                Loading messages...
              </p>
            ) : visibleMessages.length === 0 ? (
              <p className="py-8 text-center text-sm text-muted-foreground">
                No messages yet
              </p>
            ) : (
              visibleMessages.map((message) => (
                <ChatMessageItem
                  key={message.id}
                  message={message}
                  toolResults={toolResults}
                />
              ))
            )}
            {waitingForFirstChunk && (
              <Message className="mx-auto w-full items-start px-2">
                <Loader variant="text-shimmer" text="Thinking" />
              </Message>
            )}
            {error && (
              <p className="px-2 text-sm text-destructive">{error}</p>
            )}
          </ChatContainerContent>
          <div className="absolute right-4 bottom-4">
            <ScrollButton />
          </div>
        </ChatContainerRoot>
      </div>
      <ChatInput
        clientCode={clientCode}
        containerClassName="w-full max-w-3xl"
        placeholder="Ask me anything..."
        disabled={!clientCode || !sessionId}
        isLoading={isRunning}
        onSubmit={handleSubmit}
        onStop={stop}
      />
    </div>
  );
}
