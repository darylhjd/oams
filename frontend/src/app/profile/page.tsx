"use client";

import { redirectIfNotLoggedIn } from "@/routes/checks";
import { sessionStore } from "@/states/session";
import {
  Card,
  Container,
  Stack,
  TextInput,
  Title,
  createStyles,
} from "@mantine/core";

const useStyles = createStyles((theme) => ({
  card: {
    width: "100%",
  },
}));

export default function LoginPage() {
  redirectIfNotLoggedIn();

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

  return (
    <Card className={classes.card} shadow="md" withBorder>
      <Container size="30em">
        <Stack>
          <TextInput
            placeholder={session.data!.session_user.id}
            label="ID"
            disabled
          />
          <TextInput
            placeholder={session.data!.session_user.name}
            label="Name"
            disabled
          />
          <TextInput
            placeholder={session.data!.session_user.email}
            label="Email"
            disabled
          />
          <TextInput
            placeholder={session.data!.session_user.role}
            label="Role"
            disabled
          />
        </Stack>
      </Container>
    </Card>
  );
}
