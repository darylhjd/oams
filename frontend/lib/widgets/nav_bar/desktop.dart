import 'package:flutter/material.dart';
import 'package:frontend/router.dart';
import 'package:go_router/go_router.dart';

import 'nav_bar.dart';

class NavBarDesktop extends StatelessWidget {
  const NavBarDesktop({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    var theme = Theme.of(context);

    return SizedBox(
      height: 70,
      child: Container(
        padding: const EdgeInsets.all(10),
        decoration: BoxDecoration(
          color: theme.primaryColorLight,
          boxShadow: kElevationToShadow[8],
          borderRadius: const BorderRadius.all(Radius.circular(30)),
        ),
        child: const Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Logo(),
            _OptionsDesktop(),
          ],
        ),
      ),
    );
  }
}

class _OptionsDesktop extends StatelessWidget {
  const _OptionsDesktop({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        TextButton(
          child: const NavBarText("About"),
          onPressed: () => context.goNamed(Routes.about.name),
        ),
        TextButton(
          child: const NavBarText("Login"),
          onPressed: () => context.goNamed(Routes.login.name),
        ),
      ],
    );
  }
}
