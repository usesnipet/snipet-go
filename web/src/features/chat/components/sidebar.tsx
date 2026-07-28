import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { InputSearch } from "@/components/ui/input-search";
import { Link } from "@/components/ui/link";
import { Sidebar, SidebarContent, SidebarFooter, SidebarHeader } from "@/components/ui/sidebar";
import { Skeleton } from "@/components/ui/skeleton";
import { ToggleTheme } from "@/components/ui/toggle-theme";
import { useAuthStore } from "@/features/auth/store";
import { useListSession } from "@/features/session/hooks";
import { useNavigate } from "@/hooks/use-navigate";
import { applyPathParams } from "@/lib/http";
import { ROUTES } from "@/routes";
import { LogOutIcon, MessageSquarePlus, UserIcon } from "lucide-react";
import { useState } from "react";
import { useParams } from "react-router";

import { ClientCode } from "../../client/components/client-code";
import { useFindPublicByCodeClient } from "../../client/hooks";

import { SessionList } from "./session-list";

export function ClientChatSidebar() {
  const navigate = useNavigate();
  const clearAccessToken = useAuthStore(state => state.clear);
  const [sessionSearch, setSessionSearch] = useState("");
  const { clientCode: clientCodeParam } = useParams<{ clientCode: string }>();
  const clientCode = clientCodeParam ?? "";
  const { data: client, isLoading } = useFindPublicByCodeClient(clientCode);
  const { data: sessions, isLoading: isLoadingSessions, error: errorSessions } = useListSession(
    clientCode, { auth: "jwt" }
  )

  const handleLogout = () => {
    clearAccessToken();
    navigate(ROUTES.home, { replace: true });
  }

  return (
    <Sidebar>
      <SidebarHeader>
        <div className="flex items-center justify-between">
          <Link
            href="/"
            className="flex items-center gap-2 px-2 py-1 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0"
          >
            <img src="/favicon.svg" alt="Snipet" className="size-7 shrink-0" />
            <div>
              {isLoading ? <Skeleton className="h-4 w-20" /> : <p className="text-sm font-semibold">{client?.name}</p>}
              <ClientCode code={clientCode} />
            </div>
          </Link>
          <ToggleTheme />
        </div>
        <InputSearch
          value={sessionSearch}
          onChange={(event) => setSessionSearch(event.target.value)}
          placeholder="Search sessions..."
          aria-label="Search sessions"
        />
        <Button className="w-full" asChild>
          <Link href={applyPathParams(ROUTES.clientChat, { clientCode })}>
            <MessageSquarePlus />
            New chat
          </Link>
        </Button>
      </SidebarHeader>
      <SidebarContent>
        <SessionList
          clientCode={clientCode}
          search={sessionSearch}
          data={sessions?.data ?? []}
          isLoading={isLoadingSessions}
          error={errorSessions}
        />
      </SidebarContent>
      <SidebarFooter>
        <Card className="flex items-center gap-3 p-2">
          <div className="flex size-9 shrink-0 items-center justify-center rounded-full bg-muted">
            <UserIcon className="size-4" />
          </div>
          <p className="min-w-0 flex-1 truncate text-sm font-medium">Jhon Doe</p>
          <Button
            variant="ghost"
            size="icon"
            onClick={handleLogout}
            aria-label="Logout"
            title="Logout"
          >
            <LogOutIcon />
          </Button>
        </Card>
      </SidebarFooter>
    </Sidebar>
  )
}