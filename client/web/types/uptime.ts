export type URL = {
    id: number;
    label: string;
    url: string;
    interval: number;
    active: boolean;
    last_checked: string;
    created_at: string;
};

export type Stats = {
    bucket_start: string;
    total_checks: number;
    up_checks: number;
    uptime_percent: number;
}