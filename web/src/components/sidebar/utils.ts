import type { NavEntry, NavItemWithChildren, NavSubItem } from "./types";

export function isNavItemWithChildren(item: NavEntry): item is NavItemWithChildren {
  return "items" in item
}

export function isNavActive(pathname: string, href: string, exact?: boolean) {
  if (href === "/") return pathname === "/";
  if (exact) return pathname === href;
  return pathname === href || pathname.startsWith(`${href}/`);
}

export function isNavGroupActive(pathname: string, items: NavSubItem[]) {
  return items.some((item) => isNavActive(pathname, item.href, item.exact))
}