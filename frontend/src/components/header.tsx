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
  Flex,
  Space,
  DrawerContent,
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
      <Flex className={styles.items} justify="space-between" align="center">
        <Logo />
        <Items />
      </Flex>
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

  const items = session.data ? <LoggedItems /> : <GuestItems />;

  return (
    <>
      <Box visibleFrom="md">{items}</Box>

      <Burger hiddenFrom="md" opened={opened} onClick={toggle} />
      <Drawer
        size="sm"
        padding="lg"
        opened={opened}
        onClose={close}
        onClick={close}
        classNames={{
          inner: styles.systemAdminMenuDropdown
        }}
      >
        {items}
      </Drawer>
    </>
  );
}

function GuestItems() {
  return (
    <>
      <Group visibleFrom="md">
        <AboutButton />
        <LoginButton />
      </Group>

      <Stack hiddenFrom="md" gap={0}>
        <AboutButton />
        <Divider my="sm" />

        <Space h="md" />
        <Group>
          <LoginButton />
        </Group>
      </Stack>
    </>
  );
}

function LoggedItems() {
  const session = useSessionUserStore();

  return (
    <>
      <Group visibleFrom="md">
        {session.data?.session_user.role == UserRole.SystemAdmin ? (
          <SystemAdminMenu />
        ) : null}
        <AboutButton />
        <ProfileButton />
        <LogoutButton />
      </Group>

      <Stack hiddenFrom="md" gap={0}>
        {session.data?.session_user.role == UserRole.SystemAdmin ? (
          <SystemAdminMenu />
        ) : null}
        <AboutButton />
        <Divider my="sm" />

        <Space h="md" />
        <Group>
          <ProfileButton />
          <LogoutButton />
        </Group>
      </Stack>
    </>
  );
}

function AboutButton() {
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
        onClick={() => router.push(Routes.about)}
      />
    </>
  );
}

function LoginButton() {
  const pathname = usePathname();

  return (
    <Button
      onClick={async () => {
        const redirectLink = `${process.env.WEB_SERVER}${pathname}`;
        await APIClient.login(redirectLink);
      }}
    >
      Login
    </Button>
  );
}

function ProfileButton() {
  const router = useRouter();

  return (
    <Button
      color="blue"
      variant="filled"
      onClick={() => router.push(Routes.profile)}
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

function SystemAdminMenu() {
  const router = useRouter();

  return (
    <>
      <Box visibleFrom="md">
        <Menu
          width={200}
          classNames={{
            dropdown: styles.systemAdminMenuDropdown,
          }}
        >
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
        onClick={(event) => event.stopPropagation()}
      >
        <NavLink
          label="Admin Panel"
          leftSection={<IconLayoutDashboard size={16} />}
          onClick={() => router.push(Routes.adminPanel)}
        />
        <NavLink
          label="Batch Processing"
          leftSection={<IconFileDescription size={16} />}
          onClick={() => router.push(Routes.batchProcessing)}
        />
      </NavLink>
    </>
  );
}
