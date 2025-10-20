export type URLRequest = {
  label: string;
  url: string;
  active: boolean;
};

export type URLResponse = {
  id: number;
  label: string;
  url: string;
  interval: number;
  active: boolean;
  last_checked: string;
  created_at: string;
};

export type URLStats = {
  bucket_start: string;
  total_checks: number;
  up_checks: number;
  uptime_percent: number;
};
