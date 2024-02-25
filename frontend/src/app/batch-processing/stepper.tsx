"use client";

import fileProcessingStyles from "@/styles/FileProcessing.module.css";
import styles from "@/styles/BatchProcessingPage.module.css";

import { useMediaQuery } from "@mantine/hooks";
import { IS_MOBILE_MEDIA_QUERY } from "@/components/media_query";
import { useState } from "react";
import {
  Button,
  Container,
  Fieldset,
  Group,
  Stack,
  Stepper,
  StepperCompleted,
  StepperStep,
  Space,
  NumberInput,
} from "@mantine/core";
import { Completed, FilePicker, Step } from "@/components/file_processing";
import {
  BatchDataStoreType,
  BatchFileStoreType,
  defaultSem1StartWeek,
  defaultSem2StartWeek,
  useBatchDataStore,
  useBatchFilesStore,
} from "@/app/batch-processing/batch_processing_store";
import { APIClient } from "@/api/client";
import { notifications } from "@mantine/notifications";
import { getError } from "@/api/error";
import { IconCheck, IconX } from "@tabler/icons-react";
import { BatchProcessingPreviewer } from "@/app/batch-processing/previewer";

export function BatchProcessingStepper() {
  const isMobile = useMediaQuery(IS_MOBILE_MEDIA_QUERY);
  const fileStorage = useBatchFilesStore();
  const batchDataStorage = useBatchDataStore();

  const steps = getSteps(fileStorage, batchDataStorage);

  const [step, setStep] = useState(batchDataStorage.data.length == 0 ? 0 : 1);
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
          <StepperStep label="First step" description="Choose batch files">
            <FilePicker fileStorage={fileStorage} />
            <Space h="md" />
            <StartWeekSelector />
          </StepperStep>
          <StepperStep label="Second step" description="Preview batch data">
            <BatchProcessingPreviewer />
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
                fileStorage.reset();
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

function StartWeekSelector() {
  const fileStorage = useBatchFilesStore();

  return (
    <Fieldset className={styles.startWeekField} legend="Additional Settings">
      <NumberInput
        label="Start Week"
        description="Week No. of the year that the semester starts."
        clampBehavior="strict"
        value={fileStorage.startWeek}
        allowDecimal={false}
        min={1}
        max={52}
        onChange={(v) => fileStorage.setStartWeek(v as number)}
      />
      <Group mt="md" justify="center">
        <Button
          variant="outline"
          onClick={() => fileStorage.setStartWeek(defaultSem1StartWeek)}
        >
          Set Default Sem 1
        </Button>
        <Button
          variant="outline"
          onClick={() => fileStorage.setStartWeek(defaultSem2StartWeek)}
        >
          Set Default Sem 2
        </Button>
      </Group>
    </Fieldset>
  );
}

function getSteps(
  fileStorage: BatchFileStoreType,
  batchDataStorage: BatchDataStoreType,
): Step[] {
  return [
    {
      buttonText: "Preview Batch Data",
      action: async () => {
        try {
          const resp = await APIClient.batchPost(
            fileStorage.files,
            fileStorage.startWeek,
          );
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
        fileStorage.reset();
        batchDataStorage.setData([]);
        return true;
      },
    },
  ];
}
