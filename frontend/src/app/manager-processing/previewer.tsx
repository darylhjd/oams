import { StepLayout } from "@/components/file_processing";
import {
  MantineReactTable,
  MRT_DensityState,
  useMantineReactTable,
} from "mantine-react-table";
import { useManagerDataStore } from "@/stores/manager_processing";
import { ClassGroupManagerProcessingDataTableColumns } from "@/components/tabling";

export function ManagerProcessingPreviewer() {
  const managerDataStorage = useManagerDataStore();

  const table = useMantineReactTable({
    columns: ClassGroupManagerProcessingDataTableColumns,
    data: managerDataStorage.data,
    initialState: {
      density: "sm" as MRT_DensityState,
    },
  });

  return (
    <StepLayout>
      <MantineReactTable table={table} />
    </StepLayout>
  );
}
