import { User } from "@/api/user";

export type SessionResponse = {
  session: Session | null;
};

export type Session = {
  user: User;
  management_details: ManagementDetails;
};

export type ManagementDetails = {
  attendance: boolean;
  administrative: boolean;
};
