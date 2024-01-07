"use client";

import styles from "@/styles/ClassAdministrationsPage.module.css";

import { Container, Space, Text, Title } from "@mantine/core";
import { useState } from "react";
import { CoordinatingClass } from "@/api/coordinating_class";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { CoordinatingClassPicker } from "@/components/tabling";
import { Routes } from "@/routing/routes";
import { useRouter } from "next/navigation";

export default function ClassAdministrationsPage() {
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
        <Text ta="center">Select a class to manage.</Text>
        <Space h="xs" />
        <CoordinatingClassPicker
          coordinatingClasses={coordinatingClasses}
          onRowClick={(row) =>
            router.push(Routes.classAdministration + row.original.id)
          }
        />
      </Container>
    </RequestLoader>
  );
}

function PageTitle() {
  return (
    <Title order={2} ta="center">
      Class Administration
    </Title>
  );
}