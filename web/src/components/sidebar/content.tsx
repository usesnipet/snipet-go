"use client"

import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { Link } from "@/components/ui/link";
import {
  SidebarContent as SidebarContentBase, SidebarGroup, SidebarGroupContent, SidebarMenu, SidebarMenuButton,
  SidebarMenuItem, SidebarMenuSub, SidebarMenuSubButton, SidebarMenuSubItem
} from "@/components/ui/sidebar";
import { ChevronRight } from "lucide-react";
import { useLocation } from "react-router";

import { isNavActive, isNavGroupActive, isNavItemWithChildren } from "./utils";

import type { NavEntry } from "./types";

type Props = {
  navItems: NavEntry[]
}

export function SidebarContent({ navItems }: Props) {
  const { pathname } = useLocation();

  return (
    <SidebarContentBase>
      <SidebarGroup>
        <SidebarGroupContent>
          <SidebarMenu>
            {navItems.map((item) =>
              isNavItemWithChildren(item) ? (
                <Collapsible
                  key={item.title}
                  asChild
                  defaultOpen={isNavGroupActive(pathname, item.items)}
                  className="group/collapsible"
                >
                  <SidebarMenuItem>
                    <CollapsibleTrigger asChild>
                      <SidebarMenuButton
                        tooltip={item.title}
                        isActive={isNavGroupActive(pathname, item.items)}
                      >
                        <item.icon />
                        <span>{item.title}</span>
                        <ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
                      </SidebarMenuButton>
                    </CollapsibleTrigger>
                    <CollapsibleContent>
                      <SidebarMenuSub>
                        {item.items.map((subItem) => (
                          <SidebarMenuSubItem key={subItem.href}>
                            <SidebarMenuSubButton
                              asChild
                              isActive={isNavActive(pathname, subItem.href, subItem.exact)}
                            >
                              <Link href={subItem.href}>
                                <span>{subItem.title}</span>
                              </Link>
                            </SidebarMenuSubButton>
                          </SidebarMenuSubItem>
                        ))}
                      </SidebarMenuSub>
                    </CollapsibleContent>
                  </SidebarMenuItem>
                </Collapsible>
              ) : (
                <SidebarMenuItem key={item.href}>
                  <SidebarMenuButton
                    asChild
                    isActive={isNavActive(pathname, item.href, item.exact)}
                    tooltip={item.title}
                  >
                    <Link href={item.href}>
                      <item.icon />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              )
            )}
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContentBase>
  )
}
