"use client";

import { Text } from "@mantine/core";
import { Params } from "@/app/attendance-rules/[id]/layout";
import { useState } from "react";
import { CoordinatingClassGetResponse } from "@/api/coordinating_class";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";

export default function AttendanceRulePage({ params }: { params: Params }) {
  const [coordinatingClassRules, setCoordinatingClassRules] =
    useState<CoordinatingClassGetResponse>({} as CoordinatingClassGetResponse);
  const promiseFunc = async () => {
    const data = await APIClient.coordinatingClassGet(params.id);
    return setCoordinatingClassRules(data);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Text ta="center">This is the attendance rule page for {params.id}</Text>
    </RequestLoader>
  );
}
