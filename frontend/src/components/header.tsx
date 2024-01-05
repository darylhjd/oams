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
  IconBraces,
  IconCheck,
  IconChevronDown,
  IconCloudDownload,
  IconFileDescription,
  IconLayoutDashboard,
  IconReportAnalytics,
} from "@tabler/icons-react";
import { Routes } from "@/routing/routes";
import { APIClient } from "@/api/client";
import { UserRole } from "@/api/user";
import React from "react";
import {
  CanManageClassRulesAndReports,
  CanTakeAttendance,
} from "@/components/session_checker";

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
        position="right"
        padding="lg"
        opened={opened}
        onClose={close}
        onClick={close}
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
        <Group grow preventGrowOverflow>
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
        <Group grow preventGrowOverflow>
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

  if (
    !(
      session.data?.management_details.attendance ||
      session.data?.management_details.rules_and_reports
    )
  ) {
    return null;
  }

  return (
    <>
      <Box visibleFrom="md">
        <Menu width={200}>
          <MenuTarget>
            <Button
              color="red"
              variant="subtle"
              rightSection={<IconChevronDown />}
            >
              Class Management Menu
            </Button>
          </MenuTarget>
          <MenuDropdown>
            <CanTakeAttendance failNode={null}>
              <MenuItem
                leftSection={<IconCheck size={16} />}
                onClick={() => router.push(Routes.attendanceTaking)}
              >
                Attendance Taking
              </MenuItem>
            </CanTakeAttendance>
            <CanManageClassRulesAndReports failNode={null}>
              <MenuItem
                leftSection={<IconBraces size={16} />}
                onClick={() => router.push(Routes.attendanceRules)}
              >
                Attendance Rules
              </MenuItem>
            </CanManageClassRulesAndReports>
            <CanManageClassRulesAndReports failNode={null}>
              <MenuItem
                leftSection={<IconReportAnalytics size={16} />}
                onClick={() => router.push(Routes.classReports)}
              >
                Class Reports
              </MenuItem>
            </CanManageClassRulesAndReports>
          </MenuDropdown>
        </Menu>
      </Box>

      <NavLink
        color="red"
        hiddenFrom="md"
        label="Class Management Menu"
        active
        variant="subtle"
        onClick={(event) => event.stopPropagation()}
      >
        <CanTakeAttendance failNode={null}>
          <NavLink
            label="Attendance Taking"
            leftSection={<IconCheck size={16} />}
            onClick={() => router.push(Routes.attendanceTaking)}
          />
        </CanTakeAttendance>
        <CanManageClassRulesAndReports failNode={null}>
          <NavLink
            label="Attendance Rules"
            leftSection={<IconBraces size={16} />}
            onClick={() => router.push(Routes.attendanceRules)}
          />
        </CanManageClassRulesAndReports>
        <CanManageClassRulesAndReports failNode={null}>
          <NavLink
            label="Class Reports"
            leftSection={<IconReportAnalytics size={16} />}
            onClick={() => router.push(Routes.classReports)}
          />
        </CanManageClassRulesAndReports>
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
            <MenuItem
              leftSection={<IconFileDescription size={16} />}
              onClick={() => router.push(Routes.managerProcessing)}
            >
              Manager Processing
            </MenuItem>
            <MenuItem
              leftSection={<IconCloudDownload size={16} />}
              onClick={() => router.push(Routes.dataExport)}
            >
              Data Export
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
        <NavLink
          label="Data Export"
          leftSection={<IconCloudDownload size={16} />}
          onClick={() => router.push(Routes.dataExport)}
        />
      </NavLink>
    </>
  );
}
