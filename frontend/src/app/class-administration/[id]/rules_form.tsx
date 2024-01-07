import styles from "@/styles/ClassAdministrationPage.module.css";

import { UseFormReturnType } from "@mantine/form";
import { FormEvent, ReactNode } from "react";
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
  PopoverDropdown,
  PopoverTarget,
  Space,
  Text,
  Textarea,
  TextInput,
} from "@mantine/core";
import {
  CoordinatingClassRulesPostRequest,
  RuleType,
} from "@/api/coordinating_class";
import { CodeHighlight } from "@mantine/code-highlight";

export function RuleForm({
  form,
  loading,
  onSubmit,
}: {
  form: UseFormReturnType<CoordinatingClassRulesPostRequest>;
  loading: boolean;
  onSubmit: (event?: FormEvent<HTMLFormElement> | undefined) => void;
}) {
  let ruleDetailComponent: ReactNode;
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
  form: UseFormReturnType<CoordinatingClassRulesPostRequest>;
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
  form: UseFormReturnType<CoordinatingClassRulesPostRequest>;
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
        max={100}
        clampBehavior="strict"
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
  form: UseFormReturnType<CoordinatingClassRulesPostRequest>;
  loading: boolean;
}) {
  const variables = `// Array of enrollment info for a user.
// The array is sorted by start time and then end time in ascending order.
var enrollments []struct {
	ClassID       int64    
	ClassCode     string   
	ClassYear     int32    
	ClassSemester string    
	ClassType     string    
	StartTime     time.Time 
	EndTime       time.Time 
	Venue         string    
	UserID        string    
	UserName      string    
	UserEmail     string   
	Attended      bool     
}`;

  return (
    <>
      <Text c="dimmed" size="sm" ta="center">
        OAMS allows you to specify custom rules. Use the provided variables to
        form custom conditions. An alert is triggered when the condition returns{" "}
        <Code>true</Code> for a user. The language definition can be found{" "}
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
        <Popover
          position="bottom"
          withArrow
          shadow="md"
          classNames={{
            dropdown: styles.variableDropdown,
          }}
        >
          <PopoverTarget>
            <Button variant="outline" size="sm" color="gray">
              View variables
            </Button>
          </PopoverTarget>
          <PopoverDropdown>
            <CodeHighlight
              code={variables}
              language="golang"
              withCopyButton={false}
            />
          </PopoverDropdown>
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
