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
import { Badge, Text } from "@mantine/core";
import { AttendanceEntry } from "@/api/upcoming_class_group_session";
import { ClassAttendanceRule } from "@/api/class_attendance_rule";
import { CoordinatingClass } from "@/api/coordinating_class";

export const DEFAULT_PAGE_SIZE = 50;

const CreatedAtUpdatedAtDataTableColumns = [
  { accessorKey: "created_at", header: "Created At" },
  { accessorKey: "updated_at", header: "Updated At" },
];

const usersSharedColumns = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "name", header: "Name" },
];
export const UsersBatchDataTableColumns: MRT_ColumnDef<UpsertUserParams>[] = [
  ...usersSharedColumns,
];
export const UsersDataTableColumns: MRT_ColumnDef<User>[] = [
  ...usersSharedColumns,
  { accessorKey: "email", header: "Email" },
  { accessorKey: "role", header: "Role" },
  ...CreatedAtUpdatedAtDataTableColumns,
];

const classesSharedColumns = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "code", header: "Code" },
  { accessorKey: "year", header: "Year" },
  { accessorKey: "semester", header: "Semester" },
  { accessorKey: "programme", header: "Programme" },
  { accessorKey: "au", header: "AU" },
];
export const ClassesBatchDataTableColumns: MRT_ColumnDef<UpsertClassParams>[] =
  [...classesSharedColumns];
export const ClassesDataTableColumns: MRT_ColumnDef<Class>[] = [
  ...classesSharedColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const CoordinatingClassesDataTableColumns: MRT_ColumnDef<CoordinatingClass>[] =
  [...classesSharedColumns.slice(1)];

export const ClassAttendanceRulesDataTableColumns: MRT_ColumnDef<ClassAttendanceRule>[] =
  [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "class_id", header: "Class ID" },
    { accessorKey: "creator_id", header: "Creator ID" },
    { accessorKey: "title", header: "Title" },
    { accessorKey: "description", header: "Description" },
    {
      accessorKey: "rule",
      header: "Rule",
      Cell: ({ cell }: { cell: MRT_Cell<ClassAttendanceRule> }) => {
        const rule = cell.getValue<string>();
        return (
          <Text lineClamp={2} size="sm">
            {rule}
          </Text>
        );
      },
    },
    {
      accessorKey: "environment",
      header: "Environment",
      Cell: ({ cell }: { cell: MRT_Cell<ClassAttendanceRule> }) => {
        const environment = cell.getValue<JSON>();
        return (
          <Text lineClamp={2} size="sm">
            {JSON.stringify(environment)}
          </Text>
        );
      },
    },
    ...CreatedAtUpdatedAtDataTableColumns,
  ];

const classGroupsSharedColumns = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "class_id", header: "Class ID" },
  { accessorKey: "name", header: "Name" },
  { accessorKey: "class_type", header: "Class Type" },
];
export const ClassGroupsBatchDataTableColumns: MRT_ColumnDef<UpsertClassGroupParams>[] =
  [...classGroupsSharedColumns];
export const ClassGroupsDataTableColumns: MRT_ColumnDef<ClassGroup>[] = [
  ...classGroupsSharedColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

const classGroupManagersSharedColumns = [
  { accessorKey: "user_id", header: "User ID" },
  { accessorKey: "class_group_id", header: "Class Group ID" },
  { accessorKey: "managing_role", header: "Managing Role" },
];
export const ClassGroupManagersDataTableColumns: MRT_ColumnDef<ClassGroupManager>[] =
  [
    { accessorKey: "id", header: "ID" },
    ...classGroupManagersSharedColumns,
    ...CreatedAtUpdatedAtDataTableColumns,
  ];
export const ClassGroupManagersProcessingDataTableColumns: MRT_ColumnDef<UpsertClassGroupManagerParams>[] =
  [...classGroupManagersSharedColumns];

const classGroupSessionsSharedColumns = [
  { accessorKey: "class_group_id", header: "Class Group ID" },
  { accessorKey: "start_time", header: "Start Time" },
  { accessorKey: "end_time", header: "End Time" },
  { accessorKey: "venue", header: "Venue" },
];
export const ClassGroupSessionsBatchDataTableColumns: MRT_ColumnDef<UpsertClassGroupSessionParams>[] =
  [...classGroupSessionsSharedColumns];
export const ClassGroupSessionsDataTableColumns: MRT_ColumnDef<ClassGroupSession>[] =
  [
    { accessorKey: "id", header: "ID" },
    ...classGroupSessionsSharedColumns,
    ...CreatedAtUpdatedAtDataTableColumns,
  ];

function getSessionEnrollmentsSharedColumns<T extends Record<string, any>>() {
  return [
    { accessorKey: "user_id", header: "User ID" },
    {
      accessorKey: "attended",
      header: "Attended",
      Cell: ({ cell }: { cell: MRT_Cell<T> }) => {
        const attended = cell.getValue<boolean>();
        return (
          <Badge color={attended ? "green" : "red"}>
            {attended ? "Attended" : "Not attended"}
          </Badge>
        );
      },
    },
  ];
}

export const SessionEnrollmentsDataTableColumns: MRT_ColumnDef<SessionEnrollment>[] =
  [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "session_id", header: "Session ID" },
    ...getSessionEnrollmentsSharedColumns<SessionEnrollment>(),
    ...CreatedAtUpdatedAtDataTableColumns,
  ];
export const AttendanceEntriesDataTableColumns: MRT_ColumnDef<AttendanceEntry>[] =
  [
    ...getSessionEnrollmentsSharedColumns<AttendanceEntry>().slice(0, 1),
    { accessorKey: "user_name", header: "Name" },
    ...getSessionEnrollmentsSharedColumns<AttendanceEntry>().slice(1),
  ];
