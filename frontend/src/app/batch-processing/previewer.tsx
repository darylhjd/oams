import styles from "@/styles/BatchProcessingPage.module.css";

import { Tabs, TabsList, TabsPanel, TabsTab } from "@mantine/core";
import { StepLayout } from "./step_layout";
import { useMediaQuery } from "@mantine/hooks";
import { IS_MOBILE_MEDIA_QUERY } from "@/components/media_query";
import {
  ClassGroupSessionsPreviewTable,
  ClassGroupsPreviewTable,
  ClassesPreviewTable,
  UsersPreviewTable,
} from "./tables";
import { BatchData } from "@/api/batch";

export function Previewer({ batchData }: { batchData: BatchData[] }) {
  const isMobile = useMediaQuery(IS_MOBILE_MEDIA_QUERY);

  return (
    <StepLayout>
      <Tabs
        defaultValue="users"
        variant="outline"
        classNames={{
          list: styles.previewTabList,
          tab: styles.previewTabTab,
        }}
      >
        <TabsList grow justify={isMobile ? "left" : "center"}>
          <TabsTab value="users">Users</TabsTab>
          <TabsTab value="classes">Classes</TabsTab>
          <TabsTab value="classGroups">Class Groups</TabsTab>
          <TabsTab value="classGroupSessions">Class Group Sessions</TabsTab>
        </TabsList>

        <TabsPanel value="users">
          <UsersPreviewTable batches={batchData} />
        </TabsPanel>
        <TabsPanel value="classes">
          <ClassesPreviewTable batches={batchData} />
        </TabsPanel>
        <TabsPanel value="classGroups">
          <ClassGroupsPreviewTable batches={batchData} />
        </TabsPanel>
        <TabsPanel value="classGroupSessions">
          <ClassGroupSessionsPreviewTable batches={batchData} />
        </TabsPanel>
      </Tabs>
    </StepLayout>
  );
}
