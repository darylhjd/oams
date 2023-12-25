"use client";

import {
  useManagerDataStore,
  useManagerFilesStore,
} from "@/stores/manager_processing";
import { FileProcessingStepper, Step } from "@/components/file_processing";
import { APIClient } from "@/api/client";
import { notifications } from "@mantine/notifications";
import { getError } from "@/api/error";
import { IconCheck, IconX } from "@tabler/icons-react";
import { ManagerProcessingPreviewer } from "@/app/manager-processing/previewer";

export default function ManagerProcessingPage() {
  const fileStorage = useManagerFilesStore();
  const managerDataStorage = useManagerDataStore();

  const steps: Step[] = [
    {
      buttonText: "Preview Manager Data",
      action: async () => {
        try {
          const resp = await APIClient.classGroupManagersPost(
            fileStorage.files,
          );
          managerDataStorage.setData(resp.class_group_managers);
          return true;
        } catch (error) {
          notifications.show({
            title: "Manager Preview Error",
            message: getError(error),
            icon: <IconX />,
            color: "red",
          });
          return false;
        }
      },
    },
    {
      buttonText: "Confirm Manager Processing",
      action: async () => {
        try {
          await APIClient.classGroupManagersPut(managerDataStorage.data);
          notifications.show({
            title: "Success!",
            message: "All manager data has been processed!",
            icon: <IconCheck />,
            color: "teal",
          });
          return true;
        } catch (error) {
          notifications.show({
            title: "Manager Processing Error",
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
        managerDataStorage.setData([]);
        return true;
      },
    },
  ];

  return (
    <FileProcessingStepper
      fileStorage={fileStorage}
      steps={steps}
      initialStep={managerDataStorage.data.length == 0 ? 0 : 1}
      firstStepDescription="Choose manager file"
      secondStepDescription="Preview manager data"
      completionDescription="Press the 'Done' button to restart the process, or the
      'Back' button to revisit the previous steps."
      previewer={<ManagerProcessingPreviewer />}
    />
  );
}
