"use client";

import {
  UpsertClassGroupParams,
  UpsertClassGroupSessionParams,
  UpsertClassParams,
} from "@/api/batch";
import {
  ClassBatchDataTableColumns,
  ClassGroupBatchDataTableColumns,
  ClassGroupSessionBatchDataTableColumns,
  UserBatchDataTableColumns,
} from "@/components/entity_columns";
import { useBatchDataStore } from "@/stores/batch_processing";
import { MantineReactTable } from "mantine-react-table";

export function UsersPreviewTable() {
  const batchDataStorage = useBatchDataStore();

  const rows = batchDataStorage.data
    .map((batch) => batch.class_groups)
    .flat()
    .map((classGroup) => classGroup.students)
    .flat();

  return <MantineReactTable columns={UserBatchDataTableColumns} data={rows} />;
}

export function ClassesPreviewTable() {
  const batchDataStorage = useBatchDataStore();

  const rows = batchDataStorage.data.map<UpsertClassParams>((batch, index) => ({
    id: index,
    code: batch.class.code,
    year: batch.class.year,
    semester: batch.class.semester,
    programme: batch.class.programme,
    au: batch.class.au,
  }));

  return <MantineReactTable columns={ClassBatchDataTableColumns} data={rows} />;
}

export function ClassGroupsPreviewTable() {
  const batchDataStorage = useBatchDataStore();

  const rows = batchDataStorage.data
    .map((batch) => batch.class_groups)
    .flat()
    .map<UpsertClassGroupParams>((classGroup, index) => ({
      id: index,
      class_id: classGroup.class_id,
      name: classGroup.name,
      class_type: classGroup.class_type,
    }));

  return (
    <MantineReactTable columns={ClassGroupBatchDataTableColumns} data={rows} />
  );
}

export function ClassGroupSessionsPreviewTable() {
  const batchDataStorage = useBatchDataStore();

  const rows = batchDataStorage.data
    .map((batch) => batch.class_groups)
    .flat()
    .map((classGroup, index) =>
      classGroup.sessions.map<UpsertClassGroupSessionParams>((session) => ({
        class_group_id: index,
        start_time: session.start_time,
        end_time: session.end_time,
        venue: session.venue,
      })),
    )
    .flat();

  return (
    <MantineReactTable
      columns={ClassGroupSessionBatchDataTableColumns}
      data={rows}
    />
  );
}
