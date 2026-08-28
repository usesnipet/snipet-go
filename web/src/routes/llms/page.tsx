import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { CreateLlmDialog } from "@/features/llm/components/create-llm-dialog";
import { LlmTable } from "@/features/llm/components/llm-table";
import { useDialog } from "@/lib/dialog/use-dialog";
import { Plus } from "lucide-react";

export function LLMPage() {
  const { openDialog } = useDialog();

  const openCreate = () => {
    openDialog({
      component: CreateLlmDialog,
      props: {},
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
