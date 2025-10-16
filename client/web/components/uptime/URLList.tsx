"use client";

import { useState } from "react";
import { PlusCircle, ChevronLeft, ChevronRight, Loader2, InfoIcon } from "lucide-react";
import URLCard from "./URLCard";
import { AddEditDialog, DeleteDialog } from "./URLDialog";
import { useURLs } from "@/hooks/useURLs";
import useURLQueryParams from "@/hooks/useURLQueryParams";
import { Input } from "../ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import { Button } from "../ui/button";
import { URLResponse } from "@/types/url";
import { urlService } from "@/lib/services/url";
import { toast } from "sonner";
import { getErrorMessage, getValidationErrors } from "@/utils/helper";
import { ErrorInput } from "@/types/response";

export default function URLList() {
  const { query, setQuery, filter, setFilter, sortBy, setSortBy, page, setPage } = useURLQueryParams();
  const { urls, totalPages, loading, setUrls } = useURLs({ query, filter, sortBy, page });

  const initial = {
    label: "",
    url: "",
    interval: 0,
    active: true,
  }
  const [openAdd, setOpenAdd] = useState(false);
  const [openEdit, setOpenEdit] = useState(false);
  const [payload, setPayload] = useState<Partial<URLResponse>>(initial);
  const [deleteData, setDeleteData] = useState<URLResponse | null>(null);
  const [errors, setErrors] = useState<ErrorInput[]>([]);

  const handleAdd = async () => {
    const res = await urlService.create(payload)
    if (!res.success) {
      if (res.error?.code === "VALIDATION_ERROR") {
        const details = getValidationErrors(res.error.fields);
        setErrors(details);
        toast.error(getErrorMessage(res.error?.code || ""));
      } else {
        toast.error(getErrorMessage(res.error?.code || ""));
      }
    } else {
      if (!res.data) return;
      setUrls([...urls, res.data]);
      setOpenAdd(false);
      toast.success("URL berhasil ditambahkan");  
    }
  };
  const handleEdit = async () => {
    if (!payload) return;

    const res = await urlService.update(payload.id, payload)
    if (!res.success) {
      if (res.error?.code === "VALIDATION_ERROR") {
        const details = getValidationErrors(res.error.fields);
        setErrors(details);
        toast.error(getErrorMessage(res.error?.code || ""));
      } else {
        toast.error(getErrorMessage(res.error?.code || ""));
      }
    } else {
      if (!res.data) return;
      const data = res.data;
      setPayload(initial);
      setOpenEdit(false);
      setUrls(urls.map((url) => (url.id === data.id ? data : url)));
      toast.success("URL berhasil diperbarui");
    }
  };
  const handleDelete = async (id: number) => {
    const res = await urlService.delete(id)
    if (!res.success) {
      toast.error(getErrorMessage(res.error?.code || ""));
    } else {
      setUrls(urls.filter((url) => url.id !== id));
      toast.success("URL berhasil dihapus");
    }
    setDeleteData(null);
  };

  const handleToggle = async (id: number) => {
    if (!urls) return;
    let url = urls.find((url) => url.id === id);
    if (!url) return;

    const res = await urlService.update(id, { active: !url.active, label: url.label, url: url.url });
    if (!res.success) {
      toast.error(getErrorMessage(res.error?.code || ""));
    } else {
      toast.success("URL berhasil diperbarui");
      setUrls(urls.map((url) => (url.id === id ? { ...url, active: !url.active } : url)));
    }
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

          <Button onClick={() => {
            setErrors([])
            setOpenAdd(true);
            setPayload(initial);
          }} className="flex items-center gap-2">
            <PlusCircle /> Tambah URL
          </Button>
        </div>
      </div>

      {/* Grid */}
      {loading ? (
        <div className="flex flex-col items-center justify-center min-h-[50vh]">
          <Loader2 className="w-6 h-6 text-slate-500 animate-spin" />
          <p className="text-center text-slate-500 mt-2">Loading...</p>
        </div>
      ) : urls.length === 0 ? (
        <div className="flex flex-col items-center justify-center min-h-[50vh]">
          <InfoIcon className="w-6 h-6 text-slate-500 mb-2" />
          <p className="text-center text-slate-500">Tidak ada URL</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 md:gap-6">
          {urls.map((url: URLResponse) => (
            <URLCard
              key={url.id}
              data={url}
              onDelete={() => setDeleteData(url)}
              onEdit={() => {
                setErrors([])
                setPayload(url);
                setOpenEdit(true)
              }}
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
      {openAdd && <AddEditDialog 
        mode="add"
        open={openAdd}
        onClose={() => setOpenAdd(false)}
        initialData={payload}
        onInputChange={setPayload}
        onSave={handleAdd} 
        errors={errors}
      />}
      {openEdit && <AddEditDialog 
        mode="edit"
        open={openEdit}
        onClose={() => {
          setErrors([])
          setOpenEdit(false);
        }}
        initialData={payload}
        onInputChange={setPayload}
        onSave={handleEdit}
        errors={errors}
      />}
      {deleteData && <DeleteDialog
        open={!!deleteData}
        targetLabel={deleteData.label}
        onClose={() => setDeleteData(null)}
        onConfirm={() => handleDelete(deleteData.id)}
      />}
    </div>

  );
}
