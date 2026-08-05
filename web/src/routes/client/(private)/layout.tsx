"use client"

import { AnimatedOutlet } from "@/components/animated-outlet";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { ClientAdminSidebar } from "@/features/client/components/admin-sidebar";

export function ClientAdminLayout() {
  return (
    <SidebarProvider>
      <ClientAdminSidebar />
      <SidebarInset className="h-dvh overflow-hidden bg-sidebar p-4">
        <div className="border-border bg-background flex min-h-0 flex-1 overflow-hidden rounded-xl border shadow-sm">
          <AnimatedOutlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
