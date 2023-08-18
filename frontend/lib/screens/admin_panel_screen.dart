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
    Tab(text: "Placeholder"),
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
              const Placeholder(),
            ],
          ),
        ),
      ],
    );
  }
}

// _UsersSource holds the source for the users data.
class _UsersSource extends AsyncDataTableSource {
  int _rowCount = 1000;
  bool _isApproximateCount = true;

  @override
  Future<AsyncRowsResponse> getRows(int startIndex, int limit) async {
    final response = await APIClient.getUsers(startIndex, limit);

    // Last index from this request.
    final cursorEnd = startIndex + response.users.length;

    // If last index is more than current row count, set it to the last index,
    // else leave it.
    _rowCount = cursorEnd > _rowCount ? cursorEnd : _rowCount;

    // If the length of the users list is less than the limit, then we know
    // we have reached the max number.
    if (response.users.length < limit) {
      _rowCount = cursorEnd;
      _isApproximateCount = false;
    }

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

  @override
  int get rowCount => _rowCount;

  @override
  bool get isRowCountApproximate => _isApproximateCount;
}

// _UserEntities provides the paginated table to show the users data.
class _UserEntities extends StatefulWidget {
  @override
  _UserEntitiesState createState() => _UserEntitiesState();
}

// This holds the state for the _UserEntities widget.
class _UserEntitiesState extends State<_UserEntities>
    with AutomaticKeepAliveClientMixin {
  static final TableBorder _border = TableBorder(
    top: const BorderSide(color: Colors.black),
    bottom: BorderSide(color: Colors.grey[300]!),
    left: BorderSide(color: Colors.grey[300]!),
    right: BorderSide(color: Colors.grey[300]!),
    verticalInside: BorderSide(color: Colors.grey[300]!),
    horizontalInside: const BorderSide(color: Colors.grey, width: 1),
  );
  static const defaultRowsPerPage = 50;

  final _UsersSource _source = _UsersSource();
  late final List<DataColumn2> _columns;
  int _rowsPerPage = defaultRowsPerPage;

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
    return AsyncPaginatedDataTable2(
      columnSpacing: 10,
      minWidth: 800,
      border: _border,
      renderEmptyRowsInTheEnd: false,
      rowsPerPage: _rowsPerPage,
      availableRowsPerPage: const [10, defaultRowsPerPage, 100],
      onRowsPerPageChanged: (value) {
        _rowsPerPage = value!;
      },
      fixedLeftColumns: 1,
      columns: _columns,
      source: _source,
    );
  }

  @override
  bool get wantKeepAlive => true;
}
