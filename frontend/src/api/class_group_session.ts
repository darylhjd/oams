import { CreatedUpdatedAt } from "@/api/types";

export type ClassGroupSessionsGetResponse = {
  class_group_sessions: ClassGroupSession[];
};

export type ClassGroupSessionGetResponse = {
  class_group_session: ClassGroupSession;
};

export type ClassGroupSession = {
  id: number;
  class_group_id: number;
  start_time: Date;
  end_time: Date;
  venue: string;
} & CreatedUpdatedAt;
