import { ClassType } from "./class";

export type BatchPostResponse = {
  batches: BatchData[];
};

export type BatchPutRequest = BatchPostResponse;

export type BatchPutResponse = {
  class_ids: number[];
};

export type BatchData = {
  filename: string;
  file_creation_date: Date;
  class: UpsertClassParams;
  class_groups: ClassGroupData[];
};

export type ClassGroupData = UpsertClassGroupParams & {
  sessions: UpsertClassGroupSessionParams[];
  students: UpsertUserParams[];
};

export type UpsertClassGroupParams = {
  class_id: number;
  name: string;
  class_type: ClassType;
};

export type UpsertClassParams = {
  code: string;
  year: number;
  semester: string;
  programme: string;
  au: number;
};

export type UpsertClassGroupSessionParams = {
  class_group_id: number;
  start_time: Date;
  end_time: Date;
  venue: string;
};

export type UpsertUserParams = {
  id: string;
  name: string;
};
