import 'package:flutter/material.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

// Shows information about the app.
class AboutScreen extends StatelessWidget {
  static const double _mobilePadding = 10;
  static const double _desktopPadding = 20;

  const AboutScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ScreenTemplate(
      ResponsiveBreakpoints.of(context).isMobile ? _mobile() : _desktop(),
    );
  }

  Widget _mobile() {
    return ListView(
      padding: const EdgeInsets.all(_mobilePadding),
      children: const [
        _PoweredByFlutter(true),
      ],
    );
  }

  Widget _desktop() {
    return ListView(
      padding: const EdgeInsets.all(_desktopPadding),
      children: const [
        _PoweredByFlutter(false),
      ],
    );
  }
}

// Widget that shows that the webpage is created with Flutter.
class _PoweredByFlutter extends StatelessWidget {
  static const double _mobileSize = 100;
  static const double _desktopSize = 200;
  static const String _tagline = "Powered by:";
  final bool _isMobile;

  const _PoweredByFlutter(this._isMobile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return _isMobile ? _mobile(context) : _desktop(context);
  }

  Widget _mobile(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Text(
          _tagline,
          style: Theme.of(context).textTheme.labelLarge,
        ),
        const FlutterLogo(
          size: _mobileSize,
          style: FlutterLogoStyle.stacked,
        ),
      ],
    );
  }

  Widget _desktop(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Text(
          _tagline,
          style: Theme.of(context).textTheme.labelLarge,
        ),
        const FlutterLogo(
          size: _desktopSize,
          style: FlutterLogoStyle.horizontal,
        ),
      ],
    );
  }
}
