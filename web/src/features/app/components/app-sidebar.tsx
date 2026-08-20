import { SidebarContent } from "@/components/sidebar/content";
import { Link } from "@/components/ui/link";
import { Sidebar, SidebarHeader } from "@/components/ui/sidebar";
import { Skeleton } from "@/components/ui/skeleton";
import { ToggleTheme } from "@/components/ui/toggle-theme";
import { useTenantStore } from "@/features/tenant/store";
import { ROUTES } from "@/routes";
import { Boxes, Home, Settings, Users } from "lucide-react";
import { useParams } from "react-router";

import { useFindByCodeApp } from "../hooks";

import { AppCode } from "./app-code";

import type { NavEntry } from "@/components/sidebar/types";
const navItems: NavEntry[] = [
  {
    title: "Overview",
    href: ROUTES.app,
    icon: Home,
    exact: true,
  },
  {
    title: "Users",
    href: ROUTES.appUsers,
    icon: Users,
  },
  {
    title: "Sessions",
    href: ROUTES.appSessions,
    icon: Boxes,
  },
  {
    title: "Settings",
    href: ROUTES.appSettings,
    icon: Settings,
  }
]

export function AppSidebar() {
  const { appCode = "" } = useParams<{ appCode: string }>();
  const tenant = useTenantStore((state) => state.tenant);
  const { data: app, isLoading } = useFindByCodeApp(tenant?.id ?? "", appCode);

  return (
    <Sidebar>
      <SidebarHeader className="flex flex-row justify-between items-center">
        <Link
          href="/"
          className="flex items-center gap-2 px-2 py-1 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0"
        >
          <img src="/favicon.svg" alt="Snipet" className="size-7 shrink-0" />
          <div>
            {isLoading ? <Skeleton className="w-20 h-4" /> : <p className="text-sm font-semibold">{app?.name}</p>}
            <AppCode code={appCode} />
          </div>
        </Link>
        <ToggleTheme />
      </SidebarHeader>
      <SidebarContent navItems={navItems} />
    </Sidebar>
  )
}
