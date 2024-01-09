import { Panel } from "@/app/class-administration/[id]/panel";
import { Button, Center, Space, Text } from "@mantine/core";
import { APIClient } from "@/api/client";
import { saveBlobResponseAsFile } from "@/components/file_processing";

export function ReportTab({ id }: { id: number }) {
  return (
    <Panel>
      <Text ta="center">
        OAMS provides a summary report for the class. You can access it by
        clicking the button below.
      </Text>
      <Text ta="center">The report will open in a new tab.</Text>
      <Space h="lg" />
      <Center>
        <Button
          onClick={async () => {
            const response = await APIClient.coordinatingClassReportGet(id);
            saveBlobResponseAsFile(response);
          }}
        >
          Download Report
        </Button>
      </Center>
    </Panel>
  );
}
