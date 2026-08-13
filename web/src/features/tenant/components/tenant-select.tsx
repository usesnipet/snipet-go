import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Loading } from "@/components/ui/loading";
import { useNavigate } from "@/hooks/use-navigate";
import { ROUTES } from "@/routes";
import { useEffect } from "react";

import { useFindMineTenant } from "../hooks";

import { TenantCard } from "./tenant-card";

export function TenantSelect() {
  const navigate = useNavigate();
  const { data, isLoading } = useFindMineTenant();
  const tenants = data?.data ?? [];
  const singleTenant = tenants.length === 1 ? tenants[0] : null;

  useEffect(() => {
    if (!singleTenant) return;
    navigate(ROUTES.tenantHome, {
      params: { tenantSlug: singleTenant.slug },
      replace: true,
    });
  }, [singleTenant, navigate]);

  return (
    <Card className="flex flex-col gap-3">
      <CardHeader>
        <CardTitle>
          Select a tenant
        </CardTitle>
        <CardDescription>
          Choose which workspace you want to enter.
        </CardDescription>
      </CardHeader>
      <CardContent>
        {(isLoading || singleTenant) && <Loading />}
        {
          !isLoading && !singleTenant && tenants.length === 0 &&
          <p className="text-sm text-muted-foreground text-center">
            You are not a member of any tenant yet.
          </p>
        }
        {
          !isLoading && !singleTenant && tenants.length > 0 &&
          <div className="flex flex-col gap-3">
            {tenants.map((tenant) => (
              <TenantCard key={tenant.id} tenant={tenant} />
            ))}
          </div>
        }
      </CardContent>
    </Card>
  )
}
