import { Dispatch, SetStateAction, useState } from "react";
import { APIClient } from "@/api/client";
import { Panel } from "@/app/class-administration/[id]/panel";
import { RequestLoader } from "@/components/request_loader";
import { ScheduleData } from "@/api/coordinating_class";
import {
  MantineReactTable,
  MRT_DensityState,
  MRT_Row,
  useMantineReactTable,
} from "mantine-react-table";
import { CoordinatingClassScheduleTableColumns } from "@/components/columns";
import {
  ActionIcon,
  Button,
  Group,
  Modal,
  Space,
  Tooltip,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { IconAdjustments, IconX } from "@tabler/icons-react";
import { DateTimePicker } from "@mantine/dates";
import { useForm, UseFormReturnType } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { getError } from "@/api/error";
import { useRouter } from "next/navigation";
import { Routes } from "@/routing/routes";

export function ScheduleTab({ id }: { id: number }) {
  const [schedule, setSchedule] = useState<ScheduleData[]>([]);
  const promiseFunc = async () => {
    const scheduleData = await APIClient.coordinatingClassSchedulesGet(id);
    setSchedule(scheduleData.schedule);
  };

  return (
    <Panel>
      <RequestLoader promiseFunc={promiseFunc}>
        <ScheduleDisplay
          id={id}
          schedule={schedule}
          setSchedule={setSchedule}
        />
      </RequestLoader>
    </Panel>
  );
}

function ScheduleDisplay({
  id,
  schedule,
  setSchedule,
}: {
  id: number;
  schedule: ScheduleData[];
  setSchedule: Dispatch<SetStateAction<ScheduleData[]>>;
}) {
  const router = useRouter();

  const table = useMantineReactTable({
    columns: CoordinatingClassScheduleTableColumns,
    data: schedule,
    initialState: {
      density: "sm" as MRT_DensityState,
      pagination: { pageSize: 20, pageIndex: 0 },
    },
    enableRowActions: true,
    positionActionsColumn: "last",
    renderRowActions: ({ row }) => (
      <ChangeSchedule
        id={id}
        row={row}
        schedule={schedule}
        setSchedule={setSchedule}
      />
    ),
    mantineTableBodyRowProps: ({ row }) => ({
      onClick: () =>
        router.push(
          `${Routes.classAdministrations}/${id}/sessions/${row.original.class_group_session_id}`,
        ),
      style: {
        cursor: "pointer",
      },
    }),
  });

  return <MantineReactTable table={table} />;
}

function ChangeSchedule({
  id,
  row,
  schedule,
  setSchedule,
}: {
  id: number;
  row: MRT_Row<ScheduleData>;
  schedule: ScheduleData[];
  setSchedule: Dispatch<SetStateAction<ScheduleData[]>>;
}) {
  const [opened, { open, close }] = useDisclosure(false);
  const [loading, setLoading] = useState(false);
  const form: UseFormReturnType<{ start_time: Date; end_time: Date }> = useForm(
    {
      initialValues: {
        start_time: new Date(row.original.start_time),
        end_time: new Date(row.original.end_time),
      },
    },
  );

  const onSubmit = form.onSubmit(async (values) => {
    setLoading(true);
    values.start_time.setMilliseconds(0);
    values.end_time.setMilliseconds(0);

    try {
      const response = await APIClient.coordinatingClassSchedulePut(
        id,
        row.original.class_group_session_id,
        values.start_time,
        values.end_time,
      );
      schedule[row.index].start_time = response.start_time;
      schedule[row.index].end_time = response.end_time;
      form.setInitialValues({
        start_time: response.start_time,
        end_time: response.end_time,
      });
      form.resetDirty();
      setSchedule([...schedule]);
      close();
    } catch (e) {
      notifications.show({
        title: "Schedule Update Failed",
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
        title={`Change Schedule: ${row.original.class_group_name}, ${row.original.class_type}`}
      >
        <form onSubmit={onSubmit}>
          <DateTimePicker
            dropdownType="modal"
            label="Start Time"
            valueFormat="DD MMM YYYY hh:mm A"
            {...form.getInputProps("start_time")}
          />
          <Space h="md" />
          <DateTimePicker
            dropdownType="modal"
            label="End Time"
            valueFormat="DD MMM YYYY hh:mm A"
            {...form.getInputProps("end_time")}
          />

          <Space h="md" />
          <Group justify="center">
            <Button
              disabled={loading || !form.isDirty()}
              onClick={form.reset}
              color="red"
              variant="light"
            >
              Reset
            </Button>
            <Button
              disabled={loading || !form.isDirty()}
              type="submit"
              color="green"
            >
              Change
            </Button>
          </Group>
        </form>
      </Modal>
      <Tooltip label="Change session schedule">
        <ActionIcon onClick={open}>
          <IconAdjustments />
        </ActionIcon>
      </Tooltip>
    </>
  );
}
