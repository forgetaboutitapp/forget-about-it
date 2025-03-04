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
  const MainApp({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ValueListenableBuilder(
        valueListenable: Hive.box(localSettingsHiveBox).listenable(),
        builder: (context, settingsBox, w) {
          final router = GoRouter(
            redirect: (BuildContext context, GoRouterState state) async {
              if (settingsBox.get(localSettingsHiveLoginToken) == null) {
                return LoginScreen.location;
              } else {
                return state.fullPath;
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
                  token: settingsBox.get(localSettingsHiveLoginToken),
                  remoteHost: settingsBox.get(localSettingsHiveRemoteHost),
                ),
              ),
            ],
          );
          return MaterialApp.router(
            debugShowCheckedModeBanner: false,
            routerConfig: router,
          );
        });
  }
}
