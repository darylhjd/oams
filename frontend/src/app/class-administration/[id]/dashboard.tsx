import styles from "@/styles/ClassAdministrationPageDashboard.module.css";

import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { Panel } from "@/app/class-administration/[id]/panel";
import { useState } from "react";
import { AttendanceCountData } from "@/api/coordinating_class";
import { BarChart } from "@mantine/charts";
import { Center, Container, Space, Title, Text } from "@mantine/core";

export function DashboardTab({ id }: { id: number }) {
  const [data, setData] = useState<AttendanceCountData[]>([]);
  const promiseFunc = async () => {
    const response = await APIClient.coordinatingClassDashboardGet(id);
    setData(response.attendance_count);
  };

  return (
    <Panel>
      <RequestLoader promiseFunc={promiseFunc}>
        <PercentageAttendanceChart data={data} />
      </RequestLoader>
    </Panel>
  );
}

function PercentageAttendanceChart({ data }: { data: AttendanceCountData[] }) {
  const chartHeight = 300;

  const chart =
    data.length == 0 ? (
      <Center className={styles.emptyBox} h={chartHeight}>
        <Text ta="center">
          No sessions have taken place yet! Check back later.
        </Text>
      </Center>
    ) : (
      <BarChart
        className={styles.barChart}
        h={chartHeight}
        data={data}
        dataKey="class_group_name"
        series={[
          {
            name: "attended",
            color: "green",
          },
          {
            name: "not_attended",
            color: "red",
          },
        ]}
        type="percent"
        withLegend
        legendProps={{ verticalAlign: "bottom" }}
        tooltipAnimationDuration={200}
      />
    );

  return (
    <>
      <Title order={4} ta="center">
        Percentage attendance by Class Group
      </Title>
      <Space h="md" />
      <Container className={styles.barChartBox} fluid>
        {chart}
      </Container>
    </>
  );
}
