"use client"

import { AnimatedOutlet } from "@/components/animated-outlet";
import { Loading } from "@/components/ui/loading";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { TenantSidebar } from "@/features/tenant/components/tenant-sidebar";
import { useFindMineTenant } from "@/features/tenant/hooks";
import { useTenantStore } from "@/features/tenant/store";
import { useNavigate } from "@/hooks/use-navigate";
import { logger } from "@/lib/logger";
import { ROUTES } from "@/routes";
import { useEffect } from "react";
import { useParams } from "react-router";

export function Layout() {
  const { tenantSlug } = useParams<{ tenantSlug: string }>();
  if (!tenantSlug) {
    throw new Error("Tenant slug is required");
  }

  const { data, isLoading, isError } = useFindMineTenant();
  const setTenant = useTenantStore((state) => state.setTenant);
  const clearTenant = useTenantStore((state) => state.clearTenant);
  const setMember = useTenantStore((state) => state.setMember);
  const clearMember = useTenantStore((state) => state.clearMember);
  const navigate = useNavigate();

  useEffect(() => {
    if (data) {
      const tenant = data.data.find((tenant) => tenant.slug === tenantSlug);
      if (!tenant) {
        logger.error(`Tenant ${tenantSlug} not found in mine tenants`);
        navigate(ROUTES.selectTenant);
        return;
      }
      setTenant(tenant);
      const member = tenant.members?.[0];
      if (!member) {
        logger.error(`Member not found for tenant ${tenant.slug}`);
        navigate(ROUTES.selectTenant);
        return;
      }
      setMember(member);
    }
  }, [data, setTenant, setMember, navigate, tenantSlug]);

  useEffect(() => {
    if (isError) {
      clearTenant();
      clearMember();
    }
  }, [isError, clearTenant, clearMember]);


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
