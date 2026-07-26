"use client"

import { AnimatedOutlet } from "@/components/animated-outlet";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { ClientChatSidebar } from "@/features/client/components/chat-sidebar";

export function ClientChatLayout() {
  return (
    <SidebarProvider>
      <ClientChatSidebar />
      <SidebarInset className="min-h-dvh bg-sidebar p-4">
        <div className="border-border bg-background flex min-h-0 flex-1 overflow-hidden rounded-xl border shadow-sm">
          <AnimatedOutlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
