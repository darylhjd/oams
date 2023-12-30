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
  NumberInput,
  Popover,
  SegmentedControl,
  Select,
  Space,
  Text,
  Textarea,
  TextInput,
} from "@mantine/core";

export type RuleFormParams = {
  title: string;
  description: string;
  rule_type: string;
  preset_rule: string;
  consecutive_params: {
    num: number;
  };
  percentage_params: {
    percentage: number;
    from: number;
  };
  rule: string;
};

export enum RuleType {
  Simple = "simple",
  Advanced = "advanced",
}

export enum PresetRule {
  MissedConsecutiveClasses = "missed_consecutive_classes",
  MinPercentageAttendanceFromSession = "min_percentage_attendance_from_session",
}

export function RuleForm({
  form,
  loading,
  onSubmit,
}: {
  form: UseFormReturnType<RuleFormParams>;
  loading: boolean;
  onSubmit: (event?: FormEvent<HTMLFormElement> | undefined) => void;
}) {
  return (
    <form onSubmit={onSubmit}>
      <FocusTrap active>
        <TextInput
          disabled={loading}
          withAsterisk
          label="Title"
          {...form.getInputProps("title")}
          data-autofocus
        />
      </FocusTrap>
      <Space h="sm" />
      <Textarea
        disabled={loading}
        withAsterisk
        label="Description"
        {...form.getInputProps("description")}
      />
      <Space h="sm" />
      <SegmentedControl
        disabled={loading}
        fullWidth
        data={[
          { value: RuleType.Simple, label: "Simple" },
          { value: RuleType.Advanced, label: "Advanced" },
        ]}
        {...form.getInputProps("rule_type")}
      />
      <Space h="sm" />
      <>
        <SimpleRule form={form} loading={loading} />
        <AdvancedForm form={form} loading={loading} />
      </>
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
    </form>
  );
}

function SimpleRule({
  form,
  loading,
}: {
  form: UseFormReturnType<RuleFormParams>;
  loading: boolean;
}) {
  if (form.getInputProps("rule_type").value != RuleType.Simple) {
    return null;
  }

  return (
    <>
      <Select
        disabled={loading}
        label="Preset Rule"
        allowDeselect={false}
        defaultValue={PresetRule.MissedConsecutiveClasses}
        data={[
          {
            value: PresetRule.MissedConsecutiveClasses,
            label: "Missed Consecutive Classes",
          },
          {
            value: PresetRule.MinPercentageAttendanceFromSession,
            label: "Minimum Percentage Attendance From Session",
          },
        ]}
        {...form.getInputProps("preset_rule")}
      />
      <Space h="sm" />
      <Fieldset legend="Rule Details">
        <ConsecutiveClassRule form={form} loading={loading} />
        <PercentageClassRule form={form} loading={loading} />
      </Fieldset>
    </>
  );
}

function ConsecutiveClassRule({
  form,
  loading,
}: {
  form: UseFormReturnType<RuleFormParams>;
  loading: boolean;
}) {
  if (
    form.getInputProps("preset_rule").value !=
    PresetRule.MissedConsecutiveClasses
  ) {
    return null;
  }

  return (
    <>
      <Text c="dimmed" size="sm" ta="center">
        This rule triggers when a student misses at least the last few
        consecutive classes.
      </Text>
      <Space h="sm" />
      <NumberInput
        disabled={loading}
        withAsterisk
        label="Number of consecutive missed classes"
        defaultValue={2}
        allowNegative={false}
        allowDecimal={false}
        min={1}
        {...form.getInputProps("consecutive_params.num")}
      />
    </>
  );
}

function PercentageClassRule({
  form,
  loading,
}: {
  form: UseFormReturnType<RuleFormParams>;
  loading: boolean;
}) {
  if (
    form.getInputProps("preset_rule").value !=
    PresetRule.MinPercentageAttendanceFromSession
  ) {
    return null;
  }

  return (
    <>
      <Text c="dimmed" size="sm" ta="center">
        This rule triggers if a student fails to maintain the required
        percentage of attendance beginning from a certain session.
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
        {...form.getInputProps("percentage_params.from")}
      />
    </>
  );
}

function AdvancedForm({
  form,
  loading,
}: {
  form: UseFormReturnType<RuleFormParams>;
  loading: boolean;
}) {
  if (form.getInputProps("rule_type").value != RuleType.Advanced) {
    return null;
  }

  const variables = `var enrollments []struct {
	ID        int64    
	SessionID int64    
	UserID    string   
	Attended  bool     
	CreatedAt time.Time
	UpdatedAt time.Time
}`;

  return (
    <Fieldset>
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
        {...form.getInputProps("rule")}
      />
    </Fieldset>
  );
}
