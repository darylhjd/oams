"use client";

import { Text } from "@mantine/core";
import { useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { EntityLoader } from "@/app/admin-panel/entity_loader";
import { SessionEnrollment } from "@/api/session_enrollment";

export default function AdminPanelSessionEnrollmentPage({
  params,
}: {
  params: Params;
}) {
  const [sessionEnrollment, setSessionEnrollment] =
    useState<SessionEnrollment | null>(null);
  const promiseFunc = async () => {
    const data = await APIClient.sessionEnrollmentGet(params.id);
    return setSessionEnrollment(data.session_enrollment);
  };

  return (
    <EntityLoader promiseFunc={promiseFunc}>
      <SessionEnrollmentDisplay sessionEnrollment={sessionEnrollment!} />
    </EntityLoader>
  );
}

function SessionEnrollmentDisplay({
  sessionEnrollment,
}: {
  sessionEnrollment: SessionEnrollment;
}) {
  return <Text ta="center">{sessionEnrollment.id}</Text>;
}
