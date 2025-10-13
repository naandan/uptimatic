"use client";

import { useSearchParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function useURLQueryParams() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const [query, setQuery] = useState(searchParams.get("q") || "");
  const [filter, setFilter] = useState<"all" | "active" | "inactive">(
    (searchParams.get("filter") as any) || "all"
  );
  const [sortBy, setSortBy] = useState<"label" | "created_at">(
    (searchParams.get("sort") as any) || "label"
  );
  const [page, setPage] = useState(Number(searchParams.get("page") || 1));

  // Update URL ketika state berubah
  useEffect(() => {
    const params = new URLSearchParams();
    if (query) params.set("q", query);
    if (filter && filter !== "all") params.set("filter", filter);
    if (sortBy && sortBy !== "label") params.set("sort", sortBy);
    if (page && page !== 1) params.set("page", page.toString());

    router.replace(`?${params.toString()}`);
  }, [query, filter, sortBy, page, router]);

  return { query, setQuery, filter, setFilter, sortBy, setSortBy, page, setPage };
}
