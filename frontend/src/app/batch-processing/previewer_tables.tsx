"use client";

import {
  UpsertClassGroupParams,
  UpsertClassGroupSessionParams,
  UpsertClassParams,
} from "@/api/batch";
import {
  ClassesBatchDataTableColumns,
  ClassGroupsBatchDataTableColumns,
  ClassGroupSessionsBatchDataTableColumns,
  UsersBatchDataTableColumns,
} from "@/components/columns";
import { useBatchDataStore } from "@/app/batch-processing/batch_processing_store";
import {
  MantineReactTable,
  MRT_ColumnDef,
  MRT_DensityState,
  useMantineReactTable,
} from "mantine-react-table";

export function UsersPreviewTable() {
  const batchDataStorage = useBatchDataStore();

  const rows = batchDataStorage.data
    .map((batch) => batch.class_groups)
    .flat()
    .map((classGroup) => classGroup.students)
    .flat();

  return <PreviewerTable columns={UsersBatchDataTableColumns} data={rows} />;
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

  return <PreviewerTable columns={ClassesBatchDataTableColumns} data={rows} />;
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
    <PreviewerTable columns={ClassGroupsBatchDataTableColumns} data={rows} />
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
    <PreviewerTable
      columns={ClassGroupSessionsBatchDataTableColumns}
      data={rows}
    />
  );
}

function PreviewerTable<T extends Record<string, any>>({
  columns,
  data,
}: {
  columns: MRT_ColumnDef<T>[];
  data: T[];
}) {
  const table = useMantineReactTable({
    columns,
    data,
    initialState: {
      density: "sm" as MRT_DensityState,
    },
  });

  return <MantineReactTable table={table} />;
}
