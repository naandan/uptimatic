import { useEffect, useState } from "react";
import type { URL } from "@/types/uptime";

export function useURLs({ query, filter, sortBy, page }: any) {
  const [urls, setUrls] = useState<URL[]>([]);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    const params = new URLSearchParams();
    if (query) params.set("q", query);
    if (filter) params.set("filter", filter);
    if (sortBy) params.set("sort", sortBy);
    if (page) params.set("page", page.toString());

    fetch(`/api/urls?${params.toString()}`)
      .then((res) => res.json())
      .then((data) => {
        setUrls(data.urls);
        setTotalPages(data.totalPages);
      })
      .finally(() => setLoading(false));
  }, [query, filter, sortBy, page]);

  return { urls, totalPages, loading };
}
