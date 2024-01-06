import axios from "axios";
import { createClient } from "@supabase/supabase-js";
import { UserGetResponse, UsersGetResponse } from "./user";
import { BatchData, BatchPostResponse, BatchPutResponse } from "./batch";
import { FileWithPath } from "@mantine/dropzone";
import { ClassesGetResponse, ClassGetResponse } from "./class";
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
  UpcomingClassGroupSessionAttendancesGetResponse,
  UpcomingClassGroupSessionAttendancePatchResponse,
  UpcomingClassGroupSessionsGetResponse,
} from "./upcoming_class_group_session";
import { SessionResponse } from "@/api/session";
import { ClassAttendanceRulesGetResponse } from "@/api/class_attendance_rule";
import {
  CoordinatingClassesGetResponse,
  CoordinatingClassRulesGetResponse,
  CoordinatingClassRulesPostRequest,
  CoordinatingClassRulesPostResponse,
} from "@/api/coordinating_class";

export class APIClient {
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
    const { data } = await this._client.get<SessionResponse>("/session");
    return data;
  }

  static async signaturePut(id: string, newSignature: string): Promise<void> {
    await this._client.put(`/signature/${id}`, {
      signature: newSignature,
    });
  }

  static async batchPost(files: FileWithPath[]): Promise<BatchPostResponse> {
    const form = new FormData();
    files.forEach((file) => form.append("batch-attachments", file));

    const { data } = await this._client.post<BatchPostResponse>("/batch", form);
    return data;
  }

  static async batchPut(batchData: BatchData[]): Promise<BatchPutResponse> {
    const { data } = await this._client.put<BatchPutResponse>("/batch", {
      batches: batchData,
    });
    return data;
  }

  static async usersGet(
    offset: number,
    limit: number,
  ): Promise<UsersGetResponse> {
    const { data } = await this._client.get<UsersGetResponse>("/users", {
      params: {
        offset: offset,
        limit: limit,
      },
    });
    return data;
  }

  static async userGet(id: string): Promise<UserGetResponse> {
    const { data } = await this._client.get<UserGetResponse>(`/users/${id}`);
    return data;
  }

  static async classesGet(
    offset: number,
    limit: number,
  ): Promise<ClassesGetResponse> {
    const { data } = await this._client.get<ClassesGetResponse>("/classes", {
      params: {
        offset: offset,
        limit: limit,
      },
    });
    return data;
  }

  static async classGet(id: number): Promise<ClassGetResponse> {
    const { data } = await this._client.get<ClassGetResponse>(`/classes/${id}`);
    return data;
  }

  static async classAttendanceRulesGet(
    offset: number,
    limit: number,
  ): Promise<ClassAttendanceRulesGetResponse> {
    const { data } = await this._client.get<ClassAttendanceRulesGetResponse>(
      "/class-attendance-rules",
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
      "/class-groups",
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
      `/class-groups/${id}`,
    );
    return data;
  }

  static async classGroupManagersGet(
    offset: number,
    limit: number,
  ): Promise<ClassGroupManagersGetResponse> {
    const { data } = await this._client.get<ClassGroupManagersGetResponse>(
      "/class-group-managers",
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
      "/class-group-managers",
      form,
    );
    return data;
  }

  static async classGroupManagersPut(
    classGroupManagers: UpsertClassGroupManagerParams[],
  ): Promise<ClassGroupManagersPutResponse> {
    const { data } = await this._client.put<ClassGroupManagersPutResponse>(
      "/class-group-managers",
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
      "/class-group-sessions",
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
      `/class-group-sessions/${id}`,
    );
    return data;
  }

  static async sessionEnrollmentsGet(
    offset: number,
    limit: number,
  ): Promise<SessionEnrollmentsGetResponse> {
    const { data } = await this._client.get<SessionEnrollmentsGetResponse>(
      "/session-enrollments",
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
      `/session-enrollments/${id}`,
    );
    return data;
  }

  static async upcomingClassGroupSessionsGet(): Promise<UpcomingClassGroupSessionsGetResponse> {
    const { data } =
      await this._client.get<UpcomingClassGroupSessionsGetResponse>(
        "/upcoming-class-group-sessions",
      );
    return data;
  }

  static async upcomingClassGroupSessionAttendancesGet(
    id: number,
  ): Promise<UpcomingClassGroupSessionAttendancesGetResponse> {
    const { data } =
      await this._client.get<UpcomingClassGroupSessionAttendancesGetResponse>(
        `/upcoming-class-group-sessions/${id}/attendances`,
      );
    return data;
  }

  static async upcomingClassGroupSessionAttendancePatch(
    id: number,
    sessionEnrollmentId: number,
    attended: boolean,
    userId: string,
    userSignature: string,
  ): Promise<UpcomingClassGroupSessionAttendancePatchResponse> {
    const { data } =
      await this._client.patch<UpcomingClassGroupSessionAttendancePatchResponse>(
        `/upcoming-class-group-sessions/${id}/attendances/${sessionEnrollmentId}`,
        {
          attended: attended,
          user_id: userId,
          user_signature: userSignature,
        },
      );
    return data;
  }

  static async coordinatingClassesGet(): Promise<CoordinatingClassesGetResponse> {
    const { data } = await this._client.get<CoordinatingClassesGetResponse>(
      "/coordinating-classes",
    );
    return data;
  }

  static async coordinatingClassRulesGet(
    id: number,
  ): Promise<CoordinatingClassRulesGetResponse> {
    const { data } = await this._client.get<CoordinatingClassRulesGetResponse>(
      `/coordinating-classes/${id}/rules`,
    );
    return data;
  }

  static async coordinatingClassRulesPost(
    id: number,
    params: CoordinatingClassRulesPostRequest,
  ): Promise<CoordinatingClassRulesPostResponse> {
    const { data } =
      await this._client.post<CoordinatingClassRulesPostResponse>(
        `/coordinating-classes/${id}/rules`,
        params,
      );
    return data;
  }

  static async coordinatingClassReportGet(id: number) {
    return await this._client.get(`/coordinating-classes/${id}/report`, {
      responseType: "blob",
    });
  }

  static async dataExportGet() {
    return await this._client.get("/data-export", {
      responseType: "blob",
    });
  }
}
