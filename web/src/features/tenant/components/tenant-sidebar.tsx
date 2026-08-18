import { SidebarContent } from "@/components/sidebar/content";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Separator } from "@/components/ui/separator";
import { Sidebar, SidebarFooter, SidebarHeader } from "@/components/ui/sidebar";
import { useTheme } from "@/context/theme-provider";
import { useGetSystemInfo } from "@/features/app/hooks";
import { TenantCard } from "@/features/tenant/components/tenant-card";
import { useFindMineTenant } from "@/features/tenant/hooks";
import { UserCard } from "@/features/user/components/user-card";
import { useNavigate } from "@/hooks/use-navigate";
import { ROUTES } from "@/routes";
import {
  BookOpen, Bot, Building2Icon, ChevronsUpDownIcon, Cpu, Key, MoonIcon, Settings, Shield, SunIcon, Users
} from "lucide-react";
import { useMemo, useState } from "react";

import { useTenantStore } from "../store";

import type { NavEntry } from "@/components/sidebar/types";

const MAX_VISIBLE_TENANTS = 3;

export function TenantSidebar() {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();
  const { theme, setTheme } = useTheme();
  const { data } = useFindMineTenant();
  const { data: systemInfo } = useGetSystemInfo();
  const multiTenantEnabled = systemInfo?.multi_tenant_enabled ?? false;

  const tenants = data?.data ?? [];
  const currentTenant = useTenantStore(s => s.tenant);
  const visibleTenants = tenants.slice(0, MAX_VISIBLE_TENANTS);
  const hasMoreTenants = tenants.length > MAX_VISIBLE_TENANTS;


  const navItems: NavEntry[] = useMemo(() => [
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
      label: "Management",
      items: [
        ...(multiTenantEnabled ? [
          {
            title: "Members",
            items: [
              {
                title: "Members",
                href: ROUTES.members,
                exact: true,
              },
              {
                title: "Invite",
                href: ROUTES.inviteMember,
              }
            ],
            icon: Users,
          },
        ] : [
          {
            title: "Members",
            href: ROUTES.members,
            exact: true,
            icon: Users,
          },
        ]),
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
    },
  ], [multiTenantEnabled])

  const goToSelectTenant = () => {
    setOpen(false);
    navigate(ROUTES.selectTenant);
  };

  const handleToggleTheme = () => setTheme(theme === "dark" ? "light" : "dark");

  return (
    <Sidebar>
      <SidebarHeader>
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <button
              type="button"
              className="flex w-full items-center gap-2 rounded-md px-2 py-1 text-left hover:bg-muted group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0"
            >
              <Avatar size="sm">
                {currentTenant?.icon && (
                  <AvatarImage src={currentTenant.icon} alt={currentTenant.name} />
                )}
                <AvatarFallback>
                  <Building2Icon className="size-4" />
                </AvatarFallback>
              </Avatar>
              <div className="min-w-0 flex-1 group-data-[collapsible=icon]:hidden">
                <p className="truncate text-sm font-semibold leading-tight">
                  {currentTenant?.name ?? "Select tenant"}
                </p>
                {currentTenant && (
                  <p className="truncate text-xs text-muted-foreground">{currentTenant.slug}</p>
                )}
              </div>
              <ChevronsUpDownIcon className="size-4 shrink-0 text-muted-foreground group-data-[collapsible=icon]:hidden" />
            </button>
          </PopoverTrigger>
          <PopoverContent align="start" className="w-72 p-2">
            <p className="px-2 py-1 text-xs font-semibold text-muted-foreground">Tenants</p>
            <div className="flex flex-col gap-2" onClick={() => setOpen(false)}>
              {visibleTenants.length > 0 ? (
                visibleTenants.map((tenant) => (
                  <TenantCard key={tenant.id} tenant={tenant} compact />
                ))
              ) : (
                <p className="px-2 py-1 text-sm text-muted-foreground">No tenants yet.</p>
              )}
              {hasMoreTenants && (
                <Button variant="ghost" className="w-full justify-center" onClick={goToSelectTenant}>
                  Show more
                </Button>
              )}
            </div>
            <Separator className="my-2" />
            <p className="px-2 py-1 text-xs font-semibold text-muted-foreground">Theme</p>
            <Button variant="outline" className="w-full justify-start" onClick={handleToggleTheme}>
              {theme === "dark" ? <SunIcon className="size-4" /> : <MoonIcon className="size-4" />}
              {theme === "dark" ? "Light mode" : "Dark mode"}
            </Button>
          </PopoverContent>
        </Popover>
      </SidebarHeader>
      <SidebarContent navItems={navItems} />
      <SidebarFooter>
        <UserCard />
      </SidebarFooter>
    </Sidebar>
  )
}
