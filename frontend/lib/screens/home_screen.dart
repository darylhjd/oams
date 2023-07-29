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
        Divider(),
        _FinalActionCall(true),
        Divider(),
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
        Divider(),
        _FinalActionCall(false),
        Divider(),
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
            "Discover our features.",
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
        _DetailedSystemDashboard(true),
        _CloudBasedFeature(true),
      ],
    );
  }

  Widget desktop() {
    return const Column(
      children: [
        _CloudBasedFeature(false),
        _DetailedSystemDashboard(false),
        _CloudBasedFeature(false),
      ],
    );
  }
}

// _FeatureCard provides the template for a feature showcase.
class _FeatureCard extends StatelessWidget {
  static const double mobileCardMargin = 10;
  static const double mobilePadding = 10;
  static const double mobileMargin = 20;

  static const double desktopPadding = 20;
  static const double desktopMargin = 20;

  final String title;
  final String body;
  final bool isMobile;

  const _FeatureCard(this.title, this.body, this.isMobile);

  @override
  Widget build(BuildContext context) {
    return isMobile ? mobile(context) : desktop(context);
  }

  Widget mobile(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(vertical: mobileCardMargin),
      child: Container(
        padding: const EdgeInsets.all(mobilePadding),
        margin: const EdgeInsets.symmetric(vertical: mobileMargin),
        alignment: Alignment.center,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
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
        alignment: Alignment.topLeft,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
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

// _CloudBasedFeature shows the online cloud-based feature of OAMS.
class _CloudBasedFeature extends StatelessWidget {
  static const String title = "Cloud-based management system";
  static const String body =
      "Say goodbye to cumbersome paper-based attendance sheets. OAMS stores attendance data online so you can stay up-to-date wherever you are - instantly.";

  final bool isMobile;

  const _CloudBasedFeature(this.isMobile);

  @override
  Widget build(BuildContext context) {
    return _FeatureCard(title, body, isMobile);
  }
}

// _DetailedSystemDashboard shows the reports and analytics feature of OAMS.
class _DetailedSystemDashboard extends StatelessWidget {
  static const String title = "Detailed system dashboard";
  static const String body =
      "Get comprehensive access to attendance data. Our system dashboard provides you with detailed visualisations of your most important analytics - all at your fingertips.";

  final bool isMobile;

  const _DetailedSystemDashboard(this.isMobile);

  @override
  Widget build(BuildContext context) {
    return _FeatureCard(title, body, isMobile);
  }
}

// _CallToAction is the final call to the user to sign up.
class _FinalActionCall extends StatelessWidget {
  static const String callText = "Try OAMS now.";
  static const double buttonRadius = 5;
  static const String buttonText = "Sign Up/Login";

  static const double mobilePadding = 10;
  static const double mobileMargin = 40;
  static const double mobileButtonFontSize = 20;
  static const double mobileButtonVerticalPadding = 15;
  static const double mobileButtonHorizontalPadding = 10;

  static const double desktopPadding = 20;
  static const double desktopMargin = 100;
  static const double desktopButtonFontSize = 24;
  static const double desktopButtonVerticalPadding = 20;
  static const double desktopButtonHorizontalPadding = 10;

  final bool isMobile;

  const _FinalActionCall(this.isMobile);

  @override
  Widget build(BuildContext context) {
    return isMobile ? mobile(context) : desktop(context);
  }

  Widget mobile(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: mobileMargin),
      child: Column(
        children: [
          Text(
            callText,
            style: Theme.of(context).textTheme.headlineSmall,
          ),
          const SizedBox(height: mobilePadding),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              padding: const EdgeInsets.symmetric(
                vertical: mobileButtonVerticalPadding,
                horizontal: mobileButtonHorizontalPadding,
              ),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(buttonRadius),
              ),
              textStyle: const TextStyle(fontSize: mobileButtonFontSize),
            ),
            onPressed: () => context.goNamed(Routes.login.name),
            child: const Text(buttonText),
          )
        ],
      ),
    );
  }

  Widget desktop(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: desktopMargin),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            callText,
            style: Theme.of(context).textTheme.headlineLarge,
          ),
          const SizedBox(width: desktopPadding),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              padding: const EdgeInsets.symmetric(
                vertical: desktopButtonVerticalPadding,
                horizontal: desktopButtonHorizontalPadding,
              ),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(buttonRadius),
              ),
              textStyle: const TextStyle(fontSize: desktopButtonFontSize),
            ),
            onPressed: () => context.goNamed(Routes.login.name),
            child: const Text(buttonText),
          ),
        ],
      ),
    );
  }
}
