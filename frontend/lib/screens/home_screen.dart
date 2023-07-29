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
      padding: const EdgeInsets.symmetric(vertical: mobilePadding),
      children: const [
        _WelcomeBanner(true),
        _FeaturesDivider(),
        _Features(true),
      ],
    );
  }

  Widget desktop(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.symmetric(vertical: desktopPadding),
      children: const [
        _WelcomeBanner(false),
        _FeaturesDivider(),
        _Features(false),
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

// _FeaturesDivider is applicable only to the guest home screen.
// It shows the divider between the welcome banner and the bottom.
class _FeaturesDivider extends StatelessWidget {
  static const double mobilePadding = 20;
  static const double desktopPadding = 30;

  const _FeaturesDivider();

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.symmetric(
        vertical: ResponsiveBreakpoints.of(context).isMobile
            ? mobilePadding
            : desktopPadding,
      ),
      child: const Column(
        children: [
          Text(
            "Scroll down to discover more.",
            textAlign: TextAlign.center,
          ),
          Icon(Icons.arrow_downward),
        ],
      ),
    );
  }
}

// _Features is only applicable to the guest home screen.
// It contains the list of feature "selling points" of OAMS.
class _Features extends StatelessWidget {
  final bool isMobile;

  const _Features(this.isMobile);

  @override
  Widget build(BuildContext context) {
    return isMobile ? mobile(context) : desktop();
  }

  Widget mobile(BuildContext context) {
    return const Column(
      children: [
        _CloudBasedFeature(true),
      ],
    );
  }

  Widget desktop() {
    return const Row(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
      children: [
        Flexible(
          child: _CloudBasedFeature(false),
        ),
      ],
    );
  }
}

// _CloudBasedFeature shows the online cloud-based feature of OAMS.
class _CloudBasedFeature extends StatelessWidget {
  static const String headline = "Cloud-based management system";
  static const String body =
      "Say goodbye to cumbersome paper-based attendance sheets. OAMS stores attendance data online so you can stay up-to-date wherever you are - instantly.";

  static const double mobilePadding = 10;
  static const double mobileMargin = 10;

  static const double desktopPadding = 20;
  static const double desktopMargin = 20;
  static const double desktopMaxWidth = 400;

  final bool isMobile;

  const _CloudBasedFeature(this.isMobile);

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
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              headline,
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: mobilePadding),
            Text(
              body,
              style: Theme.of(context).textTheme.bodyLarge,
            ),
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
        width: desktopMaxWidth,
        alignment: Alignment.center,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              headline,
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: desktopPadding),
            Text(
              body,
              style: Theme.of(context).textTheme.bodyLarge,
            ),
          ],
        ),
      ),
    );
  }
}
