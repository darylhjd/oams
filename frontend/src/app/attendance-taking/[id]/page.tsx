"use client";

import { useState } from "react";
import { Params } from "./layout";
import { ClassGroupSession } from "@/api/class_group_session";
import { APIClient } from "@/api/client";
import { EntityLoader } from "@/components/entity_loader";
import { Text } from "@mantine/core";

export default function SessionAttendanceTakingPage({
  params,
}: {
  params: Params;
}) {
  const [session, setSession] = useState<ClassGroupSession | null>(null);
  const promiseFunc = async () => {
    const data = await APIClient.classGroupSessionGet(params.id);
    return setSession(data.class_group_session);

    // TODO: Get enrollments related to session.
  };

  return (
    <EntityLoader promiseFunc={promiseFunc}>
      <Page session={session!} />
    </EntityLoader>
  );
}

function Page({ session }: { session: ClassGroupSession }) {
  return <Text ta="center">{session.venue}</Text>;
}
