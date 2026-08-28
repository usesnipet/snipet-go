import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { Link } from "@/components/ui/link";
import { AgentList } from "@/features/agent/components/agent-list";
import { CreateAgentDialog } from "@/features/agent/components/create-agent-dialog";
import { useDialog } from "@/lib/dialog/use-dialog";
import { ROUTES } from "@/routes";
import { Play, Plus } from "lucide-react";

export function AgentsPage() {
  const { openDialog } = useDialog();

  const openCreate = () => {
    openDialog({
      component: CreateAgentDialog,
      props: {},
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
