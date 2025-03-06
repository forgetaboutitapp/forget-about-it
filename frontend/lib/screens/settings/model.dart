import 'package:app/screens/settings/services.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:http/http.dart' as http;
import 'package:riverpod_annotation/riverpod_annotation.dart';
part 'model.freezed.dart';
part 'model.g.dart';

@freezed
class RemoteSettings with _$RemoteSettings {
  const factory RemoteSettings({
    required IList<RemoteDevice>? remoteDevices,
  }) = _RemoteSettings;

  factory RemoteSettings.fromJSON(dynamic json) {
    if (json['remote-devices'] == null) {
      return RemoteSettings(remoteDevices: null);
    }
    final List<RemoteDevice> remoteList = json['remote-devices']
        .map<RemoteDevice>((e) => RemoteDevice.fromJSON(e))
        .toList();
    return RemoteSettings(remoteDevices: remoteList.toIList());
  }
}

@freezed
class RemoteDevice with _$RemoteDevice {
  const factory RemoteDevice({
    required String title,
    required DateTime dateAdded,
    required DateTime? lastUsed,
  }) = _RemoteDevice;
  factory RemoteDevice.fromJSON(dynamic json) {
    final lastUsed = json['last-used'];
    return RemoteDevice(
      title: json['title'],
      dateAdded: DateTime.fromMillisecondsSinceEpoch(
        json['date-added'] * 1000,
      ),
      lastUsed: (lastUsed == null)
          ? null
          : DateTime.fromMillisecondsSinceEpoch(
              lastUsed * 1000,
            ),
    );
  }
}

@riverpod
class RemoteSettingsNotifier extends _$RemoteSettingsNotifier {
  bool _init = false;
  @override
  RemoteSettings? build() {
    if (_init) {
      return state;
    } else {
      return null;
    }
  }

  Future<Exception?> getData(
    String remoteHost,
    String token,
    http.Client client,
  ) async {
    try {
      final newState = await getRemoteSettings(remoteHost, token, client);

      if (!_init || state != newState) {
        state = newState;
        _init = true;
      }
    } on Exception catch (e) {
      return e;
    }
    return null;
  }
}
