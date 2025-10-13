"use client";

import { useState } from "react";
import { PlusCircle, ChevronLeft, ChevronRight } from "lucide-react";
import URLCard from "./URLCard";
import { AddEditDialog, DeleteDialog } from "./URLDialog";
import { useURLs } from "@/hooks/useURLs";
import useURLQueryParams from "@/hooks/useURLQueryParams";
import { Input } from "../ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import { Button } from "../ui/button";
import { URL } from "@/types/uptime";
import { urlService } from "@/lib/services/url";

export default function URLList() {
  const { query, setQuery, filter, setFilter, sortBy, setSortBy, page, setPage } =
    useURLQueryParams();

  const { urls, totalPages, loading } = useURLs({ query, filter, sortBy, page });

  const [openAdd, setOpenAdd] = useState(false);
  const [editData, setEditData] = useState<URL | null>(null);
  const [deleteData, setDeleteData] = useState<URL | null>(null);

  const handleAdd = (data: URL) => {
    // await urlService.create(data);
    setOpenAdd(false);
  };
  const handleEdit = (data: URL) => {
    // await urlService.update(id, data);
    setEditData(null);
    setOpenAdd(false);
  };
  const handleDelete = (id: number) => {
    // await urlService.delete(id);
    setDeleteData(null);
  };

  return (
    <div className="py-12 max-w-7xl mx-auto">
      {/* Header + Filters */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
        <Input
          placeholder="Cari URL..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="w-full sm:w-1/3"
        />

        <div className="flex flex-wrap gap-3 items-center">
          <Select value={filter} onValueChange={(v: any) => setFilter(v)}>
            <SelectTrigger><SelectValue placeholder="Filter" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Semua</SelectItem>
              <SelectItem value="active">Aktif</SelectItem>
              <SelectItem value="inactive">Nonaktif</SelectItem>
            </SelectContent>
          </Select>

          <Select value={sortBy} onValueChange={(v: any) => setSortBy(v)}>
            <SelectTrigger><SelectValue placeholder="Sort By" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="label">Nama</SelectItem>
              <SelectItem value="created_at">Tanggal</SelectItem>
            </SelectContent>
          </Select>

          <Button onClick={() => setOpenAdd(true)} className="flex items-center gap-2">
            <PlusCircle /> Tambah URL
          </Button>
        </div>
      </div>

      {/* Grid */}
      {loading ? (
        <p className="text-center text-slate-500 mt-6">Loading...</p>
      ) : urls.length === 0 ? (
        <p className="text-center text-slate-500 mt-6">Tidak ada data</p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {urls.map((url: URL) => (
            <URLCard
              key={url.id}
              data={url}
              onDelete={() => setDeleteData(url)}
              onEdit={() => setEditData(url)}
              onToggle={(id) => console.log(id)}
            />
          ))}
        </div>
      )}


      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex justify-center gap-4 mt-6">
          <Button disabled={page <= 1} onClick={() => setPage(page - 1)}>
            <ChevronLeft />
          </Button>
          <span>{page} / {totalPages}</span>
          <Button disabled={page >= totalPages} onClick={() => setPage(page + 1)}>
            <ChevronRight />
          </Button>
        </div>
      )}

      {/* Dialogs */}
      <AddEditDialog open={openAdd} mode="add" onClose={() => setOpenAdd(false)} onSave={handleAdd} />
      {editData && <AddEditDialog open={!!editData} mode="edit" initialData={editData} onClose={() => setEditData(null)} onSave={handleEdit} />}
      {deleteData && <DeleteDialog open={!!deleteData} targetLabel={deleteData.label} onClose={() => setDeleteData(null)} onConfirm={() => handleDelete(deleteData.id)} />}
    </div>

  );
}
