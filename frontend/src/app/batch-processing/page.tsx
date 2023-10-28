"use client";

import { UserRole } from "@/api/user";
import { useSessionUserStore } from "@/stores/session";
import NotFoundPage from "../not-found";
import {
  Group,
  Stepper,
  Button,
  Container,
  Stack,
  StepperStep,
  StepperCompleted,
} from "@mantine/core";
import { Dispatch, SetStateAction, useState } from "react";
import styles from "@/styles/BatchProcessingPage.module.css";
import { FilePicker } from "./file_picker";
import { useBatchFiles } from "@/stores/batch_file_picker";
import { APIClient } from "@/api/client";
import { useMediaQuery } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { IconCheck, IconX } from "@tabler/icons-react";
import { getError } from "@/api/error";
import { FileWithPath } from "@mantine/dropzone";
import { BatchData } from "@/api/batch";

export default function BatchProcessingPage() {
  const isMobile = useMediaQuery(`(max-width: 62em)`);

  const session = useSessionUserStore();
  const fileStorage = useBatchFiles();
  const [batchData, setBatchData] = useState<BatchData[]>([]);
  const [step, setStep] = useState(0);
  const nextStep = () =>
    setStep((current) => Math.min(current + 1, stepDescriptions.length));
  const prevStep = () => setStep((current) => Math.max(current - 1, 0));

  const stepDescriptions = [
    {
      desc: "Preview Batch Data",
      action: async () => previewAction(fileStorage.files, setBatchData),
    },
    {
      desc: "Confirm Batch Processing",
      action: async () => batchPutAction(batchData),
    },
    { desc: "Done!", action: async () => true },
  ];

  if (session.data?.session_user.role != UserRole.SystemAdmin) {
    return <NotFoundPage />;
  }

  return (
    <>
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
              Here is your data preview.
            </StepperStep>
            <StepperCompleted>
              Completed, click back button to get to previous step
            </StepperCompleted>
          </Stepper>

          <Group justify="center" mt="xl">
            <Button variant="default" onClick={prevStep} disabled={step == 0}>
              Back
            </Button>
            <Button
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
    </>
  );
}

async function previewAction(
  files: FileWithPath[],
  setBatchData: Dispatch<SetStateAction<BatchData[]>>,
): Promise<boolean> {
  try {
    const resp = await APIClient.batchPost(files);
    setBatchData(resp.batches);
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
}

async function batchPutAction(batchData: BatchData[]): Promise<boolean> {
  try {
    const resp = await APIClient.batchPut(batchData);
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
}
