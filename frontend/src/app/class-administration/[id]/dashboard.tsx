import { Text } from "@mantine/core";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";

export function DashboardTab({ id }: { id: number }) {
  const promiseFunc = async () => {
    await APIClient.coordinatingClassDashboardGet(id);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Text ta="center">This is the dashboard tab.</Text>
    </RequestLoader>
  );
}
