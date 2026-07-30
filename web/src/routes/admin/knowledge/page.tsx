import { Page } from "@/components/page";
import { KnowledgeList } from "@/features/knowledge/components/knowledge-list";

export function AdminKnowledgePage() {
  return (
    <Page
      title="Knowledge"
      description="Browse and manage knowledge sources."
      documentTitle="Knowledge"
    >
      <KnowledgeList />
    </Page>
  );
}
