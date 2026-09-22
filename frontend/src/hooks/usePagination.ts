import { useState, useMemo } from "react";

// usePagination：列表页统一分页口径。
export function usePagination<T>(rows: T[] = [], pageSize = 8) {
  const [page, setPage] = useState(1);
  const pageRows = useMemo(
    () => rows.slice((page - 1) * pageSize, page * pageSize),
    [rows, page, pageSize]
  );
  return { page, setPage, pageSize, pageRows, total: rows.length };
}
