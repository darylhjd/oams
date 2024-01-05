"use client";

import styles from "@/styles/ClassReportsPage.module.css";

import { Container, Space, Text, Title } from "@mantine/core";
import { APIClient } from "@/api/client";
import { CoordinatingClassPicker } from "@/components/tabling";
import { Routes } from "@/routing/routes";
import { useState } from "react";
import { CoordinatingClass } from "@/api/coordinating_class";
import { useRouter } from "next/navigation";
import { RequestLoader } from "@/components/request_loader";

export default function ClassReportsPage() {
  const [coordinatingClasses, setCoordinatingClasses] = useState<
    CoordinatingClass[]
  >([]);
  const router = useRouter();

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
        <CoordinatingClassPicker
          coordinatingClasses={coordinatingClasses}
          onRowClick={(row) =>
            router.push(Routes.classReport + row.original.id)
          }
        />
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
