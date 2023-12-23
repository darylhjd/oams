import { User } from "@/api/user";

export type SessionResponse = {
  session: Session | null;
};

export type Session = {
  user: User;
  has_managed_class_groups: boolean;
};
