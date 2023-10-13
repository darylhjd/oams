import { DataTable, DataTableColumn } from "mantine-datatable";
import { useEffect, useState } from "react";

export const UserDataTableColumns = [
  { accessor: "id", title: "ID" },
  { accessor: "name" },
];
export const ClassDataTableColumns = [
  { accessor: "id", title: "ID" },
  { accessor: "code" },
  { accessor: "year" },
  { accessor: "semester" },
  { accessor: "programme" },
  { accessor: "au" },
];
export const ClassGroupDataTableColumns = [
  { accessor: "id", title: "ID" },
  { accessor: "class_id", title: "Class ID" },
  { accessor: "name" },
  { accessor: "class_type" },
];
export const ClassGroupSessionDataTableColumns = [
  { accessor: "class_group_id", title: "Class Group ID" },
  { accessor: "start_time", title: "Start Time" },
  { accessor: "end_time", title: "End Time" },
  { accessor: "venue" },
];

export function BatchDataTable({
  columns,
  records,
}: {
  columns: DataTableColumn<any>[];
  records: any[];
}) {
  const PAGE_SIZE = 50;

  const [page, setPage] = useState(1);
  const [pageRecords, setPageRecords] = useState(records.slice(0, PAGE_SIZE));

  useEffect(() => {
    const from = (page - 1) * PAGE_SIZE;
    const to = from + PAGE_SIZE;
    setPageRecords(records.slice(from, to));
  }, [page]);

  return (
    <DataTable
      height={500}
      withBorder
      withColumnBorders
      highlightOnHover
      scrollAreaProps={{ type: "auto" }}
      columns={columns}
      records={pageRecords}
      totalRecords={records.length}
      recordsPerPage={PAGE_SIZE}
      page={page}
      onPageChange={(p) => setPage(p)}
    />
  );
}

export function AsyncDataTable() {}
