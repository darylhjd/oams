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
