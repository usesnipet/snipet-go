import { Link } from "@/components/ui/link";
import { Sidebar, SidebarContent, SidebarHeader } from "@/components/ui/sidebar";
import { Skeleton } from "@/components/ui/skeleton";
import { ToggleTheme } from "@/components/ui/toggle-theme";
import { SessionList } from "@/features/session/components/session-list";
import { useParams } from "react-router";

import { useFindByCodeClient } from "../hooks";

import { ClientCode } from "./client-code";

export function ClientChatSidebar() {
  const { clientCode } = useParams<{ clientCode: string }>();
  const { data: client, isLoading } = useFindByCodeClient(clientCode);

  return (
    <Sidebar>
      <SidebarHeader className="flex flex-row justify-between items-center">
        <Link
          href="/"
          className="flex items-center gap-2 px-2 py-1 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0"
        >
          <img src="/favicon.svg" alt="Snipet" className="size-7 shrink-0" />
          <div>
            {isLoading ? <Skeleton className="w-20 h-4" /> : <p className="text-sm font-semibold">{client?.name}</p>}
            <ClientCode code={clientCode} />
          </div>
        </Link>
        <ToggleTheme />
      </SidebarHeader>
      <SidebarContent>
        <SessionList clientCode={clientCode} />
      </SidebarContent>
    </Sidebar>
  )
}