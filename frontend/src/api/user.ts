import { CreatedUpdatedAt } from "@/api/types";

export type UsersGetResponse = {
  users: User[];
};

export type UserGetResponse = {
  user: User;
};

export type User = {
  id: string;
  name: string;
  email: string;
  role: UserRole;
} & CreatedUpdatedAt;

export enum UserRole {
  User = "USER",
  SystemAdmin = "SYSTEM_ADMIN",
}
