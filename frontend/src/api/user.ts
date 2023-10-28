import { ClassType } from "./class";

export type UsersGetResponse = {
  users: User[];
};

export type UserMeResponse = {
  session_user: User;
  upcoming_class_group_sessions: UpcomingClassGroupSession[];
};

export type UpcomingClassGroupSession = {
  code: string;
  year: number;
  semester: string;
  name: string;
  class_type: ClassType;
  start_time: Date;
  end_time: Date;
};

export type User = {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  created_at: Date;
  updated_at: Date;
};

export enum UserRole {
  User = "USER",
  SystemAdmin = "SYSTEM_ADMIN",
}
