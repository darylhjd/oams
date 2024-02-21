"use client";

import styles from "@/styles/ClassAdministrationsPage.module.css";

import {
  Container,
  Group,
  Space,
  Table,
  TableScrollContainer,
  TableTbody,
  TableTd,
  TableTh,
  TableThead,
  TableTr,
  Text,
  Title,
} from "@mantine/core";
import { useState } from "react";
import { CoordinatingClass } from "@/api/coordinating_class";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { CoordinatingClassesDataTableColumns } from "@/components/columns";
import { Routes } from "@/routing/routes";
import { useRouter } from "next/navigation";
import {
  flexRender,
  MRT_Row,
  MRT_TableBodyCellValue,
  MRT_TablePagination,
  useMantineReactTable,
} from "mantine-react-table";

export default function ClassAdministrationsPage() {
  const [coordinatingClasses, setCoordinatingClasses] = useState<
    CoordinatingClass[]
  >([]);
  const router = useRouter();

  const promiseFunc = async () => {
    const data = await APIClient.coordinatingClassesGet();
    return setCoordinatingClasses(data.coordinating_classes);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Container className={styles.container} fluid>
        <PageTitle />
        <Space h="md" />
        <Text ta="center">Select a class to manage.</Text>
        <Space h="xs" />
        <CoordinatingClassPicker
          coordinatingClasses={coordinatingClasses}
          onRowClick={(row) =>
            router.push(Routes.classAdministration + row.original.id)
          }
        />
      </Container>
    </RequestLoader>
  );
}

function PageTitle() {
  return (
    <Title order={2} ta="center">
      Class Administration
    </Title>
  );
}

function CoordinatingClassPicker({
  coordinatingClasses,
  onRowClick,
}: {
  coordinatingClasses: CoordinatingClass[];
  onRowClick: (row: MRT_Row<CoordinatingClass>) => void;
}) {
  const table = useMantineReactTable({
    columns: CoordinatingClassesDataTableColumns,
    data: coordinatingClasses,
    initialState: {
      pagination: { pageSize: 5, pageIndex: 0 },
    },
  });

  return (
    <>
      <TableScrollContainer minWidth={500}>
        <Table verticalSpacing="md" highlightOnHover withTableBorder>
          <TableThead>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableTr key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableTh key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.Header ??
                            header.column.columnDef.header,
                          header.getContext(),
                        )}
                  </TableTh>
                ))}
              </TableTr>
            ))}
          </TableThead>
          <TableTbody>
            {table.getRowModel().rows.map((row) => (
              <TableTr
                key={row.id}
                className={styles.tableRow}
                onClick={() => onRowClick(row)}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableTd key={cell.id}>
                    <MRT_TableBodyCellValue cell={cell} table={table} />
                  </TableTd>
                ))}
              </TableTr>
            ))}
          </TableTbody>
        </Table>
      </TableScrollContainer>
      <Group justify="right">
        <MRT_TablePagination table={table} />
      </Group>
    </>
  );
}
