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
import { notifications } from "@mantine/notifications";
import { Dispatch, SetStateAction, useState } from "react";
import { batchesStore } from "./batches_store";
import { redirectIfNotUserRole } from "@/routes/checks";
import { UserRole } from "@/api/models";
import {
  ClassGroupSessionsTable,
  ClassGroupsTable,
  ClassesTable,
  UsersTable,
} from "./batch_tables";
import { useMediaQuery } from "@mantine/hooks";
import { MOBILE_MIN_WIDTH } from "@/components/responsive";
import { IconCheck, IconCross, IconX } from "@tabler/icons-react";

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

  tabList: {
    overflowY: "hidden",
    overflowX: "auto",
    flexWrap: "nowrap",
  },

  tabTab: {
    padding: "1em 1em",
  },

  batchPutButton: {
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
        <PreviewBatchDataButton files={files} />
      </Stack>
    </Container>
  );
}

function PreviewBatchDataButton({ files }: { files: File[] }) {
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
      Preview Batch Data
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
      <ConfirmBatchPutButton />
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
  const { classes } = useStyles();
  const batches = batchesStore();
  const isMobile = useMediaQuery(MOBILE_MIN_WIDTH);

  var batch = null;
  for (var b of batches.data!.batches) {
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
      <Tabs defaultValue="classes" variant="outline">
        <Tabs.List
          className={classes.tabList}
          position={isMobile ? "left" : "center"}
        >
          <Tabs.Tab className={classes.tabTab} value="classes">
            Classes
          </Tabs.Tab>
          <Tabs.Tab value="classGroups">Class Groups</Tabs.Tab>
          <Tabs.Tab value="classGroupSessions">Class Group Sessions</Tabs.Tab>
          <Tabs.Tab value="users">Users</Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="classes">
          <ClassesTable cls={batch.class} />
        </Tabs.Panel>
        <Tabs.Panel value="classGroups">
          <ClassGroupsTable classGroups={batch.class_groups} />
        </Tabs.Panel>
        <Tabs.Panel value="classGroupSessions">
          <ClassGroupSessionsTable classGroups={batch.class_groups} />
        </Tabs.Panel>
        <Tabs.Panel value="users">
          <UsersTable classGroups={batch.class_groups} />
        </Tabs.Panel>
      </Tabs>
    </Container>
  );
}

function ConfirmBatchPutButton() {
  const { classes } = useStyles();
  const batches = batchesStore();

  if (batches.data == null) {
    return null;
  }

  return (
    <Container className={classes.batchPutButton}>
      <Center>
        <Button
          onClick={async () => {
            notifications.show({
              id: "loading",
              title: "Processing...",
              message: "Your data is being processed. Please wait.",
              loading: true,
            });

            const result = await APIClient.batchPut(batches.data!);
            if (result == null) {
              notifications.show({
                title: "Oh no!",
                message:
                  "There was an error processing your batch data. Please try again later",
                icon: <IconCross />,
                color: "red",
              });
              return;
            }

            notifications.show({
              title: "Success!",
              message: "All batch data has been processed!",
              icon: <IconCheck />,
              color: "teal",
            });
          }}
        >
          Confirm Data Processing
        </Button>
      </Center>
    </Container>
  );
}
