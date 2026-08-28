import { CatalogList } from "@/components/catalog";
import { Loading } from "@/components/ui/loading";

import { useListKnowledge } from "../hooks";

import { KnowledgeCatalogItem } from "./knowledge-catalog-item";

export function KnowledgeList() {
  const { data, isLoading } = useListKnowledge();

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
