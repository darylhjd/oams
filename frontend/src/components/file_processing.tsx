import styles from "@/styles/FileProcessing.module.css";

import { ReactNode } from "react";
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
  Text,
} from "@mantine/core";
import { Dropzone, FileWithPath, MS_EXCEL_MIME_TYPE } from "@mantine/dropzone";
import { IconFile, IconUpload, IconX } from "@tabler/icons-react";
import { AxiosResponse } from "axios";

const MAX_FILE_SIZE = 5 * 1024 ** 2; // 5MB

export interface FileSetter {
  files: FileWithPath[];
  setFiles(files: FileWithPath[]): void;
  resetFiles(): void;
}

export type Step = {
  buttonText: string;
  action: () => Promise<boolean>;
};

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

export function StepLayout({ children }: { children: ReactNode }) {
  return (
    <>
      <Divider my="sm" />
      <Space visibleFrom="md" h="md" />
      {children}
    </>
  );
}

export function FilePicker({ fileStorage }: { fileStorage: FileSetter }) {
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
                Files should not exceed 5MB
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

function FileLister({ fileStorage }: { fileStorage: FileSetter }) {
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

function ResetFilesButton({ fileStorage }: { fileStorage: FileSetter }) {
  return (
    <Center>
      <Button
        color="red"
        variant="light"
        disabled={fileStorage.files.length == 0}
        onClick={() => fileStorage.resetFiles()}
      >
        Reset File Selection
      </Button>
    </Center>
  );
}

export function saveBlobResponseAsFile(response: AxiosResponse<any, any>) {
  // Do ugly stuff to save file :(
  const blob = new Blob([response.data]);

  const link = document.createElement("a");
  const url = window.URL.createObjectURL(blob);
  const [, filename] =
    response.headers["content-disposition"].split("filename=");
  link.href = url;
  link.setAttribute("download", filename);
  document.body.appendChild(link);

  link.click();

  // Clean up and remove the link
  link.parentNode?.removeChild(link);
  window.URL.revokeObjectURL(url);
}
