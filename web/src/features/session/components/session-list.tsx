import { ErrorFallback } from "@/components/error-fallback";
import { LoadingFallback } from "@/components/loading-fallback";
import {
  SidebarGroup, SidebarGroupContent, SidebarGroupLabel, SidebarMenu, SidebarMenuButton, SidebarMenuItem
} from "@/components/ui/sidebar";
import { MessageSquare } from "lucide-react";

import { useListSession } from "../hooks";

type SessionListProps = {
  clientCode: string;
  search?: string;
}

export const SessionList = ({ clientCode, search = "" }: SessionListProps) => {
  const { data, isLoading, error } = useListSession(clientCode);
  if (isLoading) return <LoadingFallback />;
  if (error) return <ErrorFallback error={error} />;
  if (!data) return null;

  const normalizedSearch = search.trim().toLowerCase();
  const sessions = data.data.filter((session) => {
    const name = session.metadata.name ?? "";
    return `${name} ${session.id}`.toLowerCase().includes(normalizedSearch);
  });

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
                <SidebarMenuButton tooltip={session.metadata.name ?? session.id}>
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