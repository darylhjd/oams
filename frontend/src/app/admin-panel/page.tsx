"use client";

import { UserRole } from "@/api/models";
import { MOBILE_MIN_WIDTH } from "@/components/responsive";
import { redirectIfNotUserRole } from "@/routes/checks";
import { Center, Container, Tabs, createStyles } from "@mantine/core";
import { useMediaQuery } from "@mantine/hooks";

const useStyles = createStyles((theme) => ({
  tabList: {
    overflowY: "hidden",
    overflowX: "auto",
    flexWrap: "nowrap",
  },

  tabTab: {
    padding: "1em 1em",
  },
}));

export default function AdminPanelPage() {
  const { classes } = useStyles();
  const isMobile = useMediaQuery(MOBILE_MIN_WIDTH);

  if (redirectIfNotUserRole(UserRole.SystemAdmin)) {
    return null;
  }

  return (
    <Container>
      <Tabs defaultValue="users" variant="outline">
        <Tabs.List
          className={classes.tabList}
          position={isMobile ? "left" : "center"}
        >
          <Tabs.Tab className={classes.tabTab} value="users">
            Users
          </Tabs.Tab>
          <Tabs.Tab value="classes">Classes</Tabs.Tab>
          <Tabs.Tab value="classGroups">Class Groups</Tabs.Tab>
          <Tabs.Tab value="classGroupSessions">Class Group Sessions</Tabs.Tab>
          <Tabs.Tab value="sessionEnrollments">Session Enrollments</Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="users">
          <Center>Users WIP</Center>
        </Tabs.Panel>
        <Tabs.Panel value="classes">
          <Center>Classes WIP</Center>
        </Tabs.Panel>
        <Tabs.Panel value="classGroups">
          <Center>Class Groups WIP</Center>
        </Tabs.Panel>
        <Tabs.Panel value="classGroupSessions">
          <Center>Class Group Sessions WIP</Center>
        </Tabs.Panel>
        <Tabs.Panel value="sessionEnrollments">
          <Center>Session Enrollments WIP</Center>
        </Tabs.Panel>
      </Tabs>
    </Container>
  );
}
