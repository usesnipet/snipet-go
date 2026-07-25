import { CatalogList } from "@/components/catalog";
import { Loading } from "@/components/ui/loading";

import { useListClient } from "../hooks";

import { ClientCatalogItem } from "./client-catalog-item";

export function ClientList() {
  const { data, isLoading } = useListClient();

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
