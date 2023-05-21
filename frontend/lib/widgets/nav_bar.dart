import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/router.dart';
import 'package:go_router/go_router.dart';
import 'package:responsive_framework/responsive_framework.dart';

import '../providers/session.dart';

// NavBar is a widget that shows the app's navigation bar.
class NavBar extends StatelessWidget {
  static const double height = 70;
  static const double padding = 10;
  static List<BoxShadow>? boxShadow = kElevationToShadow[8];

  static const double desktopBorderRadius = 30;

  const NavBar({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ResponsiveBreakpoints.of(context).smallerThan(DESKTOP)
        ? mobile(context)
        : desktop(context);
  }

  Widget mobile(BuildContext context) {
    return SizedBox(
      height: height,
      child: Container(
        padding: const EdgeInsets.all(padding),
        decoration: BoxDecoration(
          color: Theme.of(context).primaryColorLight,
          boxShadow: boxShadow,
        ),
        child: const Row(
          mainAxisSize: MainAxisSize.max,
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Options(true),
            Logo(),
          ],
        ),
      ),
    );
  }

  Widget desktop(BuildContext context) {
    return SizedBox(
      height: height,
      child: Container(
        padding: const EdgeInsets.all(padding),
        decoration: BoxDecoration(
            color: Theme.of(context).primaryColorLight,
            boxShadow: boxShadow,
            borderRadius:
                const BorderRadius.all(Radius.circular(desktopBorderRadius))),
        child: const Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Logo(),
            Options(false),
          ],
        ),
      ),
    );
  }
}

// Logo is a widget that shows the app logo on the navigation bar.
class Logo extends StatelessWidget {
  static const double width = 100;
  static const double height = double.infinity;

  static const String logoPath = "assets/logo.png";

  const Logo({Key? key}) : super(key: key);

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

// Options is a widget that stores all the navigation bar items except the logo.
class Options extends ConsumerWidget {
  final bool isMobile;

  const Options(this.isMobile, {Key? key}) : super(key: key);

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
        child: AboutItem(true),
      ),
    ];

    PopupMenuItem<Widget> newItem;
    if (isLoggedIn) {
      newItem = const PopupMenuItem(
        child: ProfileItem(true),
      );
    } else {
      newItem = const PopupMenuItem(
        child: LoginItem(true),
      );
    }

    items.add(newItem);

    return PopupMenuButton(
      icon: const Icon(Icons.menu),
      itemBuilder: (context) => items,
    );
  }

  Widget desktop(BuildContext context, bool isLoggedIn) {
    var children = <Widget>[
      const AboutItem(false),
    ];

    Widget newChild;
    if (isLoggedIn) {
      newChild = const ProfileItem(false);
    } else {
      newChild = const LoginItem(false);
    }

    children.add(newChild);

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: children,
    );
  }
}

// AboutItem is a widget that shows the About item on the navigation bar.
class AboutItem extends StatelessWidget {
  final Widget text = const NavBarText("About");
  final bool isMobile;

  const AboutItem(this.isMobile, {Key? key}) : super(key: key);

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

// ProfileItem is a widget that shows the Profile item on the navigation bar.
class ProfileItem extends StatelessWidget {
  final Widget text = const NavBarText("Profile");
  final bool isMobile;

  const ProfileItem(this.isMobile, {Key? key}) : super(key: key);

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
    return TextButton(
      child: text,
      onPressed: () => context.goNamed(Routes.profile.name),
    );
  }
}

// LoginItem is a widget that shows the Login item on the navigation bar.
class LoginItem extends StatelessWidget {
  final Widget text = const NavBarText("Login");
  final bool isMobile;

  const LoginItem(this.isMobile, {Key? key}) : super(key: key);

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

// NavBarText helps to standardise the text in the navigation bar.
class NavBarText extends StatelessWidget {
  static const double fontSize = 18;

  final String text;

  const NavBarText(this.text, {super.key});

  @override
  Widget build(BuildContext context) {
    return Text(
      text,
      style: const TextStyle(fontSize: fontSize),
    );
  }
}
