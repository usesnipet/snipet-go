"use client"

import { AnimatedOutlet } from "@/components/animated-outlet";
import { ClientSidebar } from "@/components/sidebar/client";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";

export function ClientLayout() {
  return (
    <SidebarProvider>
      <ClientSidebar />
      <SidebarInset className="min-h-dvh bg-sidebar p-4">
        <div className="border-border bg-background flex min-h-0 flex-1 overflow-hidden rounded-xl border shadow-sm">
          <AnimatedOutlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
