import 'package:flutter/material.dart';

import '../widgets/centered_view.dart';
import '../widgets/nav_bar.dart';

// ScreenTemplate provides a template for all screen.
// In general, a navigation bar is provided.
class ScreenTemplate extends StatelessWidget {
  final Widget child;

  const ScreenTemplate(this.child, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: CenteredView(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            const NavBar(),
            Expanded(child: child),
          ],
        ),
      ),
    );
  }
}
