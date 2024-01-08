"use client";

import styles from "@/styles/ClassAdministrationPage.module.css";

import {
  Container,
  Space,
  Tabs,
  TabsList,
  TabsPanel,
  TabsTab,
  Text,
  Title,
} from "@mantine/core";
import { Params } from "@/app/class-administration/[id]/layout";
import { ReactNode, useState } from "react";
import { CoordinatingClass } from "@/api/coordinating_class";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { ClassAttendanceRule } from "@/api/class_attendance_rule";
import { useMediaQuery } from "@mantine/hooks";
import { IS_MOBILE_MEDIA_QUERY } from "@/components/media_query";
import { RulesTab } from "@/app/class-administration/[id]/rules";
import { DashboardTab } from "@/app/class-administration/[id]/dashboard";

export default function ClassAdministrationPage({
  params,
}: {
  params: Params;
}) {
  const isMobile = useMediaQuery(IS_MOBILE_MEDIA_QUERY);

  const [coordinatingClass, setCoordinatingClass] = useState<CoordinatingClass>(
    {} as CoordinatingClass,
  );
  const promiseFunc = async () => {
    const classData = await APIClient.coordinatingClassGet(params.id);
    setCoordinatingClass(classData.coordinating_class);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Container className={styles.container} fluid>
        <CoordinatingClassDetails coordinatingClass={coordinatingClass} />
        <Space h="md" />
        <Tabs
          defaultValue="dashboard"
          classNames={{
            list: styles.tabsList,
            tab: styles.tabTab,
          }}
        >
          <TabsList grow={isMobile} justify="left">
            <TabsTab value="dashboard">Dashboard</TabsTab>
            <TabsTab value="rules">Rules</TabsTab>
          </TabsList>

          <TabsPanel value="dashboard">
            <Panel>
              <DashboardTab />
            </Panel>
          </TabsPanel>
          <TabsPanel value="rules">
            <Panel>
              <RulesTab id={params.id} />
            </Panel>
          </TabsPanel>
        </Tabs>
      </Container>
    </RequestLoader>
  );
}

function CoordinatingClassDetails({
  coordinatingClass,
}: {
  coordinatingClass: CoordinatingClass;
}) {
  return (
    <Title order={2} ta="center">
      <Text span inherit c="teal">
        {coordinatingClass.code} {coordinatingClass.year}/
        {coordinatingClass.semester}
      </Text>
    </Title>
  );
}

function Panel({ children }: { children: ReactNode }) {
  return <div className={styles.tabPanel}>{children}</div>;
}
