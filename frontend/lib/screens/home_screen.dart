import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/router.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:go_router/go_router.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

// HomeScreen shows the home screen.
class HomeScreen extends ConsumerWidget {
  static const double mobilePadding = 10;
  static const double desktopPadding = 20;

  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isLoggedIn = ref.read(sessionProvider);
    return isLoggedIn ? const _HomeScreenLoggedIn() : const _HomeScreenGuest();
  }
}

// _HomeScreenLoggedIn shows the logged-in version of the home screen.
class _HomeScreenLoggedIn extends ConsumerWidget {
  static const double mobilePadding = 10;
  static const double desktopPadding = 20;

  const _HomeScreenLoggedIn();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ScreenTemplate(
      ResponsiveBreakpoints.of(context).isMobile ? mobile() : desktop(),
    );
  }

  Widget mobile() {
    return ListView(
      padding: const EdgeInsets.all(mobilePadding),
      children: const [],
    );
  }

  Widget desktop() {
    return ListView(
      padding: const EdgeInsets.all(desktopPadding),
      children: const [],
    );
  }
}

// _HomeScreenGuest shows the guest version of the home screen.
class _HomeScreenGuest extends StatelessWidget {
  static const double mobilePadding = 10;
  static const double desktopPadding = 20;

  const _HomeScreenGuest();

  @override
  Widget build(BuildContext context) {
    return ScreenTemplate(
      ResponsiveBreakpoints.of(context).isMobile
          ? mobile(context)
          : desktop(context),
    );
  }

  Widget mobile(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(mobilePadding),
      children: const [
        _WelcomeBanner(true),
      ],
    );
  }

  Widget desktop(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.symmetric(vertical: desktopPadding),
      children: const [
        _WelcomeBanner(false),
      ],
    );
  }
}

// _WelcomeBanner is applicable only to the guest home screen.
class _WelcomeBanner extends StatelessWidget {
  static const String bannerText =
      "Welcome to the Online Attendance Management System!";
  static const double buttonRadius = 5;
  static const String buttonText = "Get started";

  static const double mobilePadding = 10;
  static const double mobileMargin = 50;
  static const double mobileButtonVerticalPadding = 20;
  static const double mobileButtonHorizontalPadding = 15;
  static const double mobileButtonFontSize = 18;

  static const double desktopPadding = 20;
  static const double desktopMargin = 100;
  static const double desktopButtonVerticalPadding = 30;
  static const double desktopButtonHorizontalPadding = 20;
  static const double desktopButtonFontSize = 24;

  final bool isMobile;

  const _WelcomeBanner(this.isMobile);

  @override
  Widget build(BuildContext context) {
    return isMobile ? mobile(context) : desktop(context);
  }

  Widget mobile(BuildContext context) {
    return Card(
      child: Container(
        padding: const EdgeInsets.all(mobilePadding),
        margin: const EdgeInsets.symmetric(vertical: mobileMargin),
        alignment: Alignment.center,
        child: Column(
          children: [
            Text(
              bannerText,
              style: Theme.of(context)
                  .textTheme
                  .headlineSmall
                  ?.copyWith(fontWeight: FontWeight.bold),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: mobilePadding),
            FilledButton(
              onPressed: () => context.goNamed(Routes.login.name),
              style: FilledButton.styleFrom(
                padding: const EdgeInsets.symmetric(
                  vertical: mobileButtonVerticalPadding,
                  horizontal: mobileButtonHorizontalPadding,
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(buttonRadius),
                ),
                textStyle: const TextStyle(fontSize: mobileButtonFontSize),
              ),
              child: const Text(buttonText),
            )
          ],
        ),
      ),
    );
  }

  Widget desktop(BuildContext context) {
    return Card(
      child: Container(
        padding: const EdgeInsets.all(desktopPadding),
        margin: const EdgeInsets.symmetric(vertical: desktopMargin),
        alignment: Alignment.center,
        child: Column(
          children: [
            Text(
              bannerText,
              style: Theme.of(context)
                  .textTheme
                  .headlineLarge
                  ?.copyWith(fontWeight: FontWeight.bold),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: desktopPadding),
            FilledButton(
              onPressed: () => context.goNamed(Routes.login.name),
              style: FilledButton.styleFrom(
                padding: const EdgeInsets.symmetric(
                  vertical: desktopButtonVerticalPadding,
                  horizontal: desktopButtonHorizontalPadding,
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(buttonRadius),
                ),
                textStyle: const TextStyle(fontSize: desktopButtonFontSize),
              ),
              child: const Text(buttonText),
            ),
          ],
        ),
      ),
    );
  }
}
