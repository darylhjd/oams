import 'package:data_table_2/data_table_2.dart';
import 'package:flutter/material.dart';
import 'package:frontend/api/client.dart';
import 'package:frontend/screens/screen_template.dart';

// Shows the admin panel screen.
class AdminPanelScreen extends StatelessWidget {
  const AdminPanelScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return ScreenTemplate(
      _EntityViewer(),
    );
  }
}

// Provides a tab view to manage the different entities.
class _EntityViewer extends StatefulWidget {
  @override
  _EntityViewerState createState() => _EntityViewerState();
}

// Holds the state and main implementation for the entity viewer.
class _EntityViewerState extends State<_EntityViewer>
    with SingleTickerProviderStateMixin {
  static const List<Widget> _tabs = [
    Tab(text: "Users"),
    Tab(text: "Classes"),
    Tab(text: "Class Groups"),
    Tab(text: "Class Group Sessions"),
    Tab(text: "Session Enrollments"),
  ];
  late final TabController _controller;

  @override
  void initState() {
    super.initState();
    _controller = TabController(length: _tabs.length, vsync: this);
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        TabBar(
          controller: _controller,
          isScrollable: true,
          tabs: _tabs,
        ),
        Expanded(
          child: TabBarView(
            controller: _controller,
            children: [
              _UserEntities(),
              _ClassEntities(),
              _ClassGroupEntities(),
              _ClassGroupSessionEntities(),
              _SessionEnrollmentEntities(),
            ],
          ),
        ),
      ],
    );
  }
}

// _DataSource provides default values for a data source.
abstract class _DataSource extends AsyncDataTableSource {
  int estimatedRowCount = 1000;
  bool isApproximateCount = true;

  @override
  int get rowCount => estimatedRowCount;

  @override
  bool get isRowCountApproximate => isApproximateCount;

  // Call this method after getting the paginated data from your data source
  // to ensure that the row estimations are up to date.
  void updateRowEstimationState(int startIndex, int limit, int dataRowCount) {
    // Last index from this request.
    final cursorEnd = startIndex + dataRowCount;

    // If last index is more than current row count, set it to the last index,
    // else leave it.
    estimatedRowCount =
        cursorEnd > estimatedRowCount ? cursorEnd : estimatedRowCount;

    // If the length of the users list is less than the limit, then we know
    // we have reached the max number.
    if (dataRowCount < limit) {
      estimatedRowCount = cursorEnd;
      isApproximateCount = false;
    }
  }
}

// This provides a base implementation for the paginated data table.
abstract class _DataTableState extends State
    with AutomaticKeepAliveClientMixin {
  static final TableBorder tableBorder = TableBorder(
    top: const BorderSide(color: Colors.black),
    bottom: BorderSide(color: Colors.grey[300]!),
    left: BorderSide(color: Colors.grey[300]!),
    right: BorderSide(color: Colors.grey[300]!),
    verticalInside: BorderSide(color: Colors.grey[300]!),
    horizontalInside: const BorderSide(color: Colors.grey, width: 1),
  );
  static const defaultNumRowsPerPage = 50;

  final _DataSource source;

  _DataTableState(this.source);

  Widget withDefaultAsyncPaginatedTable({
    required List<DataColumn2> cols,
    required double minWidth,
    required int rowsPerPage,
    required void Function(int?) onRowsPerPageChanged,
  }) {
    return AsyncPaginatedDataTable2(
      columnSpacing: 10,
      dataRowHeight: 60,
      minWidth: minWidth,
      border: tableBorder,
      renderEmptyRowsInTheEnd: false,
      availableRowsPerPage: const [
        defaultNumRowsPerPage ~/ 5,
        defaultNumRowsPerPage,
        defaultNumRowsPerPage * 2
      ],
      fixedLeftColumns: 1,
      rowsPerPage: rowsPerPage,
      onRowsPerPageChanged: onRowsPerPageChanged,
      columns: cols,
      source: source,
    );
  }

  @override
  bool get wantKeepAlive => true;
}

// Holds the source for the users data.
class _UsersSource extends _DataSource {
  @override
  Future<AsyncRowsResponse> getRows(int startIndex, int limit) async {
    final response = await APIClient.getUsers(limit, startIndex);
    updateRowEstimationState(startIndex, limit, response.users.length);

    return AsyncRowsResponse(
      response.users.length,
      response.users
          .map(
            (u) => DataRow2(
              cells: [
                DataCell(
                  Text(
                    u.id,
                    style: const TextStyle(fontWeight: FontWeight.bold),
                  ),
                ),
                DataCell(Text(u.name)),
                DataCell(Text(u.email)),
                DataCell(Text(u.role.name)),
                DataCell(Text(u.createdAt.toString())),
                DataCell(Text(u.updatedAt.toString())),
              ],
            ),
          )
          .toList(),
    );
  }
}

// Provides the paginated table to show the users data.
class _UserEntities extends StatefulWidget {
  @override
  _UserEntitiesState createState() => _UserEntitiesState();
}

// This holds the state for the _UserEntities widget.
class _UserEntitiesState extends _DataTableState {
  static const double _minWidth = 700;
  late final List<DataColumn2> _columns;
  int _rowsPerPage = _DataTableState.defaultNumRowsPerPage;

  _UserEntitiesState() : super(_UsersSource());

  @override
  void initState() {
    super.initState();
    _columns = [
      "ID",
      "Name",
      "Email",
      "Role",
      "Created At",
      "Updated At",
    ]
        .map((s) => DataColumn2(
              label:
                  Text(s, style: const TextStyle(fontWeight: FontWeight.bold)),
            ))
        .toList();
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    return withDefaultAsyncPaginatedTable(
      cols: _columns,
      minWidth: _minWidth,
      rowsPerPage: _rowsPerPage,
      onRowsPerPageChanged: (value) {
        _rowsPerPage = value!;
      },
    );
  }
}

// Holds the source for the classes data.
class _ClassesSource extends _DataSource {
  @override
  Future<AsyncRowsResponse> getRows(int startIndex, int limit) async {
    final response = await APIClient.getClasses(limit, startIndex);
    updateRowEstimationState(startIndex, limit, response.classes.length);

    return AsyncRowsResponse(
        response.classes.length,
        response.classes
            .map((c) => DataRow2(
                  cells: [
                    DataCell(
                      Text(
                        c.id.toString(),
                        style: const TextStyle(fontWeight: FontWeight.bold),
                      ),
                    ),
                    DataCell(Text(c.code)),
                    DataCell(Text(c.year.toString())),
                    DataCell(Text(c.semester)),
                    DataCell(Text(c.programme)),
                    DataCell(Text(c.au.toString())),
                    DataCell(Text(c.createdAt.toString())),
                    DataCell(Text(c.updatedAt.toString())),
                  ],
                ))
            .toList());
  }
}

// Provides the paginated table to show the classes data.
class _ClassEntities extends StatefulWidget {
  @override
  _ClassEntitiesState createState() => _ClassEntitiesState();
}

// This holds the state for the _ClassEntities widget.
class _ClassEntitiesState extends _DataTableState {
  static const double _minWidth = 850;
  late final List<DataColumn2> _columns;
  int _rowsPerPage = _DataTableState.defaultNumRowsPerPage;

  _ClassEntitiesState() : super(_ClassesSource());

  @override
  void initState() {
    super.initState();
    _columns = [
      "ID",
      "Code",
      "Year",
      "Semester",
      "Programme",
      "AU",
      "Created At",
      "Updated At",
    ]
        .map((s) => DataColumn2(
              label:
                  Text(s, style: const TextStyle(fontWeight: FontWeight.bold)),
            ))
        .toList();
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    return withDefaultAsyncPaginatedTable(
      cols: _columns,
      minWidth: _minWidth,
      rowsPerPage: _rowsPerPage,
      onRowsPerPageChanged: (value) {
        _rowsPerPage = value!;
      },
    );
  }
}

// The source for the class groups data.
class _ClassGroupsSource extends _DataSource {
  @override
  Future<AsyncRowsResponse> getRows(int startIndex, int limit) async {
    final response = await APIClient.getClassGroups(limit, startIndex);
    updateRowEstimationState(startIndex, limit, response.classGroups.length);

    return AsyncRowsResponse(
      response.classGroups.length,
      response.classGroups
          .map((c) => DataRow2(
                cells: [
                  DataCell(
                    Text(
                      c.id.toString(),
                      style: const TextStyle(fontWeight: FontWeight.bold),
                    ),
                  ),
                  DataCell(Text(c.classId.toString())),
                  DataCell(Text(c.name)),
                  DataCell(Text(c.classType.name)),
                  DataCell(Text(c.createdAt.toString())),
                  DataCell(Text(c.updatedAt.toString())),
                ],
              ))
          .toList(),
    );
  }
}

// Provides the paginated table to show the class groups data.
class _ClassGroupEntities extends StatefulWidget {
  @override
  _ClassGroupEntitiesState createState() => _ClassGroupEntitiesState();
}

// This holds the state for the _ClassGroupEntities widget.
class _ClassGroupEntitiesState extends _DataTableState {
  static const double _minWidth = 620;
  late final List<DataColumn2> _columns;
  int _rowsPerPage = _DataTableState.defaultNumRowsPerPage;

  _ClassGroupEntitiesState() : super(_ClassGroupsSource());

  @override
  void initState() {
    super.initState();
    _columns = [
      "ID",
      "Class ID",
      "Name",
      "Class Type",
      "Created At",
      "Updated At",
    ]
        .map((s) => DataColumn2(
              label:
                  Text(s, style: const TextStyle(fontWeight: FontWeight.bold)),
            ))
        .toList();
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    return withDefaultAsyncPaginatedTable(
      cols: _columns,
      minWidth: _minWidth,
      rowsPerPage: _rowsPerPage,
      onRowsPerPageChanged: (value) {
        _rowsPerPage = value!;
      },
    );
  }
}

// The source for the class group sessions data.
class _ClassGroupSessionsSource extends _DataSource {
  @override
  Future<AsyncRowsResponse> getRows(int startIndex, int limit) async {
    final response = await APIClient.getClassGroupSessions(limit, startIndex);
    updateRowEstimationState(
        startIndex, limit, response.classGroupSessions.length);

    return AsyncRowsResponse(
      response.classGroupSessions.length,
      response.classGroupSessions
          .map((c) => DataRow2(
                cells: [
                  DataCell(
                    Text(
                      c.id.toString(),
                      style: const TextStyle(fontWeight: FontWeight.bold),
                    ),
                  ),
                  DataCell(Text(c.classGroupId.toString())),
                  DataCell(Text(c.startTime.toString())),
                  DataCell(Text(c.endTime.toString())),
                  DataCell(Text(c.venue)),
                  DataCell(Text(c.createdAt.toString())),
                  DataCell(Text(c.updatedAt.toString())),
                ],
              ))
          .toList(),
    );
  }
}

// Provides the paginated table to show the class group sessions data.
class _ClassGroupSessionEntities extends StatefulWidget {
  @override
  _ClassGroupSessionEntitiesState createState() =>
      _ClassGroupSessionEntitiesState();
}

// This holds the state for the _ClassGroupSessionEntities widget.
class _ClassGroupSessionEntitiesState extends _DataTableState {
  static const double _minWidth = 1000;
  late final List<DataColumn2> _columns;
  int _rowsPerPage = _DataTableState.defaultNumRowsPerPage;

  _ClassGroupSessionEntitiesState() : super(_ClassGroupSessionsSource());

  @override
  void initState() {
    super.initState();
    _columns = [
      "ID",
      "Class Group ID",
      "Start Time",
      "End Time",
      "Venue",
      "Created At",
      "Updated At",
    ]
        .map((s) => DataColumn2(
              label:
                  Text(s, style: const TextStyle(fontWeight: FontWeight.bold)),
            ))
        .toList();
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    return withDefaultAsyncPaginatedTable(
      cols: _columns,
      minWidth: _minWidth,
      rowsPerPage: _rowsPerPage,
      onRowsPerPageChanged: (value) {
        _rowsPerPage = value!;
      },
    );
  }
}

// The source for the session enrollments data.
class _SessionEnrollmentsSource extends _DataSource {
  @override
  Future<AsyncRowsResponse> getRows(int startIndex, int limit) async {
    final response = await APIClient.getSessionEnrollments(limit, startIndex);
    updateRowEstimationState(
        startIndex, limit, response.sessionEnrollments.length);

    return AsyncRowsResponse(
      response.sessionEnrollments.length,
      response.sessionEnrollments
          .map((c) => DataRow2(
                cells: [
                  DataCell(
                    Text(
                      c.id.toString(),
                      style: const TextStyle(fontWeight: FontWeight.bold),
                    ),
                  ),
                  DataCell(Text(c.sessionId.toString())),
                  DataCell(Text(c.userId)),
                  DataCell(Text(c.attended.toString())),
                  DataCell(Text(c.createdAt.toString())),
                  DataCell(Text(c.updatedAt.toString())),
                ],
              ))
          .toList(),
    );
  }
}

// Provides the paginated table to show the session enrollments data.
class _SessionEnrollmentEntities extends StatefulWidget {
  @override
  _SessionEnrollmentEntitiesState createState() =>
      _SessionEnrollmentEntitiesState();
}

// This holds the state for the _SessionEnrollmentEntities widget.
class _SessionEnrollmentEntitiesState extends _DataTableState {
  static const double _minWidth = 650;
  late final List<DataColumn2> _columns;
  int _rowsPerPage = _DataTableState.defaultNumRowsPerPage;

  _SessionEnrollmentEntitiesState() : super(_SessionEnrollmentsSource());

  @override
  void initState() {
    super.initState();
    _columns = [
      "ID",
      "Session ID",
      "User ID",
      "Attended",
      "Created At",
      "Updated At",
    ]
        .map((s) => DataColumn2(
              label:
                  Text(s, style: const TextStyle(fontWeight: FontWeight.bold)),
            ))
        .toList();
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    return withDefaultAsyncPaginatedTable(
      cols: _columns,
      minWidth: _minWidth,
      rowsPerPage: _rowsPerPage,
      onRowsPerPageChanged: (value) {
        _rowsPerPage = value!;
      },
    );
  }
}
