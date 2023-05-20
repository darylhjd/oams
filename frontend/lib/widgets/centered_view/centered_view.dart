import 'package:flutter/material.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

import 'desktop.dart';
import 'mobile.dart';

class CenteredView extends StatelessWidget {
  final Widget child;

  const CenteredView({Key? key, required this.child}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    if (ResponsiveBreakpoints.of(context).smallerThan(DESKTOP)) {
      return CenteredViewMobile(child: child);
    }
    return CenteredViewDesktop(child: child);
  }
}
