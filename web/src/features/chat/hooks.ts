import { findMessagesSessionQueryKey, useFindMessagesSession } from "@/features/session/hooks";
import {
  messageAddedEventSchema, sseErrorEventSchema, statusChangedEventSchema
} from "@/features/session/schemas";
import { sessionService } from "@/features/session/service";
import { toast } from "@/hooks/use-toast";
import { ApiError } from "@/lib/http";
import { queryClient } from "@/lib/query-client";
import { useCallback, useEffect, useRef, useState } from "react";

import type { RuntimeMessage } from "@/features/session/schemas";

function mergeMessages(existing: RuntimeMessage[], incoming: RuntimeMessage[]): RuntimeMessage[] {
  const byId = new Map(existing.map((m) => [m.id, m]));

  for (const msg of incoming) {
    byId.set(msg.id, msg);
  }

  return Array.from(byId.values()).sort((a, b) => {
    const timeDiff = a.timestamp.getTime() - b.timestamp.getTime();
    if (timeDiff !== 0) return timeDiff;
    return a.sequence - b.sequence;
  });
}

function replaceOptimisticUserMessage(
  messages: RuntimeMessage[],
  optimisticId: string,
  confirmed: RuntimeMessage[],
): RuntimeMessage[] {
  const userConfirmed = confirmed.find((m) => m.role === "user");
  if (!userConfirmed) {
    return mergeMessages(messages, confirmed);
  }

  const withoutOptimistic = messages.filter((m) => m.id !== optimisticId);
  return mergeMessages(withoutOptimistic, confirmed);
}

export function useSessionChat(clientCode: string, sessionId: string) {
  const {
    data: historyPage,
    isLoading,
    error: historyError,
  } = useFindMessagesSession(clientCode, sessionId, {
    auth: "jwt",
    searchParams: { sort: "asc" },
  });

  const [messages, setMessages] = useState<RuntimeMessage[]>([]);
  const [isRunning, setIsRunning] = useState(false);
  const [runError, setRunError] = useState<string | null>(null);
  const abortRef = useRef<AbortController | null>(null);
  const historySyncedRef = useRef<string | null>(null);
  const prevSessionIdRef = useRef(sessionId);

  useEffect(() => {
    if (!historyPage?.data) return;
    if (historySyncedRef.current === sessionId && isRunning) return;

    const history = historyPage.data;
    setMessages((prev: RuntimeMessage[]) => {
      if (historySyncedRef.current !== sessionId) {
        historySyncedRef.current = sessionId;
        return history;
      }
      return mergeMessages(history, prev);
    });
  }, [historyPage, sessionId, isRunning]);

  useEffect(() => {
    return () => {
      abortRef.current?.abort();
    };
  }, []);

  useEffect(() => {
    if (prevSessionIdRef.current === sessionId) return;
    prevSessionIdRef.current = sessionId;

    abortRef.current?.abort();
    abortRef.current = null;
    setIsRunning(false);
    setRunError(null);
    historySyncedRef.current = null;
    setMessages([]);
  }, [sessionId]);

  const sendMessage = useCallback(
    async (text: string) => {
      const message = text.trim();
      if (!clientCode || !sessionId || !message || isRunning) return;

      setIsRunning(true);
      setRunError(null);

      const optimisticId = crypto.randomUUID();
      const optimistic: ChatMessage = {
        id: optimisticId,
        sequence: Number.MAX_SAFE_INTEGER,
        role: "user",
        content: message,
        created_at: new Date(),
      };
      setMessages((prev) => [...prev, optimistic]);

      const controller = new AbortController();
      abortRef.current = controller;

      try {
        await sessionService.run(
          clientCode,
          sessionId,
          { message },
          {
            auth: "jwt",
            signal: controller.signal,
            onEvent: (event, data) => {
              switch (event) {
                case "message_added": {
                  const parsed = messageAddedEventSchema.safeParse(data);
                  if (!parsed.success) return;
                  const incoming = parsed.data.messages.map(normalizeRuntimeMessage);
                  setMessages((prev) =>
                    replaceOptimisticUserMessage(prev, optimisticId, incoming),
                  );
                  break;
                }
                case "status_changed": {
                  const parsed = statusChangedEventSchema.safeParse(data);
                  if (!parsed.success) return;
                  if (
                    parsed.data.status === "failed" ||
                    parsed.data.status === "cancelled"
                  ) {
                    setRunError(parsed.data.error_message ?? `Run ${parsed.data.status}`);
                  }
                  break;
                }
                case "error": {
                  const parsed = sseErrorEventSchema.safeParse(data);
                  const errMessage = parsed.success
                    ? parsed.data.message
                    : "Run failed";
                  setRunError(errMessage);
                  break;
                }
                case "close":
                  break;
              }
            },
          },
        );
      } catch (error) {
        if (error instanceof DOMException && error.name === "AbortError") {
          return;
        }
        const errMessage =
          error instanceof ApiError ? error.message : "Failed to send message";
        setRunError(errMessage);
        setMessages((prev) => prev.filter((m) => m.id !== optimisticId));
        toast({
          title: "Failed to send message",
          description: errMessage,
          variant: "destructive",
        });
      } finally {
        if (abortRef.current === controller) {
          abortRef.current = null;
        }
        setIsRunning(false);
        void queryClient.invalidateQueries({
          queryKey: findMessagesSessionQueryKey(clientCode, sessionId),
        });
      }
    },
    [clientCode, sessionId, isRunning],
  );

  return {
    messages,
    isLoading,
    isRunning,
    error: historyError?.message ?? runError,
    sendMessage,
  };
}
