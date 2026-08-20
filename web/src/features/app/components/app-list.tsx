import { CatalogList } from "@/components/catalog";
import { Loading } from "@/components/ui/loading";
import { useTenantStore } from "@/features/tenant/store";

import { useListApp } from "../hooks";

import { AppCatalogItem } from "./app-catalog-item";

export function AppList() {
  const tenant = useTenantStore((state) => state.tenant);
  const { data, isLoading } = useListApp(tenant?.id ?? "");

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
      emptyMessage="No apps yet."
      renderItem={app => <AppCatalogItem app={app} />}
    />
  );
}
