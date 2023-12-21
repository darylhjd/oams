import { ClassType } from "./class_group";
import { ManagingRole } from "./class_group_manager";

export type AttendanceTakingGetResponse = {
  upcoming_class_group_sessions: UpcomingClassGroupSession[];
};

export type UpcomingClassGroupSession = {
  id: number;
  start_time: Date;
  end_time: Date;
  venue: string;
  class_type: ClassType;
  managing_role: ManagingRole | null;
};
