import styles from "@/styles/ClassAdministrationPageRules.module.css";

import { ClassAttendanceRule } from "@/api/class_attendance_rule";
import {
  Accordion,
  AccordionControl,
  AccordionItem,
  AccordionPanel,
  Button,
  Center,
  Divider,
  Group,
  Modal,
  Space,
  Switch,
  Text,
  Title,
} from "@mantine/core";
import { Dispatch, SetStateAction, useState } from "react";
import { CodeHighlightTabs } from "@mantine/code-highlight";
import { useDisclosure } from "@mantine/hooks";
import { useForm, UseFormReturnType } from "@mantine/form";
import {
  CoordinatingClassRulesPostRequest,
  RuleType,
} from "@/api/coordinating_class";
import { APIClient } from "@/api/client";
import { notifications } from "@mantine/notifications";
import { getError } from "@/api/error";
import { IconCheck, IconX } from "@tabler/icons-react";
import { RuleForm } from "@/app/class-administration/[id]/rules_form";
import { RequestLoader } from "@/components/request_loader";
import { Panel } from "@/app/class-administration/[id]/panel";

export function RulesTab({ id }: { id: number }) {
  const [rules, setRules] = useState<ClassAttendanceRule[]>([]);
  const promiseFunc = async () => {
    const rulesData = await APIClient.coordinatingClassRulesGet(id);
    setRules(rulesData.rules);
  };

  return (
    <Panel>
      <RequestLoader promiseFunc={promiseFunc}>
        <CreateRuleButton id={id} rules={rules} setRules={setRules} />
        <RuleDisplay id={id} rules={rules} setRules={setRules} />
      </RequestLoader>
    </Panel>
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
  const form: UseFormReturnType<CoordinatingClassRulesPostRequest> = useForm({
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
      const resp = await APIClient.coordinatingClassRulesPost(id, values);
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

function RuleDisplay({
  id,
  rules,
  setRules,
}: {
  id: number;
  rules: ClassAttendanceRule[];
  setRules: Dispatch<SetStateAction<ClassAttendanceRule[]>>;
}) {
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
        <RuleLabel rule={rule} />
      </AccordionControl>
      <AccordionPanel>
        <RuleContent
          classId={id}
          rule={rule}
          rules={rules}
          setRules={setRules}
        />
      </AccordionPanel>
    </AccordionItem>
  ));

  return (
    <Accordion className={styles.accordion} variant="contained">
      {items}
    </Accordion>
  );
}

function RuleLabel({ rule }: { rule: ClassAttendanceRule }) {
  return (
    <Group justify="space-between">
      <div>
        <Title order={4}>{rule.title}</Title>
        <Space h="xs" />
        <Text size="sm" c="dimmed" lineClamp={3}>
          {rule.description}
        </Text>
      </div>
      <div>
        <Text size="sm" c="dimmed" className={styles.creatorLabel}>
          Created By: {rule.creator_id}
        </Text>
      </div>
    </Group>
  );
}

function RuleContent({
  classId,
  rule,
  rules,
  setRules,
}: {
  classId: number;
  rule: ClassAttendanceRule;
  rules: ClassAttendanceRule[];
  setRules: Dispatch<SetStateAction<ClassAttendanceRule[]>>;
}) {
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
              rule.environment,
              (k, v) => (k == "env_type" ? undefined : v),
              4,
            ),
            language: "json",
          },
        ]}
      />
      <Space h="lg" />
      <RuleActions
        classId={classId}
        rule={rule}
        rules={rules}
        setRules={setRules}
      />
    </>
  );
}

function RuleActions({
  classId,
  rule,
  rules,
  setRules,
}: {
  classId: number;
  rule: ClassAttendanceRule;
  rules: ClassAttendanceRule[];
  setRules: Dispatch<SetStateAction<ClassAttendanceRule[]>>;
}) {
  const [checked, setChecked] = useState(rule.active);
  const [opened, { open, close }] = useDisclosure(false);

  return (
    <Group justify="left">
      <Switch
        size="md"
        checked={checked}
        color="teal"
        label={checked ? "Enabled" : "Disabled"}
        thumbIcon={
          checked ? (
            <IconCheck size={16} color="green" />
          ) : (
            <IconX size={16} color="red" />
          )
        }
        onChange={async (event) => {
          try {
            const response = await APIClient.coordinatingClassRulePatch(
              classId,
              rule.id,
              event.currentTarget.checked,
            );
            setChecked(response.active);
          } catch (error) {
            notifications.show({
              title: "Rule Toggle Error",
              message: getError(error),
              icon: <IconX />,
              color: "red",
            });
          }
        }}
      />
      <>
        <Modal opened={opened} onClose={close} title="Delete Rule" centered>
          <Text ta="center">
            You are about to delete this rule. Are you sure?
          </Text>
          <Space h="lg" />
          <Group justify="center">
            <Button variant="outline" onClick={close}>
              Cancel
            </Button>
            <Button
              color="red"
              onClick={async () => {
                try {
                  await APIClient.coordinatingClassRuleDelete(classId, rule.id);
                  close();
                  setRules([...rules.filter((r) => r.id != rule.id)]);
                } catch (error) {
                  notifications.show({
                    title: "Rule Delete Error",
                    message: getError(error),
                    icon: <IconX />,
                    color: "red",
                  });
                }
              }}
            >
              Yes, I&apos;m sure!
            </Button>
          </Group>
        </Modal>
        <Button variant="outline" color="red" onClick={open}>
          Delete Rule
        </Button>
      </>
    </Group>
  );
}
