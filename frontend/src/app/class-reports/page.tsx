"use client";

import styles from "@/styles/ClassReportsPage.module.css";

import { Container, Space, Text, Title } from "@mantine/core";
import { APIClient } from "@/api/client";
import { CoordinatingClassPicker } from "@/components/tabling";
import { useState } from "react";
import { CoordinatingClass } from "@/api/coordinating_class";
import { RequestLoader } from "@/components/request_loader";
import { saveBlobResponseAsFile } from "@/components/file_processing";

export default function ClassReportsPage() {
  const [coordinatingClasses, setCoordinatingClasses] = useState<
    CoordinatingClass[]
  >([]);
  const promiseFunc = async () => {
    const data = await APIClient.coordinatingClassesGet();
    return setCoordinatingClasses(data.coordinating_classes);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Container className={styles.container} fluid>
        <PageTitle />
        <Space h="md" />
        <Text ta="center">Choose the class to generate reports.</Text>
        <Space h="xs" />
        <ReportGenerator coordinatingClasses={coordinatingClasses} />
      </Container>
    </RequestLoader>
  );
}

function PageTitle() {
  return (
    <Title order={2} ta="center">
      Class Reports
    </Title>
  );
}

function ReportGenerator({
  coordinatingClasses,
}: {
  coordinatingClasses: CoordinatingClass[];
}) {
  return (
    <CoordinatingClassPicker
      coordinatingClasses={coordinatingClasses}
      onRowClick={async (row) => {
        const response = await APIClient.coordinatingClassReportGet(
          row.original.id,
        );
        saveBlobResponseAsFile(response);
      }}
    />
  );
}
