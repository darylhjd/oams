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
import { CoordinatingClass, ScheduleData } from "@/api/coordinating_class";

export const DEFAULT_PAGE_SIZE = 50;

function dateCell(date: Date) {
  const d = new Date(date);
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, "0");
  const day = d.getDate().toString().padStart(2, "0");
  const hour = d.getHours().toString().padStart(2, "0");
  const minute = d.getMinutes().toString().padStart(2, "0");
  const second = d.getSeconds().toString().padStart(2, "0");
  return <Text>{`${day}/${month}/${year}, ${hour}:${minute}:${second}`}</Text>;
}

export const UsersBatchDataTableColumns: MRT_ColumnDef<UpsertUserParams>[] = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "name", header: "Name" },
];
export const UsersDataTableColumns: MRT_ColumnDef<User>[] = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "name", header: "Name" },
  { accessorKey: "email", header: "Email" },
  { accessorKey: "role", header: "Role" },
  {
    accessorKey: "created_at",
    header: "Created At",
    Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
  },
  {
    accessorKey: "updated_at",
    header: "Updated At",
    Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
  },
];

export const ClassesBatchDataTableColumns: MRT_ColumnDef<UpsertClassParams>[] =
  [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "code", header: "Code" },
    { accessorKey: "year", header: "Year" },
    { accessorKey: "semester", header: "Semester" },
    { accessorKey: "programme", header: "Programme" },
    { accessorKey: "au", header: "AU" },
  ];

export const ClassesDataTableColumns: MRT_ColumnDef<Class>[] = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "code", header: "Code" },
  { accessorKey: "year", header: "Year" },
  { accessorKey: "semester", header: "Semester" },
  { accessorKey: "programme", header: "Programme" },
  { accessorKey: "au", header: "AU" },
  {
    accessorKey: "created_at",
    header: "Created At",
    Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
  },
  {
    accessorKey: "updated_at",
    header: "Updated At",
    Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
  },
];

export const CoordinatingClassesDataTableColumns: MRT_ColumnDef<CoordinatingClass>[] =
  [
    { accessorKey: "code", header: "Code" },
    { accessorKey: "year", header: "Year" },
    { accessorKey: "semester", header: "Semester" },
    { accessorKey: "programme", header: "Programme" },
    { accessorKey: "au", header: "AU" },
  ];

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
    {
      accessorKey: "active",
      header: "Active",
      Cell: ({ cell }: { cell: MRT_Cell<ClassAttendanceRule> }) => {
        const active = cell.getValue<boolean>();
        return (
          <Badge color={active ? "green" : "red"}>
            {active ? "Active" : "Inactive"}
          </Badge>
        );
      },
    },
    {
      accessorKey: "created_at",
      header: "Created At",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
    {
      accessorKey: "updated_at",
      header: "Updated At",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
  ];

export const ClassGroupsBatchDataTableColumns: MRT_ColumnDef<UpsertClassGroupParams>[] =
  [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "class_id", header: "Class ID" },
    { accessorKey: "name", header: "Name" },
    { accessorKey: "class_type", header: "Class Type" },
  ];
export const ClassGroupsDataTableColumns: MRT_ColumnDef<ClassGroup>[] = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "class_id", header: "Class ID" },
  { accessorKey: "name", header: "Name" },
  { accessorKey: "class_type", header: "Class Type" },
  {
    accessorKey: "created_at",
    header: "Created At",
    Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
  },
  {
    accessorKey: "updated_at",
    header: "Updated At",
    Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
  },
];

export const ClassGroupManagersDataTableColumns: MRT_ColumnDef<ClassGroupManager>[] =
  [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "user_id", header: "User ID" },
    { accessorKey: "class_group_id", header: "Class Group ID" },
    { accessorKey: "managing_role", header: "Managing Role" },
    {
      accessorKey: "created_at",
      header: "Created At",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
    {
      accessorKey: "updated_at",
      header: "Updated At",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
  ];
export const ClassGroupManagersProcessingDataTableColumns: MRT_ColumnDef<UpsertClassGroupManagerParams>[] =
  [
    { accessorKey: "user_id", header: "User ID" },
    { accessorKey: "class_group_id", header: "Class Group ID" },
    { accessorKey: "managing_role", header: "Managing Role" },
  ];

export const ClassGroupSessionsBatchDataTableColumns: MRT_ColumnDef<UpsertClassGroupSessionParams>[] =
  [
    { accessorKey: "class_group_id", header: "Class Group ID" },
    {
      accessorKey: "start_time",
      header: "Start Time",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
    {
      accessorKey: "end_time",
      header: "End Time",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
    { accessorKey: "venue", header: "Venue" },
  ];
export const ClassGroupSessionsDataTableColumns: MRT_ColumnDef<ClassGroupSession>[] =
  [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "class_group_id", header: "Class Group ID" },
    {
      accessorKey: "start_time",
      header: "Start Time",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
    {
      accessorKey: "end_time",
      header: "End Time",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
    { accessorKey: "venue", header: "Venue" },
    {
      accessorKey: "created_at",
      header: "Created At",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
    {
      accessorKey: "updated_at",
      header: "Updated At",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
  ];

export const SessionEnrollmentsDataTableColumns: MRT_ColumnDef<SessionEnrollment>[] =
  [
    { accessorKey: "id", header: "ID" },
    { accessorKey: "session_id", header: "Session ID" },
    { accessorKey: "user_id", header: "User ID" },
    {
      accessorKey: "attended",
      header: "Attended",
      Cell: ({ cell }) => {
        const attended = cell.getValue<boolean>();
        return (
          <Badge color={attended ? "green" : "red"}>
            {attended ? "Attended" : "Not attended"}
          </Badge>
        );
      },
    },
    {
      accessorKey: "created_at",
      header: "Created At",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
    {
      accessorKey: "updated_at",
      header: "Updated At",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
  ];
export const AttendanceEntriesDataTableColumns: MRT_ColumnDef<AttendanceEntry>[] =
  [
    { accessorKey: "user_id", header: "User ID" },
    { accessorKey: "user_name", header: "Name" },
    {
      accessorKey: "attended",
      header: "Attended",
      Cell: ({ cell }) => {
        const attended = cell.getValue<boolean>();
        return (
          <Badge color={attended ? "green" : "red"}>
            {attended ? "Attended" : "Not attended"}
          </Badge>
        );
      },
    },
  ];

export const CoordinatingClassScheduleTableColumns: MRT_ColumnDef<ScheduleData>[] =
  [
    { accessorKey: "class_group_name", header: "Class Group Name" },
    { accessorKey: "class_type", header: "Class Type" },
    {
      accessorKey: "start_time",
      header: "Start Time",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
    {
      accessorKey: "end_time",
      header: "End Time",
      Cell: ({ cell }) => dateCell(cell.getValue<Date>()),
    },
    { accessorKey: "venue", header: "Venue" },
  ];
