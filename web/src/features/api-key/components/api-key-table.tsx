import { DataTable } from "@/components/data-table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { DateFormat } from "@/components/ui/date";
import {
  DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger
} from "@/components/ui/dropdown-menu";
import { CalendarClock, MoreHorizontal, RefreshCw, Trash2 } from "lucide-react";

import type { DataTableColumn } from "@/components/data-table";
import type { ApiKey } from "../schemas";

type ApiKeyTableProps = {
  data: ApiKey[]
  onUpdateExpiration: (apiKey: ApiKey) => void
  onRoll: (apiKey: ApiKey) => void
  onDelete: (apiKey: ApiKey) => void
}

export function ApiKeyTable({
  data,
  onUpdateExpiration,
  onRoll,
  onDelete,
}: ApiKeyTableProps) {
  const columns: DataTableColumn<ApiKey>[] = [
    {
      id: "name",
      header: "Name",
      cell: (row) => (
        <div className="flex flex-col gap-0.5">
          <span className="font-medium">{row.name}</span>
          <span className="font-mono text-xs text-muted-foreground">{row.key_id}</span>
        </div>
      ),
    },
    {
      id: "status",
      header: "Status",
      cell: (row) => (
        <Badge variant={row.active ? "default" : "secondary"}>
          {row.active ? "Active" : "Disabled"}
        </Badge>
      ),
    },
    {
      id: "expires_at",
      header: "Expires",
      cell: (row) => (
        <DateFormat className="text-muted-foreground" emptyValue="Never" date={row.expires_at} />
      ),
    },
    {
      id: "created_at",
      header: "Created",
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
            <DropdownMenuItem onClick={() => onUpdateExpiration(row)}>
              <CalendarClock />
              Update expiration
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => onRoll(row)}>
              <RefreshCw />
              Roll key
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              onClick={() => onDelete(row)}
            >
              <Trash2 />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ),
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={data}
      getRowKey={(row) => row.id}
      emptyMessage="No API keys yet."
    />
  )
}
