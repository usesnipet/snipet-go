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

import type { Message as ChatMessage } from "@/features/session/schemas";
import type { ChatInputSubmit } from "./chat-input";

type ChatLocationState = {
  initialMessage?: string;
};

function isVisibleMessage(message: ChatMessage) {
  return message.role === "user" || message.role === "assistant";
}

function ChatMessageItem({ message }: { message: ChatMessage }) {
  const isUser = message.role === "user";

  return (
    <Message
      className={cn(
        "mx-auto w-full max-w-3xl flex-col gap-1 px-2",
        isUser ? "items-end" : "items-start",
      )}
    >
      <MessageContent
        markdown={!isUser}
        className={cn(
          "max-w-[90%]",
          !isUser && "bg-transparent",
        )}
      >
        {message.content}
      </MessageContent>
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

  const { messages, isLoading, isRunning, error, sendMessage, stop } =
    useSessionChat(clientCode, sessionId);

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
                <ChatMessageItem key={message.id} message={message} />
              ))
            )}
            {isRunning && (
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
