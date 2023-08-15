import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/client.dart';
import 'package:frontend/api/models.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:frontend/widgets/dialogs.dart';

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
    Tab(text: "Tab 1"),
    Tab(text: "Tab 2"),
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
            children: const [
              _UserEntities(),
              Placeholder(),
            ],
          ),
        ),
      ],
    );
  }
}

final _usersProvider = FutureProvider<GetUsersResponse>((ref) async {
  try {
    return await APIClient.getUsers();
  } catch (e) {
    rethrow;
  }
});

class _UserEntities extends ConsumerWidget {
  const _UserEntities();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final usersAsync = ref.watch(_usersProvider);

    return usersAsync.when(
      data: (data) {
        return _dataTable(data);
      },
      error: (error, stackTrace) => const InvalidSession(),
      loading: () => const Center(child: CircularProgressIndicator()),
    );
  }

  Widget _dataTable(GetUsersResponse data) {
    final rows = List<DataRow>.from(
      data.users.map(
        (u) => DataRow(
          cells: [
            DataCell(Text(u.id)),
          ],
        ),
      ),
    );
    return SingleChildScrollView(
      child: DataTable(
        columns: const [
          DataColumn(label: Text("ID")),
        ],
        rows: rows,
      ),
    );
  }
}
