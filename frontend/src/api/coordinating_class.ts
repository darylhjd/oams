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
