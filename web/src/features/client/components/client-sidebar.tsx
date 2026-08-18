import { SidebarContent } from "@/components/sidebar/content";
import { Link } from "@/components/ui/link";
import { Sidebar, SidebarHeader } from "@/components/ui/sidebar";
import { Skeleton } from "@/components/ui/skeleton";
import { ToggleTheme } from "@/components/ui/toggle-theme";
import { ROUTES } from "@/routes";
import { Boxes, Home, Settings, Users } from "lucide-react";
import { useParams } from "react-router";

import { useFindPublicByCodeClient } from "../hooks";

import { ClientCode } from "./client-code";

import type { NavEntry } from "@/components/sidebar/types";
const navItems: NavEntry[] = [
  {
    title: "Overview",
    href: ROUTES.client,
    icon: Home,
    exact: true,
  },
  {
    title: "Users",
    href: ROUTES.clientUsers,
    icon: Users,
  },
  {
    title: "Sessions",
    href: ROUTES.clientSessions,
    icon: Boxes,
  },
  {
    title: "Settings",
    href: ROUTES.clientSettings,
    icon: Settings,
  }
]

export function ClientSidebar() {
  const { clientCode = "" } = useParams<{ clientCode: string }>();
  const { data: client, isLoading } = useFindPublicByCodeClient(clientCode);

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
      <SidebarContent navItems={navItems} />
    </Sidebar>
  )
}