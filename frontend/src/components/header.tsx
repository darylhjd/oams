"use client";

import {
  ActionIcon,
  Button,
  Center,
  Container,
  Flex,
  Image,
  Menu,
  Text,
  createStyles,
} from "@mantine/core";
import {
  IconLogin,
  IconLogout,
  IconMenu2,
  IconUserCircle,
} from "@tabler/icons-react";
import { Desktop, Mobile } from "./responsive";
import { useRouter } from "next/navigation";
import { Routes } from "@/routes/routes";
import { sessionStore } from "@/states/session";
import { APIClient } from "@/api/client";
import { UserRole } from "@/api/models";

const useStyles = createStyles((theme) => ({
  container: {
    position: "sticky",
    top: 0,
    zIndex: 1,
    backgroundColor: "white",
    padding: "0.29em 0em",
    borderBottom: "1px solid black",
    boxShadow: "0em 0.1em 1em -0.1em rgba(0,0,0,0.4)",

    [theme.fn.smallerThan("md")]: {
      padding: "0.6em 0em",
    },
  },

  centeredContainer: {
    padding: "0em 1em",
    width: "100%",
    maxWidth: "80em",
  },

  logo: {
    width: "9em",
    height: "auto",
    padding: "0.5em 0em",
    marginRight: "0.7em",

    [theme.fn.smallerThan("md")]: {
      width: "7em",
      padding: "0.25em 0em",
      marginRight: "0",
    },
  },

  desktopOptions: {
    width: "100%",
  },
}));

// Header stores the navigation bar and shows a horizontal divider bottom border.
export default function Header() {
  const { classes } = useStyles();

  return (
    <Container className={classes.container} fluid={true}>
      <Center>
        <NavBar />
      </Center>
    </Container>
  );
}

// This shows the navigation bar.
function NavBar() {
  const { classes } = useStyles();

  return (
    <nav className={classes.centeredContainer}>
      <Flex align="center" justify="space-between">
        <Logo />
        <Options />
      </Flex>
    </nav>
  );
}

function Logo() {
  const { classes } = useStyles();
  const router = useRouter();

  return (
    <Button
      className={classes.logo}
      variant="subtle"
      onClick={() => router.push("/")}
    >
      <Image src="logo.png" alt="OAMS Logo" fit="contain" />
    </Button>
  );
}

function Options() {
  const router = useRouter();
  const session = sessionStore();
  const { classes } = useStyles();

  return (
    <>
      <Mobile>
        <Menu position="bottom-end" width={150}>
          <Menu.Target>
            <Button leftIcon={<IconMenu2 />} variant="subtle">
              Menu
            </Button>
          </Menu.Target>

          <Menu.Dropdown>
            <Menu.Item onClick={() => router.push(Routes.about)}>
              <AboutButton />
            </Menu.Item>
            {session.userMe != null &&
            session.userMe.session_user.role == UserRole.SystemAdmin ? (
              <>
                <Menu.Label>Admin Controls</Menu.Label>
                <Menu.Item onClick={() => router.push(Routes.adminPanel)}>
                  <AdminPanelButton />
                </Menu.Item>
                <Menu.Item onClick={() => router.push(Routes.batchProcessing)}>
                  <BatchProcessingButton />
                </Menu.Item>
              </>
            ) : null}
            {session.userMe == null ? (
              <Menu.Item
                icon={<IconLogin stroke={1} />}
                onClick={() => router.push(Routes.login)}
              >
                <LoginButton />
              </Menu.Item>
            ) : (
              <>
                <Menu.Label>Account</Menu.Label>
                <Menu.Item
                  icon={<IconUserCircle stroke={1} />}
                  onClick={() => router.push(Routes.profile)}
                >
                  <ProfileButton />
                </Menu.Item>
                <Menu.Item
                  icon={<IconLogout stroke={1} />}
                  onClick={async () => {
                    session.invalidate();
                    await APIClient.logout();
                    location.href = "/";
                  }}
                >
                  <LogoutButton />
                </Menu.Item>
              </>
            )}
          </Menu.Dropdown>
        </Menu>
      </Mobile>

      <Desktop>
        <Flex
          className={classes.desktopOptions}
          align="center"
          justify="space-between"
        >
          <div>
            {session.userMe != null &&
            session.userMe.session_user.role == UserRole.SystemAdmin ? (
              <>
                <AdminPanelButton />
                <BatchProcessingButton />
              </>
            ) : null}
            <AboutButton />
          </div>
          <div>
            {session.userMe == null ? <LoginButton /> : <ProfileButton />}
          </div>
        </Flex>
      </Desktop>
    </>
  );
}

function AdminPanelButton() {
  const router = useRouter();

  return (
    <>
      <Mobile>
        <Text c="yellow">Admin Panel</Text>
      </Mobile>

      <Desktop>
        <Button
          variant="subtle"
          color="yellow"
          onClick={() => router.push(Routes.adminPanel)}
        >
          Admin Panel
        </Button>
      </Desktop>
    </>
  );
}

function BatchProcessingButton() {
  const router = useRouter();

  return (
    <>
      <Mobile>
        <Text c="yellow">Batch Processing</Text>
      </Mobile>

      <Desktop>
        <Button
          variant="subtle"
          color="yellow"
          onClick={() => router.push(Routes.batchProcessing)}
        >
          Batch Processing
        </Button>
      </Desktop>
    </>
  );
}

function AboutButton() {
  const router = useRouter();

  return (
    <>
      <Mobile>
        <Text c="cyan">About</Text>
      </Mobile>

      <Desktop>
        <Button
          variant="subtle"
          color="cyan"
          onClick={() => router.push(Routes.about)}
        >
          About
        </Button>
      </Desktop>
    </>
  );
}

function LoginButton() {
  const router = useRouter();

  return (
    <>
      <Mobile>
        <Text c="blue">Login</Text>
      </Mobile>

      <Desktop>
        <Button onClick={() => router.push(Routes.login)}>Login</Button>
      </Desktop>
    </>
  );
}

function ProfileButton() {
  const router = useRouter();
  const session = sessionStore();

  return (
    <>
      <Mobile>
        <Text>Your Profile</Text>
      </Mobile>

      <Desktop>
        <Menu position="bottom-end" width={150}>
          <Menu.Target>
            <ActionIcon size="lg">
              <IconUserCircle size="4em" />
            </ActionIcon>
          </Menu.Target>

          <Menu.Dropdown>
            <Menu.Item
              icon={<IconUserCircle stroke={1} />}
              onClick={() => router.push(Routes.profile)}
            >
              <Text>Your Profile</Text>
            </Menu.Item>
            <Menu.Item
              icon={<IconLogout color="red" />}
              onClick={async () => {
                session.invalidate();
                await APIClient.logout();
                location.href = "/";
              }}
            >
              <LogoutButton />
            </Menu.Item>
          </Menu.Dropdown>
        </Menu>
      </Desktop>
    </>
  );
}

function LogoutButton() {
  return <Text c="red">Logout</Text>;
}
