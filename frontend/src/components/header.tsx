'use client'

import { Button, Center, Container, Flex, Image, Menu, Space, Text, createStyles } from "@mantine/core";
import { IconLogin, IconMenu2 } from "@tabler/icons-react";
import { Desktop, Mobile } from "./responsive";

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
      <Mobile>
        <Flex align='center' justify='space-between'>
          <Logo />
          <MobileDropDownMenu />
        </Flex>
      </Mobile>   

      <Desktop>
        <Flex align='center' justify='space-between'>
          <Logo />
          <Flex align='center'>
            <AboutButton />
            <Space w='md' />
            <LoginButton />
          </Flex>
        </Flex>
      </Desktop>
    </nav>
  )
}

function Logo() {
  const { classes } = useStyles()

  return (
    <Button 
      className={classes.logo} 
      variant='subtle' 
      component='a' href={`${process.env.WEB_SERVER_HOST}:${process.env.WEB_SERVER_PORT}`}>
      <Image src='logo.png' alt='OAMS Logo' fit='contain' />
    </Button>
  )
}

function AboutButton() {
  return (
    <>
      <Mobile>
        <Text c='cyan'>About</Text>
      </Mobile>

      <Desktop>
        <Button variant='subtle' color='cyan' component='a' href='/about'>
          About
        </Button>  
      </Desktop>
    </>
  )
}

function LoginButton() {
  return (
    <>
      <Mobile>
        <Text c='blue'>Login</Text>
      </Mobile>

      <Desktop>
        <Button component='a' href='/login'>
          Login
        </Button>
      </Desktop>
    </>
  )
}

function MobileDropDownMenu() {
  return (
    <Menu position='bottom-end' width={150}>
      <Menu.Target>
        <Button leftIcon={<IconMenu2 />} variant='subtle'>
          Menu
        </Button>
      </Menu.Target>

      <Menu.Dropdown>
        <Menu.Item component='a' href='/about'>
          <AboutButton />
        </Menu.Item>
        <Menu.Item icon={<IconLogin stroke={1}/>} component='a' href='/login'>
          <LoginButton />
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  )
}
