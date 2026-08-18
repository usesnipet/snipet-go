import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { AgentList } from "@/features/agent/components/agent-list";
import { CreateAgentDialog } from "@/features/agent/components/create-agent-dialog";
import { useTenantStore } from "@/features/tenant/store";
import { useDialog } from "@/lib/dialog/use-dialog";
import { Plus } from "lucide-react";

export function AgentsPage() {
  const tenant = useTenantStore((state) => state.tenant);
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
