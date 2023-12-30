import { ClassAttendanceRule } from "@/api/class_attendance_rule";

export type CoordinatingClassesGetResponse = {
  coordinating_classes: CoordinatingClass[];
};

export type CoordinatingClassGetResponse = {
  coordinating_class: CoordinatingClass;
  rules: ClassAttendanceRule[];
};

export type CoordinatingClass = {
  id: number;
  code: string;
  year: number;
  semester: string;
  programme: string;
  au: number;
};

export type CoordinatingClassPostResponse = {
  rule: ClassAttendanceRule;
};

export enum RuleType {
  MissedConsecutiveClasses = "missed_consecutive_classes",
  MinPercentageAttendanceFromSession = "min_percentage_attendance_from_session",
  Advanced = "advanced",
}

export type CoordinatingClassPostRequest = {
  title: string;
  description: string;
  rule_type: string;
  consecutive_params: MissedConsecutiveClassParams;
  percentage_params: MinPercentageAttendanceFromSessionParams;
  advanced_params: AdvancedParams;
};

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
