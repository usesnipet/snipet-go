import { DataTable } from "@/components/data-table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { DateFormat } from "@/components/ui/date";
import {
  DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger
} from "@/components/ui/dropdown-menu";
import { useTenantStore } from "@/features/tenant/store";
import { useDialog } from "@/lib/dialog";
import { MoreHorizontal, XCircle } from "lucide-react";

import { useFilterInvitation } from "../hooks";

import { CancelInvitationDialog } from "./cancel-invitation-dialog";

import type { DataTableColumn, DataTablePagination } from "@/components/data-table";
import type { Invitation, InvitationStatusFilter } from "../schemas";
const statusVariant: Record<InvitationStatusFilter, "default" | "secondary" | "destructive"> = {
  pending: "secondary",
  accepted: "default",
  declined: "destructive",
  expired: "destructive",
};

const statusLabel: Record<InvitationStatusFilter, string> = {
  pending: "Pending",
  accepted: "Accepted",
  declined: "Declined",
  expired: "Expired",
};

function getDisplayStatus(invitation: Invitation): InvitationStatusFilter {
  if (invitation.status === "pending" && invitation.expires_at < new Date()) {
    return "expired";
  }
  return invitation.status;
}

type InvitationTableProps = {
  status?: InvitationStatusFilter
}

export function InvitationTable({ status }: InvitationTableProps) {
  const tenant = useTenantStore((state) => state.tenant);
  const { openDialog } = useDialog();

  const useInvitationTableQuery = (pagination: DataTablePagination) =>
    useFilterInvitation(tenant?.id ?? "", { searchParams: { ...pagination, status } });

  const openCancel = (invitation: Invitation) => {
    if (!tenant) return;
    openDialog({
      component: CancelInvitationDialog,
      props: { tenantId: tenant.id, invitation },
    });
  };

  const columns: DataTableColumn<Invitation>[] = [
    {
      id: "email",
      header: "Email",
      cell: (row) => <span className="font-medium">{row.email}</span>,
    },
    {
      id: "role",
      header: "Role",
      cell: (row) => (
        <Badge variant={row.role === "admin" ? "default" : "secondary"}>
          {row.role === "admin" ? "Admin" : "Member"}
        </Badge>
      ),
    },
    {
      id: "status",
      header: "Status",
      cell: (row) => {
        const displayStatus = getDisplayStatus(row);
        return <Badge variant={statusVariant[displayStatus]}>{statusLabel[displayStatus]}</Badge>;
      },
    },
    {
      id: "expires_at",
      header: "Expires",
      cell: (row) => (
        <DateFormat className="text-muted-foreground" date={row.expires_at} />
      ),
    },
    {
      id: "created_at",
      header: "Sent",
      cell: (row) => (
        <DateFormat className="text-muted-foreground" date={row.created_at} />
      ),
    },
    {
      id: "actions",
      header: <span className="sr-only">Actions</span>,
      headerClassName: "w-12",
      className: "text-right",
      cell: (row) => {
        if (row.status !== "pending") return null;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon-sm">
                <MoreHorizontal />
                <span className="sr-only">Open menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem
                className="text-destructive focus:text-destructive"
                onClick={() => openCancel(row)}
              >
                <XCircle />
                Cancel invitation
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ];

  return (
    <DataTable
      columns={columns}
      useQuery={useInvitationTableQuery}
      getRowKey={(row) => row.id}
      emptyMessage="No invitations yet."
    />
  )
}
