import 'package:flutter/material.dart';
import 'package:frontend/router.dart';
import 'package:frontend/widgets/nav_bar/mobile.dart';
import 'package:go_router/go_router.dart';
import 'package:responsive_framework/responsive_framework.dart';

import 'desktop.dart';

class NavBar extends StatelessWidget {
  const NavBar({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    if (ResponsiveBreakpoints.of(context).smallerThan(DESKTOP)) {
      return const NavBarMobile();
    }
    return const NavBarDesktop();
  }
}

class Logo extends StatelessWidget {
  const Logo({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: double.infinity,
      width: 100,
      child: IconButton(
        icon: Image.asset('assets/logo.png'),
        onPressed: () => context.goNamed(Routes.index.name),
      ),
    );
  }
}

class NavBarText extends StatelessWidget {
  final String text;

  const NavBarText(this.text, {super.key});

  @override
  Widget build(BuildContext context) {
    return Text(
      text,
      style: const TextStyle(fontSize: 18),
    );
  }
}
