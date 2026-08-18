import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { AgentList } from "@/features/agent/components/agent-list";
import { CreateAgentDialog } from "@/features/agent/components/create-agent-dialog";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useDialog } from "@/lib/dialog/use-dialog";
import { Plus } from "lucide-react";
import { useParams } from "react-router";

export function AgentsPage() {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const { openDialog } = useDialog();

  const openCreate = () => {
    if (!tenant) return;
    openDialog({
      component: CreateAgentDialog,
      props: { tenantId: tenant.id },
    });
  };

  return (
    <Page
      title="Agents"
      description="Create and manage AI agents."
      documentTitle="Agents"
    >
      <PageActions>
        <Button onClick={openCreate}>
          <Plus />
          Create Agent
        </Button>
      </PageActions>
      <AgentList />
    </Page>
  );
}
