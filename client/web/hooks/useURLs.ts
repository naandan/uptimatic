import { useEffect, useState } from "react";
import type { URL } from "@/types/uptime";
import { urlService } from "@/lib/services/url";

export function useURLs({ query, filter, sortBy, page }: any) {
  const [urls, setUrls] = useState<URL[]>([]);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    const params = new URLSearchParams();
    if (query) params.set("q", query);
    if (filter) params.set("active", filter);
    if (sortBy) params.set("sort", sortBy);
    if (page) params.set("page", page.toString());

    urlService.list(params.toString()).then((data) => {
      setUrls(data.data);
      setTotalPages(data.meta.total_page);
      setLoading(false);
    })

  }, [query, filter, sortBy, page]);

  return { urls, totalPages, loading, setUrls };
}
