export type LoginResponse = {
  redirect_url: string;
};

export type UserMeResponse = {
  session_user: User;
  upcoming_class_group_sessions: UpcomingClassGroupSession[];
};

export type User = {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  createdAt: Date;
  updatedAt: Date;
};

export enum UserRole {
  User = "USER",
  SystemAdmin = "SYSTEM_ADMIN",
}

export type UpcomingClassGroupSession = {
  code: string;
  year: number;
  semester: string;
  name: string;
  classType: ClassType;
  start_time: Date;
  end_time: Date;
};

export enum ClassType {
  Lecture = "LEC",
  Tutorial = "TUT",
  Lab = "LAB",
}

export type BatchPostResponse = {
  batches: BatchData[];
};

export type BatchData = {
  filename: string;
  file_creation_date: Date;
  class: UpsertClassParams;
  class_groups: ClassGroupData[];
};

export type ClassGroupData = UpsertClassGroupParams & {
  sessions: UpsertClassGroupSessionParams[];
  students: UpsertUserParams[];
};

export type UpsertClassGroupParams = {
  class_id: number;
  name: string;
  class_type: ClassType;
};

export type UpsertClassParams = {
  code: string;
  year: number;
  semester: string;
  programme: string;
  au: number;
};

export type UpsertClassGroupSessionParams = {
  class_group_id: number;
  start_time: Date;
  end_time: Date;
  venue: string;
};

export type UpsertUserParams = {
  id: string;
  name: string;
};
