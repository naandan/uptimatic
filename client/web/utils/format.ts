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