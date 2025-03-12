import 'package:app/network/interfaces.dart';
import 'package:app/screens/settings/services.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'model.freezed.dart';
part 'model.g.dart';

@freezed
class RemoteSettings with _$RemoteSettings {
  const RemoteSettings({
    required this.remoteDevices,
  });
  @override
  final IList<RemoteDevice>? remoteDevices;

  static RemoteSettings fromJSON(dynamic json) {
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
  @override
  final DateTime dateAdded;
  @override
  final DateTime? lastUsed;
  @override
  final String title;
  const RemoteDevice({
    required this.title,
    required this.dateAdded,
    required this.lastUsed,
  });
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

  Future<Exception?> getData(FetchData remoteServer) async {
    try {
      final newState = await getRemoteSettings(remoteServer);

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
