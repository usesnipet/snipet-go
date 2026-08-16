import { UserPlus } from "lucide-react";
import { useState } from "react";
import { useParams } from "react-router";

import { Page, PageActions } from "@/components/page";
import { Button } from "@/components/ui/button";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { InvitationTable } from "@/features/invitation/components/invitation-table";
import { InviteMemberDialog } from "@/features/invitation/components/invite-member-dialog";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useDialog } from "@/lib/dialog";

import type { InvitationStatusFilter } from "@/features/invitation/schemas";

export function InviteMembersPage() {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const { openDialog } = useDialog();
  const [status, setStatus] = useState<InvitationStatusFilter>("pending");

  const openInvite = () => {
    if (!tenant) return;
    openDialog({
      component: InviteMemberDialog,
      props: { tenantId: tenant.id },
    });
  };

  return (
    <Page
      title="Invitations"
      description="Invite new members and manage pending, expired, accepted, and declined invitations."
      documentTitle="Invitations"
    >
      <PageActions>
        <Button onClick={openInvite} disabled={!tenant}>
          <UserPlus />
          Invite member
        </Button>
      </PageActions>
      <div className="flex min-h-0 flex-1 flex-col gap-3">
        <Tabs value={status} onValueChange={(value) => setStatus(value as InvitationStatusFilter)}>
          <TabsList>
            <TabsTrigger value="pending">Pending</TabsTrigger>
            <TabsTrigger value="expired">Expired</TabsTrigger>
            <TabsTrigger value="accepted">Accepted</TabsTrigger>
            <TabsTrigger value="declined">Declined</TabsTrigger>
          </TabsList>
        </Tabs>
        <InvitationTable key={status} status={status} />
      </div>
    </Page>
  )
}
