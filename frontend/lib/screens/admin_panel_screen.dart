import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
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
class _EntityViewer extends ConsumerStatefulWidget {
  @override
  _EntityViewerState createState() => _EntityViewerState();
}

// Holds the state and main implementation for the entity viewer.
class _EntityViewerState extends ConsumerState
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
              Placeholder(),
              Placeholder(),
            ],
          ),
        ),
      ],
    );
  }
}
