"use client";

import styles from "@/styles/DataExportPage.module.css";

import { Button, Center, Container, Space, Text, Title } from "@mantine/core";
import { APIClient } from "@/api/client";

export default function DataExportPage() {
  return (
    <Container className={styles.container} fluid>
      <PageTitle />
      <Details />
      <DataExportButton />
    </Container>
  );
}

function PageTitle() {
  return (
    <Title order={2} ta="center">
      Data Export Service
    </Title>
  );
}

function Details() {
  return (
    <Container className={styles.details} fluid>
      <Text ta="center">
        OAMS includes support for data export from our systems. The data is
        meant to be fed to other systems for integration or used for further
        data analysis.
      </Text>
      <Space h="md" />
      <Text ta="center">
        No sensitive information (i.e., signature data) from our systems will be
        exported in this process.
      </Text>
    </Container>
  );
}

function DataExportButton() {
  return (
    <Center>
      <Button
        onClick={async () => {
          const response = await APIClient.dataExportGet();

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
        }}
      >
        Start Data Export
      </Button>
    </Center>
  );
}
