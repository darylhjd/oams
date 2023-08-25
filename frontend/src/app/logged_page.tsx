'use client'

import { Center, Container, Divider, Flex, Space, Stack, Text, createStyles } from '@mantine/core'
import { Calendar, dayjsLocalizer } from 'react-big-calendar'
import dayjs from 'dayjs'

import 'react-big-calendar/lib/addons/dragAndDrop/styles.css'
import 'react-big-calendar/lib/css/react-big-calendar.css'

const localizer = dayjsLocalizer(dayjs)

const useStyles = createStyles((theme) => ({
  calendar: {
    flexGrow: 1,
    height: '35em',
  },  

  previewer: {
    width: '22em',
  },

  previewStack: {
    height: '100%',
    borderStyle: 'solid',
    borderRadius: '0.5em',
  },

  previewHeader: {
    padding: '0.1em 0',
  },

  mainPreview: {
    flexGrow: 1,
    padding: '0.5em 0',
  },

  previewFooter: {
    padding: '0.1em 0',
  }
}))

export default function LoggedHomePage() {
  const { classes } = useStyles()

  return (
    <Container fluid={true}>
      <Flex justify='space-between'>
        <Calendar className={classes.calendar}
          localizer={localizer}
          defaultDate={new Date()}
        />
        <Space w='md' />
        <Previewer />
      </Flex>
    </Container>
  )
}

function Previewer() {
  const { classes } = useStyles()

  return (
    <div className={classes.previewer}>
      <Stack className={classes.previewStack} justify='space-between' spacing={0}>
        <Center className={classes.previewHeader}>Selected Day's Events</Center>
        <Divider />
        <Center className={classes.mainPreview}>Event Previews</Center>
        <Divider />
        <Center className={classes.previewFooter}>Legend</Center>
      </Stack>
    </div>
  )
}
