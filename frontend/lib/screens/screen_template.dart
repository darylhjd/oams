import 'package:flutter/material.dart';
import 'package:frontend/widgets/nav_bar.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

// Provides a template for all screens.
class ScreenTemplate extends StatelessWidget {
  final Widget _child;

  const ScreenTemplate(this._child, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _CenteredView(
        Column(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            const NavBar(),
            Expanded(child: _child),
          ],
        ),
      ),
    );
  }
}

// Allows the screens to be centered.
class _CenteredView extends StatelessWidget {
  static const double _mobileHorizontalPadding = 10;
  static const double _mobileVerticalPadding = 5;
  static const double _desktopHorizontalPadding = 20;
  static const double _desktopVerticalPadding = 10;
  static const double _desktopMaxWidth = 1500;

  final Widget _child;

  const _CenteredView(this._child, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ResponsiveBreakpoints.of(context).isMobile
        ? _mobile(context)
        : _desktop(context);
  }

  Widget _mobile(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: _mobileHorizontalPadding,
        vertical: _mobileVerticalPadding,
      ),
      alignment: Alignment.center,
      child: _child,
    );
  }

  Widget _desktop(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: _desktopHorizontalPadding,
        vertical: _desktopVerticalPadding,
      ),
      alignment: Alignment.center,
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: _desktopMaxWidth),
        child: _child,
      ),
    );
  }
}
