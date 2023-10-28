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
import { useState } from "react";
import styles from "@/styles/BatchProcessingPage.module.css";
import { FilePicker } from "./file_picker";
import { useBatchFiles } from "@/stores/batch_file_picker";
import { APIClient } from "@/api/client";
import { useMediaQuery } from "@mantine/hooks";

export default function BatchProcessingPage() {
  const isMobile = useMediaQuery(`(max-width: 62em)`);

  const session = useSessionUserStore();
  const fileStorage = useBatchFiles();

  const stepDescriptions = [
    {
      desc: "Preview Batch Data",
      action: async () => APIClient.batchPost(fileStorage.files),
    },
    { desc: "Confirm Batch Processing", action: async () => {} },
    { desc: "Done!", action: async () => {} },
  ];

  const [step, setStep] = useState(0);
  const nextStep = () =>
    setStep((current) => Math.min(current + 1, stepDescriptions.length));
  const prevStep = () => setStep((current) => Math.max(current - 1, 0));

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
            <Button variant="default" onClick={prevStep}>
              Back
            </Button>
            <Button
              onClick={async () => {
                await stepDescriptions[step].action();

                nextStep();
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
