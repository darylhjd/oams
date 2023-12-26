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
} from "@mantine/core";
import { usePathname, useRouter } from "next/navigation";
import { useSessionUserStore } from "@/stores/session";
import { useDisclosure } from "@mantine/hooks";
import {
  IconCheck,
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
          inner: styles.systemAdminMenuDrawer,
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
  return (
    <>
      <Group visibleFrom="md">
        <SystemAdminMenu />
        <ClassGroupManagerMenu />
        <AboutButton />
        <ProfileButton />
        <LogoutButton />
      </Group>

      <Stack hiddenFrom="md" gap={0}>
        <SystemAdminMenu />
        <ClassGroupManagerMenu />
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

function ClassGroupManagerMenu() {
  const router = useRouter();
  const session = useSessionUserStore();

  const isClassGroupManager =
    session.data?.management_details.has_managed_class_groups;
  if (!isClassGroupManager) {
    return null;
  }

  return (
    <>
      <Box visibleFrom="md">
        <Menu
          width={200}
          classNames={{
            dropdown: styles.menuDropdown,
          }}
        >
          <MenuTarget>
            <Button
              color="red"
              variant="subtle"
              rightSection={<IconChevronDown />}
            >
              Class Group Manager Menu
            </Button>
          </MenuTarget>
          <MenuDropdown>
            <MenuItem
              leftSection={<IconCheck size={16} />}
              onClick={() => router.push(Routes.attendanceTaking)}
            >
              Attendance Taking
            </MenuItem>
          </MenuDropdown>
        </Menu>
      </Box>

      <NavLink
        color="red"
        hiddenFrom="md"
        label="Class Group Manager Menu"
        active
        variant="subtle"
        onClick={(event) => event.stopPropagation()}
      >
        <NavLink
          label="Attendance Taking"
          leftSection={<IconCheck size={16} />}
          onClick={() => router.push(Routes.attendanceTaking)}
        />
      </NavLink>
    </>
  );
}

function SystemAdminMenu() {
  const router = useRouter();
  const session = useSessionUserStore();

  const isSystemAdmin = session.data?.user.role == UserRole.SystemAdmin;
  if (!isSystemAdmin) {
    return null;
  }

  return (
    <>
      <Box visibleFrom="md">
        <Menu
          width={200}
          classNames={{
            dropdown: styles.menuDropdown,
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
            <MenuItem
              leftSection={<IconFileDescription size={16} />}
              onClick={() => router.push(Routes.managerProcessing)}
            >
              Manager Processing
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
        <NavLink
          label="Manager Processing"
          leftSection={<IconFileDescription size={16} />}
          onClick={() => router.push(Routes.managerProcessing)}
        />
      </NavLink>
    </>
  );
}
