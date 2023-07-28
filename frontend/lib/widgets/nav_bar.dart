import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/client.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/router.dart';
import 'package:go_router/go_router.dart';
import 'package:responsive_framework/responsive_framework.dart';

// NavBar is a widget that shows the app's navigation bar.
class NavBar extends StatelessWidget {
  static const double height = 70;
  static const double padding = 10;
  static List<BoxShadow>? boxShadow = kElevationToShadow[8];
  static const double borderRadius = 10;

  const NavBar({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: height,
      child: Container(
        padding: const EdgeInsets.all(padding),
        decoration: BoxDecoration(
          color: Theme.of(context).primaryColorLight,
          boxShadow: boxShadow,
          borderRadius: const BorderRadius.all(
            Radius.circular(borderRadius),
          ),
        ),
        child: ResponsiveBreakpoints.of(context).isMobile
            ? mobile(context)
            : desktop(context),
      ),
    );
  }

  Widget mobile(BuildContext context) {
    return const Row(
      mainAxisSize: MainAxisSize.max,
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        _Options(true),
        _Logo(),
      ],
    );
  }

  Widget desktop(BuildContext context) {
    return const Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        _Logo(),
        _Options(false),
      ],
    );
  }
}

// _Logo is a widget that shows the app logo on the navigation bar.
class _Logo extends StatelessWidget {
  static const double width = 100;
  static const double height = double.infinity;

  static const String logoPath = "assets/logo.png";

  const _Logo({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: width,
      height: height,
      child: IconButton(
        icon: Image.asset(logoPath),
        onPressed: () => context.goNamed(Routes.index.name),
      ),
    );
  }
}

// _Options is a widget that stores all the navigation bar items except the logo.
class _Options extends ConsumerWidget {
  final bool isMobile;

  const _Options(this.isMobile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isLoggedIn = ref.read(sessionProvider);
    return isMobile
        ? mobile(context, isLoggedIn)
        : desktop(context, isLoggedIn);
  }

  Widget mobile(BuildContext context, bool isLoggedIn) {
    var items = <PopupMenuItem<Widget>>[
      const PopupMenuItem(
        padding: EdgeInsets.zero,
        child: _AboutItem(true),
      ),
    ];

    if (isLoggedIn) {
      items.addAll(
        [
          const PopupMenuItem(
            padding: EdgeInsets.zero,
            child: _ProfileItem(true),
          ),
          const PopupMenuItem(
            padding: EdgeInsets.zero,
            child: _LogoutItem(true),
          )
        ],
      );
    } else {
      items.add(
        const PopupMenuItem(
          padding: EdgeInsets.zero,
          child: _LoginItem(true),
        ),
      );
    }

    return PopupMenuButton(
      icon: const Icon(Icons.menu),
      itemBuilder: (context) => items,
    );
  }

  Widget desktop(BuildContext context, bool isLoggedIn) {
    var children = <Widget>[
      const _AboutItem(false),
    ];

    Widget newChild;
    if (isLoggedIn) {
      newChild = const _ProfileItem(false);
    } else {
      newChild = const _LoginItem(false);
    }

    children.add(newChild);

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: children,
    );
  }
}

// _AboutItem is a widget that shows the About item on the navigation bar.
class _AboutItem extends StatelessWidget {
  final Widget text = const _NavBarText("About");
  final bool isMobile;

  const _AboutItem(this.isMobile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return isMobile ? mobile(context) : desktop(context);
  }

  Widget mobile(BuildContext context) {
    return PopupMenuItem<Widget>(
      onTap: () => context.goNamed(Routes.about.name),
      child: text,
    );
  }

  Widget desktop(BuildContext context) {
    return TextButton(
      onPressed: () => context.goNamed(Routes.about.name),
      child: text,
    );
  }
}

// _ProfileItem is a widget that shows the Profile item on the navigation bar.
class _ProfileItem extends StatelessWidget {
  final Widget text = const _NavBarText("Profile");
  final bool isMobile;

  const _ProfileItem(this.isMobile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return isMobile ? mobile(context) : desktop(context);
  }

  Widget mobile(BuildContext context) {
    return PopupMenuItem<Widget>(
      child: text,
      onTap: () => context.goNamed(Routes.profile.name),
    );
  }

  Widget desktop(BuildContext context) {
    var items = <PopupMenuItem<Widget>>[
      PopupMenuItem(
        onTap: () => context.goNamed(Routes.profile.name),
        child: text,
      ),
      const PopupMenuItem(
        padding: EdgeInsets.zero,
        child: _LogoutItem(false),
      ),
    ];

    return PopupMenuButton(
      icon: const Icon(Icons.account_circle_rounded),
      itemBuilder: (context) => items,
    );
  }
}

// _LoginItem is a widget that shows the Login item on the navigation bar.
class _LoginItem extends StatelessWidget {
  final Widget text = const _NavBarText("Login");
  final bool isMobile;

  const _LoginItem(this.isMobile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return isMobile ? mobile(context) : desktop(context);
  }

  Widget mobile(BuildContext context) {
    return PopupMenuItem<Widget>(
      child: text,
      onTap: () => context.goNamed(Routes.login.name),
    );
  }

  Widget desktop(BuildContext context) {
    return TextButton(
      child: text,
      onPressed: () => context.goNamed(Routes.login.name),
    );
  }
}

// _LogoutItem is a widget that allows user to logout on click.
class _LogoutItem extends ConsumerWidget {
  static const String text = "Logout";
  final bool isMobile;

  const _LogoutItem(this.isMobile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return PopupMenuItem(
      onTap: () => logoutAction(context, ref),
      child: Text.rich(
        TextSpan(
          children: [
            const WidgetSpan(
              child: Icon(
                Icons.logout,
                color: Colors.red,
              ),
            ),
            WidgetSpan(
              child: Text(
                text,
                style: _navBarTextStyle(context)?.copyWith(color: Colors.red),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> logoutAction(BuildContext context, WidgetRef ref) async {
    final ok = await APIClient.logout();
    if (ok && context.mounted) {
      ref.read(sessionProvider.notifier).update((_) => false);
      ref.invalidate(sessionUserProvider);
      context.goNamed(Routes.index.name);
    }
  }
}

// _NavBarText helps to standardise the text in the navigation bar.
class _NavBarText extends StatelessWidget {
  final String text;

  const _NavBarText(this.text, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Text(
      text,
      style: _navBarTextStyle(context),
    );
  }
}

// _navBarTextStyle is a helper function to return the default text style
// for text in teh navigation bar.
TextStyle? _navBarTextStyle(BuildContext context) {
  return Theme.of(context).textTheme.titleMedium;
}
