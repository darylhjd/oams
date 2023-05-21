import 'package:flutter/material.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

// CenteredView allows the screens to be centered.
class CenteredView extends StatelessWidget {
  static const double mobileHorizontalPadding = 10;
  static const double mobileVerticalPadding = 5;
  static const double desktopHorizontalPadding = 20;
  static const double desktopVerticalPadding = 10;
  static const double desktopMaxWidth = 1100;

  final Widget child;

  const CenteredView({Key? key, required this.child}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ResponsiveBreakpoints.of(context).smallerThan(DESKTOP)
        ? mobile(context)
        : desktop(context);
  }

  Widget mobile(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
          horizontal: mobileHorizontalPadding, vertical: mobileVerticalPadding),
      alignment: Alignment.center,
      child: child,
    );
  }

  Widget desktop(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
          horizontal: desktopHorizontalPadding,
          vertical: desktopVerticalPadding),
      alignment: Alignment.center,
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: desktopMaxWidth),
        child: child,
      ),
    );
  }
}
