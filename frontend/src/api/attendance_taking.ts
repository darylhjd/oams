import { ClassGroupSession } from "./class_group_session";

export type AttendanceTakingGetResponse = {
  upcoming_class_group_sessions: ClassGroupSession[];
};
