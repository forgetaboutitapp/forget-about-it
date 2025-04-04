import 'dart:async';
import 'dart:convert';

import 'package:app/network/interfaces.dart';
import 'package:app/screens/general-display/show_error.dart';
import 'package:app/screens/settings/add_algorithm.dart';
import 'package:app/screens/settings/models/model.dart';
import 'package:app/screens/settings/qr_dialog.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:flutter_settings_ui/flutter_settings_ui.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:intl/intl.dart';

class SettingsScreen extends HookConsumerWidget {
  static const location = '/settings';
  final bool curDarkMode;
  final FetchData remoteServer;
  final GenericFilepicker filepicker;
  final FutureOr<void> Function(bool) switchDarkMode;
  const SettingsScreen({
    super.key,
    required this.curDarkMode,
    required this.remoteServer,
    required this.switchDarkMode,
    required this.filepicker,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final localCurDarkMode = useState(curDarkMode);
    final stillWaitingForRemoteServer = useState(true);

    final remoteSettings = ref.watch(remoteSettingsNotifierProvider);
    useEffect(() {
      refresh(stillWaitingForRemoteServer, ref, context);
      return null;
    }, []);
    return Scaffold(
      appBar: AppBar(
        title: Text('Settings'),
      ),
      body: stillWaitingForRemoteServer.value
          ? Center(child: CircularProgressIndicator())
          : SettingsList(
              sections: [
                SettingsSection(
                  title: Text('Local Setting'),
                  tiles: [
                    SettingsTile.switchTile(
                      initialValue: false,
                      onToggle: null,
                      title: Text('Enable Notifications'),
                    ),
                    SettingsTile.navigation(
                      enabled: false,
                      title: Text('Daily Notification Time'),
                      value: null,
                      onPressed: null,
                    ),
                    SettingsTile.switchTile(
                      initialValue: localCurDarkMode.value,
                      onToggle: (v) {
                        localCurDarkMode.value = v;
                        switchDarkMode(v);
                      },
                      title: Text('Dark Theme'),
                    ),
                    SettingsTile.switchTile(
                      initialValue: true,
                      onToggle: (v) => {},
                      title: Text('Cutesy Content'),
                    ),
                  ],
                ),
                SettingsSection(
                  title: Text('General Setting'),
                  tiles: [
                    ...?remoteSettings?.remoteAlgorithms?.map(
                      (e) => SettingsTile(
                        title: Text(e.algorithmName),
                        description: Text(
                            'Author: ${e.authorName}\nLicense: ${e.license}\nVersion: ${e.version}\nAdded: ${e.timeAdded}'),
                        leading: Radio(
                          value: e.algorithmID,
                          groupValue: remoteSettings.defaultAlgorithm ?? 0,
                          onChanged: (int? v) async {
                            stillWaitingForRemoteServer.value = true;
                            try {
                              await remoteServer
                                  .setDefaultAlgorithm(e.algorithmID);
                            } on Exception catch (e) {
                              if (context.mounted) {
                                showError(context, e.toString());
                              }
                            }
                            if (context.mounted) {
                              refresh(
                                  stillWaitingForRemoteServer, ref, context);
                            }
                          },
                        ),
                        trailing: IconButton(
                          icon: Icon(Icons.delete),
                          onPressed:
                              remoteSettings.defaultAlgorithm == e.algorithmID
                                  ? null
                                  : () async {
                                      try {
                                        final res = await remoteServer
                                            .removeAlgorithm(e.algorithmName);
                                        if (res != null && res != '') {
                                          final resJson = jsonDecode(res);
                                          resJson.keys.contains('err');
                                          if (context.mounted) {
                                            showError(context,
                                                'Error: ${resJson['err']}');
                                          }
                                        }
                                        if (context.mounted) {
                                          refresh(stillWaitingForRemoteServer,
                                              ref, context);
                                        }
                                      } on Exception catch (e) {
                                        if (context.mounted) {
                                          showError(context, e.toString());
                                        }
                                      }
                                    },
                        ),
                      ),
                    ),
                    SettingsTile(
                      title: Text('Add Algorithm'),
                      trailing: Icon(Icons.add),
                      onPressed: (context) async {
                        final needsRefresh = await showDialog(
                          context: context,
                          builder: (context) => AddAlgorithm(
                              algos: remoteSettings?.remoteAlgorithms
                                  ?.toSet()
                                  .toISet(),
                              remoteServer: remoteServer,
                              filepicker: filepicker),
                        );
                        if (context.mounted && needsRefresh == true) {
                          refresh(stillWaitingForRemoteServer, ref, context);
                        }
                      },
                    )
                  ],
                ),
                remoteSettings == null
                    ? SettingsSection(
                        title: Text('Other Devices'),
                        tiles: [
                          SettingsTile(
                            title: Text('Cannot reach remote server... '),
                          ),
                        ],
                      )
                    : SettingsSection(
                        title: Text('Other Devices'),
                        tiles: [
                          ...?remoteSettings.remoteDevices?.indexed.map(
                            (tuple) {
                              final e = tuple.$2;
                              final i = tuple.$1;
                              final lastUsed = e.lastUsed;
                              return SettingsTile(
                                title: Text(e.title),
                                value: Row(
                                  children: [
                                    Padding(
                                      padding: const EdgeInsets.fromLTRB(
                                          0, 0, 16.0, 0),
                                      child: Text(
                                          'Date Added: ${DateFormat.yMd().format(e.dateAdded)}'),
                                    ),
                                    if (lastUsed != null)
                                      Text(
                                          'Last Used: ${DateFormat.yMd().format(lastUsed)}'),
                                  ],
                                ),
                                trailing: IconButton(
                                  icon: Icon(Icons.delete),
                                  onPressed: i == 0
                                      ? null
                                      : () async {
                                          try {
                                            final res = await remoteServer
                                                .removeLogin(e.loginId);
                                            if (res != null && res != '') {
                                              final resJson = jsonDecode(res);
                                              resJson.keys.contains('err');
                                              if (context.mounted) {
                                                showError(context,
                                                    'Error: ${resJson['err']}');
                                              }
                                            }
                                            if (context.mounted) {
                                              refresh(
                                                  stillWaitingForRemoteServer,
                                                  ref,
                                                  context);
                                            }
                                          } on Exception catch (e) {
                                            if (context.mounted) {
                                              showError(context, e.toString());
                                            }
                                          }
                                        },
                                ),
                              );
                            },
                          ),
                          SettingsTile(
                            title: Text('Connect Other Devices'),
                            trailing: Icon(Icons.add),
                            onPressed: (context) async {
                              await showDialog(
                                context: context,
                                builder: (context) =>
                                    QRDialog(remoteServer: remoteServer),
                              );
                              if (context.mounted) {
                                refresh(
                                    stillWaitingForRemoteServer, ref, context);
                              }
                            },
                          ),
                        ],
                      ),
                SettingsSection(
                  title: Text('Licensess'),
                  tiles: [
                    SettingsTile(
                      title: Text('Show Licenses'),
                      onPressed: (context) => showLicensePage(
                          applicationLegalese: 'AGPL-V3',
                          context: context,
                          applicationName: 'Forget About It'),
                    )
                  ],
                ),
              ],
            ),
    );
  }

  void refresh(ValueNotifier<bool> stillWaitingForRemoteServer, WidgetRef ref,
      BuildContext context) {
    () async {
      if (stillWaitingForRemoteServer.value != true) {
        stillWaitingForRemoteServer.value = true;
      }
      final p = await ref
          .read(remoteSettingsNotifierProvider.notifier)
          .getData(remoteServer);
      stillWaitingForRemoteServer.value = false;
      if (p != null) {
        stillWaitingForRemoteServer.value = false;
        if (context.mounted) {
          ScaffoldMessenger.of(context)
              .showSnackBar(SnackBar(content: Text(p.toString())));
        }
      }
    }();
  }
}
