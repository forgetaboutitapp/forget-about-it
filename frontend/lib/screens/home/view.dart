import 'package:app/network/interfaces.dart';
import 'package:app/screens/bulk-edit/view.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../settings/view.dart';

class HomeScreen extends HookConsumerWidget {
  static String location = '/';
  final FetchData remoteServer;
  const HomeScreen({super.key, required this.remoteServer});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Forget About It'),
        actions: [
          PopupMenuButton(
            // https://stackoverflow.com/a/75375750 for the strategy to use a function as value
            onSelected: (fn) => fn(),
            itemBuilder: (context) => [
              PopupMenuItem(
                value: () => context.go(BulkEditScreen.location),
                child: Text(
                  'Bulk Edit',
                ),
              ),
              PopupMenuItem(
                value: () => context.go(SettingsScreen.location),
                child: Text(
                  'Settings',
                ),
              ),
            ],
          )
        ],
      ),
    );
  }
}
