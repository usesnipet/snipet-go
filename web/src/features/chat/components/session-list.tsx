import { ErrorFallback } from "@/components/error-fallback";
import { LoadingFallback } from "@/components/loading-fallback";
import {
  SidebarGroup, SidebarGroupContent, SidebarGroupLabel, SidebarMenu, SidebarMenuButton, SidebarMenuItem
} from "@/components/ui/sidebar";
import { useNavigate } from "@/hooks/use-navigate";
import { ROUTES } from "@/routes";
import { MessageSquare } from "lucide-react";

import type { Session } from "../../session/schemas";
type SessionListProps = {
  clientCode: string;
  data: Session[];
  isLoading: boolean;
  error: Error;
  search?: string;
}

export const SessionList = ({ clientCode, search = "", data, error, isLoading }: SessionListProps) => {
  const navigate = useNavigate();

  if (isLoading) return <LoadingFallback />;
  if (error) return <ErrorFallback error={error} />;
  if (!data) return null;

  const normalizedSearch = search.trim().toLowerCase();
  const sessions = data.filter((session) => {
    const name = session.metadata.name ?? "";
    return `${name} ${session.id}`.toLowerCase().includes(normalizedSearch);
  });

  const handleClick = (sessionId: string) => {
    console.log("handleClick", sessionId);
    navigate(ROUTES.clientChatSession, { params: { clientCode, sessionId } });
  }

  return (
    <SidebarGroup>
      <SidebarGroupLabel>Sessions</SidebarGroupLabel>
      <SidebarGroupContent>
        {sessions.length === 0 ? (
          <p className="px-2 py-4 text-center text-xs text-muted-foreground">
            No sessions found.
          </p>
        ) : (
          <SidebarMenu>
            {sessions.map((session) => (
              <SidebarMenuItem key={session.id}>
                <SidebarMenuButton
                  tooltip={session.metadata.name ?? session.id}
                  onClick={() => handleClick(session.id)}
                >
                  <MessageSquare />
                  <span>{session.metadata.name ?? `Session ${session.id.slice(0, 8)}`}</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        )}
      </SidebarGroupContent>
    </SidebarGroup>
  )
}