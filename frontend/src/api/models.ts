export type LoginResponse = {
  redirect_url: string;
}

export type UserMeResponse = {
  session_user: User;
  upcoming_class_group_sessions: UpcomingClassGroupSession[];
}

export type User = {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  createdAt: Date;
  updatedAt: Date;
}

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
}

export enum ClassType {
  Lecture = "LEC",
  Tutorial = "TUT",
  Lab = "LAB",
}
