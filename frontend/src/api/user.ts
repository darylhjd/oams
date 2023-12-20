export type UsersGetResponse = {
  users: User[];
};

export type UserGetResponse = {
  user: User;
};

export type UserMeResponse = {
  session_user: User;
  has_managed_class_groups: boolean;
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
