import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { CreateLlmDialog } from "@/features/llm/components/create-llm-dialog";
import { LlmTable } from "@/features/llm/components/llm-table";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useDialog } from "@/lib/dialog/use-dialog";
import { Plus } from "lucide-react";
import { useParams } from "react-router";

export function LLMPage() {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const { openDialog } = useDialog();

  const openCreate = () => {
    if (!tenant) return;
    openDialog({
      component: CreateLlmDialog,
      props: { tenantId: tenant.id },
    });
  };

  return (
    <Page
      title="LLMs"
      description="Create and manage language model provider configurations."
      documentTitle="LLMs"
    >
      <PageActions>
        <Button onClick={openCreate}>
          <Plus />
          Create LLM
        </Button>
      </PageActions>
      <LlmTable />
    </Page>
  );
}
