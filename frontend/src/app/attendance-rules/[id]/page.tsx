"use client";

import styles from "@/styles/AttendanceRulePage.module.css";

import {
  Accordion,
  AccordionControl,
  AccordionItem,
  AccordionPanel,
  Button,
  Center,
  Container,
  Divider,
  FocusTrap,
  Group,
  Modal,
  Space,
  Text,
  Textarea,
  TextInput,
  Title,
} from "@mantine/core";
import { Params } from "@/app/attendance-rules/[id]/layout";
import { useState } from "react";
import {
  CoordinatingClass,
  CoordinatingClassGetResponse,
} from "@/api/coordinating_class";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { ClassAttendanceRule } from "@/api/class_attendance_rule";
import { useDisclosure } from "@mantine/hooks";
import { useForm } from "@mantine/form";

export default function AttendanceRulePage({ params }: { params: Params }) {
  const [data, setData] = useState<CoordinatingClassGetResponse>(
    {} as CoordinatingClassGetResponse,
  );
  const promiseFunc = async () => {
    const data = await APIClient.coordinatingClassGet(params.id);
    return setData(data);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Container className={styles.container} fluid>
        <CoordinatingClassDetails coordinatingClass={data.coordinating_class} />
        <Divider my="md" />
        <CreateRuleButton />
        <RuleDisplay rules={data.rules} />
      </Container>
    </RequestLoader>
  );
}

function CoordinatingClassDetails({
  coordinatingClass,
}: {
  coordinatingClass: CoordinatingClass;
}) {
  return (
    <Title order={2} ta="center">
      <Text span inherit c="teal">
        {coordinatingClass.code} {coordinatingClass.year}/
        {coordinatingClass.semester}
      </Text>{" "}
      - Attendance Rules
    </Title>
  );
}

function CreateRuleButton() {
  const [loading, setLoading] = useState(false);

  const [opened, { open, close }] = useDisclosure(false);
  const form = useForm({
    initialValues: {
      title: "",
      description: "",
      rule: "",
    },
    validate: {
      title: (value) => (value.length == 0 ? "Title cannot be empty" : null),
      description: (value) =>
        value.length == 0 ? "Description cannot be empty" : null,
      rule: (value) => (value.length == 0 ? "Rule cannot be empty" : null),
    },
  });

  return (
    <>
      <Modal
        opened={opened}
        onClose={close}
        centered
        title="Create New Rule"
        size="lg"
        overlayProps={{
          backgroundOpacity: 0.55,
          blur: 3,
        }}
      >
        <form
          onSubmit={form.onSubmit(async (values) => {
            setLoading(true);
            console.log(values);
            await new Promise((f) => setTimeout(f, 1000));
            setLoading(false);
          })}
        >
          <FocusTrap active>
            <TextInput
              label="Title"
              {...form.getInputProps("title")}
              disabled={loading}
              data-autofocus
            />
            <Textarea
              label="Description"
              disabled={loading}
              {...form.getInputProps("description")}
            />
            <Textarea
              label="Rule"
              disabled={loading}
              {...form.getInputProps("rule")}
              autosize
              minRows={6}
              maxRows={15}
            />
          </FocusTrap>
          <Space h="sm" />
          <Group justify="center">
            <Button
              onClick={form.reset}
              color="red"
              variant="light"
              disabled={loading}
            >
              Reset
            </Button>
            <Button type="submit" color="green" loading={loading}>
              Create
            </Button>
          </Group>
        </form>
      </Modal>
      <Center className={styles.createRuleButton}>
        <Button onClick={open}>Create New Rule</Button>
      </Center>
    </>
  );
}

function RuleDisplay({ rules }: { rules: ClassAttendanceRule[] }) {
  if (rules.length == 0) {
    return (
      <Text className={styles.noRulesText} ta="center">
        No rules for this class have been defined.
      </Text>
    );
  }

  const items = rules.map((rule, idx) => (
    <AccordionItem value={idx.toString()} key={idx}>
      <AccordionControl>
        <AccordionLabel title={rule.title} description={rule.description} />
      </AccordionControl>
      <AccordionPanel>
        <AccordionContent rule={rule} />
      </AccordionPanel>
    </AccordionItem>
  ));

  return (
    <Accordion className={styles.accordion} variant="contained">
      {items}
    </Accordion>
  );
}

function AccordionLabel({
  title,
  description,
}: {
  title: string;
  description: string;
}) {
  return (
    <>
      <Title order={4}>{title}</Title>
      <Space h="xs" />
      <Text size="sm" c="dimmed" lineClamp={3}>
        {description}
      </Text>
    </>
  );
}

function AccordionContent({ rule }: { rule: ClassAttendanceRule }) {
  return (
    <>
      <Divider className={styles.divider} />
      {rule.rule}
    </>
  );
}
