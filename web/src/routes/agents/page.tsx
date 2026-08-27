import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { Link } from "@/components/ui/link";
import { AgentList } from "@/features/agent/components/agent-list";
import { CreateAgentDialog } from "@/features/agent/components/create-agent-dialog";
import { useTenantStore } from "@/features/tenant/store";
import { useDialog } from "@/lib/dialog/use-dialog";
import { ROUTES } from "@/routes";
import { Play, Plus } from "lucide-react";

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
        <Button variant="outline" asChild>
          <Link href={ROUTES.agentPlayground}>
            <Play />
            Playground
          </Link>
        </Button>
        <Button onClick={openCreate}>
          <Plus />
          Create Agent
        </Button>
      </PageActions>
      <AgentList />
    </Page>
  );
}
