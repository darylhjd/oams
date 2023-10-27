import axios from "axios";
import {
  BatchPostResponse,
  BatchPutRequest,
  BatchPutResponse,
  ClassGroupSessionsGetResponse,
  ClassGroupsGetResponse,
  ClassManagersGetResponse,
  ClassesGetResponse,
  LoginResponse,
  SessionEnrollmentsGetResponse,
  UserMeResponse,
  UsersGetResponse,
} from "./models";

export class APIClient {
  static _client = axios.create({
    withCredentials: true,
    baseURL: `${process.env.API_SERVER_HOST}/api/v1`,
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

  static async getLoginUrl(returnTo: string): Promise<string> {
    if (returnTo.length == 0) {
      returnTo = `${process.env.WEB_SERVER_HOST}:${process.env.WEB_SERVER_PORT}`;
    }

    try {
      const { data } = await this._client.get<LoginResponse>(this._loginPath, {
        params: {
          redirect_url: returnTo,
        },
      });
      return data.redirect_url;
    } catch (error) {
      throw error;
    }
  }

  static async logout(): Promise<boolean> {
    try {
      await this._client.get(this._logoutPath);
      return true;
    } catch (error) {
      return false;
    }
  }

  static async batchPost(files: File[]): Promise<BatchPostResponse | null> {
    const form = new FormData();
    for (var file of files) {
      form.append("attachments", file);
    }

    try {
      const { data } = await this._client.post<BatchPostResponse>(
        this._batchPath,
        form,
      );
      return data;
    } catch (error) {
      return null;
    }
  }

  static async batchPut(
    req: BatchPutRequest,
  ): Promise<BatchPutResponse | null> {
    try {
      const { data } = await this._client.put<BatchPutResponse>(
        this._batchPath,
        req,
      );
      return data;
    } catch (error) {
      return null;
    }
  }

  static async usersGet(
    offset: number,
    limit: number,
  ): Promise<UsersGetResponse | null> {
    try {
      const { data } = await this._client.get<UsersGetResponse>(
        this._usersPath,
        {
          params: {
            offset: offset,
            limit: limit,
          },
        },
      );
      return data;
    } catch (error) {
      return null;
    }
  }

  static async getUserMe(): Promise<UserMeResponse | null> {
    try {
      const { data } = await this._client.get<UserMeResponse>(this._userMePath);
      return data;
    } catch (error) {
      return null;
    }
  }

  static async classesGet(
    offset: number,
    limit: number,
  ): Promise<ClassesGetResponse | null> {
    try {
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
    } catch (error) {
      return null;
    }
  }

  static async classManagersGet(
    offset: number,
    limit: number,
  ): Promise<ClassManagersGetResponse | null> {
    try {
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
    } catch (error) {
      return null;
    }
  }

  static async classGroupsGet(
    offset: number,
    limit: number,
  ): Promise<ClassGroupsGetResponse | null> {
    try {
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
    } catch (error) {
      return null;
    }
  }

  static async classGroupSessionsGet(
    offset: number,
    limit: number,
  ): Promise<ClassGroupSessionsGetResponse | null> {
    try {
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
    } catch (error) {
      return null;
    }
  }

  static async sessionEnrollmentsGet(
    offset: number,
    limit: number,
  ): Promise<SessionEnrollmentsGetResponse | null> {
    try {
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
    } catch (error) {
      return null;
    }
  }
}
