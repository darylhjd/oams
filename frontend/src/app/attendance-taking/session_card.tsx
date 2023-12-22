import { UpcomingClassGroupSession } from "@/api/attendance_taking";
import { useRouter } from "next/navigation";
import { Paper, Space, Text } from "@mantine/core";
import styles from "@/styles/AttendanceTakingPage.module.css";
import { Routes } from "@/routing/routes";
import { UserRole } from "@/api/user";

export function SessionCard({
  session,
}: {
  session: UpcomingClassGroupSession;
}) {
  const router = useRouter();

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

  const isOngoing = new Date() >= startDatetime;

  return (
    <Paper
      withBorder
      p="xs"
      className={styles.sessionCard}
      component="button"
      onClick={() => router.push(Routes.attendanceTakingSession + session.id)}
    >
      <Text ta="left">
        {session.code}{" "}
        <Text span size="sm" c="dimmed">
          {session.year}/{session.semester}
        </Text>
      </Text>
      <Text ta="left" size="xs">
        Group Name: {session.name}
      </Text>
      <Text ta="left" size="xs">
        Class Type: {session.class_type}
      </Text>
      <Text ta="left" size="xs">
        Venue: {session.venue}
      </Text>
      <Space h="xs" />
      <Text ta="left" size="xs">
        {date}
        <br />
        {startTime} - {endTime}
      </Text>
      <Text ta="left" size="xs" c={isOngoing ? "green" : "orange"}>
        {isOngoing ? "ONGOING" : "STARTING"}
      </Text>
      <Space h="xs" />
      <Text ta="left" size="xs" c="dimmed">
        {session.managing_role ? session.managing_role : UserRole.SystemAdmin}
      </Text>
    </Paper>
  );
}
