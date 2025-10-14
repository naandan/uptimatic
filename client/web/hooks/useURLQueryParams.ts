"use client";

import { useState } from "react";

export default function useURLQueryParams() {

  const [query, setQuery] = useState("");
  const [filter, setFilter] = useState<"all" | "active" | "inactive">("all");
  const [sortBy, setSortBy] = useState<"label" | "created_at">("label");
  const [page, setPage] = useState(1);

  return { query, setQuery, filter, setFilter, sortBy, setSortBy, page, setPage };
}
