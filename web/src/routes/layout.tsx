"use client"

import { AnimatedOutlet } from "@/components/animated-outlet";
import { Loading } from "@/components/ui/loading";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { TenantSidebar } from "@/features/tenant/components/tenant-sidebar";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useTenantStore } from "@/features/tenant/store";
import { useEffect } from "react";
import { useParams } from "react-router";

export function Layout() {
  const { tenantSlug } = useParams<{ tenantSlug: string }>();
  if (!tenantSlug) {
    throw new Error("Tenant slug is required");
  }

  const { data, isError, isLoading } = useFindBySlugTenant(tenantSlug);
  const setTenant = useTenantStore((state) => state.setTenant);
  const clearTenant = useTenantStore((state) => state.clearTenant);

  useEffect(() => {
    if (data) setTenant(data);
  }, [data, setTenant]);

  useEffect(() => {
    if (isError) clearTenant();
  }, [isError, clearTenant]);


  if (isLoading) {
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
