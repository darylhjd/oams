import {
  UpsertClassGroupParams,
  UpsertClassGroupSessionParams,
  UpsertClassParams,
  UpsertUserParams,
} from "@/api/batch";
import { Class } from "@/api/class";
import { ClassGroup } from "@/api/class_group";
import { ClassGroupSession } from "@/api/class_group_session";
import {
  ClassGroupManager,
  UpsertClassGroupManagerParams,
} from "@/api/class_group_manager";
import { SessionEnrollment } from "@/api/session_enrollment";
import { User } from "@/api/user";
import { MRT_Cell, MRT_ColumnDef } from "mantine-react-table";
import { Badge } from "@mantine/core";
import { AttendanceEntry } from "@/api/attendance_taking";
import { ClassAttendanceRule } from "@/api/class_attendance_rule";
import { CoordinatingClass } from "@/api/coordinating_class";

export const DEFAULT_PAGE_SIZE = 50;

const CreatedAtUpdatedAtDataTableColumns = [
  { accessorKey: "created_at", header: "Created At" },
  { accessorKey: "updated_at", header: "Updated At" },
];

const userSharedColumns = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "name", header: "Name" },
];
export const UserBatchDataTableColumns: MRT_ColumnDef<UpsertUserParams>[] = [
  ...userSharedColumns,
];
export const UserDataTableColumns: MRT_ColumnDef<User>[] = [
  ...userSharedColumns,
  { accessorKey: "email", header: "Email" },
  { accessorKey: "role", header: "Role" },
  ...CreatedAtUpdatedAtDataTableColumns,
];

const classSharedColumns = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "code", header: "Code" },
  { accessorKey: "year", header: "Year" },
  { accessorKey: "semester", header: "Semester" },
  { accessorKey: "programme", header: "Programme" },
  { accessorKey: "au", header: "AU" },
];
export const ClassBatchDataTableColumns: MRT_ColumnDef<UpsertClassParams>[] = [
  ...classSharedColumns,
];
export const ClassDataTableColumns: MRT_ColumnDef<Class>[] = [
  ...classSharedColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const CoordinatingClassDataTableColumns: MRT_ColumnDef<CoordinatingClass>[] =
  [...classSharedColumns.slice(1)];

export const ClassAttendanceRuleDataTableColumns: MRT_ColumnDef<ClassAttendanceRule>[] =
  [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "class_id", header: "Class ID" },
    { accessorKey: "title", header: "Title" },
    { accessorKey: "description", header: "Description" },
    { accessorKey: "rule", header: "Rule" },
    ...CreatedAtUpdatedAtDataTableColumns,
  ];

const classGroupSharedColumns = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "class_id", header: "Class ID" },
  { accessorKey: "name", header: "Name" },
  { accessorKey: "class_type", header: "Class Type" },
];
export const ClassGroupBatchDataTableColumns: MRT_ColumnDef<UpsertClassGroupParams>[] =
  [...classGroupSharedColumns];
export const ClassGroupDataTableColumns: MRT_ColumnDef<ClassGroup>[] = [
  ...classGroupSharedColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

const classGroupManagerSharedColumns = [
  { accessorKey: "user_id", header: "User ID" },
  { accessorKey: "class_group_id", header: "Class Group ID" },
  { accessorKey: "managing_role", header: "Managing Role" },
];
export const ClassGroupManagerDataTableColumns: MRT_ColumnDef<ClassGroupManager>[] =
  [
    { accessorKey: "id", header: "ID" },
    ...classGroupManagerSharedColumns,
    ...CreatedAtUpdatedAtDataTableColumns,
  ];
export const ClassGroupManagerProcessingDataTableColumns: MRT_ColumnDef<UpsertClassGroupManagerParams>[] =
  [...classGroupManagerSharedColumns];

const classGroupSessionSharedColumns = [
  { accessorKey: "class_group_id", header: "Class Group ID" },
  { accessorKey: "start_time", header: "Start Time" },
  { accessorKey: "end_time", header: "End Time" },
  { accessorKey: "venue", header: "Venue" },
];
export const ClassGroupSessionBatchDataTableColumns: MRT_ColumnDef<UpsertClassGroupSessionParams>[] =
  [...classGroupSessionSharedColumns];
export const ClassGroupSessionDataTableColumns: MRT_ColumnDef<ClassGroupSession>[] =
  [
    { accessorKey: "id", header: "ID" },
    ...classGroupSessionSharedColumns,
    ...CreatedAtUpdatedAtDataTableColumns,
  ];

const sessionEnrollmentsSharedColumns = [
  { accessorKey: "user_id", header: "User ID" },
  {
    accessorKey: "attended",
    header: "Attended",
    Cell: ({ cell }: { cell: MRT_Cell<any> }) => {
      const attended = cell.getValue<boolean>();
      return (
        <Badge color={attended ? "green" : "red"}>
          {attended ? "Attended" : "Not attended"}
        </Badge>
      );
    },
  },
];
export const SessionEnrollmentsDataTableColumns: MRT_ColumnDef<SessionEnrollment>[] =
  [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "session_id", header: "Session ID" },
    ...sessionEnrollmentsSharedColumns,
    ...CreatedAtUpdatedAtDataTableColumns,
  ];
export const AttendanceTakingDataTableColumns: MRT_ColumnDef<AttendanceEntry>[] =
  [
    ...sessionEnrollmentsSharedColumns.slice(0, 1),
    { accessorKey: "user_name", header: "Name" },
    ...sessionEnrollmentsSharedColumns.slice(1),
  ];
