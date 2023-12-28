"use client";

import styles from "@/styles/AttendanceRulePage.module.css";

import {
  Accordion,
  AccordionControl,
  AccordionItem,
  AccordionPanel,
  Container,
  Divider,
  Space,
  Text,
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
