"use client";

import {
  Center,
  Container,
  Divider,
  Flex,
  Space,
  Stack,
  createStyles,
} from "@mantine/core";
import { Calendar, dayjsLocalizer } from "react-big-calendar";
import dayjs from "dayjs";

import "react-big-calendar/lib/addons/dragAndDrop/styles.css";
import "react-big-calendar/lib/css/react-big-calendar.css";
import { sessionStore } from "@/states/session";
import { ClassType } from "@/api/models";
import { useMediaQuery } from "@mantine/hooks";
import { MOBILE_MIN_WIDTH } from "@/components/responsive";

const localizer = dayjsLocalizer(dayjs);

const useStyles = createStyles((theme) => ({
  calendar: {
    flexGrow: 1,
    height: "35em",
  },

  tutorialEvent: {
    backgroundColor: "green",
  },

  lectureEvent: {
    backgroundColor: "blue",
  },

  labEvent: {
    backgroundColor: "chocolate",
  },

  previewer: {
    width: "22em",

    [theme.fn.smallerThan("md")]: {
      width: "100%",
    },
  },

  previewStack: {
    height: "100%",
    borderStyle: "solid",
    borderRadius: "0.5em",
  },

  previewHeader: {
    padding: "0.1em 0",
  },

  mainPreview: {
    flexGrow: 1,
    padding: "0.5em 0",
  },

  previewFooter: {
    padding: "0.1em 0",
  },
}));

export default function LoggedHomePage() {
  const layout = useMediaQuery(MOBILE_MIN_WIDTH) ? (
    <Stack>
      <CalendarPreview />
      <Previewer />
    </Stack>
  ) : (
    <Flex justify="space-between">
      <CalendarPreview />
      <Space w="md" />
      <Previewer />
    </Flex>
  );

  return <Container fluid={true}>{layout}</Container>;
}

function CalendarPreview() {
  const { classes } = useStyles();
  const session = sessionStore();

  const events = session.data!.upcoming_class_group_sessions.map((session) => ({
    title: `${session.code} ${session.class_type}`,
    start: new Date(session.start_time), // Since typescript parses date from API as a string.
    end: new Date(session.end_time),
    allDay: false,
    resource: session,
  }));

  return (
    <Calendar
      className={classes.calendar}
      localizer={localizer}
      defaultDate={new Date()}
      events={events}
      eventPropGetter={(event, start, end, isSelected) => {
        var className = null;
        switch (event.resource.class_type) {
          case ClassType.Lab:
            className = classes.labEvent;
            break;
          case ClassType.Lecture:
            className = classes.lectureEvent;
            break;
          case ClassType.Tutorial:
            className = classes.tutorialEvent;
            break;
        }
        return {
          className: className,
        };
      }}
      formats={{
        agendaHeaderFormat: ({ start, end }, _, localizer) =>
          localizer!.format(start, "YYYY/MM/DD") +
          " - " +
          localizer!.format(end, "YYYY/MM/DD"),
      }}
    />
  );
}

function Previewer() {
  const { classes } = useStyles();

  return (
    <div className={classes.previewer}>
      <Stack
        className={classes.previewStack}
        justify="space-between"
        spacing={0}
      >
        <Center className={classes.previewHeader}>
          Selected Day&apos;s Events
        </Center>
        <Divider />
        <Center className={classes.mainPreview}>Event Previews</Center>
        <Divider />
        <Center className={classes.previewFooter}>Legend</Center>
      </Stack>
    </div>
  );
}
