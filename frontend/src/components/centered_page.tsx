'use client'

import { Center, createStyles } from "@mantine/core"

const useStyles = createStyles((theme) => ({
  constrained: {
    width: '100%',
    maxWidth: '80em',
  }
}))

// Helps to constrain the contents into the middle of the screen.
export default function CenteredScreen({ children }: { children: React.ReactNode }) {
  const { classes } = useStyles()
  
  return (
    <Center>
      <div className={classes.constrained}>
        <Center>
          {children}
        </Center>
      </div>
    </Center>
  )
}
