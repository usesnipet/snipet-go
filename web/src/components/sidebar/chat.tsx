

import { Link } from "../ui/link";
import { Sidebar, SidebarHeader } from "../ui/sidebar";
import { ToggleTheme } from "../ui/toggle-theme";

export function ChatSidebar() {

  return (
    <Sidebar>
      <SidebarHeader className="flex flex-row justify-between items-center">
        <Link
          href="/"
          className="flex items-center gap-2 px-2 py-1 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0"
        >
          <img src="/favicon.svg" alt="Snipet" className="size-7 shrink-0" />
          <p className="text-sm font-semibold">Snipet</p>
        </Link>
        <ToggleTheme />
      </SidebarHeader>
    </Sidebar>
  )
}