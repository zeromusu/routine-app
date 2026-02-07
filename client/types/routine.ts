export interface Routine {
  id: number;
  title: string;
  interval: string;
  created_at: string;
  updated_at: string;
}

export interface APIResponse<T> {
  success: boolean;
  data: T;
  error?: {
    code: string;
    message: string;
  };
}
