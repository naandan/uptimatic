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
import { URLRequest, urlService } from "@/lib/services/url";
import { toast } from "sonner";

export default function URLList() {
  const { query, setQuery, filter, setFilter, sortBy, setSortBy, page, setPage } =
    useURLQueryParams();

  const { urls, totalPages, loading, setUrls } = useURLs({ query, filter, sortBy, page });

  const [openAdd, setOpenAdd] = useState(false);
  const [editData, setEditData] = useState<URL | null>(null);
  const [deleteData, setDeleteData] = useState<URL | null>(null);

  const handleAdd = (data: URLRequest) => {
    urlService.create(data)
    .then((res) => {
      setUrls([...urls, res.data]);
      setOpenAdd(false);
    })
    .catch((err) => {
      console.error(err);
      toast.error(err.response.data.message);
    });
  };
  const handleEdit = (data: URLRequest) => {
    if (!editData) return;
    urlService.update(editData.id, data)
    .then((res) => {
      setEditData(null);
      setUrls(urls.map((url) => (url.id === res.id ? res : url)));
      toast.success("URL berhasil diperbarui");
    })
    .catch((err) => {
      console.error(err);
      toast.error(err.response.data.message);
    });
  };
  const handleDelete = (id: number) => {
    urlService.delete(id)
    .then(() => {
      setUrls(urls.filter((url) => url.id !== id));
      toast.success("URL berhasil dihapus");
    })
    .catch((err) => {
      console.error(err);
      toast.error(err.response.data.message);
    });
    setDeleteData(null);
  };

  const handleToggle = (id: number) => {
    if (!urls) return;
    let url = urls.find((url) => url.id === id);
    if (!url) return;

    urlService.update(id, { active: !url.active, label: url.label, url: url.url });
    setUrls(urls.map((url) => (url.id === id ? { ...url, active: !url.active } : url)));
  }

  return (
    <div className="py-12 min-h-screen max-w-7xl mx-auto px-4">
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
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 md:gap-6">
          {urls.map((url: URL) => (
            <URLCard
              key={url.id}
              data={url}
              onDelete={() => setDeleteData(url)}
              onEdit={() => setEditData(url)}
              onToggle={(id) => handleToggle(id)}
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
