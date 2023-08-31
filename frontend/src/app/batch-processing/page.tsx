"use client";

import { APIClient } from "@/api/client";
import { UserRole } from "@/api/models";
import { buildRedirectUrlQueryParamsString } from "@/routes/checks";
import { Routes } from "@/routes/routes";
import { sessionStore } from "@/states/session";
import {
  Button,
  Center,
  Container,
  FileButton,
  List,
  Stack,
  Title,
  createStyles,
} from "@mantine/core";
import { getURL } from "next/dist/shared/lib/utils";
import { useRouter } from "next/navigation";
import { Dispatch, SetStateAction, useEffect, useState } from "react";
import { batchStore } from "./batch_store";

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
  const [files, setFiles] = useState<File[]>([]);
  const clearFiles = () => setFiles([]);

  const router = useRouter();
  const session = sessionStore();
  const batches = batchStore();

  useEffect(() => {
    if (session.data == null) {
      router.replace(
        `${Routes.login}?${buildRedirectUrlQueryParamsString(getURL())}`,
      );
    } else if (session.data!.session_user.role != UserRole.SystemAdmin) {
      router.replace(Routes.home);
    }
  }, [router, session]);

  if (
    session.data == null ||
    session.data!.session_user.role != UserRole.SystemAdmin
  ) {
    return null;
  }

  return (
    <>
      <Container>
        <Center>
          <Stack>
            <Center>
              <p>Upload your batch files here.</p>
            </Center>
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
  const batches = batchStore();

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
  const batches = batchStore();

  if (batches.data == null) {
    return null;
  }

  return (
    <Center>
      {batches.data.batches.map((batch) => (
        <>
          <p>{batch.filename}</p>
          {batch.class_groups.map((classGroup) => (
            <p>{classGroup.name}</p>
          ))}
        </>
      ))}
    </Center>
  );
}
