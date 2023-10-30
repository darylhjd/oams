export type ClassManagersGetResponse = {
  class_managers: ClassManager[];
};

export type ClassManager = {
  id: number;
  user_id: string;
  class_id: number;
  managing_role: ManagingRole;
  created_at: Date;
  updated_at: Date;
};

export enum ManagingRole {
  CourseCoordinator = "COURSE_COORDINATOR",
}
