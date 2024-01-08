import { ClassAttendanceRule } from "@/api/class_attendance_rule";
import { Class } from "@/api/class";
import { ManagingRole } from "@/api/class_group_manager";
import { ClassGroup } from "@/api/class_group";
import { ClassGroupSession } from "@/api/class_group_session";
import { SessionEnrollment } from "@/api/session_enrollment";

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

export type CoordinatingClassDashboardGetResponse = {
  data: CoordinatingClassDashboardReportData;
};

export type CoordinatingClassDashboardReportData = {
  class: Class;
  rules: ClassAttendanceRule[];
  managers: ClassGroupManagerReportData;
  class_groups: ClassGroupReportData;
};

export type ClassGroupManagerReportData = {
  user_id: string;
  user_name: string;
  class_group_name: string;
  managing_role: ManagingRole;
};

export type ClassGroupReportData = {
  class_group: ClassGroup;
  class_group_session: ClassGroupSession;
  session_enrollment: SessionEnrollment;
};
