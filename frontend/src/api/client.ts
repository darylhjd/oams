import axios from "axios";
import { LoginResponse } from "./login";
import { UserMeResponse, UsersGetResponse } from "./user";
import { BatchData, BatchPostResponse, BatchPutResponse } from "./batch";
import { FileWithPath } from "@mantine/dropzone";
import { ClassesGetResponse } from "./class";
import { ClassManagersGetResponse } from "./class_manager";
import { ClassGroupsGetResponse } from "./class_group";
import { ClassGroupSessionsGetResponse } from "./class_group_session";
import { SessionEnrollmentsGetResponse } from "./session_enrollment";

export class APIClient {
  static _client = axios.create({
    withCredentials: true,
    baseURL: `${process.env.API_SERVER}/api/v1`,
  });

  static readonly _loginPath = "/login";
  static readonly _logoutPath = "/logout";
  static readonly _batchPath = "/batch";
  static readonly _usersPath = "/users";
  static readonly _userMePath = "/users/me";
  static readonly _classesPath = "/classes";
  static readonly _classManagersPath = "/class-managers";
  static readonly _classGroupsPath = "/class-groups";
  static readonly _classGroupSessionsPath = "/class-group-sessions";
  static readonly _sessionEnrollmentsPath = "/session-enrollments";

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

  static async usersGet(
    offset: number,
    limit: number,
  ): Promise<UsersGetResponse> {
    const { data } = await this._client.get<UsersGetResponse>(this._usersPath, {
      params: {
        offset: offset,
        limit: limit,
      },
    });
    return data;
  }

  static async userMe(): Promise<UserMeResponse> {
    const { data } = await this._client.get<UserMeResponse>(this._userMePath);
    return data;
  }

  static async classesGet(
    offset: number,
    limit: number,
  ): Promise<ClassesGetResponse> {
    const { data } = await this._client.get<ClassesGetResponse>(
      this._classesPath,
      {
        params: {
          offset: offset,
          limit: limit,
        },
      },
    );
    return data;
  }

  static async classManagersGet(
    offset: number,
    limit: number,
  ): Promise<ClassManagersGetResponse> {
    const { data } = await this._client.get<ClassManagersGetResponse>(
      this._classManagersPath,
      {
        params: {
          offset: offset,
          limit: limit,
        },
      },
    );
    return data;
  }

  static async classGroupsGet(
    offset: number,
    limit: number,
  ): Promise<ClassGroupsGetResponse> {
    const { data } = await this._client.get<ClassGroupsGetResponse>(
      this._classGroupsPath,
      {
        params: {
          offset: offset,
          limit: limit,
        },
      },
    );
    return data;
  }

  static async classGroupSessionsGet(
    offset: number,
    limit: number,
  ): Promise<ClassGroupSessionsGetResponse> {
    const { data } = await this._client.get<ClassGroupSessionsGetResponse>(
      this._classGroupSessionsPath,
      {
        params: {
          offset: offset,
          limit: limit,
        },
      },
    );
    return data;
  }

  static async sessionEnrollmentsGet(
    offset: number,
    limit: number,
  ): Promise<SessionEnrollmentsGetResponse> {
    const { data } = await this._client.get<SessionEnrollmentsGetResponse>(
      this._sessionEnrollmentsPath,
      {
        params: {
          offset: offset,
          limit: limit,
        },
      },
    );
    return data;
  }
}
