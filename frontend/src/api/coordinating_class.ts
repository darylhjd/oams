import { ClassAttendanceRule } from "@/api/class_attendance_rule";

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
