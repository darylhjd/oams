"use client";

import { useState } from "react";
import { Params } from "./layout";
import { ClassGroupSession } from "@/api/class_group_session";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { Text } from "@mantine/core";

export default function SessionAttendanceTakingPage({
  params,
}: {
  params: Params;
}) {
  const [session, setSession] = useState<ClassGroupSession>(
    {} as ClassGroupSession,
  );
  const promiseFunc = async () => {
    const data = await APIClient.classGroupSessionGet(params.id);
    return setSession(data.class_group_session);

    // TODO: Get enrollments related to session.
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Page session={session} />
    </RequestLoader>
  );
}

function Page({ session }: { session: ClassGroupSession }) {
  return <Text ta="center">{session.venue}</Text>;
}
