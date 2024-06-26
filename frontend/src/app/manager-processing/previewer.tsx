import { StepLayout } from "@/components/file_processing";
import {
  MantineReactTable,
  MRT_DensityState,
  useMantineReactTable,
} from "mantine-react-table";
import { useManagerDataStore } from "@/app/manager-processing/manager_processing_store";
import { ClassGroupManagersProcessingDataTableColumns } from "@/components/columns";

export function ManagerProcessingPreviewer() {
  const managerDataStorage = useManagerDataStore();

  const table = useMantineReactTable({
    columns: ClassGroupManagersProcessingDataTableColumns,
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
