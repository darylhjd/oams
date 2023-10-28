import { Tabs, TabsList, TabsPanel } from "@mantine/core";
import styles from "@/styles/BatchProcessingPage.module.css";
import { StepLayout } from "./step_layout";

export function Previewer() {
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
        <TabsList grow justify="center">
          <Tabs.Tab value="classes">Classes</Tabs.Tab>
          <Tabs.Tab value="classGroups">Class Groups</Tabs.Tab>
          <Tabs.Tab value="classGroupSessions">Class Group Sessions</Tabs.Tab>
          <Tabs.Tab value="users">Users</Tabs.Tab>
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
