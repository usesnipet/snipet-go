import { CatalogList } from "@/components/catalog";
import { Loading } from "@/components/ui/loading";
import { useTenantStore } from "@/features/tenant/store";

import { useListClient } from "../hooks";

import { ClientCatalogItem } from "./client-catalog-item";

export function ClientList() {
  const tenant = useTenantStore((state) => state.tenant);
  const { data, isLoading } = useListClient(tenant?.id ?? "");

  if (isLoading) {
    return (
      <div className="flex flex-1 items-center justify-center py-12">
        <Loading />
      </div>
    );
  }

  return (
    <CatalogList
      items={data?.data ?? []}
      emptyMessage="No clients yet."
      renderItem={client => <ClientCatalogItem client={client} />}
    />
  );
}
