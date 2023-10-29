import { UpsertUserParams } from "@/api/batch";
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

export const ClassBatchDataTableColumns = [
  { accessor: "id", title: "ID" },
  { accessor: "code" },
  { accessor: "year" },
  { accessor: "semester" },
  { accessor: "programme" },
  { accessor: "au" },
];
export const ClassDataTableColumns = [
  ...ClassBatchDataTableColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const ClassManagerDataTableColumns = [
  { accessor: "id", title: "ID" },
  { accessor: "user_id", title: "User ID" },
  { accessor: "class_id", title: "Class ID" },
  { accessor: "managing_role", title: "Managing Role" },
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const ClassGroupBatchDataTableColumns = [
  { accessor: "id", title: "ID" },
  { accessor: "class_id", title: "Class ID" },
  { accessor: "name" },
  { accessor: "class_type" },
];
export const ClassGroupDataTableColumns = [
  ...ClassGroupBatchDataTableColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const ClassGroupSessionBatchDataTableColumns = [
  { accessor: "class_group_id", title: "Class Group ID" },
  { accessor: "start_time", title: "Start Time" },
  { accessor: "end_time", title: "End Time" },
  { accessor: "venue" },
];
export const ClassGroupSessionDataTableColumns = [
  { accessor: "id", title: "ID" },
  ...ClassGroupSessionBatchDataTableColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const SessionEnrollmentsDataTableColumns = [
  { accessor: "id", title: "ID" },
  { accessor: "session_id", title: "Session ID" },
  { accessor: "user_id", title: "User ID" },
  {
    accessor: "attended",
    render: (record: SessionEnrollment) => record.attended.toString(),
  },
  ...CreatedAtUpdatedAtDataTableColumns,
];
