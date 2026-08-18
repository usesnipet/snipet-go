import { CatalogList } from "@/components/catalog";
import { Loading } from "@/components/ui/loading";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useParams } from "react-router";

import { useListKnowledge } from "../hooks";

import { KnowledgeCatalogItem } from "./knowledge-catalog-item";

export function KnowledgeList() {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const { data, isLoading } = useListKnowledge(tenant?.id ?? "");

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
      emptyMessage="No knowledge sources yet."
      renderItem={(knowledge) => <KnowledgeCatalogItem knowledge={knowledge} />}
    />
  );
}
