"use client";

import { APIClient } from "@/api/client";
import {
  Button,
  Center,
  Container,
  Divider,
  FileButton,
  List,
  NativeSelect,
  Stack,
  Tabs,
  Text,
  Title,
  createStyles,
} from "@mantine/core";
import { Dispatch, SetStateAction, useState } from "react";
import { batchesStore } from "./batches_store";
import { redirectIfNotUserRole } from "@/routes/checks";
import { UserRole } from "@/api/models";

const useStyles = createStyles((theme) => ({
  fileContainer: {
    padding: "2em 0",
  },

  listItem: {
    listStyleType: "none",
    margin: 0,
    padding: 0,
  },

  batchChooserButton: {
    padding: "1.5em 0",
  },
}));

export default function BatchProcessingPage() {
  const [files, setFiles] = useState<File[]>([]);
  const clearFiles = () => setFiles([]);

  const batches = batchesStore();

  if (redirectIfNotUserRole(UserRole.SystemAdmin)) {
    return null;
  }

  return (
    <>
      <Container>
        <Center>
          <Stack>
            <Text align="center">Upload your batch files here.</Text>
            <ChooseFilesButton onChange={setFiles} />
            <ResetFilesButton
              files={files}
              onClick={() => {
                clearFiles();
                batches.invalidate();
              }}
            />
          </Stack>
        </Center>
      </Container>
      <SelectedFilesList files={files} />
      <Divider my="md" />
      <BatchData />
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
  const batches = batchesStore();

  return (
    <Button
      disabled={files.length == 0}
      color="green"
      onClick={async () => {
        const data = await APIClient.batchPost(files);
        batches.setData(data);
      }}
    >
      Process Files
    </Button>
  );
}

function BatchData() {
  const batches = batchesStore();
  const [filename, setFilename] = useState("");

  if (batches.data == null) {
    return (
      <>
        <Title align="center" order={6}>
          Batch Data
        </Title>
        <Text align="center">Choose some files to process!</Text>
      </>
    );
  }

  return (
    <>
      <Title align="center" order={6}>
        Batch Data
      </Title>
      <BatchChooserMenu onChange={setFilename} />
      <BatchTabViewer filename={filename} />
    </>
  );
}

function BatchChooserMenu({
  onChange,
}: {
  onChange: Dispatch<SetStateAction<string>>;
}) {
  const { classes } = useStyles();
  const batches = batchesStore();

  if (batches.data == null) {
    return null;
  }

  var data = batches.data.batches.map((batch) => ({
    value: batch.filename,
    label: batch.filename,
  }));
  data.unshift({ value: "", label: "Select a file" });

  return (
    <Center>
      <Container className={classes.batchChooserButton}>
        <NativeSelect
          onChange={(event) => onChange(event.currentTarget.value)}
          data={data}
          label="Select uploaded file"
          variant="filled"
        />
      </Container>
    </Center>
  );
}

function BatchTabViewer({ filename }: { filename: string }) {
  const batches = batchesStore();

  var batch = null;
  for (var b of batches.data!.batches) {
    console.log(b.filename);
    if (b.filename == filename) {
      batch = b;
      break;
    }
  }

  if (batch == null) {
    return null;
  }

  return (
    <Container>
      <Tabs defaultValue="classes">
        <Tabs.List>
          <Tabs.Tab value="classes">Classes</Tabs.Tab>
          <Tabs.Tab value="classGroups">Class Groups</Tabs.Tab>
          <Tabs.Tab value="classGroupSessions">Class Group Sessions</Tabs.Tab>
          <Tabs.Tab value="users">Users</Tabs.Tab>
          <Tabs.Tab value="sessionEnrollments">Session Enrollments</Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="classes">{batch.filename}</Tabs.Panel>
        <Tabs.Panel value="classGroups">Class Groups Content</Tabs.Panel>
        <Tabs.Panel value="classGroupSessions">
          Class Group Sessions Content
        </Tabs.Panel>
        <Tabs.Panel value="users">Users Content</Tabs.Panel>
        <Tabs.Panel value="sessionEnrollments">
          Session Enrollments Content
        </Tabs.Panel>
      </Tabs>
    </Container>
  );
}
