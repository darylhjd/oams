'use client'

import { Center, createStyles } from "@mantine/core";

const useStyles = createStyles((theme) => ({
  expand: {
    width: '100%',
    height: '100%',
  }
}))

export default function LoginPage() {
  const { classes } = useStyles()

  return (
    <Center>
      <div className={classes.expand}>
        Hi there.
      </div>
    </Center>
  )
}