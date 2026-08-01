import { findMessagesSessionQueryKey, useFindMessagesSession } from "@/features/session/hooks";
import {
  SSE_EVENT,
  executionMessageAddedEventSchema,
  executionStatusChangedEventSchema,
  executionTurnCompletedEventSchema,
  sseErrorEventSchema,
} from "@/features/session/schemas";
import { sessionService } from "@/features/session/service";
import { toast } from "@/hooks/use-toast";
import { ApiError } from "@/lib/http";
import { queryClient } from "@/lib/query-client";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import type { Message } from "@/features/session/schemas";

const EMPTY_MESSAGES: Message[] = [];

type PendingMessage = {
  message: Message;
  knownIds: Set<string>;
};

function isAbortError(error: unknown): boolean {
  return error instanceof Error && error.name === "AbortError";
}

/**
 * The run stream never echoes the user message back, so the persisted copy is only
 * recognizable by its content: any user message that was not around when we sent it.
 */
function isPendingConfirmed(messages: Message[], pending: PendingMessage): boolean {
  return messages.some(
    (message) =>
      message.role === "user" &&
      message.content === pending.message.content &&
      !pending.knownIds.has(message.id),
  );
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

  const [streamed, setStreamed] = useState<Message[]>(EMPTY_MESSAGES);
  const [pending, setPending] = useState<PendingMessage | null>(null);
  const [isRunning, setIsRunning] = useState(false);
  const [turn, setTurn] = useState<number | null>(null);
  const [runError, setRunError] = useState<string | null>(null);
  const [renderedSessionId, setRenderedSessionId] = useState(sessionId);

  const abortRef = useRef<AbortController | null>(null);
  const isRunningRef = useRef(false);
  const sessionIdRef = useRef(sessionId);
  const messagesRef = useRef<Message[]>(EMPTY_MESSAGES);

  if (renderedSessionId !== sessionId) {
    setRenderedSessionId(sessionId);
    setStreamed(EMPTY_MESSAGES);
    setPending(null);
    setIsRunning(false);
    setTurn(null);
    setRunError(null);
  }

  useEffect(() => {
    sessionIdRef.current = sessionId;
    return () => {
      abortRef.current?.abort();
      abortRef.current = null;
      isRunningRef.current = false;
    };
  }, [sessionId]);

  const history = historyPage?.data;

  const messages = useMemo(() => {
    const byId = new Map<string, Message>();

    for (const message of history ?? EMPTY_MESSAGES) {
      byId.set(message.id, message);
    }
    for (const message of streamed) {
      if (!byId.has(message.id)) byId.set(message.id, message);
    }

    const persisted = Array.from(byId.values());
    if (!pending || isPendingConfirmed(persisted, pending)) return persisted;

    return [...persisted, pending.message];
  }, [history, streamed, pending]);

  useEffect(() => {
    messagesRef.current = messages;
  }, [messages]);

  const sendMessage = useCallback(
    async (text: string) => {
      const message = text.trim();
      if (!clientCode || !sessionId || !message || isRunningRef.current) return;

      const runSessionId = sessionId;
      const isSameSession = () => sessionIdRef.current === runSessionId;

      isRunningRef.current = true;
      setIsRunning(true);
      setRunError(null);
      setTurn(null);
      setPending({
        message: {
          id: crypto.randomUUID(),
          sequence: -1,
          role: "user",
          content: message,
          timestamp: new Date(),
        },
        knownIds: new Set(messagesRef.current.map((m) => m.id)),
      });

      const controller = new AbortController();
      abortRef.current = controller;

      try {
        await sessionService.run(
          clientCode,
          runSessionId,
          { message },
          {
            auth: "jwt",
            signal: controller.signal,
            onEvent: (event, data) => {
              if (controller.signal.aborted || !isSameSession()) return;

              switch (event) {
                case SSE_EVENT.EXECUTION_MESSAGE_ADDED: {
                  const parsed = executionMessageAddedEventSchema.safeParse(data);
                  if (!parsed.success) return;
                  const incoming = parsed.data.message;
                  setStreamed((prev) =>
                    prev.some((m) => m.id === incoming.id) ? prev : [...prev, incoming],
                  );
                  break;
                }
                case SSE_EVENT.EXECUTION_STATUS_CHANGED: {
                  const parsed = executionStatusChangedEventSchema.safeParse(data);
                  if (!parsed.success) return;
                  if (
                    parsed.data.status === "failed" ||
                    parsed.data.status === "max_turns"
                  ) {
                    setRunError(
                      parsed.data.error_message ?? `Run ${parsed.data.status}`,
                    );
                  }
                  break;
                }
                case SSE_EVENT.EXECUTION_TURN_COMPLETED: {
                  const parsed = executionTurnCompletedEventSchema.safeParse(data);
                  if (!parsed.success) return;
                  setTurn(parsed.data.turn);
                  break;
                }
                case SSE_EVENT.ERROR: {
                  const parsed = sseErrorEventSchema.safeParse(data);
                  setRunError(parsed.success ? parsed.data.message : "Run failed");
                  break;
                }
              }
            },
          },
        );
      } catch (error) {
        // The message is already persisted when the stream is aborted, so the refetch
        // below is what replaces the optimistic copy.
        if (isAbortError(error)) return;

        const errMessage =
          error instanceof ApiError ? error.message : "Failed to send message";
        if (isSameSession()) {
          setRunError(errMessage);
          setPending(null);
        }
        toast({
          title: "Failed to send message",
          description: errMessage,
          variant: "destructive",
        });
      } finally {
        if (abortRef.current === controller) abortRef.current = null;
        isRunningRef.current = false;
        if (isSameSession()) setIsRunning(false);
        void queryClient.invalidateQueries({
          queryKey: findMessagesSessionQueryKey(clientCode, runSessionId),
        });
      }
    },
    [clientCode, sessionId],
  );

  const stop = useCallback(() => {
    abortRef.current?.abort();
  }, []);

  return {
    messages,
    isLoading,
    isRunning,
    turn,
    error: historyError?.message ?? runError,
    sendMessage,
    stop,
  };
}
