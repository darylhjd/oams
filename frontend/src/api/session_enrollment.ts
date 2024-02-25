import { CreatedUpdatedAt } from "@/api/types";

export type SessionEnrollmentsGetResponse = {
  session_enrollments: SessionEnrollment[];
};

export type SessionEnrollmentGetResponse = {
  session_enrollment: SessionEnrollment;
};

export type SessionEnrollment = {
  id: number;
  session_id: number;
  user_id: string;
  attended: boolean;
} & CreatedUpdatedAt;
