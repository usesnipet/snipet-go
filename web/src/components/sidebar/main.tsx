import { useGetSystemInfo } from "@/features/app/hooks";
import { useLogout } from "@/features/auth/hooks";
import { ROUTES } from "@/routes";
import { BookOpen, Bot, Cpu, Key, LogOutIcon, Settings, Shield, Users } from "lucide-react";

import { Button } from "../ui/button";
import { Link } from "../ui/link";
import { Sidebar, SidebarFooter, SidebarHeader } from "../ui/sidebar";
import { Skeleton } from "../ui/skeleton";
import { ToggleTheme } from "../ui/toggle-theme";

import { SidebarContent } from "./content";

import type { NavEntry } from "./types";
const navItems: NavEntry[] = [
  {
    title: "Home",
    href: ROUTES.tenantHome,
    icon: Shield,
    exact: true,
  },
  {
    title: "Clients",
    href: ROUTES.clients,
    icon: Users,
  },
  {
    title: "Agent",
    href: ROUTES.agent,
    icon: Bot,
  },
  {
    title: "LLMs",
    href: ROUTES.llms,
    icon: Cpu,
  },
  {
    title: "Knowledge",
    href: ROUTES.knowledge,
    icon: BookOpen,
  },
  {
    title: "API Keys",
    href: ROUTES.apiKeys,
    icon: Key,
  },
  {
    title: "Settings",
    href: ROUTES.settings,
    icon: Settings,
  },
]

export function MainSidebar() {
  const { data: systemInfo, isLoading } = useGetSystemInfo();
  const { mutate: logout } = useLogout();


  return (
    <Sidebar>
      <SidebarHeader className="flex flex-row justify-between items-center">
        <Link
          href="/"
          className="flex items-center gap-2 px-2 py-1 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0"
        >
          <img src="/favicon.svg" alt="Snipet" className="size-7 shrink-0" />
          <div>
            <p className="text-sm font-semibold">Snipet</p>
            <span className="inline-flex gap-1 text-xs text-muted-foreground">
              <strong>Version:</strong>
              {isLoading ? <Skeleton className="w-20 h-4" /> : <p>{systemInfo?.version}</p>}
            </span>
          </div>
        </Link>
        <ToggleTheme />
      </SidebarHeader>
      <SidebarContent navItems={navItems} />
      <SidebarFooter>
        <Button variant="destructive" onClick={() => logout({})}>
          <LogOutIcon className="size-4" />
          Logout
        </Button>
      </SidebarFooter>
    </Sidebar>
  )
}