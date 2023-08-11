import 'package:flutter/material.dart';
import 'package:frontend/router.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:go_router/go_router.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

// Shows the guest version of the home screen.
class HomeScreenGuest extends StatelessWidget {
  static const double _mobilePadding = 10;
  static const double _desktopPadding = 20;

  const HomeScreenGuest({super.key});

  @override
  Widget build(BuildContext context) {
    return ScreenTemplate(
      ResponsiveBreakpoints.of(context).isMobile
          ? _mobile(context)
          : _desktop(context),
    );
  }

  Widget _mobile(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.symmetric(vertical: _mobilePadding),
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

  Widget _desktop(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.symmetric(vertical: _desktopPadding),
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

// Shows the welcome banner on the guest home screen.
class _WelcomeBanner extends StatelessWidget {
  static const String _bannerText =
      "Welcome to the Online Attendance Management System!";
  static const double _buttonRadius = 5;
  static const String _buttonText = "Get started";

  static const double _mobilePadding = 10;
  static const double _mobileMargin = 50;
  static const double _mobileButtonVerticalPadding = 20;
  static const double _mobileButtonHorizontalPadding = 15;
  static const double _mobileButtonFontSize = 18;

  static const double _desktopPadding = 20;
  static const double _desktopMargin = 100;
  static const double _desktopButtonVerticalPadding = 30;
  static const double _desktopButtonHorizontalPadding = 20;
  static const double _desktopButtonFontSize = 24;

  final bool _isMobile;

  const _WelcomeBanner(this._isMobile);

  @override
  Widget build(BuildContext context) {
    return _isMobile ? _mobile(context) : _desktop(context);
  }

  Widget _mobile(BuildContext context) {
    return Card(
      child: Container(
        padding: const EdgeInsets.all(_mobilePadding),
        margin: const EdgeInsets.symmetric(vertical: _mobileMargin),
        alignment: Alignment.center,
        child: Column(
          children: [
            Text(
              _bannerText,
              style: Theme.of(context)
                  .textTheme
                  .headlineSmall
                  ?.copyWith(fontWeight: FontWeight.bold),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: _mobilePadding),
            FilledButton(
              onPressed: () => context.goNamed(Routes.login.name),
              style: FilledButton.styleFrom(
                padding: const EdgeInsets.symmetric(
                  vertical: _mobileButtonVerticalPadding,
                  horizontal: _mobileButtonHorizontalPadding,
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(_buttonRadius),
                ),
                textStyle: const TextStyle(fontSize: _mobileButtonFontSize),
              ),
              child: const Text(_buttonText),
            )
          ],
        ),
      ),
    );
  }

  Widget _desktop(BuildContext context) {
    return Card(
      child: Container(
        padding: const EdgeInsets.all(_desktopPadding),
        margin: const EdgeInsets.symmetric(vertical: _desktopMargin),
        alignment: Alignment.center,
        child: Column(
          children: [
            Text(
              _bannerText,
              style: Theme.of(context)
                  .textTheme
                  .headlineLarge
                  ?.copyWith(fontWeight: FontWeight.bold),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: _desktopPadding),
            FilledButton(
              onPressed: () => context.goNamed(Routes.login.name),
              style: FilledButton.styleFrom(
                padding: const EdgeInsets.symmetric(
                  vertical: _desktopButtonVerticalPadding,
                  horizontal: _desktopButtonHorizontalPadding,
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(_buttonRadius),
                ),
                textStyle: const TextStyle(fontSize: _desktopButtonFontSize),
              ),
              child: const Text(_buttonText),
            ),
          ],
        ),
      ),
    );
  }
}

// Shows the divider between the welcome banner and the rest of the content on
// the guest home screen.
class _FeaturesDivider extends StatelessWidget {
  static const double _mobilePadding = 20;
  static const double _desktopPadding = 30;

  const _FeaturesDivider();

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.symmetric(
        vertical: ResponsiveBreakpoints.of(context).isMobile
            ? _mobilePadding
            : _desktopPadding,
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

// Contains the list of features of OAMS on the guest home screen.
class _Features extends StatelessWidget {
  final bool _isMobile;

  const _Features(this._isMobile);

  @override
  Widget build(BuildContext context) {
    return _isMobile ? _mobile(context) : _desktop();
  }

  Widget _mobile(BuildContext context) {
    return const Column(
      children: [
        _CloudBasedFeature(true),
        _DetailedSystemDashboard(true),
      ],
    );
  }

  Widget _desktop() {
    return const Column(
      children: [
        _CloudBasedFeature(false),
        _DetailedSystemDashboard(false),
      ],
    );
  }
}

// Provides the template for a feature showcase.
class _FeatureCard extends StatelessWidget {
  static const double _mobileCardMargin = 10;
  static const double _mobilePadding = 10;
  static const double _mobileMargin = 20;

  static const double _desktopPadding = 20;
  static const double _desktopMargin = 20;

  final String _title;
  final String _body;
  final bool _isMobile;

  const _FeatureCard(this._title, this._body, this._isMobile);

  @override
  Widget build(BuildContext context) {
    return _isMobile ? _mobile(context) : _desktop(context);
  }

  Widget _mobile(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(vertical: _mobileCardMargin),
      child: Container(
        padding: const EdgeInsets.all(_mobilePadding),
        margin: const EdgeInsets.symmetric(vertical: _mobileMargin),
        alignment: Alignment.center,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              _title,
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: _mobilePadding),
            Text(
              _body,
              style: Theme.of(context).textTheme.bodyLarge,
            ),
          ],
        ),
      ),
    );
  }

  Widget _desktop(BuildContext context) {
    return Card(
      child: Container(
        padding: const EdgeInsets.all(_desktopPadding),
        margin: const EdgeInsets.symmetric(vertical: _desktopMargin),
        alignment: Alignment.topLeft,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              _title,
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: _desktopPadding),
            Text(
              _body,
              style: Theme.of(context).textTheme.bodyLarge,
            ),
          ],
        ),
      ),
    );
  }
}

// Shows the online cloud-based feature of OAMS.
class _CloudBasedFeature extends StatelessWidget {
  static const String _title = "Cloud-based management system";
  static const String _body =
      "Say goodbye to cumbersome paper-based attendance sheets. OAMS stores attendance data online so you can stay up-to-date wherever you are - instantly.";

  final bool _isMobile;

  const _CloudBasedFeature(this._isMobile);

  @override
  Widget build(BuildContext context) {
    return _FeatureCard(_title, _body, _isMobile);
  }
}

// Shows the reports and analytics feature of OAMS.
class _DetailedSystemDashboard extends StatelessWidget {
  static const String _title = "Detailed system dashboard";
  static const String _body =
      "Get comprehensive access to attendance data. Our system dashboard provides you with detailed visualisations of your most important analytics - all at your fingertips.";

  final bool _isMobile;

  const _DetailedSystemDashboard(this._isMobile);

  @override
  Widget build(BuildContext context) {
    return _FeatureCard(_title, _body, _isMobile);
  }
}

// The final call to the user to sign up.
class _FinalActionCall extends StatelessWidget {
  static const String _callText = "Try OAMS now.";
  static const double _buttonRadius = 5;
  static const String _buttonText = "Sign Up/Login";

  static const double _mobilePadding = 10;
  static const double _mobileMargin = 40;
  static const double _mobileButtonFontSize = 20;
  static const double _mobileButtonVerticalPadding = 15;
  static const double _mobileButtonHorizontalPadding = 10;

  static const double _desktopPadding = 20;
  static const double _desktopMargin = 100;
  static const double _desktopButtonFontSize = 24;
  static const double _desktopButtonVerticalPadding = 20;
  static const double _desktopButtonHorizontalPadding = 10;

  final bool _isMobile;

  const _FinalActionCall(this._isMobile);

  @override
  Widget build(BuildContext context) {
    return _isMobile ? _mobile(context) : _desktop(context);
  }

  Widget _mobile(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: _mobileMargin),
      child: Column(
        children: [
          Text(
            _callText,
            style: Theme.of(context).textTheme.headlineSmall,
          ),
          const SizedBox(height: _mobilePadding),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              padding: const EdgeInsets.symmetric(
                vertical: _mobileButtonVerticalPadding,
                horizontal: _mobileButtonHorizontalPadding,
              ),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(_buttonRadius),
              ),
              textStyle: const TextStyle(fontSize: _mobileButtonFontSize),
            ),
            onPressed: () => context.goNamed(Routes.login.name),
            child: const Text(_buttonText),
          )
        ],
      ),
    );
  }

  Widget _desktop(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: _desktopMargin),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            _callText,
            style: Theme.of(context).textTheme.headlineLarge,
          ),
          const SizedBox(width: _desktopPadding),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              padding: const EdgeInsets.symmetric(
                vertical: _desktopButtonVerticalPadding,
                horizontal: _desktopButtonHorizontalPadding,
              ),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(_buttonRadius),
              ),
              textStyle: const TextStyle(fontSize: _desktopButtonFontSize),
            ),
            onPressed: () => context.goNamed(Routes.login.name),
            child: const Text(_buttonText),
          ),
        ],
      ),
    );
  }
}
