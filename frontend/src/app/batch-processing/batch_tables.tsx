import { Table } from "@mantine/core";
import { ClassGroupData, UpsertClassParams } from "@/api/models";
import { DataTable } from "mantine-datatable";

export function ClassesTable({ cls }: { cls: UpsertClassParams }) {
  return (
    <DataTable
      withBorder
      withColumnBorders
      highlightOnHover
      columns={[
        { accessor: "id", title: "ID" },
        { accessor: "code" },
        { accessor: "year" },
        { accessor: "semester" },
        { accessor: "programme" },
        { accessor: "au" },
      ]}
      records={[
        {
          id: 0,
          code: cls.code,
          year: cls.year,
          semester: cls.semester,
          programme: cls.programme,
          au: cls.au,
        },
      ]}
    />
  );
}

export function ClassGroupsTable({
  classGroups,
}: {
  classGroups: ClassGroupData[];
}) {
  const rows = classGroups.map((classGroup, index) => ({
    id: index,
    classId: classGroup.class_id,
    name: classGroup.name,
    classType: classGroup.class_type,
  }));

  return (
    <DataTable
      withBorder
      withColumnBorders
      highlightOnHover
      columns={[
        { accessor: "id", title: "ID" },
        { accessor: "classId", title: "Class ID" },
        { accessor: "name" },
        { accessor: "classType" },
      ]}
      records={rows}
    />
  );
}

export function ClassGroupSessionsTable({
  classGroups,
}: {
  classGroups: ClassGroupData[];
}) {
  const rows = classGroups
    .map((group, index) => {
      for (let i = 0; i < group.sessions.length; i++) {
        group.sessions[i].class_group_id = index;
      }

      return group.sessions;
    })
    .flat();

  return (
    <DataTable
      withBorder
      withColumnBorders
      highlightOnHover
      columns={[
        { accessor: "class_group_id", title: "Class Group ID" },
        { accessor: "start_time", title: "Start Time" },
        { accessor: "end_time", title: "End Time" },
        { accessor: "venue" },
      ]}
      records={rows}
    />
  );
}

export function UsersTable({ classGroups }: { classGroups: ClassGroupData[] }) {
  const rows = classGroups.map((group) => group.students).flat();

  return (
    <DataTable
      withBorder
      withColumnBorders
      highlightOnHover
      columns={[{ accessor: "id", title: "ID" }, { accessor: "name" }]}
      records={rows}
    />
  );
}
