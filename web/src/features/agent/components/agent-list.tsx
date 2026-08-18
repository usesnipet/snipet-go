import { CatalogList } from "@/components/catalog";
import { Loading } from "@/components/ui/loading";
import { useTenantStore } from "@/features/tenant/store";

import { useListAgent } from "../hooks";

import { AgentCatalogItem } from "./agent-catalog-item";

export function AgentList() {
  const tenant = useTenantStore((state) => state.tenant);
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
