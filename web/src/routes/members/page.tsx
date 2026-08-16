import { UserPlus } from "lucide-react";
import { useParams } from "react-router";

import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { MemberTable } from "@/features/member/components/member-table";
import { useNavigate } from "@/hooks/use-navigate";
import { ROUTES } from "@/routes";

export function MembersPage() {
  const navigate = useNavigate();
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();

  const goToInvite = () => {
    navigate(ROUTES.inviteMember, { params: { tenantSlug } });
  };

  return (
    <Page
      title="Members"
      description="Manage who has access to this tenant and their roles."
      documentTitle="Members"
    >
      <PageActions>
        <Button onClick={goToInvite}>
          <UserPlus />
          Invite member
        </Button>
      </PageActions>
      <MemberTable />
    </Page>
  )
}
