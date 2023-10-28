import request from "axios";

export function getError(error: any): string {
  if (!request.isAxiosError(error)) {
    return "";
  }

  return error.response?.data["error"];
}
