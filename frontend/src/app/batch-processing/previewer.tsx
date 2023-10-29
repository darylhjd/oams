import { Tabs, TabsList, TabsPanel, TabsTab } from "@mantine/core";
import styles from "@/styles/BatchProcessingPage.module.css";
import { StepLayout } from "./step_layout";
import { useMediaQuery } from "@mantine/hooks";
import { IS_MOBILE_MEDIA_QUERY } from "@/components/media_query";

export function Previewer() {
  const isMobile = useMediaQuery(IS_MOBILE_MEDIA_QUERY);

  return (
    <StepLayout>
      <Tabs
        defaultValue="classes"
        variant="outline"
        classNames={{
          list: styles.previewTabList,
          tab: styles.previewTabTab,
        }}
      >
        <TabsList grow justify={isMobile ? "left" : "center"}>
          <TabsTab value="classes">Classes</TabsTab>
          <TabsTab value="classGroups">Class Groups</TabsTab>
          <TabsTab value="classGroupSessions">Class Group Sessions</TabsTab>
          <TabsTab value="users">Users</TabsTab>
        </TabsList>

        <TabsPanel value="classes">Classes Table goes here.</TabsPanel>
        <TabsPanel value="classGroups">Class Groups Table goes here.</TabsPanel>
        <TabsPanel value="classGroupSessions">
          Class Group Sessions Table goes here.
        </TabsPanel>
        <TabsPanel value="users">Users Table goes here.</TabsPanel>
      </Tabs>
    </StepLayout>
  );
}
