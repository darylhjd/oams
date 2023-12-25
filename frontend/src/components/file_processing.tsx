import styles from "@/styles/FileProcessing.module.css";

import React, { useState } from "react";
import {
  Button,
  Center,
  Container,
  Divider,
  Group,
  List,
  ListItem,
  rem,
  Space,
  Stack,
  Stepper,
  StepperCompleted,
  StepperStep,
  Text,
} from "@mantine/core";
import { Dropzone, MS_EXCEL_MIME_TYPE } from "@mantine/dropzone";
import { IconFile, IconUpload, IconX } from "@tabler/icons-react";
import { FileStoreType } from "@/stores/file_store";
import { useMediaQuery } from "@mantine/hooks";
import { IS_MOBILE_MEDIA_QUERY } from "@/components/media_query";

const MAX_FILE_SIZE = 5 * 1024 ** 2; // 5MB

export type Step = {
  buttonText: string;
  action: () => Promise<boolean>;
};

export function FileProcessingStepper({
  fileStorage,
  steps,
  initialStep,
  firstStepDescription,
  secondStepDescription,
  completionDescription,
  previewer,
}: {
  fileStorage: FileStoreType;
  steps: Step[];
  initialStep: number;
  firstStepDescription: string;
  secondStepDescription: string;
  completionDescription: string;
  previewer: React.ReactNode;
}) {
  const isMobile = useMediaQuery(IS_MOBILE_MEDIA_QUERY);

  const [step, setStep] = useState(initialStep);
  const nextStep = () =>
    setStep((current) => Math.min(current + 1, steps.length));
  const prevStep = () => setStep((current) => Math.max(current - 1, 0));

  return (
    <Container className={styles.processor} fluid>
      <Stack>
        <Stepper
          active={step}
          onStepClick={setStep}
          allowNextStepsSelect={false}
          orientation={isMobile ? "vertical" : "horizontal"}
        >
          <StepperStep label="First step" description={firstStepDescription}>
            <FilePicker fileStorage={fileStorage} />
          </StepperStep>
          <StepperStep label="Second step" description={secondStepDescription}>
            {previewer}
          </StepperStep>
          <StepperCompleted>
            <Completed description={completionDescription} />
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
                fileStorage.clearFiles();
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

export function Completed({ description }: { description: string }) {
  return (
    <StepLayout>
      <Container size="md">
        <Text ta="center">
          Processing complete!
          <br />
          <br />
          {description}
        </Text>
      </Container>
    </StepLayout>
  );
}

export function StepLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <Divider my="sm" />
      <Space visibleFrom="md" h="md" />
      {children}
    </>
  );
}

export function FilePicker({ fileStorage }: { fileStorage: FileStoreType }) {
  return (
    <StepLayout>
      <Container>
        <Dropzone
          onDrop={(files) => fileStorage.setFiles(files)}
          maxSize={MAX_FILE_SIZE}
          accept={MS_EXCEL_MIME_TYPE}
        >
          <Group
            justify="center"
            gap="xl"
            mih={220}
            style={{ pointerEvents: "none" }}
          >
            <Dropzone.Accept>
              <IconUpload
                style={{
                  width: rem(52),
                  height: rem(52),
                  color: "var(--mantine-color-blue-6)",
                }}
                stroke={1.5}
              />
            </Dropzone.Accept>
            <Dropzone.Reject>
              <IconX
                style={{
                  width: rem(52),
                  height: rem(52),
                  color: "var(--mantine-color-red-6)",
                }}
                stroke={1.5}
              />
            </Dropzone.Reject>
            <Dropzone.Idle>
              <IconFile
                style={{
                  width: rem(52),
                  height: rem(52),
                  color: "var(--mantine-color-dimmed)",
                }}
                stroke={1.5}
              />
            </Dropzone.Idle>

            <div>
              <Text size="xl" inline>
                Drag files here or click to select files
              </Text>
              <Text size="sm" c="dimmed" inline mt={7}>
                Attach as many files as you like, each file should not exceed
                5MB
              </Text>
            </div>

            <FileLister fileStorage={fileStorage} />
          </Group>
        </Dropzone>
        <Space h="md" />
        <ResetFilesButton fileStorage={fileStorage} />
      </Container>
    </StepLayout>
  );
}

function FileLister({ fileStorage }: { fileStorage: FileStoreType }) {
  if (fileStorage.files.length == 0) {
    return null;
  }

  return (
    <List className={styles.fileList}>
      {fileStorage.files.map((file) => (
        <ListItem key={file.name}>
          <>
            <Text visibleFrom="md" size="md">
              {file.name}
            </Text>
            <Text hiddenFrom="md" size="sm">
              {file.name}
            </Text>
          </>
        </ListItem>
      ))}
    </List>
  );
}

function ResetFilesButton({ fileStorage }: { fileStorage: FileStoreType }) {
  return (
    <Center>
      <Button
        color="red"
        variant="light"
        disabled={fileStorage.files.length == 0}
        onClick={() => fileStorage.clearFiles()}
      >
        Reset File Selection
      </Button>
    </Center>
  );
}
