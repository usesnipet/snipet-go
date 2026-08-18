import { CatalogList } from "@/components/catalog";
import { Loading } from "@/components/ui/loading";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useParams } from "react-router";

import { useListAgent } from "../hooks";

import { AgentCatalogItem } from "./agent-catalog-item";

export function AgentList() {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const { data, isLoading } = useListAgent(tenant?.id ?? "");

  if (isLoading) {
    return (
      <div className="flex flex-1 items-center justify-center py-12">
        <Loading />
      </div>
    );
  }

  return (
    <CatalogList
      size="md"
      items={data?.data ?? []}
      emptyMessage="No agents yet."
      renderItem={(agent) => <AgentCatalogItem agent={agent} />}
    />
  );
}
