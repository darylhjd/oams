import { Text } from "@mantine/core";
import { useState } from "react";
import { CoordinatingClassDashboardReportData } from "@/api/coordinating_class";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";

export function DashboardTab({ id }: { id: number }) {
  const [reportData, setReportData] = useState(
    {} as CoordinatingClassDashboardReportData,
  );
  const promiseFunc = async () => {
    const response = await APIClient.coordinatingClassDashboardGet(id);
    setReportData(response.data);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Text ta="center">This is the dashboard tab.</Text>
    </RequestLoader>
  );
}
