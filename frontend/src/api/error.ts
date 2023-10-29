import { isAxiosError } from "axios";

export type ErrorResponse = {
  error: string;
};

export function getError(error: any): string {
  if (!isAxiosError(error) || !error.response) {
    return "";
  }

  return (error.response.data as ErrorResponse).error;
}
