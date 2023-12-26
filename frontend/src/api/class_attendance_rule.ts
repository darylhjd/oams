export type ClassAttendanceRulesGetResponse = {
  class_attendance_rules: ClassAttendanceRule[];
};

export type ClassAttendanceRule = {
  id: number;
  class_id: number;
  title: string;
  description: string;
  rule: string;
  created_at: Date;
  updated_at: Date;
};
