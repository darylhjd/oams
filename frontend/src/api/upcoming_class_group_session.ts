import { ClassType } from "./class_group";
import { ManagingRole } from "./class_group_manager";

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

export type UpcomingClassGroupSessionsGetResponse = {
  upcoming_class_group_sessions: UpcomingClassGroupSession[];
};

export type UpcomingClassGroupSessionAttendancesGetResponse = {
  upcoming_class_group_session: UpcomingClassGroupSession;
  attendance_entries: AttendanceEntry[];
};

export type AttendanceEntry = {
  id: number;
  session_id: number;
  user_id: string;
  user_name: string;
  attended: boolean;
};

export type UpcomingClassGroupSessionAttendancePatchResponse = {
  attended: boolean;
};
