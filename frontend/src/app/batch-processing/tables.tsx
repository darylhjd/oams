import { BatchData } from "@/api/batch";
import { UserBatchDataTableColumns } from "@/components/entity_columns";
import { MantineReactTable } from "mantine-react-table";

export function UsersDataTable({ batches }: { batches: BatchData[] }) {
  const rows = batches
    .map((batch) => batch.class_groups)
    .flat()
    .map((classGroup) => classGroup.students)
    .flat();

  return <MantineReactTable columns={UserBatchDataTableColumns} data={rows} />;
}
