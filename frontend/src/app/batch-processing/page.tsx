"use client";

import styles from "@/styles/BatchProcessingPage.module.css";

import {
  Group,
  Stepper,
  Button,
  Container,
  Stack,
  StepperStep,
  StepperCompleted,
  Text,
} from "@mantine/core";
import { useState } from "react";
import { FilePicker } from "./file_picker";
import {
  useBatchDataStore,
  useBatchFilesStore,
} from "@/stores/batch_processing";
import { APIClient } from "@/api/client";
import { useMediaQuery } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { IconCheck, IconX } from "@tabler/icons-react";
import { getError } from "@/api/error";
import { Previewer } from "./previewer";
import { StepLayout } from "./step_layout";
import { IS_MOBILE_MEDIA_QUERY } from "@/components/media_query";

export default function BatchProcessingPage() {
  const isMobile = useMediaQuery(IS_MOBILE_MEDIA_QUERY);

  const fileStorage = useBatchFilesStore();
  const batchDataStorage = useBatchDataStore();

  const [step, setStep] = useState(batchDataStorage.data.length != 0 ? 1 : 0);
  const nextStep = () =>
    setStep((current) => Math.min(current + 1, stepDescriptions.length));
  const prevStep = () => setStep((current) => Math.max(current - 1, 0));

  const stepDescriptions = [
    {
      desc: "Preview Batch Data",
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
      desc: "Confirm Batch Processing",
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
    { desc: "Done!", action: async () => true },
  ];

  return (
    <Container className={styles.processor} fluid>
      <Stack>
        <Stepper
          active={step}
          onStepClick={setStep}
          allowNextStepsSelect={false}
          orientation={isMobile ? "vertical" : "horizontal"}
        >
          <StepperStep label="First step" description="Choose batch files">
            <FilePicker />
          </StepperStep>
          <StepperStep label="Second step" description="Preview batch data">
            <Previewer />
          </StepperStep>
          <StepperCompleted>
            <Completed />
          </StepperCompleted>
        </Stepper>

        <Group justify="center" mt="xl">
          <Button variant="default" onClick={prevStep} disabled={step == 0}>
            Back
          </Button>
          <Button
            disabled={fileStorage.files.length == 0}
            onClick={async () => {
              if (await stepDescriptions[step].action()) {
                nextStep();
              }

              if (step == stepDescriptions.length - 1) {
                setStep(0);
                fileStorage.clearFiles();
              }
            }}
          >
            {stepDescriptions[step].desc}
          </Button>
        </Group>
      </Stack>
    </Container>
  );
}

function Completed() {
  return (
    <StepLayout>
      <Container size="md">
        <Text ta="center">
          Processing complete!
          <br />
          <br />
          Press the &apos;Done&apos; button to restart the process, or the
          &apos;Back&apos; button to revisit the previous steps.
        </Text>
      </Container>
    </StepLayout>
  );
}
