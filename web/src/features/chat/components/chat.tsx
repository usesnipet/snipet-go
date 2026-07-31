import { ScrollArea } from "@/components/ui/scroll-area";
import { useNavigate } from "@/hooks/use-navigate";
import { cn } from "@/lib/utils";
import { ROUTES } from "@/routes";
import { useCallback, useEffect, useRef } from "react";
import { useLocation, useParams } from "react-router";

import { useSessionChat } from "../hooks";

import { ChatInput } from "./chat-input";

import type { ChatMessage } from "../hooks";
import type { ChatInputSubmit } from "./chat-input";

type ChatLocationState = {
  initialMessage?: string;
};

function isVisibleMessage(message: ChatMessage) {
  return message.role === "user" || message.role === "assistant" || message.role === "final";
}

function MessageBubble({ message }: { message: ChatMessage }) {
  const isUser = message.role === "user";

  return (
    <div className={cn("flex w-full", isUser ? "justify-end" : "justify-start")}>
      <div
        className={cn(
          "max-w-[85%] rounded-md px-3 py-2 text-sm whitespace-pre-wrap wrap-break-word",
          isUser
            ? "bg-primary text-primary-foreground"
            : "bg-muted text-foreground",
        )}
      >
        {message.content}
      </div>
    </div>
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

  const { messages, isLoading, isRunning, error, sendMessage } = useSessionChat(
    clientCode,
    sessionId,
  );

  const visibleMessages = messages.filter(isVisibleMessage);

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
    <div className="flex flex-col gap-4 p-4 h-full items-center">
      <ScrollArea className="flex-1 w-full max-w-2xl min-h-0">
        <div className="flex flex-col gap-3 pr-4 pb-2">
          {isLoading && visibleMessages.length === 0 ? (
            <p className="text-sm text-muted-foreground text-center py-8">
              Loading messages...
            </p>
          ) : visibleMessages.length === 0 ? (
            <p className="text-sm text-muted-foreground text-center py-8">
              No messages yet
            </p>
          ) : (
            visibleMessages.map((message) => (
              <MessageBubble key={message.id} message={message} />
            ))
          )}
          {isRunning && (
            <p className="text-sm text-muted-foreground">Thinking...</p>
          )}
          {error && (
            <p className="text-sm text-destructive">{error}</p>
          )}
        </div>
      </ScrollArea>
      <ChatInput
        clientCode={clientCode}
        containerclassname="w-full max-w-2xl"
        placeholder="Ask me anything..."
        disabled={!clientCode || !sessionId || isRunning}
        onSubmit={handleSubmit}
      />
    </div>
  );
}
