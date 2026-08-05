"use client"

import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { Link } from "@/components/ui/link";
import {
  SidebarContent as SidebarContentBase, SidebarGroup, SidebarGroupContent, SidebarMenu, SidebarMenuButton,
  SidebarMenuItem, SidebarMenuSub, SidebarMenuSubButton, SidebarMenuSubItem
} from "@/components/ui/sidebar";
import { applyPathParams } from "@/lib/http";
import { ChevronRight } from "lucide-react";
import { useCallback, useMemo } from "react";
import { useLocation, useParams } from "react-router";

import { isNavActive, isNavGroupActive, isNavItemWithChildren } from "./utils";

import type { NavEntry } from "./types";
type Props = {
  navItems: NavEntry[]
}

export function SidebarContent({ navItems }: Props) {
  const { pathname } = useLocation();
  const params = useParams();

  const buildPath = useCallback((path: string) => {
    const pathParams = Object.fromEntries(
      Object.entries(params).filter((entry): entry is [string, string] => entry[1] != null),
    );
    return applyPathParams(path, pathParams);
  }, [params])

  const transformedNavItems = useMemo(() => {
    const transform = (item: NavEntry): NavEntry => {
      if (isNavItemWithChildren(item)) {
        return {
          ...item,
          items: item.items.map(sub => ({
            ...sub,
            href: buildPath(sub.href),
          })),
        }
      }
      return {
        ...item,
        href: buildPath(item.href),
      }
    }
    return navItems.map((item) => transform(item))
  }, [navItems, buildPath]);

  return (
    <SidebarContentBase>
      <SidebarGroup>
        <SidebarGroupContent>
          <SidebarMenu>
            {transformedNavItems.map((item) =>
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
