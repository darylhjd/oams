export type ClassGroupManagersGetResponse = {
  class_group_managers: ClassGroupManager[];
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
  created_at: Date;
  updated_at: Date;
};

export enum ManagingRole {
  CourseCoordinator = "COURSE_COORDINATOR",
}
