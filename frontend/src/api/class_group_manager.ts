import { CreatedUpdatedAt } from "@/api/types";

export type ClassGroupManagersGetResponse = {
  class_group_managers: ClassGroupManager[];
};

export type ClassGroupManagerGetResponse = {
  manager: ClassGroupManager;
};

export type ClassGroupManagerPatchResponse = {
  manager: ClassGroupManager;
};

export type ClassGroupManagersPostResponse = {
  class_group_managers: UpsertClassGroupManagerParams[];
};

export type UpsertClassGroupManagerParams = {
  user_id: string;
  class_group_id: number;
  managing_role: ManagingRole;
};

export type ClassGroupManagersPutResponse = {
  class_group_managers: ClassGroupManager[];
};

export type ClassGroupManager = {
  id: number;
  user_id: string;
  class_group_id: number;
  managing_role: ManagingRole;
} & CreatedUpdatedAt;

export enum ManagingRole {
  TeachingAssistant = "TEACHING_ASSISTANT",
  CourseCoordinator = "COURSE_COORDINATOR",
}
