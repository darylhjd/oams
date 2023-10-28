"use client";

import { Dropzone, MS_EXCEL_MIME_TYPE } from "@mantine/dropzone";
import { IconFile, IconUpload, IconX } from "@tabler/icons-react";
import {
  Button,
  Center,
  Container,
  Group,
  List,
  ListItem,
  Space,
  Text,
  rem,
} from "@mantine/core";
import { useBatchFiles } from "@/stores/batch_file_picker";
import styles from "@/styles/BatchProcessingPage.module.css";
import { StepLayout } from "./step_layout";

const MAX_FILE_SIZE = 5 * 1024 ** 2; // 5MB

export function FilePicker() {
  const fileStorage = useBatchFiles();

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

            <FileLister />
          </Group>
        </Dropzone>
        <Space h="md" />
        <ResetFilesButton />
      </Container>
    </StepLayout>
  );
}

function FileLister() {
  const fileStorage = useBatchFiles();

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

function ResetFilesButton() {
  const fileStorage = useBatchFiles();

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
