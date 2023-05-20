import 'package:flutter/material.dart';
import 'package:frontend/router.dart';
import 'package:go_router/go_router.dart';

import 'nav_bar.dart';

class NavBarMobile extends StatelessWidget {
  const NavBarMobile({Key? key}) : super(key: key);

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
        ),
        child: const Row(
          mainAxisSize: MainAxisSize.max,
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            _OptionsMobile(),
            Logo(),
          ],
        ),
      ),
    );
  }
}

class _OptionsMobile extends StatelessWidget {
  const _OptionsMobile({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton(
      icon: const Icon(Icons.menu),
      itemBuilder: (BuildContext context) {
        return <PopupMenuEntry<Widget>>[
          PopupMenuItem<Widget>(
            child: const NavBarText("About"),
            onTap: () => context.goNamed(Routes.about.name),
          ),
          PopupMenuItem<Widget>(
            child: const NavBarText("Login"),
            onTap: () => context.goNamed(Routes.login.name),
          ),
        ];
      },
    );
  }
}
