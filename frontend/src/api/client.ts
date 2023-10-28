import axios from "axios";
import { LoginResponse } from "./login";
import { UserMeResponse } from "./user";
import { BatchData, BatchPostResponse, BatchPutResponse } from "./batch";
import { FileWithPath } from "@mantine/dropzone";

export class APIClient {
  static _client = axios.create({
    withCredentials: true,
    baseURL: `${process.env.API_SERVER}/api/v1`,
  });

  static readonly _loginPath = "/login";
  static readonly _logoutPath = "/logout";
  static readonly _batchPath = "/batch";
  static readonly _userMePath = "/users/me";

  static async login(redirectUrl: string = ""): Promise<string> {
    const { data } = await this._client.get<LoginResponse>(this._loginPath, {
      params: {
        redirect_url: redirectUrl ? redirectUrl : process.env.WEB_SERVER,
      },
    });
    return data.redirect_url;
  }

  static async logout(): Promise<boolean> {
    await this._client.get(this._logoutPath);
    return true;
  }

  static async batchPost(files: FileWithPath[]): Promise<BatchPostResponse> {
    const form = new FormData();
    files.forEach((file) => form.append("attachments", file));

    const { data } = await this._client.post<BatchPostResponse>(
      this._batchPath,
      form,
    );
    return data;
  }

  static async batchPut(batchData: BatchData[]): Promise<BatchPutResponse> {
    const { data } = await this._client.put<BatchPutResponse>(this._batchPath, {
      batches: batchData,
    });
    return data;
  }

  static async userMe(): Promise<UserMeResponse> {
    const { data } = await this._client.get<UserMeResponse>(this._userMePath);
    return data;
  }
}
