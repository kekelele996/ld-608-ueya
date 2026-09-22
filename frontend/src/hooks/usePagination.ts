import { useState, useMemo } from "react";

// usePagination is a tiny client-side pagination used by list pages.
export function usePagination<T>(rows: T[] = [], pageSize = 8) {
  const [page, setPage] = useState(1);
  const pageRows = useMemo(
    () => rows.slice((page - 1) * pageSize, page * pageSize),
    [rows, page, pageSize],
  );
  return { page, setPage, pageSize, pageRows, total: rows.length };
}
