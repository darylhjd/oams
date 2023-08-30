"use client";

import { buildRedirectUrlQueryParamsString } from "@/routes/checks";
import { Routes } from "@/routes/routes";
import { sessionStore } from "@/states/session";
import {
  Card,
  Container,
  Stack,
  TextInput,
  Title,
  createStyles,
} from "@mantine/core";
import { getURL } from "next/dist/shared/lib/utils";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

const useStyles = createStyles((theme) => ({
  card: {
    width: "100%",
  },
}));

export default function LoginPage() {
  const router = useRouter();
  const session = sessionStore();

  useEffect(() => {
    if (session.data == null) {
      router.replace(
        `${Routes.login}?${buildRedirectUrlQueryParamsString(getURL())}`,
      );
    }
  });

  if (session.data == null) {
    return null;
  }

  return (
    <Stack align="center">
      <Header />
      <Information />
    </Stack>
  );
}

function Header() {
  return <Title>Your Profile</Title>;
}

function Information() {
  const session = sessionStore();
  const { classes } = useStyles();

  const fields = [
    [session.data!.session_user.id, "ID"],
    [session.data!.session_user.name, "Name"],
    [session.data!.session_user.email, "Email"],
    [session.data!.session_user.role, "Role"],
  ];

  return (
    <Card className={classes.card} shadow="md" withBorder>
      <Container size="30em">
        <Stack>
          {fields.map((field) => (
            <TextInput
              placeholder={field[0]}
              label={field[1]}
              key={field[1]}
              disabled
            />
          ))}
        </Stack>
      </Container>
    </Card>
  );
}
