import {
  UpsertClassGroupParams,
  UpsertClassGroupSessionParams,
  UpsertClassParams,
  UpsertUserParams,
} from "@/api/batch";
import { Class } from "@/api/class";
import { ClassGroup } from "@/api/class_group";
import { ClassGroupSession } from "@/api/class_group_session";
import { ClassManager } from "@/api/class_manager";
import { SessionEnrollment } from "@/api/session_enrollment";
import { User } from "@/api/user";
import { MRT_ColumnDef } from "mantine-react-table";

const CreatedAtUpdatedAtDataTableColumns = [
  { accessorKey: "created_at", header: "Created At" },
  { accessorKey: "updated_at", header: "Updated At" },
];

export const UserBatchDataTableColumns = <MRT_ColumnDef<UpsertUserParams>[]>[
  { accessorKey: "id", header: "ID" },
  { accessorKey: "name", header: "Name" },
];

export const UserDataTableColumns = <MRT_ColumnDef<User>[]>[
  ...UserBatchDataTableColumns,
  { accessorKey: "email", header: "Email" },
  { accessorKey: "role", header: "Role" },
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const ClassBatchDataTableColumns = <MRT_ColumnDef<UpsertClassParams>[]>[
  { accessorKey: "id", header: "ID" },
  { accessorKey: "code", header: "Code" },
  { accessorKey: "year", header: "Year" },
  { accessorKey: "semester", header: "Semester" },
  { accessorKey: "programme", header: "Programme" },
  { accessorKey: "au", header: "AU" },
];
export const ClassDataTableColumns = <MRT_ColumnDef<Class>[]>[
  ...ClassBatchDataTableColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const ClassManagerDataTableColumns = <MRT_ColumnDef<ClassManager>[]>[
  { accessorKey: "id", header: "ID" },
  { accessorKey: "user_id", header: "User ID" },
  { accessorKey: "class_id", header: "Class ID" },
  { accessorKey: "managing_role", header: "Managing Role" },
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const ClassGroupBatchDataTableColumns = <
  MRT_ColumnDef<UpsertClassGroupParams>[]
>[
  { accessorKey: "id", header: "ID" },
  { accessorKey: "class_id", header: "Class ID" },
  { accessorKey: "name", header: "Name" },
  { accessorKey: "class_type", header: "Class Type" },
];
export const ClassGroupDataTableColumns = <MRT_ColumnDef<ClassGroup>[]>[
  ...ClassGroupBatchDataTableColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const ClassGroupSessionBatchDataTableColumns = <
  MRT_ColumnDef<UpsertClassGroupSessionParams>[]
>[
  { accessorKey: "class_group_id", header: "Class Group ID" },
  { accessorKey: "start_time", header: "Start Time" },
  { accessorKey: "end_time", header: "End Time" },
  { accessorKey: "venue", header: "Venue" },
];
export const ClassGroupSessionDataTableColumns = <
  MRT_ColumnDef<ClassGroupSession>[]
>[
  { accessorKey: "id", header: "ID" },
  ...ClassGroupSessionBatchDataTableColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const SessionEnrollmentsDataTableColumns = <
  MRT_ColumnDef<SessionEnrollment>[]
>[
  { accessorKey: "id", header: "ID" },
  { accessorKey: "session_id", header: "Session ID" },
  { accessorKey: "user_id", header: "User ID" },
  { accessorKey: "attended", header: "Attended" },
  ...CreatedAtUpdatedAtDataTableColumns,
];
