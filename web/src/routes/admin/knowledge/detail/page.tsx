import { ErrorFallback } from "@/components/error-fallback";
import { LoadingFallback } from "@/components/loading-fallback";
import { Page, PageActions, PageLeftActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { KnowledgeDetails } from "@/features/knowledge/components/knowledge-details";
import { KnowledgeItemsTable } from "@/features/knowledge/components/knowledge-items-table";
import { UpdateKnowledgeDialog } from "@/features/knowledge/components/update-knowledge-dialog";
import { useKnowledge, useSyncKnowledge } from "@/features/knowledge/hooks";
import { useNavigate } from "@/hooks/use-navigate";
import { useDialog } from "@/lib/dialog";
import { ArrowLeftIcon, PencilIcon, RefreshCw, RotateCcw } from "lucide-react";
import { useParams } from "react-router";

export function AdminKnowledgeDetailPage() {
  const { id = "" } = useParams<{ id: string }>();
  const { openDialog } = useDialog();
  const navigate = useNavigate();
  const { data: knowledge, isLoading, error } = useKnowledge(id);
  const { mutate: sync, isPending: isSyncing } = useSyncKnowledge();

  // useEffect(() => {
  //   if (!knowledge?.last_synced_at) return;
  //   queryClient.invalidateQueries({
  //     queryKey: listKnowledgeItemsQueryKey(knowledge.id),
  //   });
  // }, [knowledge.id, knowledge?.last_synced_at]);

  if (isLoading) {
    return <LoadingFallback />;
  }

  if (error || !knowledge) {
    return <ErrorFallback error={error ?? new Error("Knowledge not found")} />;
  }

  const isSyncBusy =
    knowledge.sync_status === "pending" || knowledge.sync_status === "in_progress";

  const openEdit = () => {
    openDialog({
      component: UpdateKnowledgeDialog,
      props: { knowledge },
    });
  };

  return (
    <Page
      title={knowledge.name}
      description={knowledge.description || `Driver: ${knowledge.driver}`}
      documentTitle={knowledge.name}
    >
      <PageLeftActions>
        <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
          <ArrowLeftIcon />
        </Button>
      </PageLeftActions>
      <PageActions>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            disabled={isSyncBusy || isSyncing}
            onClick={() => sync({ id: knowledge.id })}
          >
            <RefreshCw />
            Sync
          </Button>
          <Button
            variant="outline"
            size="sm"
            disabled={isSyncBusy || isSyncing}
            onClick={() => sync({ id: knowledge.id, force: true })}
          >
            <RotateCcw />
            Full resync
          </Button>
          <Button variant="outline" size="sm" onClick={openEdit}>
            <PencilIcon />
            Edit
          </Button>
        </div>
      </PageActions>

      <div className="grid min-h-0 grid-cols-1 gap-4 md:grid-cols-[1fr_3fr]">
        <KnowledgeDetails knowledge={knowledge} />
        <KnowledgeItemsTable />
      </div>
    </Page>
  );
}
