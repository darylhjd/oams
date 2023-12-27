"use client";

import { Text } from "@mantine/core";
import { Params } from "@/app/attendance-rules/[id]/layout";
import { useState } from "react";
import { AttendanceRuleGetResponse } from "@/api/attendance_rule";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";

export default function AttendanceRulePage({ params }: { params: Params }) {
  const [coordinatingClassRules, setCoordinatingClassRules] =
    useState<AttendanceRuleGetResponse>({} as AttendanceRuleGetResponse);
  const promiseFunc = async () => {
    const data = await APIClient.attendanceRuleGet(params.id);
    return setCoordinatingClassRules(data);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Text ta="center">This is the attendance rule page for {params.id}</Text>
    </RequestLoader>
  );
}
