"use client";

import { Button, Center, Container, Image, Group } from "@mantine/core";
import styles from "@/styles/Header.module.css";
import { useRouter } from "next/navigation";

export default function Header() {
  return (
    <Container className={styles.header} fluid>
      <Items />
    </Container>
  );
}

function Items() {
  return (
    <Center>
      <Group className={styles.items} justify="space-between" align="center">
        <Logo />
        <GroupItems />
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
      onClick={() => router.push("/")}
    >
      <div>
        <Image src="logo.png" fit="contain" />
      </div>
    </Button>
  );
}

function GroupItems() {
  return (
    <Group gap={5}>
      <Button>Login</Button>
    </Group>
  );
}
