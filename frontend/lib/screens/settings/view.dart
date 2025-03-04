import 'package:app/screens/settings/qr_dialog.dart';
import 'package:app/state/login.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:flutter_settings_ui/flutter_settings_ui.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:http/http.dart' as http;

class SettingsScreen extends HookConsumerWidget {
  static const location = '/settings';
  final http.Client client;
  const SettingsScreen({super.key, required this.client});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final potentialToken = ref.watch(loginProvider);
    if (potentialToken is! LoggedIn) {
      Navigator.pushReplacementNamed(context, '/login');
      return Container();
    }
    final token = potentialToken.id;
    var remoteDevices = useState([('Real Root', '1/1/2020', '')].toIList());
    return Scaffold(
      appBar: AppBar(
        title: Text('Settings'),
      ),
      body: SettingsList(
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
                initialValue: true,
                onToggle: null,
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
                value:
                    Text('This is the FSRS algorithm currently in use by Anki'),
              ),
              SettingsTile(
                title: Text('Add Algorithm'),
                trailing: Icon(Icons.add),
              )
            ],
          ),
          SettingsSection(
            title: Text('Other Devices'),
            tiles: [
              ...remoteDevices.value.map(
                (e) => SettingsTile(
                  title: Text(e.$1),
                  value: Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.fromLTRB(0, 0, 16.0, 0),
                        child: Text(e.$2),
                      ),
                      Text(e.$3)
                    ],
                  ),
                  trailing: Icon(
                    Icons.delete,
                  ),
                ),
              ),
              SettingsTile(
                title: Text('Connect Other Devices'),
                trailing: Icon(Icons.add),
                onPressed: (context) async {
                  final p = await showDialog(
                    context: context,
                    builder: (context) => QRDialog(
                        token: token,
                        client: client,
                        remoteHost: 'http://127.0.0.1:8080'),
                  );
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
}
