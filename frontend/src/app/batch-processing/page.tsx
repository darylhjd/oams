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
import { IconCheck, IconX } from "@tabler/icons-react";
import { v4 as uuidv4 } from "uuid";

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
      <BatchTabViewer />
      <ConfirmBatchPutButton />
    </>
  );
}

function BatchTabViewer() {
  const { classes } = useStyles();
  const batches = batchesStore();
  const isMobile = useMediaQuery(MOBILE_MIN_WIDTH);

  if (batches.data == null) {
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
          <ClassesTable batches={batches.data} />
        </Tabs.Panel>
        <Tabs.Panel value="classGroups">
          <ClassGroupsTable batches={batches.data} />
        </Tabs.Panel>
        <Tabs.Panel value="classGroupSessions">
          <ClassGroupSessionsTable batches={batches.data} />
        </Tabs.Panel>
        <Tabs.Panel value="users">
          <UsersTable batches={batches.data} />
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
        <Button color="green" onClick={confirmDataProcessingAction}>
          Confirm Data Processing
        </Button>
      </Center>
    </Container>
  );
}

async function confirmDataProcessingAction() {
  const batches = batchesStore();

  const loadingId = uuidv4();
  notifications.show({
    id: loadingId,
    title: "Processing...",
    message: "Your data is being processed. Please wait.",
    loading: true,
  });

  const result = await APIClient.batchPut({ batches: batches.data! });
  if (result == null) {
    notifications.show({
      title: "Oh no!",
      message:
        "There was an error processing your batch data. Please try again later",
      icon: <IconX />,
      color: "red",
    });
    return;
  }

  notifications.hide(loadingId);
  notifications.show({
    title: "Success!",
    message: "All batch data has been processed!",
    icon: <IconCheck />,
    color: "teal",
  });
}
