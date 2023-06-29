import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_web_plugins/flutter_web_plugins.dart';
import 'package:frontend/env/env.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/router.dart';
import 'package:go_router/go_router.dart';
import 'package:responsive_framework/breakpoint.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

// main is the entry point to the app.
void main() {
  usePathUrlStrategy();
  checkEnvVars();
  runApp(
    const ProviderScope(
      child: OAMSFrontend(),
    ),
  );
}

// OAMSFrontend is the MaterialApp for OAM's frontend.
class OAMSFrontend extends ConsumerWidget {
  static const String title = "OAMS";

  static const double mobileMin = 0;
  static const double mobileMax = 800;
  static const double desktopMin = mobileMax + 1;
  static const double desktopMax = double.infinity;

  const OAMSFrontend({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return FutureBuilder(
      future: ref.watch(userInfoProvider.future),
      builder: (context, snapshot) {
        return snapshot.connectionState == ConnectionState.done
            ? _materialApp(ref.watch(routerProvider))
            : _loading();
      },
    );
  }

  Widget _materialApp(GoRouter router) {
    return MaterialApp.router(
      theme: ThemeData.light(useMaterial3: true),
      title: title,
      routerConfig: router,
      builder: (context, child) {
        return ResponsiveBreakpoints.builder(
          child: child!,
          breakpoints: [
            const Breakpoint(start: mobileMin, end: mobileMax, name: MOBILE),
            const Breakpoint(start: desktopMin, end: desktopMax, name: DESKTOP),
          ],
        );
      },
    );
  }

  Widget _loading() {
    return const Center(
      child: CircularProgressIndicator(),
    );
  }
}
