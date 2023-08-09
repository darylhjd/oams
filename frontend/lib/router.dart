import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/client.dart';
import 'package:frontend/env/env.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/screens/about_screen.dart';
import 'package:frontend/screens/home_screen.dart';
import 'package:frontend/screens/login_screen.dart';
import 'package:frontend/screens/profile_screen.dart';
import 'package:go_router/go_router.dart';

typedef RouteInfo = ({String name, String path});

class Routes {
  static const RouteInfo index = (name: "index", path: "/");
  static const RouteInfo about = (name: "about", path: "about");
  static const RouteInfo login = (name: "login", path: "login");
  static const RouteInfo profile = (name: "profile", path: "profile");
}

final routerProvider = Provider<GoRouter>((ref) => _indexRoute(ref));

GoRouter _indexRoute(ProviderRef<GoRouter> ref) {
  return GoRouter(
    routes: [
      GoRoute(
        name: Routes.index.name,
        path: Routes.index.path,
        pageBuilder: (context, state) => NoTransitionPage(
          key: state.pageKey,
          child: const HomeScreen(),
        ),
        routes: [
          _aboutRoute(),
          _loginRoute(ref),
          _profileRoute(ref),
        ],
      ),
    ],
  );
}

GoRoute _aboutRoute() {
  return GoRoute(
    name: Routes.about.name,
    path: Routes.about.path,
    pageBuilder: (context, state) => NoTransitionPage(
      key: state.pageKey,
      child: const AboutScreen(),
    ),
  );
}

GoRoute _loginRoute(ProviderRef<GoRouter> ref) {
  final isLoggedIn = ref.watch(sessionProvider);

  return GoRoute(
    name: Routes.login.name,
    path: Routes.login.path,
    pageBuilder: (context, state) {
      return NoTransitionPage(
        key: state.pageKey,
        child: LoginScreen(
          redirectUrl:
              state.uri.queryParameters[APIClient.loginRedirectUrlParam],
        ),
      );
    },
    redirect: (context, state) =>
        isLoggedIn ? state.namedLocation(Routes.index.name) : null,
  );
}

GoRoute _profileRoute(ProviderRef<GoRouter> ref) {
  final isLoggedIn = ref.watch(sessionProvider);

  return GoRoute(
    name: Routes.profile.name,
    path: Routes.profile.path,
    pageBuilder: (context, state) {
      return NoTransitionPage(
        key: state.pageKey,
        child: const ProfileScreen(),
      );
    },
    redirect: (context, state) {
      var to = state.namedLocation(
        Routes.login.name,
        queryParameters: {
          APIClient.loginRedirectUrlParam:
              "${webServerHost()}:${webServerPort()}${state.fullPath!}",
        },
      );
      return isLoggedIn ? null : to;
    },
  );
}
