"use client";

import styles from "@/styles/Header.module.css";

import {
  Button,
  Center,
  Container,
  Image,
  Group,
  Burger,
  Drawer,
  Box,
  Stack,
  Menu,
  MenuTarget,
  MenuDropdown,
  MenuItem,
  NavLink,
  Divider,
} from "@mantine/core";
import { usePathname, useRouter } from "next/navigation";
import { useSessionUserStore } from "@/stores/session";
import { useDisclosure } from "@mantine/hooks";
import {
  IconChevronDown,
  IconFileDescription,
  IconLayoutDashboard,
} from "@tabler/icons-react";
import { Routes } from "@/routing/routes";
import { APIClient } from "@/api/client";
import { UserRole } from "@/api/user";

export default function Header() {
  return (
    <Container className={styles.header} fluid>
      <NavBarItems />
    </Container>
  );
}

function NavBarItems() {
  return (
    <Center>
      <Group className={styles.items} justify="space-between">
        <Logo />
        <Items />
      </Group>
    </Center>
  );
}

function Logo() {
  const router = useRouter();

  return (
    <Button
      className={styles.logo}
      variant="subtle"
      onClick={() => router.push(Routes.index)}
    >
      <div>
        <Image src="/logo.png" alt="OAMS" fit="contain" />
      </div>
    </Button>
  );
}

function Items() {
  const session = useSessionUserStore();
  const [opened, { close, toggle }] = useDisclosure();

  const items = session.data ? (
    <LoggedItems close={close} />
  ) : (
    <GuestItems close={close} />
  );

  return (
    <>
      <Box visibleFrom="md">{items}</Box>

      <Burger hiddenFrom="md" opened={opened} onClick={toggle} />
      <Drawer opened={opened} onClose={close}>
        {items}
      </Drawer>
    </>
  );
}

function GuestItems({ close }: { close: () => void }) {
  return (
    <>
      <Group visibleFrom="md">
        <AboutButton close={close} />
        <LoginButton />
      </Group>

      <Stack hiddenFrom="md">
        <AboutButton close={close} />
        <LoginButton />
      </Stack>
    </>
  );
}

function LoggedItems({ close }: { close: () => void }) {
  const session = useSessionUserStore();

  return (
    <>
      <Group visibleFrom="md">
        {session.data?.session_user.role == UserRole.SystemAdmin ? (
          <SystemAdminMenu close={close} />
        ) : null}
        <AboutButton close={close} />
        <ProfileButton close={close} />
        <LogoutButton />
      </Group>

      <Stack hiddenFrom="md">
        {session.data?.session_user.role == UserRole.SystemAdmin ? (
          <SystemAdminMenu close={close} />
        ) : null}
        <Divider />
        <AboutButton close={close} />
        <ProfileButton close={close} />
        <LogoutButton />
      </Stack>
    </>
  );
}

function AboutButton({ close }: { close: () => void }) {
  const router = useRouter();

  return (
    <>
      <Button
        visibleFrom="md"
        variant="subtle"
        onClick={() => router.push(Routes.about)}
      >
        About
      </Button>

      <NavLink
        hiddenFrom="md"
        label="About"
        active
        variant="subtle"
        onClick={() => {
          close();
          router.push(Routes.about);
        }}
      />
    </>
  );
}

function LoginButton() {
  const router = useRouter();
  const pathname = usePathname();

  return (
    <Button
      onClick={async () => {
        const redirectLink = `${process.env.WEB_SERVER}${pathname}`;
        const loginLink = await APIClient.login(redirectLink);
        router.push(loginLink);
      }}
    >
      Login
    </Button>
  );
}

function ProfileButton({ close }: { close: () => void }) {
  const router = useRouter();

  return (
    <Button
      color="blue"
      variant="filled"
      onClick={() => {
        close();
        router.push(Routes.profile);
      }}
    >
      Profile
    </Button>
  );
}

function LogoutButton() {
  return (
    <Button
      color="red"
      onClick={async () => {
        await APIClient.logout();
        location.href = Routes.index;
      }}
    >
      Logout
    </Button>
  );
}

function SystemAdminMenu({ close }: { close: () => void }) {
  const router = useRouter();

  return (
    <>
      <Box visibleFrom="md">
        <Menu width={200}>
          <MenuTarget>
            <Button
              color="orange"
              variant="subtle"
              rightSection={<IconChevronDown />}
            >
              System Admin Menu
            </Button>
          </MenuTarget>
          <MenuDropdown>
            <MenuItem
              leftSection={<IconLayoutDashboard size={16} />}
              onClick={() => router.push(Routes.adminPanel)}
            >
              Admin Panel
            </MenuItem>
            <MenuItem
              leftSection={<IconFileDescription size={16} />}
              onClick={() => router.push(Routes.batchProcessing)}
            >
              Batch Processing
            </MenuItem>
          </MenuDropdown>
        </Menu>
      </Box>

      <NavLink
        color="orange"
        hiddenFrom="md"
        label="System Admin Menu"
        active
        variant="subtle"
      >
        <NavLink
          label="Admin Panel"
          leftSection={<IconLayoutDashboard size={16} />}
          onClick={() => {
            close();
            router.push(Routes.adminPanel);
          }}
        />
        <NavLink
          label="Batch Processing"
          leftSection={<IconFileDescription size={16} />}
          onClick={() => {
            close();
            router.push(Routes.batchProcessing);
          }}
        />
      </NavLink>
    </>
  );
}
