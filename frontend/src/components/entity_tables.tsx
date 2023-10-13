import { DataTable, DataTableColumn } from "mantine-datatable";
import { useEffect, useState } from "react";

const CreatedAtUpdatedAtDataTableColumns = [
  { accessor: "created_at", title: "Created At" },
  { accessor: "updated_at", title: "Updated At" },
];

export const UserBatchDataTableColumns = [
  { accessor: "id", title: "ID" },
  { accessor: "name" },
];
export const UserDataTableColumns = [
  ...UserBatchDataTableColumns,
  { accessor: "email" },
  { accessor: "role" },
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const ClassBatchDataTableColumns = [
  { accessor: "id", title: "ID" },
  { accessor: "code" },
  { accessor: "year" },
  { accessor: "semester" },
  { accessor: "programme" },
  { accessor: "au" },
];
export const ClassDataTableColumns = [
  ...ClassBatchDataTableColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const ClassGroupBatchDataTableColumns = [
  { accessor: "id", title: "ID" },
  { accessor: "class_id", title: "Class ID" },
  { accessor: "name" },
  { accessor: "class_type" },
];
export const ClassGroupDataTableColumns = [
  ...ClassGroupBatchDataTableColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const ClassGroupSessionBatchDataTableColumns = [
  { accessor: "class_group_id", title: "Class Group ID" },
  { accessor: "start_time", title: "Start Time" },
  { accessor: "end_time", title: "End Time" },
  { accessor: "venue" },
];
export const ClassGroupSessionDataTableColumns = [
  { accessor: "id", title: "ID" },
  ...ClassGroupSessionBatchDataTableColumns,
  ...CreatedAtUpdatedAtDataTableColumns,
];

export const SessionEnrollmentsDataTableColumns = [
  { accessor: "id", title: "ID" },
  { accessor: "session_id", title: "Session ID" },
  { accessor: "user_id", title: "User ID" },
  { accessor: "attended" },
  ...CreatedAtUpdatedAtDataTableColumns,
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

export function AsyncDataTable({
  columns,
}: {
  columns: DataTableColumn<any>[];
}) {
  const PAGE_SIZE = 100;

  const [page, setPage] = useState(1);
  const [pageRecords, setPageRecords] = useState([]);
  const [fetching, setFetching] = useState(true);
  const [totalRecords, setTotalRecords] = useState(0);

  useEffect(() => {}, [page]);

  return (
    <DataTable
      height={1000}
      withBorder
      withColumnBorders
      highlightOnHover
      scrollAreaProps={{ type: "auto" }}
      columns={columns}
      records={pageRecords}
      totalRecords={totalRecords}
      recordsPerPage={PAGE_SIZE}
      page={page}
      onPageChange={(p) => setPage(p)}
    />
  );
}
