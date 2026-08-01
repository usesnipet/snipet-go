import { useCreateSession } from "@/features/session/hooks";
import { useNavigate } from "@/hooks/use-navigate";
import { ROUTES } from "@/routes";
import { useParams } from "react-router";

import { ChatInput } from "./chat-input";

import type { ChatInputSubmit } from "./chat-input";
export function NewChat() {
  const { clientCode: clientCodeParam } = useParams<{ clientCode: string }>();
  const clientCode = clientCodeParam ?? "";
  const navigate = useNavigate();
  const { mutate: createSession, isPending } = useCreateSession({ auth: "jwt" });

  const handleCreateSession = ({ message, agentId }: ChatInputSubmit) => {
    if (!clientCode || !agentId) return;

    createSession({
      clientCode,
      data: {
        agent_id: agentId,
        metadata: {
          name: message.slice(0, 80),
        },
      },
    }, {
      onSuccess: (createdSession) => {
        navigate(ROUTES.clientChatSession, {
          params: {
            clientCode,
            sessionId: createdSession.id,
          },
          state: { initialMessage: message },
        })
      }
    });
  };

  return (
    <div className="flex h-full flex-col items-center justify-center gap-4 p-4">
      <h1 className="text-2xl font-bold">New Chat</h1>
      <p className="text-sm text-muted-foreground">Start a new chat with me</p>
      <ChatInput
        clientCode={clientCode}
        containerClassName="w-full max-w-2xl"
        placeholder="Ask me anything..."
        disabled={!clientCode || isPending}
        isLoading={isPending}
        onSubmit={handleCreateSession}
      />
    </div>
  );
}
