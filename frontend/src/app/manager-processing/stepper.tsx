"use client";

import { APIClient } from "@/api/client";
import { notifications } from "@mantine/notifications";
import { getError } from "@/api/error";
import { IconCheck, IconX } from "@tabler/icons-react";
import { Completed, FilePicker, Step } from "@/components/file_processing";
import {
  ManagerDataStoreType,
  ManagerFileStoreType,
  useManagerDataStore,
  useManagerFilesStore,
} from "@/app/manager-processing/manager_processing_store";
import { useState } from "react";
import {
  Anchor,
  Button,
  Container,
  Group,
  Stack,
  Stepper,
  StepperCompleted,
  StepperStep,
  Text,
} from "@mantine/core";
import fileProcessingStyles from "@/styles/FileProcessing.module.css";
import { ManagerProcessingPreviewer } from "@/app/manager-processing/previewer";
import { useMediaQuery } from "@mantine/hooks";
import { IS_MOBILE_MEDIA_QUERY } from "@/components/media_query";

export function ManagerProcessingStepper() {
  const isMobile = useMediaQuery(IS_MOBILE_MEDIA_QUERY);
  const fileStorage = useManagerFilesStore();
  const managerDataStorage = useManagerDataStore();

  const steps = getSteps(fileStorage, managerDataStorage);

  const [step, setStep] = useState(managerDataStorage.data.length == 0 ? 0 : 1);
  const nextStep = () =>
    setStep((current) => Math.min(current + 1, steps.length));
  const prevStep = () => setStep((current) => Math.max(current - 1, 0));

  return (
    <Container className={fileProcessingStyles.processor} fluid>
      <Stack>
        <Stepper
          active={step}
          onStepClick={setStep}
          allowNextStepsSelect={false}
          orientation={isMobile ? "vertical" : "horizontal"}
        >
          <StepperStep label="First step" description="Choose manager file">
            <FilePicker fileStorage={fileStorage} />
            <Text mt="md" ta="center">
              Unsure about the file formatting?{" "}
              <Anchor href="/manager-processing/template.xlsx">
                Download the file template.
              </Anchor>
            </Text>
          </StepperStep>
          <StepperStep label="Second step" description="Preview manager data">
            <ManagerProcessingPreviewer />
          </StepperStep>
          <StepperCompleted>
            <Completed description="Press the 'Done' button to restart the process, or the 'Back' button to revisit the previous steps." />
          </StepperCompleted>
        </Stepper>

        <Group justify="center" mt="xl">
          <Button variant="default" onClick={prevStep} disabled={step == 0}>
            Back
          </Button>
          <Button
            disabled={fileStorage.files.length == 0}
            onClick={async () => {
              if (await steps[step].action()) {
                nextStep();
              }

              if (step == steps.length - 1) {
                setStep(0);
                fileStorage.resetFiles();
              }
            }}
          >
            {steps[step].buttonText}
          </Button>
        </Group>
      </Stack>
    </Container>
  );
}

function getSteps(
  fileStorage: ManagerFileStoreType,
  managerDataStorage: ManagerDataStoreType,
): Step[] {
  return [
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
        fileStorage.resetFiles();
        managerDataStorage.setData([]);
        return true;
      },
    },
  ];
}
