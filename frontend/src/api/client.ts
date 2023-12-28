import axios from "axios";
import { createClient } from "@supabase/supabase-js";
import { UserGetResponse, UsersGetResponse } from "./user";
import { BatchData, BatchPostResponse, BatchPutResponse } from "./batch";
import { FileWithPath } from "@mantine/dropzone";
import { ClassGetResponse, ClassesGetResponse } from "./class";
import {
  ClassGroupManagersGetResponse,
  ClassGroupManagersPostResponse,
  ClassGroupManagersPutResponse,
  UpsertClassGroupManagerParams,
} from "./class_group_manager";
import { ClassGroupGetResponse, ClassGroupsGetResponse } from "./class_group";
import {
  ClassGroupSessionGetResponse,
  ClassGroupSessionsGetResponse,
} from "./class_group_session";
import {
  SessionEnrollmentGetResponse,
  SessionEnrollmentsGetResponse,
} from "./session_enrollment";
import {
  UpcomingClassGroupSessionGetResponse,
  UpcomingClassGroupSessionsGetResponse,
  UpcomingClassGroupSessionPostResponse,
} from "./upcoming_class_group_session";
import { SessionResponse } from "@/api/session";
import { ClassAttendanceRulesGetResponse } from "@/api/class_attendance_rule";
import {
  CoordinatingClassGetResponse,
  CoordinatingClassesGetResponse,
} from "@/api/coordinating_class";

export class APIClient {
  static readonly _sessionPath = "/session";
  static readonly _signaturePath = "/signature/";
  static readonly _batchPath = "/batch";
  static readonly _usersPath = "/users";
  static readonly _userPath = "/users/";
  static readonly _classesPath = "/classes";
  static readonly _classPath = "/classes/";
  static readonly _classAttendanceRulesPath = "/class-attendance-rules";
  static readonly _classGroupsPath = "/class-groups";
  static readonly _classGroupPath = "/class-groups/";
  static readonly _classGroupManagersPath = "/class-group-managers";
  static readonly _classGroupSessionsPath = "/class-group-sessions";
  static readonly _classGroupSessionPath = "/class-group-sessions/";
  static readonly _sessionEnrollmentsPath = "/session-enrollments";
  static readonly _sessionEnrollmentPath = "/session-enrollments/";
  static readonly _upcomingClassGroupSessionsPath =
    "/upcoming-class-group-sessions";
  static readonly _upcomingClassGroupSessionPath =
    "/upcoming-class-group-sessions/";
  static readonly _coordinatingClassesPath = "/coordinating-classes";
  static readonly _coordinatingClassPath = "/coordinating-classes/";

  static readonly _supabase = createClient(
    process.env.SUPABASE_URL!,
    process.env.SUPABASE_KEY!,
  );
  static _client = axios.create({
    baseURL: `${process.env.API_SERVER}/api/v1`,
  });

  static async login(redirectUrl: string = "") {
    await this._supabase.auth.signInWithOAuth({
      provider: "azure",
      options: {
        redirectTo: redirectUrl,
        scopes: `${process.env.AZURE_LOGIN_SCOPE} email`,
      },
    });
  }

  static async logout(): Promise<boolean> {
    await this._supabase.auth.signOut();
    return true;
  }

  static async loadSessionToken() {
    const session = (await this._supabase.auth.getSession()).data.session;
    if (session == null) {
      return;
    }

    this._client.defaults.headers.common["Authorization"] =
      `Bearer ${session.provider_token}`;
  }

  static async sessionGet(): Promise<SessionResponse> {
    const { data } = await this._client.get<SessionResponse>(this._sessionPath);
    return data;
  }

  static async signaturePut(id: string, newSignature: string): Promise<void> {
    await this._client.put(this._signaturePath + id, {
      signature: newSignature,
    });
  }

  static async batchPost(files: FileWithPath[]): Promise<BatchPostResponse> {
    const form = new FormData();
    files.forEach((file) => form.append("batch-attachments", file));

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

  static async userGet(id: string): Promise<UserGetResponse> {
    const { data } = await this._client.get<UserGetResponse>(
      this._userPath + id,
    );
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

  static async classGet(id: number): Promise<ClassGetResponse> {
    const { data } = await this._client.get<ClassGetResponse>(
      this._classPath + id,
    );
    return data;
  }

  static async classAttendanceRulesGet(
    offset: number,
    limit: number,
  ): Promise<ClassAttendanceRulesGetResponse> {
    const { data } = await this._client.get<ClassAttendanceRulesGetResponse>(
      this._classAttendanceRulesPath,
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

  static async classGroupGet(id: number): Promise<ClassGroupGetResponse> {
    const { data } = await this._client.get<ClassGroupGetResponse>(
      this._classGroupPath + id,
    );
    return data;
  }

  static async classGroupManagersGet(
    offset: number,
    limit: number,
  ): Promise<ClassGroupManagersGetResponse> {
    const { data } = await this._client.get<ClassGroupManagersGetResponse>(
      this._classGroupManagersPath,
      {
        params: {
          offset: offset,
          limit: limit,
        },
      },
    );
    return data;
  }

  static async classGroupManagersPost(
    files: FileWithPath[],
  ): Promise<ClassGroupManagersPostResponse> {
    const form = new FormData();
    files.forEach((file) => form.append("manager-attachments", file));

    const { data } = await this._client.post<ClassGroupManagersPostResponse>(
      this._classGroupManagersPath,
      form,
    );
    return data;
  }

  static async classGroupManagersPut(
    classGroupManagers: UpsertClassGroupManagerParams[],
  ): Promise<ClassGroupManagersPutResponse> {
    const { data } = await this._client.put<ClassGroupManagersPutResponse>(
      this._classGroupManagersPath,
      {
        class_group_managers: classGroupManagers,
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

  static async classGroupSessionGet(
    id: number,
  ): Promise<ClassGroupSessionGetResponse> {
    const { data } = await this._client.get<ClassGroupSessionGetResponse>(
      this._classGroupSessionPath + id,
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

  static async sessionEnrollmentGet(
    id: number,
  ): Promise<SessionEnrollmentGetResponse> {
    const { data } = await this._client.get<SessionEnrollmentGetResponse>(
      this._sessionEnrollmentPath + id,
    );
    return data;
  }

  static async upcomingClassGroupSessionsGet(): Promise<UpcomingClassGroupSessionsGetResponse> {
    const { data } =
      await this._client.get<UpcomingClassGroupSessionsGetResponse>(
        this._upcomingClassGroupSessionsPath,
      );
    return data;
  }

  static async upcomingClassGroupSessionGet(
    id: number,
  ): Promise<UpcomingClassGroupSessionGetResponse> {
    const { data } =
      await this._client.get<UpcomingClassGroupSessionGetResponse>(
        this._upcomingClassGroupSessionPath + id,
      );
    return data;
  }

  static async upcomingClassGroupSessionPost(
    id: number,
    sessionEnrollmentId: number,
    userId: string,
    attended: boolean,
    userSignature: string,
  ): Promise<UpcomingClassGroupSessionPostResponse> {
    const { data } =
      await this._client.post<UpcomingClassGroupSessionPostResponse>(
        this._upcomingClassGroupSessionPath + id,
        {
          id: sessionEnrollmentId,
          user_id: userId,
          attended: attended,
          user_signature: userSignature,
        },
      );
    return data;
  }

  static async coordinatingClassesGet(): Promise<CoordinatingClassesGetResponse> {
    const { data } = await this._client.get<CoordinatingClassesGetResponse>(
      this._coordinatingClassesPath,
    );
    return data;
  }

  static async coordinatingClassGet(
    id: number,
  ): Promise<CoordinatingClassGetResponse> {
    const { data } = await this._client.get<CoordinatingClassGetResponse>(
      this._coordinatingClassPath + id,
    );
    return data;
  }
}
