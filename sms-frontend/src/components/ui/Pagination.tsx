import {
  ChevronLeft,
  ChevronRight,
} from "lucide-react";

import { Button } from "./Button";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./Select";

interface PaginationProps {
  page: number;
  totalPages: number;
  pageSize: number;
  totalCount: number;
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
}

export function Pagination({
  page,
  totalPages,
  pageSize,
  totalCount,
  onPageChange,
  onPageSizeChange,
}: PaginationProps) {
  const safeTotalPages = Math.max(
    totalPages,
    1
  );

  const start =
    totalCount === 0
      ? 0
      : (page - 1) * pageSize + 1;

  const end = Math.min(
    page * pageSize,
    totalCount
  );

  return (
    <div className="mt-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">

      <div className="text-sm text-muted-foreground">
        Showing {start}-{end} of {totalCount} students
      </div>

      <div className="flex flex-wrap items-center justify-end gap-3">

        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <span>Rows per page</span>

          <Select
            value={String(pageSize)}
            onValueChange={(value) =>
              onPageSizeChange(
                Number(value)
              )
            }
          >
            <SelectTrigger className="w-[80px]">
              <SelectValue />
            </SelectTrigger>

            <SelectContent>
              <SelectItem value="5">
                5
              </SelectItem>

              <SelectItem value="10">
                10
              </SelectItem>

              <SelectItem value="20">
                20
              </SelectItem>

              <SelectItem value="50">
                50
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="flex items-center gap-2">

          <Button
            variant="outline"
            size="sm"
            onClick={() =>
              onPageChange(page - 1)
            }
            disabled={page <= 1}
          >
            <ChevronLeft />
            Previous
          </Button>

          <span className="min-w-20 text-center text-sm font-medium">
            Page {page} of {safeTotalPages}
          </span>

          <Button
            variant="outline"
            size="sm"
            onClick={() =>
              onPageChange(page + 1)
            }
            disabled={
              page >= safeTotalPages
            }
          >
            Next
            <ChevronRight />
          </Button>

        </div>
      </div>
    </div>
  );
}
