import 'dart:async';

import 'package:app/network/interfaces.dart';
import 'package:app/screens/settings/model.dart';
import 'package:app/screens/settings/qr_dialog.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:flutter_settings_ui/flutter_settings_ui.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:intl/intl.dart';

class SettingsScreen extends HookConsumerWidget {
  static const location = '/settings';
  final bool curDarkMode;
  final FetchData remoteServer;
  final FutureOr<void> Function(bool) switchDarkMode;
  const SettingsScreen(
      {super.key,
      required this.curDarkMode,
      required this.remoteServer,
      required this.switchDarkMode,
      requird});

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
                    SettingsTile(
                      title: Text('Remote Server'),
                      value: Text('None'),
                    ),
                    SettingsTile(
                      title: Text('Algorithm In Use'),
                      value: Text('FSRS'),
                    ),
                    SettingsTile(
                      title: Text('FSRS'),
                      value: Text(
                          'This is the FSRS algorithm currently in use by Anki'),
                    ),
                    SettingsTile(
                      title: Text('Add Algorithm'),
                      trailing: Icon(Icons.add),
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
                          ...?remoteSettings.remoteDevices?.map(
                            (e) {
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
                                trailing: Icon(
                                  Icons.delete,
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
      stillWaitingForRemoteServer.value = true;
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
