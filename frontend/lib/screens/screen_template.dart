import 'package:flutter/material.dart';

import '../widgets/centered_view/centered_view.dart';
import '../widgets/nav_bar/nav_bar.dart';

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
