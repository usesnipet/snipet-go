import { UserPlus } from "lucide-react";
import { useParams } from "react-router";

import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { useGetSystemInfo } from "@/features/app/hooks";
import { CreateMemberDialog } from "@/features/member/components/create-member-dialog";
import { MemberTable } from "@/features/member/components/member-table";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useNavigate } from "@/hooks/use-navigate";
import { useDialog } from "@/lib/dialog";
import { ROUTES } from "@/routes";

export function MembersPage() {
  const navigate = useNavigate();
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const { data: systemInfo } = useGetSystemInfo();
  const { openDialog } = useDialog();

  const goToInvite = () => {
    navigate(ROUTES.inviteMember, { params: { tenantSlug } });
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
