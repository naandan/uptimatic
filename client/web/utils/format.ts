export const FormatDateTimeGMT7 = (date: Date) => {
    const utcDate = new Date(date.toUTCString());
    const gmt7Date = new Date(utcDate.getTime() + 7 * 60 * 60 * 1000);
    return gmt7Date.toLocaleString("en-US", {
        day: "numeric",
        month: "long",
        year: "numeric",
        hour: "numeric",
        minute: "numeric",
        second: "numeric",
    });
};

export const FormatDateGMT7 = (date: Date) => {
    const utcDate = new Date(date.toUTCString());
    const gmt7Date = new Date(utcDate.getTime() + 7 * 60 * 60 * 1000);
    return gmt7Date.toLocaleString("en-US", {
        day: "numeric",
        month: "long",
        year: "numeric",
    });
};
