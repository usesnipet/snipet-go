import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { CreateKnowledgeDialog } from "@/features/knowledge/components/create-knowledge-dialog";
import { KnowledgeList } from "@/features/knowledge/components/knowledge-list";
import { useDialog } from "@/lib/dialog/use-dialog";
import { Plus } from "lucide-react";

export function AdminKnowledgePage() {
  const { openDialog } = useDialog();

  const openCreate = () => {
    openDialog({
      component: CreateKnowledgeDialog,
      props: {},
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
