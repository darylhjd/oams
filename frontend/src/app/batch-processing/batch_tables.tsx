import { BatchData } from "@/api/models";
import { DataTable, DataTableColumn } from "mantine-datatable";

export function ClassesTable({ batches }: { batches: BatchData[] }) {
  const rows = batches.map((batch, index) => ({
    id: index,
    code: batch.class.code,
    year: batch.class.year,
    semester: batch.class.semester,
    programme: batch.class.programme,
    au: batch.class.au,
  }));

  return (
    <StyledDataTable
      columns={[
        { accessor: "id", title: "ID" },
        { accessor: "code" },
        { accessor: "year" },
        { accessor: "semester" },
        { accessor: "programme" },
        { accessor: "au" },
      ]}
      records={rows}
    />
  );
}

export function ClassGroupsTable({ batches }: { batches: BatchData[] }) {
  const rows = batches
    .map((batch) => batch.class_groups)
    .flat()
    .map((classGroup, index) => ({
      id: index,
      classId: classGroup.class_id,
      name: classGroup.name,
      classType: classGroup.class_type,
    }));

  return (
    <StyledDataTable
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

export function ClassGroupSessionsTable({ batches }: { batches: BatchData[] }) {
  const rows = batches
    .map((batch) => batch.class_groups)
    .flat()
    .map((classGroup) => classGroup.sessions)
    .flat();

  return (
    <StyledDataTable
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

export function UsersTable({ batches }: { batches: BatchData[] }) {
  const rows = batches
    .map((batch) => batch.class_groups)
    .flat()
    .map((classGroup) => classGroup.students)
    .flat();

  return (
    <StyledDataTable
      columns={[{ accessor: "id", title: "ID" }, { accessor: "name" }]}
      records={rows}
    />
  );
}

function StyledDataTable({
  columns,
  records,
}: {
  columns: DataTableColumn<any>[];
  records: any;
}) {
  return (
    <DataTable
      height={500}
      withBorder
      withColumnBorders
      highlightOnHover
      scrollAreaProps={{ type: "auto" }}
      columns={columns}
      records={records}
    />
  );
}
