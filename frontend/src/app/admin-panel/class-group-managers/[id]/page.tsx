"use client";

import styles from "@/styles/AdminPanelClassGroupManager.module.css";

import { Params } from "@/app/admin-panel/class-group-managers/[id]/layout";
import {
  Anchor,
  Button,
  Container,
  Fieldset,
  Group,
  Modal,
  NativeSelect,
  Space,
  Text,
} from "@mantine/core";
import { useState } from "react";
import { ClassGroupManager } from "@/api/class_group_manager";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { ManagingRole } from "@/api/class_group_manager";
import { Routes } from "@/routing/routes";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { notifications } from "@mantine/notifications";
import { IconX } from "@tabler/icons-react";
import { getError } from "@/api/error";

export default function AdminPanelClassGroupManagerPage({
  params,
}: {
  params: Params;
}) {
  const [manager, setManager] = useState<ClassGroupManager>(
    {} as ClassGroupManager,
  );
  const promiseFunc = async () => {
    const data = await APIClient.classGroupManagerGet(params.id);
    return setManager(data.manager);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <ClassGroupManagerDisplay manager={manager} />
    </RequestLoader>
  );
}

function ClassGroupManagerDisplay({ manager }: { manager: ClassGroupManager }) {
  return (
    <Container className={styles.container} fluid>
      <ManagerDetails manager={manager} />
      <Space h="xl" />
      <ManagerSettings manager={manager} />
      <Space h="xl" />
      <ManagerDangerZone manager={manager} />
    </Container>
  );
}

function ManagerDetails({ manager }: { manager: ClassGroupManager }) {
  return (
    <Fieldset legend="Details">
      <Text ta="center">ID: {manager.id}</Text>
      <Text ta="center">User: {manager.user_id}</Text>
      <Text ta="center">
        Class Group:{" "}
        <Anchor
          href={`${Routes.adminPanelClassGroup}/${manager.class_group_id}`}
        >
          {manager.class_group_id}
        </Anchor>
      </Text>
    </Fieldset>
  );
}

function ManagerSettings({ manager }: { manager: ClassGroupManager }) {
  const [loading, setLoading] = useState(false);
  const form = useForm({
    initialValues: {
      role: manager.managing_role,
    },
  });

  const changeRoleOnSubmit = form.onSubmit(async (values) => {
    setLoading(true);
    try {
      await APIClient.classGroupManagerPatch(manager.id, values.role);
      form.resetDirty();
    } catch (e) {
      notifications.show({
        title: "Could not update role",
        message: getError(e),
        icon: <IconX />,
        color: "red",
      });
    }
    setLoading(false);
  });

  return (
    <Fieldset legend="Settings">
      <form onSubmit={changeRoleOnSubmit}>
        <NativeSelect
          label="Managing Role"
          description="Set managing role"
          data={[
            ManagingRole.CourseCoordinator,
            ManagingRole.TeachingAssistant,
          ]}
          {...form.getInputProps("role")}
        />
        <Space h="md" />
        <Group justify="center">
          <Button
            type="submit"
            variant="filled"
            disabled={!form.isDirty()}
            loading={loading}
          >
            Confirm
          </Button>
        </Group>
      </form>
    </Fieldset>
  );
}

function ManagerDangerZone({ manager }: { manager: ClassGroupManager }) {
  const [opened, { open, close }] = useDisclosure(false);

  const deleteManager = async () => {
    close();
  };

  return (
    <Fieldset className={styles.danger} legend="DANGER ZONE">
      <Group justify="space-between">
        <div>
          <Text>Delete Manager</Text>
          <Text c="dimmed">
            This user will no longer be a manager for this class group.
          </Text>
        </div>
        <div>
          <Modal
            opened={opened}
            onClose={close}
            centered
            title="Deleting Class Group Manager"
          >
            <Text ta="center">This action is irreversible.</Text>
            <Text ta="center">Are you sure?</Text>
            <Space h="md" />
            <Group justify="center">
              <Button onClick={close} variant="light">
                Cancel
              </Button>
              <Button color="red" variant="filled" onClick={deleteManager}>
                Confirm
              </Button>
            </Group>
          </Modal>
          <Button color="red" variant="light" onClick={open}>
            Delete
          </Button>
        </div>
      </Group>
    </Fieldset>
  );
}
