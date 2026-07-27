import { useCreateSession } from "@/features/session/hooks";
import { useParams } from "react-router";

import { ChatInput, type ChatInputSubmit } from "./chat-input";

export function NewChat() {
  const { clientCode: clientCodeParam } = useParams<{ clientCode: string }>();
  const clientCode = clientCodeParam ?? "";
  const { mutate: createSession, isPending } = useCreateSession();

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
    });
  };

  return (
    <div className="flex flex-col gap-4 p-4 h-full items-center justify-center">
      <h1 className="text-2xl font-bold">New Chat</h1>
      <p className="text-sm text-muted-foreground">Start a new chat with me</p>
      <ChatInput
        containerclassname="w-full max-w-2xl"
        placeholder="Ask me anything..."
        disabled={!clientCode || isPending}
        onSubmit={handleCreateSession}
      />
    </div>
  );
}
