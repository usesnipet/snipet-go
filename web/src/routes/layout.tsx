"use client"

import { AnimatedOutlet } from "@/components/animated-outlet";
import { Loading } from "@/components/ui/loading";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { TenantSidebar } from "@/features/tenant/components/tenant-sidebar";
import { useCurrentTenant } from "@/features/tenant/hooks";

export function Layout() {
  const { isLoading: isLoadingTenant } = useCurrentTenant();

  if (isLoadingTenant) {
    return (
      <div className="flex h-dvh items-center justify-center">
        <Loading />
      </div>
    );
  }

  return (
    <SidebarProvider>
      <TenantSidebar />
      <SidebarInset className="h-dvh overflow-hidden bg-sidebar p-4">
        <div className="border-border bg-background flex min-h-0 flex-1 overflow-hidden rounded-xl border shadow-sm">
          <AnimatedOutlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
