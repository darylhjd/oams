import { CreatedUpdatedAt } from "@/api/types";

export type ClassGroupsGetResponse = {
  class_groups: ClassGroup[];
};

export type ClassGroupGetResponse = {
  class_group: ClassGroup;
};

export type ClassGroup = {
  id: number;
  class_id: number;
  name: string;
  class_type: ClassType;
} & CreatedUpdatedAt;

export enum ClassType {
  Lecture = "LEC",
  Tutorial = "TUT",
  Lab = "LAB",
}
