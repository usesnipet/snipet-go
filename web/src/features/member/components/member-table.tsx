import { DataTable } from "@/components/data-table";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { DateFormat } from "@/components/ui/date";
import {
  DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger
} from "@/components/ui/dropdown-menu";
import { useTenantStore } from "@/features/tenant/store";
import { useDialog } from "@/lib/dialog";
import { MoreHorizontal, Shield, UserIcon, UserMinus } from "lucide-react";

import { useFilterMember } from "../hooks";

import { RemoveMemberDialog } from "./remove-member-dialog";
import { UpdateMemberRoleDialog } from "./update-member-role-dialog";

import type { DataTableColumn, DataTablePagination } from "@/components/data-table";
import type { Member } from "../schemas";
function useMemberTableQuery(pagination: DataTablePagination) {
  const tenant = useTenantStore((state) => state.tenant);
  return useFilterMember(tenant?.id ?? "", { searchParams: pagination });
}

export function MemberTable() {
  const tenant = useTenantStore((state) => state.tenant);
  const { openDialog } = useDialog();

  const openUpdateRole = (member: Member) => {
    if (!tenant) return;
    openDialog({
      component: UpdateMemberRoleDialog,
      props: { tenantId: tenant.id, member },
    });
  };

  const openRemove = (member: Member) => {
    if (!tenant) return;
    openDialog({
      component: RemoveMemberDialog,
      props: { tenantId: tenant.id, member },
    });
  };

  const columns: DataTableColumn<Member>[] = [
    {
      id: "name",
      header: "Name",
      cell: (row) => (
        <div className="flex items-center gap-3">
          <Avatar size="sm">
            {row.user?.picture && <AvatarImage src={row.user.picture} alt={row.user.name} />}
            <AvatarFallback>
              <UserIcon className="size-4" />
            </AvatarFallback>
          </Avatar>
          <div className="flex flex-col gap-0.5">
            <span className="font-medium">{row.user?.name}</span>
            <span className="text-xs text-muted-foreground">{row.user?.email}</span>
          </div>
        </div>
      ),
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
      cell: (row) => (
        <Badge variant={row.is_active ? "default" : "secondary"}>
          {row.is_active ? "Active" : "Inactive"}
        </Badge>
      ),
    },
    {
      id: "created_at",
      header: "Joined",
      cell: (row) => (
        <DateFormat className="text-muted-foreground" date={row.created_at} />
      ),
    },
    {
      id: "actions",
      header: <span className="sr-only">Actions</span>,
      headerClassName: "w-12",
      className: "text-right",
      cell: (row) => (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon-sm">
              <MoreHorizontal />
              <span className="sr-only">Open menu</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => openUpdateRole(row)}>
              <Shield />
              Update role
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              onClick={() => openRemove(row)}
            >
              <UserMinus />
              Remove
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ),
    },
  ];

  return (
    <DataTable
      columns={columns}
      useQuery={useMemberTableQuery}
      getRowKey={(row) => row.id}
      emptyMessage="No members yet."
    />
  )
}
