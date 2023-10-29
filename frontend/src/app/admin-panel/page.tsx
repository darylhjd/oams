"use client";

import { UserRole } from "@/api/user";
import { useSessionUserStore } from "@/stores/session";
import { Container, TabsList, Tabs, TabsPanel, TabsTab } from "@mantine/core";
import NotFoundPage from "../not-found";
import styles from "@/styles/AdminPage.module.css";
import { useMediaQuery } from "@mantine/hooks";
import { IS_MOBILE_MEDIA_QUERY } from "@/components/media_query";

export default function AdminPanelPage() {
  const session = useSessionUserStore();
  const isMobile = useMediaQuery(IS_MOBILE_MEDIA_QUERY);

  if (session.data?.session_user.role != UserRole.SystemAdmin) {
    return <NotFoundPage />;
  }

  return (
    <Container className={styles.container} fluid>
      <Tabs
        defaultValue="classes"
        classNames={{
          list: styles.entityTabList,
          tab: styles.entityTabTab,
        }}
      >
        <TabsList grow justify={isMobile ? "left" : "center"}>
          <TabsTab value="users">Users</TabsTab>
          <TabsTab value="classes">Classes</TabsTab>
          <TabsTab value="classManagers">Class Managers</TabsTab>
          <TabsTab value="classGroups">Class Groups</TabsTab>
          <TabsTab value="classGroupSessions">Class Group Sessions</TabsTab>
          <TabsTab value="sessionEnrollments">Session Enrollments</TabsTab>
        </TabsList>

        <TabsPanel value="users">Users Table goes here.</TabsPanel>
        <TabsPanel value="classes">Classes Table goes here.</TabsPanel>
        <TabsPanel value="classManagers">
          Class Managers Table goes here.
        </TabsPanel>
        <TabsPanel value="classGroups">Class Groups Table goes here.</TabsPanel>
        <TabsPanel value="classGroupSessions">
          Class Group Sessions Table goes here.
        </TabsPanel>
        <TabsPanel value="sessionEnrollments">
          Session Enrollments Table goes here.
        </TabsPanel>
      </Tabs>
    </Container>
  );
}
