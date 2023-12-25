"use client";

import styles from "@/styles/SessionAttendanceTaking.module.css";

import { Dispatch, SetStateAction, useEffect, useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import {
  Box,
  Button,
  Container,
  FocusTrap,
  Group,
  Modal,
  Paper,
  PasswordInput,
  Space,
  Text,
  Title,
  Tooltip,
} from "@mantine/core";
import {
  AttendanceEntry,
  AttendanceTakingGetResponse,
  UpcomingClassGroupSession,
} from "@/api/attendance_taking";
import {
  MantineReactTable,
  MRT_PaginationState,
  MRT_Row,
} from "mantine-react-table";
import {
  AttendanceTakingDataTableColumns,
  DEFAULT_PAGE_SIZE,
} from "@/components/tabling";
import { IconHelp, IconX } from "@tabler/icons-react";
import { useDisclosure } from "@mantine/hooks";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { signatureInputForm } from "@/components/signature_form";

const UPDATE_INTERVAL_MS = 5000;

export default function SessionAttendanceTakingPage({
  params,
}: {
  params: Params;
}) {
  const [attendance, setAttendance] = useState<AttendanceTakingGetResponse>(
    {} as AttendanceTakingGetResponse,
  );
  const promiseFunc = async () => {
    const data = await APIClient.attendanceTakingGet(params.id);
    return setAttendance(data);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Container className={styles.container} fluid>
        <PageTitle />
        <Space h="md" />
        <SessionInfo session={attendance.upcoming_class_group_session} />
        <Space h="md" />
        <AttendanceTaker
          id={params.id}
          attendanceEntries={attendance.attendance_entries}
        />
      </Container>
    </RequestLoader>
  );
}

function PageTitle() {
  const label = `The attendance sheet is updated automatically every ${
    UPDATE_INTERVAL_MS / 1000
  } seconds.`;
  return (
    <Group gap="xs" justify="center">
      <Title order={3} ta="center">
        Attendance Taking
      </Title>
      <Tooltip
        label={label}
        events={{
          hover: true,
          focus: false,
          touch: true,
        }}
      >
        <IconHelp size={15} />
      </Tooltip>
    </Group>
  );
}

function SessionInfo({ session }: { session: UpcomingClassGroupSession }) {
  const startDatetime = new Date(session.start_time);
  const endDatetime = new Date(session.end_time);

  const date = startDatetime.toLocaleString(undefined, {
    day: "numeric",
    month: "numeric",
    year: "numeric",
  });
  const startTime = startDatetime.toLocaleString(undefined, {
    hour: "2-digit",
    minute: "2-digit",
  });
  const endTime = endDatetime.toLocaleString(undefined, {
    hour: "2-digit",
    minute: "2-digit",
  });

  return (
    <Paper withBorder p="xs" className={styles.sessionInfo}>
      <Text ta="center" c="green">
        {session.code} {session.year}/{session.semester}
      </Text>
      <Text ta="center" size="sm" c="orange">
        {session.name} {session.class_type} @ {session.venue}
      </Text>
      <Text ta="center" c="yellow">
        {date}, {startTime} - {endTime}
      </Text>
    </Paper>
  );
}

function AttendanceTaker({
  id,
  attendanceEntries,
}: {
  id: number;
  attendanceEntries: AttendanceEntry[];
}) {
  const [paginationState, setPaginationState] = useState<MRT_PaginationState>({
    pageIndex: 0,
    pageSize: DEFAULT_PAGE_SIZE,
  });
  const [data, setData] = useState(attendanceEntries);

  useEffect(() => {
    const interval = setInterval(async () => {
      const response = await APIClient.attendanceTakingGet(id);
      setData(response.attendance_entries);
    }, UPDATE_INTERVAL_MS);

    return () => clearInterval(interval);
  }, []);

  return (
    <Box className={styles.table}>
      <MantineReactTable
        columns={AttendanceTakingDataTableColumns}
        data={data}
        state={{ pagination: paginationState }}
        onPaginationChange={setPaginationState}
        enableRowActions
        positionActionsColumn="last"
        renderRowActions={({ row }) => (
          <SignAttendance id={id} row={row} data={data} setData={setData} />
        )}
      />
    </Box>
  );
}

function SignAttendance({
  id,
  row,
  data,
  setData,
}: {
  id: number;
  row: MRT_Row<AttendanceEntry>;
  data: AttendanceEntry[];
  setData: Dispatch<SetStateAction<AttendanceEntry[]>>;
}) {
  const [loading, setLoading] = useState(false);

  const [opened, { open, close }] = useDisclosure(false);
  const form = useForm(signatureInputForm());

  if (row.original.attended) {
    return null;
  }

  return (
    <>
      <Modal
        opened={opened}
        onClose={() => {
          close();
          form.reset();
        }}
        centered
        title="Sign Attendance"
        overlayProps={{
          backgroundOpacity: 0.55,
          blur: 3,
        }}
      >
        <form
          onSubmit={form.onSubmit(async (values) => {
            setLoading(true);
            try {
              const resp = await APIClient.attendanceTakingPost(
                id,
                row.original.id,
                row.original.user_id,
                true,
                values.signature,
              );
              data[row.index].attended = resp.attended;
              setData([...data]);
              close();
              form.reset();
            } catch (error) {
              notifications.show({
                title: "Wrong signature",
                message: "You entered the wrong signature. Please try again.",
                icon: <IconX />,
                color: "red",
              });
            }
            setLoading(false);
          })}
        >
          <FocusTrap active>
            <PasswordInput
              label="Signature"
              {...form.getInputProps("signature")}
              data-autofocus
            />
          </FocusTrap>
          <Space h="sm" />
          <Group justify="center">
            <Button type="submit" color="green" loading={loading}>
              Confirm Attendance
            </Button>
          </Group>
        </form>
      </Modal>
      <Button onClick={open} variant="outline">
        Sign Attendance
      </Button>
    </>
  );
}
