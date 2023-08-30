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
import { MOBILE_MIN_WIDTH } from "./responsive";
import { useRouter } from "next/navigation";
import { Routes } from "@/routes/routes";
import { sessionStore } from "@/states/session";
import { APIClient } from "@/api/client";
import { UserRole } from "@/api/models";
import { useMediaQuery } from "@mantine/hooks";

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

const mobileProp = {
  mobile: false,
};

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
  const session = sessionStore();
  const { classes } = useStyles();

  if (useMediaQuery(MOBILE_MIN_WIDTH)) {
    return (
      <Menu position="bottom-end" width={150}>
        <Menu.Target>
          <Button leftIcon={<IconMenu2 />} variant="subtle">
            Menu
          </Button>
        </Menu.Target>

        <Menu.Dropdown>
          <AboutButton mobile />
          {session.data != null &&
          session.data.session_user.role == UserRole.SystemAdmin ? (
            <>
              <Menu.Label>Admin Controls</Menu.Label>
              <AdminPanelButton mobile />
              <BatchProcessingButton mobile />
            </>
          ) : null}
          {session.data == null ? (
            <LoginButton mobile />
          ) : (
            <>
              <Menu.Label>Account</Menu.Label>
              <ProfileButton mobile />
              <LogoutButton />
            </>
          )}
        </Menu.Dropdown>
      </Menu>
    );
  }

  return (
    <Flex
      className={classes.desktopOptions}
      align="center"
      justify="space-between"
    >
      <div>
        {session.data != null &&
        session.data.session_user.role == UserRole.SystemAdmin ? (
          <>
            <AdminPanelButton />
            <BatchProcessingButton />
          </>
        ) : null}
        <AboutButton />
      </div>
      <div>{session.data == null ? <LoginButton /> : <ProfileButton />}</div>
    </Flex>
  );
}

const AdminPanelButton = ({ mobile }: { mobile: boolean }) => {
  const router = useRouter();

  if (mobile) {
    return (
      <Menu.Item onClick={() => router.push(Routes.adminPanel)}>
        <Text c="yellow">Admin Panel</Text>
      </Menu.Item>
    );
  }

  return (
    <Button
      variant="subtle"
      color="yellow"
      onClick={() => router.push(Routes.adminPanel)}
    >
      Admin Panel
    </Button>
  );
};

AdminPanelButton.defaultProps = mobileProp;

const BatchProcessingButton = ({ mobile }: { mobile: boolean }) => {
  const router = useRouter();

  if (mobile) {
    return (
      <Menu.Item onClick={() => router.push(Routes.batchProcessing)}>
        <Text c="yellow">Batch Processing</Text>
      </Menu.Item>
    );
  }

  return (
    <Button
      variant="subtle"
      color="yellow"
      onClick={() => router.push(Routes.batchProcessing)}
    >
      Batch Processing
    </Button>
  );
};

BatchProcessingButton.defaultProps = mobileProp;

const AboutButton = ({ mobile }: { mobile: boolean }) => {
  const router = useRouter();

  if (mobile) {
    return (
      <Menu.Item onClick={() => router.push(Routes.about)}>
        <Text c="cyan">About</Text>
      </Menu.Item>
    );
  }

  return (
    <Button
      variant="subtle"
      color="cyan"
      onClick={() => router.push(Routes.about)}
    >
      About
    </Button>
  );
};

AboutButton.defaultProps = mobileProp;

const LoginButton = ({ mobile }: { mobile: boolean }) => {
  const router = useRouter();

  if (mobile) {
    return (
      <Menu.Item
        icon={<IconLogin color="darkblue" />}
        onClick={() => router.push(Routes.login)}
      >
        <Text c="blue">Login</Text>
      </Menu.Item>
    );
  }

  return (
    <Button color="blue" onClick={() => router.push(Routes.login)}>
      Login
    </Button>
  );
};

LoginButton.defaultProps = mobileProp;

const ProfileButton = ({ mobile }: { mobile: boolean }) => {
  const router = useRouter();

  if (mobile) {
    return (
      <Menu.Item
        icon={<IconUserCircle stroke={1} />}
        onClick={() => router.push(Routes.profile)}
      >
        <Text>Your Profile</Text>
      </Menu.Item>
    );
  }

  return (
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
        <LogoutButton />
      </Menu.Dropdown>
    </Menu>
  );
};

ProfileButton.defaultProps = mobileProp;

const LogoutButton = () => {
  const session = sessionStore();

  return (
    <Menu.Item
      icon={<IconLogout color="red" />}
      onClick={async () => {
        session.invalidate();
        await APIClient.logout();
        location.href = "/";
      }}
    >
      <Text c="red">Logout</Text>
    </Menu.Item>
  );
};
