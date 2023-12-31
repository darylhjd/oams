import { UseFormReturnType } from "@mantine/form";
import React, { FormEvent } from "react";
import {
  Anchor,
  Button,
  Center,
  Code,
  Fieldset,
  FocusTrap,
  Group,
  NativeSelect,
  NumberInput,
  Popover,
  Space,
  Text,
  Textarea,
  TextInput,
} from "@mantine/core";
import {
  CoordinatingClassPostRequest,
  RuleType,
} from "@/api/coordinating_class";

export function RuleForm({
  form,
  loading,
  onSubmit,
}: {
  form: UseFormReturnType<CoordinatingClassPostRequest>;
  loading: boolean;
  onSubmit: (event?: FormEvent<HTMLFormElement> | undefined) => void;
}) {
  let ruleDetailComponent: React.ReactNode;
  switch (form.getInputProps("rule_type").value) {
    case RuleType.MissedConsecutiveClasses:
      ruleDetailComponent = (
        <ConsecutiveClassRule form={form} loading={loading} />
      );
      break;
    case RuleType.MinPercentageAttendanceFromSession:
      ruleDetailComponent = (
        <PercentageClassRule form={form} loading={loading} />
      );
      break;
    case RuleType.Advanced:
      ruleDetailComponent = <AdvancedForm form={form} loading={loading} />;
      break;
    default:
      ruleDetailComponent = null;
  }

  return (
    <form onSubmit={onSubmit}>
      <FocusTrap active>
        <div>
          <TextInput
            disabled={loading}
            withAsterisk
            label="Title"
            {...form.getInputProps("title")}
            data-autofocus
          />
          <Space h="sm" />
          <Textarea
            disabled={loading}
            withAsterisk
            label="Description"
            {...form.getInputProps("description")}
          />
          <Space h="sm" />
          <NativeSelect
            disabled={loading}
            label="Rule Type"
            data={[
              {
                value: RuleType.MissedConsecutiveClasses.toString(),
                label: "Missed Consecutive Classes",
              },
              {
                value: RuleType.MinPercentageAttendanceFromSession.toString(),
                label: "Minimum Percentage Attendance",
              },
              {
                value: RuleType.Advanced.toString(),
                label: "Advanced",
              },
            ]}
            onChange={(value) =>
              form.setFieldValue("rule_type", Number(value.currentTarget.value))
            }
          />
          <Space h="sm" />
          <Fieldset legend="Rule Details">{ruleDetailComponent}</Fieldset>
          <Space h="sm" />
          <Group justify="center">
            <Button
              disabled={loading}
              onClick={form.reset}
              color="red"
              variant="light"
            >
              Reset
            </Button>
            <Button type="submit" color="green" loading={loading}>
              Create
            </Button>
          </Group>
        </div>
      </FocusTrap>
    </form>
  );
}

function ConsecutiveClassRule({
  form,
  loading,
}: {
  form: UseFormReturnType<CoordinatingClassPostRequest>;
  loading: boolean;
}) {
  return (
    <>
      <Text c="dimmed" size="sm" ta="center">
        This rule triggers when a student misses the last few consecutive
        classes.
      </Text>
      <Space h="sm" />
      <NumberInput
        disabled={loading}
        withAsterisk
        label="Number of consecutive missed classes"
        defaultValue={4}
        allowNegative={false}
        allowDecimal={false}
        min={1}
        {...form.getInputProps("consecutive_params.consecutive_classes")}
      />
    </>
  );
}

function PercentageClassRule({
  form,
  loading,
}: {
  form: UseFormReturnType<CoordinatingClassPostRequest>;
  loading: boolean;
}) {
  return (
    <>
      <Text c="dimmed" size="sm" ta="center">
        This rule triggers if a student fails to maintain the required
        percentage of attendance beginning from a certain session. Sessions
        before will be considered in the total number of sessions.
      </Text>
      <Space h="sm" />
      <NumberInput
        disabled={loading}
        withAsterisk
        label="Minimum required attendance percentage"
        defaultValue={70}
        allowNegative={false}
        decimalScale={2}
        fixedDecimalScale
        suffix="%"
        {...form.getInputProps("percentage_params.percentage")}
      />
      <Space h="sm" />
      <NumberInput
        disabled={loading}
        withAsterisk
        label="From session"
        defaultValue={4}
        allowNegative={false}
        allowDecimal={false}
        min={1}
        {...form.getInputProps("percentage_params.from_session")}
      />
    </>
  );
}

function AdvancedForm({
  form,
  loading,
}: {
  form: UseFormReturnType<CoordinatingClassPostRequest>;
  loading: boolean;
}) {
  const variables = `var enrollments []struct {
	ID        int64    
	SessionID int64    
	UserID    string   
	Attended  bool     
	CreatedAt time.Time
	UpdatedAt time.Time
}`;

  return (
    <>
      <Text c="dimmed" size="sm" ta="center">
        OAMS allows you to specify custom rules. Use the provided variables to
        form custom conditions to trigger alerts. The language definition can be
        found{" "}
        <Anchor
          href="https://expr-lang.org/docs/language-definition"
          target="_blank"
        >
          here
        </Anchor>
        .
      </Text>
      <Space h="xs" />
      <Center>
        <Popover position="bottom" withArrow shadow="md">
          <Popover.Target>
            <Button variant="outline" size="sm" color="gray">
              View variables
            </Button>
          </Popover.Target>
          <Popover.Dropdown>
            <Code block>{variables}</Code>
          </Popover.Dropdown>
        </Popover>
      </Center>
      <Space h="sm" />
      <Textarea
        disabled={loading}
        withAsterisk
        autosize
        label="Rule"
        minRows={5}
        maxRows={5}
        {...form.getInputProps("advanced_params.rule")}
      />
    </>
  );
}
