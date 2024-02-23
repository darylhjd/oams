import axios from "axios";
import { UserGetResponse, UsersGetResponse } from "./user";
import { BatchData, BatchPostResponse, BatchPutResponse } from "./batch";
import { FileWithPath } from "@mantine/dropzone";
import { ClassesGetResponse, ClassGetResponse } from "./class";
import {
  ClassGroupManagerGetResponse,
  ClassGroupManagerPatchResponse,
  ClassGroupManagersGetResponse,
  ClassGroupManagersPostResponse,
  ClassGroupManagersPutResponse,
  ManagingRole,
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
  CoordinatingClassDashboardGetResponse,
  CoordinatingClassesGetResponse,
  CoordinatingClassGetResponse,
  CoordinatingClassRulePatchResponse,
  CoordinatingClassRulesGetResponse,
  CoordinatingClassRulesPostRequest,
  CoordinatingClassRulesPostResponse,
  CoordinatingClassSchedulePutResponse,
  CoordinatingClassSchedulesGetResponse,
} from "@/api/coordinating_class";
import { LoginResponse } from "@/api/login";

export class APIClient {
  static _client = axios.create({
    baseURL: `${process.env.API_SERVER}/api/v1`,
    withCredentials: true,
  });

  static async login(redirectUrl: string = ""): Promise<string> {
    const { data } = await this._client.get<LoginResponse>("/login", {
      params: {
        redirect_url: redirectUrl ? redirectUrl : process.env.WEB_SERVER,
      },
    });
    return data.redirect_url;
  }

  static async logout() {
    await this._client.get("/logout");
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

  static async classGroupManagerGet(
    id: number,
  ): Promise<ClassGroupManagerGetResponse> {
    const { data } = await this._client.get<ClassGroupManagerGetResponse>(
      `/class-group-managers/${id}`,
    );
    return data;
  }

  static async classGroupManagerPatch(
    id: number,
    role: ManagingRole,
  ): Promise<ClassGroupManagerPatchResponse> {
    const { data } = await this._client.patch<ClassGroupManagerPatchResponse>(
      `/class-group-managers/${id}`,
      {
        managing_role: role,
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

  static async coordinatingClassGet(
    id: number,
  ): Promise<CoordinatingClassGetResponse> {
    const { data } = await this._client.get<CoordinatingClassGetResponse>(
      `/coordinating-classes/${id}`,
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

  static async coordinatingClassRulePatch(
    classId: number,
    ruleId: number,
    active: boolean,
  ): Promise<CoordinatingClassRulePatchResponse> {
    const { data } =
      await this._client.patch<CoordinatingClassRulePatchResponse>(
        `/coordinating-classes/${classId}/rules/${ruleId}`,
        {
          active,
        },
      );
    return data;
  }

  static async coordinatingClassRuleDelete(classId: number, ruleId: number) {
    await this._client.delete(
      `/coordinating-classes/${classId}/rules/${ruleId}`,
    );
  }

  static async coordinatingClassReportGet(id: number) {
    return await this._client.get(`/coordinating-classes/${id}/report`, {
      responseType: "blob",
    });
  }

  static async coordinatingClassDashboardGet(
    id: number,
  ): Promise<CoordinatingClassDashboardGetResponse> {
    const { data } =
      await this._client.get<CoordinatingClassDashboardGetResponse>(
        `/coordinating-classes/${id}/dashboard`,
      );
    return data;
  }

  static async coordinatingClassSchedulesGet(
    id: number,
  ): Promise<CoordinatingClassSchedulesGetResponse> {
    const { data } =
      await this._client.get<CoordinatingClassSchedulesGetResponse>(
        `/coordinating-classes/${id}/schedule`,
      );
    return data;
  }

  static async coordinatingClassSchedulePut(
    classId: number,
    sessionId: number,
    start_time: Date,
    end_time: Date,
  ): Promise<CoordinatingClassSchedulePutResponse> {
    const { data } =
      await this._client.put<CoordinatingClassSchedulePutResponse>(
        `/coordinating-classes/${classId}/schedule/${sessionId}`,
        {
          start_time: start_time.getTime(),
          end_time: end_time.getTime(),
        },
      );
    return data;
  }

  static async dataExportGet() {
    return await this._client.get("/data-export", {
      responseType: "blob",
    });
  }
}
