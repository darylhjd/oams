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
  Modal,
  Space,
  Text,
  Title,
} from "@mantine/core";
import { Params } from "@/app/attendance-rules/[id]/layout";
import React, { Dispatch, SetStateAction, useState } from "react";
import {
  CoordinatingClass,
  CoordinatingClassPostRequest,
  RuleType,
} from "@/api/coordinating_class";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { ClassAttendanceRule } from "@/api/class_attendance_rule";
import { useDisclosure } from "@mantine/hooks";
import { useForm, UseFormReturnType } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { IconX } from "@tabler/icons-react";
import { getError } from "@/api/error";
import { RuleForm } from "@/app/attendance-rules/[id]/rule_form";
import { CodeHighlightTabs } from "@mantine/code-highlight";

export default function AttendanceRulePage({ params }: { params: Params }) {
  const [coordinatingClass, setCoordinatingClass] = useState<CoordinatingClass>(
    {} as CoordinatingClass,
  );
  const [rules, setRules] = useState<ClassAttendanceRule[]>([]);
  const promiseFunc = async () => {
    const data = await APIClient.coordinatingClassGet(params.id);
    setCoordinatingClass(data.coordinating_class);
    return setRules(data.rules);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Container className={styles.container} fluid>
        <CoordinatingClassDetails coordinatingClass={coordinatingClass} />
        <Divider my="md" />
        <CreateRuleButton id={params.id} rules={rules} setRules={setRules} />
        <RuleDisplay rules={rules} />
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

function CreateRuleButton({
  id,
  rules,
  setRules,
}: {
  id: number;
  rules: ClassAttendanceRule[];
  setRules: Dispatch<SetStateAction<ClassAttendanceRule[]>>;
}) {
  const [opened, { open, close }] = useDisclosure(false);

  const [loading, setLoading] = useState(false);
  const form: UseFormReturnType<CoordinatingClassPostRequest> = useForm({
    initialValues: {
      title: "",
      description: "",
      rule_type: RuleType.MissedConsecutiveClasses as number,
      consecutive_params: {
        consecutive_classes: 4,
      },
      percentage_params: {
        percentage: 70,
        from_session: 4,
      },
      advanced_params: { rule: "" },
    },
    validate: {
      title: (value) => (value.length == 0 ? "Title cannot be empty" : null),
      description: (value) =>
        value.length == 0 ? "Description cannot be empty" : null,
    },
  });

  const formSubmit = form.onSubmit(async (values) => {
    setLoading(true);
    try {
      const resp = await APIClient.coordinatingClassPost(id, values);
      close();
      form.reset();
      rules.push(resp.rule);
      setRules([...rules]);
    } catch (e) {
      notifications.show({
        title: "Rule Validation Failed",
        message: getError(e),
        icon: <IconX />,
        color: "red",
      });
    }
    setLoading(false);
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
        <RuleForm form={form} loading={loading} onSubmit={formSubmit} />
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
        No rules have been defined for this class.
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
      <CodeHighlightTabs
        code={[
          {
            fileName: "Rule",
            code: rule.rule,
            language: "typescript",
          },
          {
            fileName: "Environment",
            code: JSON.stringify(
              rule.environment.data,
              (k, v) => (k == "env_type" ? undefined : v),
              4,
            ),
            language: "json",
          },
        ]}
      />
    </>
  );
}
