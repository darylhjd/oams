'use client'

import { APIClient } from "@/api/client";
import { redirectIfAlreadyLoggedIn } from "@/routes/checks";
import { Routes } from "@/routes/routes";
import { Button, Center, Container, Image, Stack, createStyles } from "@mantine/core";
import { useSearchParams } from "next/navigation";

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
  redirectIfAlreadyLoggedIn(Routes.home)

  const { classes } = useStyles()

  const redirectUrl = useSearchParams().get("redirect_url") ?? ""
  return (
    <Center>
      <Container className={classes.container}>
        <Button className={classes.image} variant='light' component='a' onClick={() => loginAction(redirectUrl)}>
          <Stack>
            <Image src='microsoft_logo.png' fit='contain' alt='Microsoft Logo'/>
            <Center>Login with Microsoft</Center>
          </Stack>
        </Button>
      </Container>
    </Center>
  )
}

async function loginAction(redirectUrl: string) {
  const uri = await APIClient.getLoginUrl(redirectUrl)
  location.replace(uri)
}
