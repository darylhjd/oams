"use client";

import styles from "@/styles/SessionAttendanceTaking.module.css";

import { useEffect, useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { Box, Button, Container, Group, Title, Tooltip } from "@mantine/core";
import {
  AttendanceTakingGetResponse,
  UpcomingClassGroupSession,
} from "@/api/attendance_taking";
import { SessionEnrollment } from "@/api/session_enrollment";
import { MantineReactTable, MRT_PaginationState } from "mantine-react-table";
import {
  AttendanceTakingDataTableColumns,
  DEFAULT_PAGE_SIZE,
} from "@/components/tabling";
import { IconHelp } from "@tabler/icons-react";

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
        <SessionInfo session={attendance.upcoming_class_group_session} />
        <AttendanceTaker
          id={params.id}
          enrollments={attendance.enrollment_data}
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
    <Group className={styles.title} gap="xs" justify="center">
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
  return null;
}

function AttendanceTaker({
  id,
  enrollments,
}: {
  id: number;
  enrollments: SessionEnrollment[];
}) {
  const [paginationState, setPaginationState] = useState<MRT_PaginationState>({
    pageIndex: 0,
    pageSize: DEFAULT_PAGE_SIZE,
  });
  const [data, setData] = useState(enrollments);

  useEffect(() => {
    const interval = setInterval(async () => {
      const response = await APIClient.attendanceTakingGet(id);
      setData(response.enrollment_data);
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
        renderRowActions={({ row }) =>
          row.original.attended ? null : (
            <TakeAttendanceButton
              onClick={async () => {
                const response = await APIClient.sessionEnrollmentPatch(
                  row.original.id,
                  true,
                );
                data[row.index] = response.session_enrollment;
                setData([...data]);
              }}
            />
          )
        }
      />
    </Box>
  );
}

function TakeAttendanceButton({ onClick }: { onClick: () => void }) {
  return (
    <Button variant="outline" onClick={onClick}>
      Take attendance
    </Button>
  );
}
