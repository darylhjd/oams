import { ClassAttendanceRule } from "@/api/class_attendance_rule";
import { ClassGroupSession } from "@/api/class_group_session";
import { ClassType } from "@/api/class_group";

export type CoordinatingClassesGetResponse = {
  coordinating_classes: CoordinatingClass[];
};

export type CoordinatingClass = {
  id: number;
  code: string;
  year: number;
  semester: string;
  programme: string;
  au: number;
};

export type CoordinatingClassGetResponse = {
  coordinating_class: CoordinatingClass;
};

export type CoordinatingClassRulesGetResponse = {
  rules: ClassAttendanceRule[];
};

export type CoordinatingClassRulesPostRequest = {
  title: string;
  description: string;
  rule_type: RuleType;
  consecutive_params: MissedConsecutiveClassParams;
  percentage_params: MinPercentageAttendanceFromSessionParams;
  advanced_params: AdvancedParams;
};

export enum RuleType {
  MissedConsecutiveClasses = 0,
  MinPercentageAttendanceFromSession = 1,
  Advanced = 2,
}

export type MissedConsecutiveClassParams = {
  consecutive_classes: number;
};

export type MinPercentageAttendanceFromSessionParams = {
  percentage: number;
  from_session: number;
};

export type AdvancedParams = {
  rule: string;
};

export type CoordinatingClassRulesPostResponse = {
  rule: ClassAttendanceRule;
};

export type CoordinatingClassRulePatchResponse = {
  active: boolean;
};

export type CoordinatingClassDashboardGetResponse = {
  result: boolean;
  attendance_count: AttendanceCountData[];
};

export type CoordinatingClassSchedulesGetResponse = {
  schedule: ScheduleData[];
};

export type ScheduleData = {
  class_group_name: string;
  class_type: ClassType;
  class_group_session_id: number;
  start_time: Date;
  end_time: Date;
  venue: string;
};

export type AttendanceCountData = {
  class_group_name: string;
  attended: number;
  not_attended: number;
};
