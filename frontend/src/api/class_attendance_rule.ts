import { CreatedUpdatedAt } from "@/api/types";

export type ClassAttendanceRulesGetResponse = {
  class_attendance_rules: ClassAttendanceRule[];
};

export type ClassAttendanceRule = {
  id: number;
  class_id: number;
  creator_id: string;
  title: string;
  description: string;
  rule: string;
  environment: JSON;
  active: boolean;
} & CreatedUpdatedAt;
