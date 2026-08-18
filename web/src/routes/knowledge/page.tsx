import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { CreateKnowledgeDialog } from "@/features/knowledge/components/create-knowledge-dialog";
import { KnowledgeList } from "@/features/knowledge/components/knowledge-list";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useDialog } from "@/lib/dialog/use-dialog";
import { Plus } from "lucide-react";
import { useParams } from "react-router";

export function KnowledgePage() {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const { openDialog } = useDialog();

  const openCreate = () => {
    if (!tenant) return;
    openDialog({
      component: CreateKnowledgeDialog,
      props: { tenantId: tenant.id },
    });
  };

  return (
    <Page
      title="Knowledge"
      description="Browse and manage knowledge sources."
      documentTitle="Knowledge"
    >
      <PageActions>
        <Button onClick={openCreate}>
          <Plus />
          Create Knowledge
        </Button>
      </PageActions>
      <KnowledgeList />
    </Page>
  );
}
