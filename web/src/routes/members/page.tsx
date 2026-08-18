import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { useGetSystemInfo } from "@/features/app/hooks";
import { CreateMemberDialog } from "@/features/member/components/create-member-dialog";
import { MemberTable } from "@/features/member/components/member-table";
import { useTenantStore } from "@/features/tenant/store";
import { useNavigate } from "@/hooks/use-navigate";
import { useDialog } from "@/lib/dialog";
import { ROUTES } from "@/routes";
import { UserPlus } from "lucide-react";

export function MembersPage() {
  const navigate = useNavigate();
  const tenant = useTenantStore((state) => state.tenant);
  const { data: systemInfo } = useGetSystemInfo();
  const { openDialog } = useDialog();

  const goToInvite = () => {
    if (!tenant) return;
    navigate(ROUTES.inviteMember, { params: { tenantSlug: tenant.slug } });
  };

  const openCreate = () => {
    if (!tenant) return;
    openDialog({
      component: CreateMemberDialog,
      props: { tenantId: tenant.id },
    });
  };

  return (
    <Page
      title="Members"
      description="Manage who has access to this tenant and their roles."
      documentTitle="Members"
    >
      <PageActions>
        {systemInfo?.multi_tenant_enabled ? (
          <Button onClick={goToInvite}>
            <UserPlus />
            Invite member
          </Button>
        ) : (
          <Button onClick={openCreate} disabled={!tenant}>
            <UserPlus />
            Create member
          </Button>
        )}
      </PageActions>
      <MemberTable />
    </Page>
  )
}
