import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue
} from "@/components/ui/select";
import { useListAgent } from "@/features/agent/hooks";
import { useCreateSession } from "@/features/session/hooks";
import { useEffect, useState } from "react";
import { useParams } from "react-router";

import { ChatInput } from "./chat-input";

export function NewChat() {
  const { clientCode: clientCodeParam } = useParams<{ clientCode: string }>();
  const clientCode = clientCodeParam ?? "";
  const { data: agentsPage, isLoading: isLoadingAgents } = useListAgent();
  const agents = agentsPage?.data ?? [];
  const [agentId, setAgentId] = useState("");
  const { mutate: createSession, isPending } = useCreateSession();

  useEffect(() => {
    if (agentId || !agentsPage?.data?.length) return;
    setAgentId(agentsPage.data[0].id);
  }, [agentId, agentsPage?.data]);

  const handleCreateSession = (initialMessage: string) => {
    if (!clientCode || !agentId) return;

    createSession({
      clientCode,
      data: {
        agent_id: agentId,
        metadata: {
          name: initialMessage.slice(0, 80),
        },
      },
    });
  };

  const canSubmit = !!clientCode && !!agentId && !isPending;

  return (
    <div className="flex flex-col gap-4 p-4 h-full items-center justify-center">
      <h1 className="text-2xl font-bold">New Chat</h1>
      <p className="text-sm text-muted-foreground">Start a new chat with me</p>
      <div className="flex w-full max-w-2xl flex-col gap-3">
        <Select
          value={agentId || undefined}
          onValueChange={setAgentId}
          disabled={isLoadingAgents || agents.length === 0}
        >
          <SelectTrigger className="w-full">
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
        <ChatInput
          containerclassname="w-full"
          placeholder="Ask me anything..."
          disabled={!canSubmit}
          onSubmit={handleCreateSession}
        />
      </div>
    </div>
  );
}
