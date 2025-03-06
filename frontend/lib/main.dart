import 'package:app/data/constants.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hive_ce_flutter/adapters.dart';

import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:http/http.dart' as http;
import 'screens/login/view.dart';
import 'screens/home/view.dart';
import 'screens/settings/view.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Hive.initFlutter();
  await Hive.openBox<dynamic>(localSettingsHiveBox);
  runApp(
    ProviderScope(
      child: MainApp(),
    ),
  );
}

class MainApp extends HookConsumerWidget {
  MainApp({super.key});
  final router = GoRouter(
    refreshListenable: Hive.box(localSettingsHiveBox).listenable(),
    redirect: (BuildContext context, GoRouterState state) async {
      final settingsBox = Hive.box(localSettingsHiveBox);
      if (settingsBox.get(localSettingsHiveLoginToken) == null) {
        return LoginScreen.location;
      } else if (state.fullPath != LoginScreen.location) {
        return state.fullPath;
      } else {
        return HomeScreen.location;
      }
    },
    routes: [
      GoRoute(
        path: HomeScreen.location,
        builder: (context, state) => HomeScreen(),
      ),
      GoRoute(
        path: LoginScreen.location,
        builder: (context, state) => LoginScreen(client: http.Client()),
      ),
      GoRoute(
        path: SettingsScreen.location,
        builder: (context, state) => SettingsScreen(
          client: http.Client(),
          token:
              Hive.box(localSettingsHiveBox).get(localSettingsHiveLoginToken),
          remoteHost:
              Hive.box(localSettingsHiveBox).get(localSettingsHiveRemoteHost),
          curDarkMode: Hive.box(localSettingsHiveBox)
              .get(localSettingsHiveDarkTheme, defaultValue: false),
          switchDarkMode: (_) async {
            final v = Hive.box(localSettingsHiveBox)
                .get(localSettingsHiveDarkTheme, defaultValue: false);
            await Hive.box(localSettingsHiveBox)
                .put(localSettingsHiveDarkTheme, !v);
          },
        ),
      ),
    ],
  );
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ValueListenableBuilder(
        valueListenable: Hive.box(localSettingsHiveBox).listenable(),
        builder: (context, settingsBox, w) {
          return MaterialApp.router(
            debugShowCheckedModeBanner: false,
            routerConfig: router,
            theme: settingsBox.get(localSettingsHiveDarkTheme) == true
                ? ThemeData.dark(useMaterial3: true)
                : ThemeData.light(useMaterial3: true),
          );
        });
  }
}
