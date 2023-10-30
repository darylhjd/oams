export type ClassGroupsGetResponse = {
  class_groups: ClassGroup[];
};

export type ClassGroup = {
  id: number;
  class_id: number;
  name: string;
  class_type: ClassType;
  created_at: Date;
  updated_at: Date;
};

export enum ClassType {
  Lecture = "LEC",
  Tutorial = "TUT",
  Lab = "LAB",
}
