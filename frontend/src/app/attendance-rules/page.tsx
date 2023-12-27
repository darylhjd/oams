"use client";

import styles from "@/styles/AttendanceRules.module.css";

import { Container, Space, Text, Title } from "@mantine/core";
import { useState } from "react";
import { CoordinatingClass } from "@/api/attendance_rule";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import {
  MantineReactTable,
  MRT_DensityState,
  useMantineReactTable,
} from "mantine-react-table";
import { CoordinatingClassDataTableColumns } from "@/components/tabling";
import { useRouter } from "next/navigation";
import { Routes } from "@/routing/routes";

export default function AttendanceRulesPage() {
  const [coordinatingClasses, setCoordinatingClasses] = useState<
    CoordinatingClass[]
  >([]);
  const promiseFunc = async () => {
    const data = await APIClient.attendanceRulesGet();
    return setCoordinatingClasses(data.coordinating_classes);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Container className={styles.container} fluid>
        <PageTitle />
        <Space h="md" />
        <Text ta="center">Click each row to view associated rules.</Text>
        <Space h="xs" />
        <CoordinatingClassPicker coordinatingClasses={coordinatingClasses} />
      </Container>
    </RequestLoader>
  );
}

function PageTitle() {
  return (
    <Title order={2} ta="center">
      Attendance Rules
    </Title>
  );
}

function CoordinatingClassPicker({
  coordinatingClasses,
}: {
  coordinatingClasses: CoordinatingClass[];
}) {
  const router = useRouter();

  const table = useMantineReactTable({
    columns: CoordinatingClassDataTableColumns,
    data: coordinatingClasses,
    initialState: {
      density: "lg" as MRT_DensityState,
    },
    mantineTableBodyRowProps: ({ row }) => ({
      onClick: (_) => router.push(Routes.attendanceRule + row.original.id),
      style: { cursor: "pointer" },
    }),
  });

  return <MantineReactTable table={table} />;
}
