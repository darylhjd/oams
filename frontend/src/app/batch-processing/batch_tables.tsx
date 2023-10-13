import { BatchData } from "@/api/models";
import {
  BatchDataTable,
  ClassBatchDataTableColumns,
  ClassGroupBatchDataTableColumns,
  ClassGroupSessionBatchDataTableColumns,
  UserBatchDataTableColumns,
} from "@/components/entity_tables";

export function ClassesTable({ batches }: { batches: BatchData[] }) {
  const rows = batches.map((batch, index) => ({
    id: index,
    code: batch.class.code,
    year: batch.class.year,
    semester: batch.class.semester,
    programme: batch.class.programme,
    au: batch.class.au,
  }));

  return <BatchDataTable columns={ClassBatchDataTableColumns} records={rows} />;
}

export function ClassGroupsTable({ batches }: { batches: BatchData[] }) {
  const rows = batches
    .map((batch) => batch.class_groups)
    .flat()
    .map((classGroup, index) => ({
      id: index,
      class_id: classGroup.class_id,
      name: classGroup.name,
      class_type: classGroup.class_type,
    }));

  return (
    <BatchDataTable columns={ClassGroupBatchDataTableColumns} records={rows} />
  );
}

export function ClassGroupSessionsTable({ batches }: { batches: BatchData[] }) {
  const rows = batches
    .map((batch) => batch.class_groups)
    .flat()
    .map((classGroup) => classGroup.sessions)
    .flat();

  return (
    <BatchDataTable
      columns={ClassGroupSessionBatchDataTableColumns}
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

  return <BatchDataTable columns={UserBatchDataTableColumns} records={rows} />;
}
