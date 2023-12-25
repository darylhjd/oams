"use client";

import {
  useBatchDataStore,
  useBatchFilesStore,
} from "@/stores/batch_processing";
import { APIClient } from "@/api/client";
import { notifications } from "@mantine/notifications";
import { IconCheck, IconX } from "@tabler/icons-react";
import { getError } from "@/api/error";
import { Step, FileProcessingStepper } from "@/components/file_processing";
import { BatchProcessingPreviewer } from "@/app/batch-processing/previewer";

export default function BatchProcessingPage() {
  const fileStorage = useBatchFilesStore();
  const batchDataStorage = useBatchDataStore();

  const steps: Step[] = [
    {
      buttonText: "Preview Batch Data",
      action: async () => {
        try {
          const resp = await APIClient.batchPost(fileStorage.files);
          batchDataStorage.setData(resp.batches);
          return true;
        } catch (error) {
          notifications.show({
            title: "Batch Preview Error",
            message: getError(error),
            icon: <IconX />,
            color: "red",
          });
          return false;
        }
      },
    },
    {
      buttonText: "Confirm Batch Processing",
      action: async () => {
        try {
          await APIClient.batchPut(batchDataStorage.data);
          notifications.show({
            title: "Success!",
            message: "All batch data has been processed!",
            icon: <IconCheck />,
            color: "teal",
          });
          return true;
        } catch (error) {
          notifications.show({
            title: "Batch Processing Error",
            message: getError(error),
            icon: <IconX />,
            color: "red",
          });
          return false;
        }
      },
    },
    {
      buttonText: "Done!",
      action: async () => {
        fileStorage.clearFiles();
        batchDataStorage.setData([]);
        return true;
      },
    },
  ];

  return (
    <FileProcessingStepper
      fileStorage={fileStorage}
      steps={steps}
      initialStep={batchDataStorage.data.length == 0 ? 0 : 1}
      firstStepDescription="Choose batch files"
      secondStepDescription="Preview batch data"
      completionDescription="Press the 'Done' button to restart the process, or the
      'Back' button to revisit the previous steps."
      previewer={<BatchProcessingPreviewer />}
    />
  );
}
