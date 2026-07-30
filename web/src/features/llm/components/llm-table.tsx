import { DataTable } from "@/components/data-table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useDialog } from "@/lib/dialog";
import { PencilIcon, TrashIcon } from "lucide-react";

import { useListLlm } from "../hooks";

import { DeleteLlmDialog } from "./delete-llm-dialog";
import { UpdateLlmDialog } from "./update-llm-dialog";

import type { DataTableColumn, DataTablePagination } from "@/components/data-table";
import type { Llm } from "../schemas";

function useLlmTableQuery(pagination: DataTablePagination) {
  return useListLlm({ searchParams: pagination })
}

export function LlmTable() {
  const { openDialog } = useDialog();

  const openEdit = (llm: Llm) => {
    openDialog({
      component: UpdateLlmDialog,
      props: { llm },
    });
  };

  const openDelete = (llm: Llm) => {
    openDialog({
      component: DeleteLlmDialog,
      props: { llm },
    });
  };

  const columns: DataTableColumn<Llm>[] = [
    {
      id: "name",
      header: "Name",
      cell: (row) => <span className="font-medium">{row.name}</span>,
    },
    {
      id: "provider",
      header: "Provider",
      cell: (row) => (
        <Badge variant="outline" className="font-normal">
          {row.provider}
        </Badge>
      ),
    },
    {
      id: "actions",
      header: <span className="sr-only">Actions</span>,
      headerClassName: "w-24",
      className: "text-right",
      cell: (row) => (
        <div className="flex justify-end">
          <Button
            variant="ghost"
            size="icon"
            aria-label="Edit LLM"
            onClick={() => openEdit(row)}
          >
            <PencilIcon className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            aria-label="Delete LLM"
            onClick={() => openDelete(row)}
          >
            <TrashIcon className="size-4" />
          </Button>
        </div>
      ),
    },
  ];

  return (
    <DataTable
      columns={columns}
      useQuery={useLlmTableQuery}
      getRowKey={(row) => row.id}
      emptyMessage="No LLMs yet."
    />
  );
}
