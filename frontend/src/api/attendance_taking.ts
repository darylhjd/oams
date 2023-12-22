import { ClassType } from "./class_group";
import { ManagingRole } from "./class_group_manager";
import { SessionEnrollment } from "@/api/session_enrollment";

export type UpcomingClassGroupSession = {
  id: number;
  start_time: Date;
  end_time: Date;
  venue: string;
  code: string;
  year: number;
  name: string;
  semester: string;
  class_type: ClassType;
  managing_role: ManagingRole | null;
};

export type AttendanceTakingGetsResponse = {
  upcoming_class_group_sessions: UpcomingClassGroupSession[];
};

export type AttendanceTakingGetResponse = {
  upcoming_class_group_session: UpcomingClassGroupSession;
  enrollment_data: SessionEnrollment[];
};
