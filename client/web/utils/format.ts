export const formatYearGMT7 = (dateStr: string) => {
  const date = new Date(dateStr);
  return date.toLocaleDateString("id-ID", {
    year: "numeric",
  });
};

export const formatMonthGMT7 = (dateStr: string) => {
  const date = new Date(dateStr);
  return date.toLocaleDateString("id-ID", {
    month: "long",
    year: "numeric",
  });
};

export const formatDateGMT7 = (dateStr: string) => {
  const date = new Date(dateStr);
  return date.toLocaleDateString("id-ID", {
    day: "2-digit",
    month: "long",
    year: "numeric",
  });
};

export const formatTimeGMT7 = (dateStr: string) => {
  const date = new Date(dateStr);
  return date.toLocaleTimeString("en-GB", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });
};

export const formatDateTimeGMT7 = (dateStr: string) => {
  return `${formatDateGMT7(dateStr)} ${formatTimeGMT7(dateStr)}`;
};

export const formatHourShort = (iso: string) =>
  new Date(iso).toLocaleTimeString("id-ID", {
    hour: "2-digit",
    minute: "2-digit",
  });

export const formatDayShort = (iso: string) =>
  new Date(iso).toLocaleDateString("id-ID", { day: "2-digit", month: "short" });

export const formatMonthShort = (iso: string) =>
  new Date(iso).toLocaleDateString("id-ID", {
    month: "short",
    year: "numeric",
  });
