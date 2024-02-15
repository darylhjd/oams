import { useState } from "react";
import { APIClient } from "@/api/client";
import { Panel } from "@/app/class-administration/[id]/panel";
import { RequestLoader } from "@/components/request_loader";
import { Text } from "@mantine/core";
import { ScheduleData } from "@/api/coordinating_class";

export function ScheduleTab({ id }: { id: number }) {
  const [schedule, setSchedule] = useState<ScheduleData[]>([]);
  const promiseFunc = async () => {
    const scheduleData = await APIClient.coordinatingClassSchedulesGet(id);
    setSchedule(scheduleData.schedule);
  };

  return (
    <Panel>
      <RequestLoader promiseFunc={promiseFunc}>
        <Text ta="center">Coming Soon</Text>
      </RequestLoader>
    </Panel>
  );
}
