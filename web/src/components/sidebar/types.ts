import type { LucideIcon } from "lucide-react"

export type NavSubItem = {
  title: string
  href: string
  /** When true, only exact pathname matches (avoids /llm matching /llm/playground). */
  exact?: boolean
}

export type NavItem = {
  title: string
  href: string
  icon: LucideIcon
  exact?: boolean
}

export type NavItemWithChildren = {
  title: string
  icon: LucideIcon
  items: NavSubItem[]
}

export type NavEntry = NavItem | NavItemWithChildren
