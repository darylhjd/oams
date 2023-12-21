"use client";

import { Text } from "@mantine/core";
import { useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { EntityLoader } from "@/components/entity_loader";
import { Class } from "@/api/class";

export default function AdminPanelClassPage({ params }: { params: Params }) {
  const [cls, setClass] = useState<Class | null>(null);
  const promiseFunc = async () => {
    const data = await APIClient.classGet(params.id);
    return setClass(data.class);
  };

  return (
    <EntityLoader promiseFunc={promiseFunc}>
      <ClassDisplay cls={cls!} />
    </EntityLoader>
  );
}

function ClassDisplay({ cls }: { cls: Class }) {
  return <Text ta="center">{cls.id}</Text>;
}
