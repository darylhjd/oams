import 'package:flutter/material.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

// AboutScreen shows information about the app.
class AboutScreen extends StatelessWidget {
  static const double mobilePadding = 10;
  static const double desktopPadding = 20;

  const AboutScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ScreenTemplate(
      ResponsiveBreakpoints.of(context).isMobile ? mobile() : desktop(),
    );
  }

  Widget mobile() {
    return ListView(
      padding: const EdgeInsets.all(mobilePadding),
      children: const [
        _PoweredByFlutter(true),
      ],
    );
  }

  Widget desktop() {
    return ListView(
      padding: const EdgeInsets.all(desktopPadding),
      children: const [
        _PoweredByFlutter(false),
      ],
    );
  }
}

// _PoweredByFlutter is the widget that shows the webpage is created with Flutter.
class _PoweredByFlutter extends StatelessWidget {
  static const double mobileSize = 100;
  static const double desktopSize = 200;
  static const String tagline = "Powered by:";
  final bool isMobile;

  const _PoweredByFlutter(this.isMobile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return isMobile ? mobile(context) : desktop(context);
  }

  Widget mobile(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Text(
          tagline,
          style: Theme.of(context).textTheme.labelLarge,
        ),
        const FlutterLogo(
          size: mobileSize,
          style: FlutterLogoStyle.stacked,
        ),
      ],
    );
  }

  Widget desktop(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Text(
          tagline,
          style: Theme.of(context).textTheme.labelLarge,
        ),
        const FlutterLogo(
          size: desktopSize,
          style: FlutterLogoStyle.horizontal,
        ),
      ],
    );
  }
}
