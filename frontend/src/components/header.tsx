'use client'

import { Button, Center, Container, Flex, Image, Menu, Space, Text, createStyles } from "@mantine/core";
import { IconLogin, IconMenu2, IconUserCircle } from "@tabler/icons-react";
import { Desktop, Mobile } from "./responsive";
import { useRouter } from "next/navigation";
import { Routes } from "@/routes/routes";
import { sessionStore } from "@/states/session";

const useStyles = createStyles((theme) => ({
  container: {
    position: 'sticky',
    top: 0,
    backgroundColor: 'white',
    padding: '0.29em 0em',
    borderBottom: '1px solid black',

    [theme.fn.smallerThan('md')]: {
      padding: '0.6em 0em',
    }
  },

  centeredContainer: {
    padding: '0em 1em',
    width: '100%',
    maxWidth: '80em',
  },

  logo: {
    width: '9em',
    height: 'auto',
    padding: '0.5em 0em',
    marginRight: '0.7em',

    [theme.fn.smallerThan('md')]: {
      width: '7em',
      padding: '0.25em 0em',
      marginRight: '0',
    }
  }
}))

// Header stores the navigation bar and shows a horizontal divider bottom border.
export default function Header() {
  const { classes } = useStyles()

  return (
    <Container className={classes.container} fluid={true}>
      <Center>
        <NavBar />
      </Center>
    </Container>
  )
}

// This shows the navigation bar.
function NavBar() {
  const { classes } = useStyles()

  return (
    <nav className={classes.centeredContainer}>
      <Flex align='center' justify='space-between'>
        <Logo />
        <Options />
      </Flex>
    </nav>
  )
}

function Logo() {
  const { classes } = useStyles()
  const router = useRouter()

  return (
    <Button
      className={classes.logo}
      variant='subtle'
      onClick={() => router.push("/")}>
      <Image src='logo.png' alt='OAMS Logo' fit='contain' />
    </Button>
  )
}

function Options() {
  const router = useRouter()
  const session = sessionStore()

  return (
    <>
      <Mobile>
        <Menu position='bottom-end' width={150}>
          <Menu.Target>
            <Button leftIcon={<IconMenu2 />} variant='subtle'>
              Menu
            </Button>
          </Menu.Target>

          <Menu.Dropdown>
            <Menu.Item onClick={() => router.push(Routes.about)}>
              <AboutButton />
            </Menu.Item>
            {session.user == null ?
              <Menu.Item icon={<IconLogin stroke={1} />} onClick={() => router.push(Routes.login)}>
                <LoginButton />
              </Menu.Item> :
              <Menu.Item icon={<IconUserCircle stroke={1} />} onClick={() => router.push(Routes.profile)}>
                <ProfileButton />
              </Menu.Item>
            }
          </Menu.Dropdown>
        </Menu>
      </Mobile>

      <Desktop>
        <Flex align='center'>
          <AboutButton />
          <Space w='md' />
          {session.user == null ? <LoginButton /> : <ProfileButton />}
        </Flex>
      </Desktop>
    </>
  )
}

function AboutButton() {
  const router = useRouter()

  return (
    <>
      <Mobile>
        <Text c='cyan'>About</Text>
      </Mobile>

      <Desktop>
        <Button variant='subtle' color='cyan' onClick={() => router.push(Routes.about)}>
          About
        </Button>
      </Desktop>
    </>
  )
}

function LoginButton() {
  const router = useRouter()

  return (
    <>
      <Mobile>
        <Text c='blue'>Login</Text>
      </Mobile>

      <Desktop>
        <Button onClick={() => router.push(Routes.login)}>
          Login
        </Button>
      </Desktop>
    </>
  )
}

function ProfileButton() {
  const router = useRouter()

  return (
    <>
      <Mobile>
        <Text>Your Profile</Text>
      </Mobile>

      <Desktop>
        <Button onClick={() => router.push(Routes.profile)} variant='outline'>
          Your Profile
        </Button>
      </Desktop>
    </>
  )
}
