"use client";

import { APIClient } from "@/api/client";
import { UserRole } from "@/api/models";
import { redirectIfNotUserRole } from "@/routes/checks";
import {
  Button,
  Center,
  Container,
  FileButton,
  List,
  Space,
  Stack,
  Title,
  createStyles,
} from "@mantine/core";
import { Dispatch, SetStateAction, useState } from "react";

const useStyles = createStyles((theme) => ({
  fileContainer: {
    padding: "2em 0",
  },

  listItem: {
    listStyleType: "none",
    margin: 0,
    padding: 0,
  },
}));

export default function BatchProcessingPage() {
  redirectIfNotUserRole(UserRole.SystemAdmin);

  const [files, setFiles] = useState<File[]>([]);
  const clearFiles = () => setFiles([]);

  return (
    <>
      <Container>
        <Center>
          <Stack>
            <Center>
              <p>Upload your batch files here.</p>
            </Center>
            <ChooseFilesButton onChange={setFiles} />
            <ResetFilesButton files={files} onClick={clearFiles} />
            <Space h="md" />
          </Stack>
        </Center>
      </Container>
      <SelectedFilesList files={files} />
    </>
  );
}

function ChooseFilesButton({
  onChange,
}: {
  onChange: Dispatch<SetStateAction<File[]>>;
}) {
  return (
    <FileButton onChange={onChange} accept="xlsx" multiple>
      {(props) => <Button {...props}>Choose files</Button>}
    </FileButton>
  );
}

function ResetFilesButton({
  files,
  onClick,
}: {
  files: File[];
  onClick: () => void;
}) {
  return (
    <Button disabled={files.length == 0} color="red" onClick={onClick}>
      Reset
    </Button>
  );
}

function SelectedFilesList({ files }: { files: File[] }) {
  const { classes } = useStyles();

  return (
    <Container className={classes.fileContainer}>
      <Stack align="center">
        <Title order={6}>Selected Files</Title>
        {files.length == 0 ? (
          <>No files selected.</>
        ) : (
          <List>
            {files.map((file, index) => (
              <List.Item className={classes.listItem} key={index}>
                {file.name}
              </List.Item>
            ))}
          </List>
        )}
        <ProcessFilesButton files={files} />
      </Stack>
    </Container>
  );
}

function ProcessFilesButton({ files }: { files: File[] }) {
  return (
    <Button
      disabled={files.length == 0}
      color="green"
      onClick={() => APIClient.batchPost(files)}
    >
      Process Files
    </Button>
  );
}
