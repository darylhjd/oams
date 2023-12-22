"use client";

import styles from "@/styles/SessionAttendanceTaking.module.css";

import { useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { Box, Button, Center, Container, Text, Title } from "@mantine/core";
import {
  AttendanceTakingGetResponse,
  UpcomingClassGroupSession,
} from "@/api/attendance_taking";
import { SessionEnrollment } from "@/api/session_enrollment";
import {
  MantineReactTable,
  MRT_PaginationState,
  MRT_Row,
} from "mantine-react-table";
import {
  AttendanceTakingDataTableColumns,
  DEFAULT_PAGE_SIZE,
} from "@/components/tabling";

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
        <AttendanceTaker enrollments={attendance.enrollment_data} />
      </Container>
    </RequestLoader>
  );
}

function PageTitle() {
  return (
    <Center className={styles.title}>
      <Title visibleFrom="md" order={3}>
        Attendance Taking
      </Title>
      <Title hiddenFrom="md" order={4}>
        Attendance Taking
      </Title>
    </Center>
  );
}

function SessionInfo({ session }: { session: UpcomingClassGroupSession }) {
  return <Text ta="center">{session.venue}</Text>;
}

function AttendanceTaker({
  enrollments,
}: {
  enrollments: SessionEnrollment[];
}) {
  const [paginationState, setPaginationState] = useState<MRT_PaginationState>({
    pageIndex: 0,
    pageSize: DEFAULT_PAGE_SIZE,
  });
  const [data, setData] = useState(enrollments);

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
              onClick={() => {
                data[row.index].attended = true;
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
