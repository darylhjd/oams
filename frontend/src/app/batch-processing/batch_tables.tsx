import { Table } from "@mantine/core";
import { ClassGroupData, UpsertClassParams } from "@/api/models";
import { DataTable } from "mantine-datatable";

export function ClassesTable({ cls }: { cls: UpsertClassParams }) {
  const row = (
    <tr key={cls.code + cls.year.toString() + cls.semester}>
      <td>0</td>
      <td>{cls.code}</td>
      <td>{cls.year}</td>
      <td>{cls.semester}</td>
      <td>{cls.programme}</td>
      <td>{cls.au}</td>
    </tr>
  );

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
    .flat()
    .map((session) => (
      <tr key={session.class_group_id + session.start_time.toString()}>
        <td>{session.class_group_id}</td>
        <td>{session.start_time.toString()}</td>
        <td>{session.end_time.toString()}</td>
        <td>{session.venue}</td>
      </tr>
    ));

  return (
    <Table>
      <thead>
        <tr>
          <th>Class Group ID</th>
          <th>Start Time</th>
          <th>End Time</th>
          <th>Venue</th>
        </tr>
      </thead>
      <tbody>{rows}</tbody>
    </Table>
  );
}

export function UsersTable({ classGroups }: { classGroups: ClassGroupData[] }) {
  const rows = classGroups
    .map((group) => group.students)
    .flat()
    .map((student) => (
      <tr key={student.id}>
        <td>{student.id}</td>
        <td>{student.name}</td>
      </tr>
    ));

  return (
    <Table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Name</th>
        </tr>
      </thead>
      <tbody>{rows}</tbody>
    </Table>
  );
}
