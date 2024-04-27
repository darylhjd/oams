"use client";

import { Paper, Text, Space, Title } from "@mantine/core";
import { Params } from "./layout";
import { useState } from "react";
import {
  CoordinatingClassScheduleGetResponse,
  ScheduleData,
} from "@/api/coordinating_class";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { AttendanceEntry } from "@/api/upcoming_class_group_session";
import { MantineReactTable, useMantineReactTable } from "mantine-react-table";
import { AttendanceEntriesDataTableColumns } from "@/components/columns";

export default function ScheduleAttendancePage({ params }: { params: Params }) {
  const [attendanceData, setAttendanceData] =
    useState<CoordinatingClassScheduleGetResponse>(
      {} as CoordinatingClassScheduleGetResponse,
    );
  const promiseFunc = async () => {
    const data = await APIClient.coordinatingClassScheduleGet(
      params.id,
      params.sessionId,
    );
    console.log(data);
    setAttendanceData(data);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Title ta="center" order={2} py={10}>
        Session Attendance
      </Title>
      <SessionData session={attendanceData.session} />
      <AttendanceData attendance={attendanceData.attendance_entries} />
    </RequestLoader>
  );
}

function SessionData({ session }: { session: ScheduleData }) {
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
    <Paper withBorder p="xs">
      <Text ta="center">Group Name: {session.class_group_name}</Text>
      <Text ta="center">Class Type: {session.class_type}</Text>
      <Text ta="center">Venue: {session.venue}</Text>
      <Space h="xs" />
      <Text ta="center">
        {date}
        <br />
        {startTime} - {endTime}
      </Text>
    </Paper>
  );
}

function AttendanceData({ attendance }: { attendance: AttendanceEntry[] }) {
  const table = useMantineReactTable({
    columns: AttendanceEntriesDataTableColumns,
    data: attendance,
  });

  return <MantineReactTable table={table} />;
}
