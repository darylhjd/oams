import { useState } from "react";
import { APIClient } from "@/api/client";
import { Panel } from "@/app/class-administration/[id]/panel";
import { RequestLoader } from "@/components/request_loader";
import { ScheduleData } from "@/api/coordinating_class";
import {
  MantineReactTable,
  MRT_DensityState,
  useMantineReactTable,
} from "mantine-react-table";
import { CoordinatingClassScheduleTableColumns } from "@/components/tabling";

export function ScheduleTab({ id }: { id: number }) {
  const [schedule, setSchedule] = useState<ScheduleData[]>([]);
  const promiseFunc = async () => {
    const scheduleData = await APIClient.coordinatingClassSchedulesGet(id);
    setSchedule(scheduleData.schedule);
  };

  return (
    <Panel>
      <RequestLoader promiseFunc={promiseFunc}>
        <ScheduleDisplay schedule={schedule} />
      </RequestLoader>
    </Panel>
  );
}

function ScheduleDisplay({ schedule }: { schedule: ScheduleData[] }) {
  const table = useMantineReactTable({
    columns: CoordinatingClassScheduleTableColumns,
    data: schedule,
    initialState: {
      density: "sm" as MRT_DensityState,
      pagination: { pageSize: 20, pageIndex: 0 },
    },
  });

  return <MantineReactTable table={table} />;
}
