import 'package:flutter/material.dart';
import 'package:frontend/screens/about_screen.dart';
import 'package:frontend/screens/home_screen.dart';
import 'package:frontend/screens/login_screen.dart';
import 'package:go_router/go_router.dart';

typedef RouteInfo = ({String name, String path});

class Routes {
  static const RouteInfo index = (name: "index", path: "/");
  static const RouteInfo about = (name: "about", path: "about");
  static const RouteInfo login = (name: "login", path: "login");
}

final router = GoRouter(
  routes: [
    _noTransitionRoute(
      Routes.index,
      page: const HomeScreen(),
      childRoutes: [
        _noTransitionRoute(Routes.about, page: const AboutScreen()),
        _noTransitionRoute(Routes.login, builder: (context, state) {
          final queryParams = state.queryParameters;
          return LoginScreen(queryParams: queryParams);
        }),
      ],
    ),
  ],
);

// _noTransitionRoute is a helper function to build a transition-less route.
// Be sure to provide at least one of page or builder.
// If both are provided, then it defaults to using page.
GoRoute _noTransitionRoute(RouteInfo info,
    {Widget? page,
    Widget Function(BuildContext, GoRouterState)? builder,
    List<GoRoute> childRoutes = const []}) {
  return GoRoute(
    name: info.name,
    path: info.path,
    builder: builder,
    pageBuilder: (context, state) => NoTransitionPage<void>(
      key: state.pageKey,
      child: page ?? builder!(context, state),
    ),
    routes: childRoutes,
  );
}
