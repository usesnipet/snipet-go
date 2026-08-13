import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent } from "@/components/ui/card";
import { useNavigate } from "@/hooks/use-navigate";
import { cn } from "@/lib/utils";
import { ROUTES } from "@/routes";
import { Building2Icon } from "lucide-react";

import type { Tenant } from "../schemas";

export function TenantCard({ tenant, compact = false }: { tenant: Tenant; compact?: boolean }) {
  const navigate = useNavigate();

  const goToTenant = () => {
    navigate(ROUTES.tenantHome, { params: { tenantSlug: tenant.slug } });
  };

  return (
    <Card
      className={cn(
        "cursor-pointer transition-colors hover:bg-accent",
        compact && "border-none shadow-none",
      )}
      onClick={goToTenant}
    >
      <CardContent className={cn("flex items-center gap-3 p-4", compact && "gap-2 p-2")}>
        <Avatar size={compact ? "default" : "lg"}>
          {tenant.icon && <AvatarImage src={tenant.icon} alt={tenant.name} />}
          <AvatarFallback>
            <Building2Icon className="size-4" />
          </AvatarFallback>
        </Avatar>
        <div className="min-w-0 flex-1">
          <p className={cn("truncate font-semibold leading-tight", compact ? "text-xs" : "text-sm")}>
            {tenant.name}
          </p>
          <p className={cn("truncate text-muted-foreground", compact ? "text-xs" : "text-sm")}>
            {tenant.slug}
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
