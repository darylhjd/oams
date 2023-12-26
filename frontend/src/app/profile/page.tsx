"use client";

import styles from "@/styles/ProfilePage.module.css";

import { useSessionUserStore } from "@/stores/session";
import {
  Box,
  Button,
  Container,
  Group,
  Paper,
  PasswordInput,
  Space,
  Text,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { APIClient } from "@/api/client";
import { notifications } from "@mantine/notifications";
import { IconCheck } from "@tabler/icons-react";
import { useState } from "react";
import { signatureInputForm } from "@/components/signature_form";

export default function ProfilePage() {
  const session = useSessionUserStore();
  const user = session.data!.user;

  return (
    <Container fluid>
      <Paper className={styles.paper} radius="md" shadow="xs" withBorder p="xl">
        <Text ta="center" size="xl" fw={1000}>
          User ID: {user.id}
        </Text>
        <Space h="md" />
        <Text ta="center" size="sm">
          {user.role} • {user.email}
        </Text>
        <Space h="xs" />
        <Text c="dimmed" fs="italic" ta="center" size="sm">
          {user.name ? user.name : "No name registered"}
        </Text>
      </Paper>
      <SignatureUpdater userId={user.id} />
    </Container>
  );
}

function SignatureUpdater({ userId }: { userId: string }) {
  const [loading, setLoading] = useState(false);
  const form = useForm(signatureInputForm());

  return (
    <Box className={styles.signatureUpdater}>
      <form
        onSubmit={form.onSubmit(async (values) => {
          setLoading(true);
          await APIClient.signaturePut(userId, values.signature);
          form.reset();
          setLoading(false);
          notifications.show({
            title: "Update successful!",
            message: "Your signature has been updated successfully.",
            icon: <IconCheck />,
            color: "green",
          });
        })}
      >
        <PasswordInput
          label="New Attendance Signature"
          description="This signature is used for signing your attendance."
          placeholder="New attendance signature"
          {...form.getInputProps("signature")}
        />
        <Text fs="italic" size="sm">
          Note: Your default signature is your User ID. We recommend that you
          change your signature from the default.
        </Text>
        <Space h="xs" />
        <Group justify="flex-end">
          <Button type="submit" loading={loading}>
            Update Signature
          </Button>
        </Group>
      </form>
    </Box>
  );
}
