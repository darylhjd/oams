'use client'

import { Button, Center, Container, Image, Stack, createStyles } from "@mantine/core";

const useStyles = createStyles((theme) => ({
  container: {
    paddingTop: '10em',
  },

  image: {
    height: 'auto',
    width: '13em',
    padding: '1em 1em',
  }
}))

export default function LoginPage() {
  const { classes } = useStyles()

  return (
    <Container className={classes.container}>
      <Button className={classes.image} variant='light'>
        <Stack>
          <Image src='microsoft_logo.png' fit='contain' alt='Microsoft Logo'/>
          <Center>Login with Microsoft</Center>
        </Stack>
      </Button>
    </Container>
  )
}
