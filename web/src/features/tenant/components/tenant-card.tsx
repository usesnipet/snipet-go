import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent } from "@/components/ui/card";
import { useNavigate } from "@/hooks/use-navigate";
import { ROUTES } from "@/routes";
import { Building2Icon } from "lucide-react";

import type { Tenant } from "../schemas";

export function TenantCard({ tenant }: { tenant: Tenant }) {
  const navigate = useNavigate();

  const goToTenant = () => {
    navigate(ROUTES.tenantHome, { params: { tenantSlug: tenant.slug } });
  };

  return (
    <Card
      className="cursor-pointer transition-colors hover:bg-accent"
      onClick={goToTenant}
    >
      <CardContent className="flex items-center gap-3 p-4">
        <Avatar size="lg">
          {tenant.icon && <AvatarImage src={tenant.icon} alt={tenant.name} />}
          <AvatarFallback>
            <Building2Icon className="size-4" />
          </AvatarFallback>
        </Avatar>
        <div className="min-w-0 flex-1">
          <p className="truncate text-sm font-semibold leading-tight">{tenant.name}</p>
          <p className="truncate text-sm text-muted-foreground">{tenant.slug}</p>
        </div>
      </CardContent>
    </Card>
  );
}
