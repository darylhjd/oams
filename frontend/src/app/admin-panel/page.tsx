"use client";

import styles from "@/styles/AdminPage.module.css";

import { Container, TabsList, Tabs, TabsPanel, TabsTab } from "@mantine/core";
import { useMediaQuery } from "@mantine/hooks";
import { IS_MOBILE_MEDIA_QUERY } from "@/components/media_query";
import {
  ClassGroupSessionsTable,
  ClassGroupsTable,
  ClassGroupManagersTable,
  ClassesTable,
  SessionEnrollmentsTable,
  UsersTable,
  ClassAttendanceRulesTable,
} from "./tables";

export default function AdminPanelPage() {
  const isMobile = useMediaQuery(IS_MOBILE_MEDIA_QUERY);

  return (
    <Container className={styles.container} fluid>
      <Tabs
        defaultValue="users"
        classNames={{
          list: styles.entityTabList,
          tab: styles.entityTabTab,
        }}
      >
        <TabsList grow justify={isMobile ? "left" : "center"}>
          <TabsTab value="users">Users</TabsTab>
          <TabsTab value="classes">Classes</TabsTab>
          <TabsTab value="classAttendanceRules">Class Attendance Rules</TabsTab>
          <TabsTab value="classGroups">Class Groups</TabsTab>
          <TabsTab value="classGroupManagers">Class Group Managers</TabsTab>
          <TabsTab value="classGroupSessions">Class Group Sessions</TabsTab>
          <TabsTab value="sessionEnrollments">Session Enrollments</TabsTab>
        </TabsList>

        <TabsPanel value="users">
          <UsersTable />
        </TabsPanel>
        <TabsPanel value="classes">
          <ClassesTable />
        </TabsPanel>
        <TabsPanel value="classAttendanceRules">
          <ClassAttendanceRulesTable />
        </TabsPanel>
        <TabsPanel value="classGroups">
          <ClassGroupsTable />
        </TabsPanel>
        <TabsPanel value="classGroupManagers">
          <ClassGroupManagersTable />
        </TabsPanel>
        <TabsPanel value="classGroupSessions">
          <ClassGroupSessionsTable />
        </TabsPanel>
        <TabsPanel value="sessionEnrollments">
          <SessionEnrollmentsTable />
        </TabsPanel>
      </Tabs>
    </Container>
  );
}
