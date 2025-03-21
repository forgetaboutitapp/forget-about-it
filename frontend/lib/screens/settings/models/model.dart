import 'package:app/network/interfaces.dart';
import 'package:app/screens/settings/models/remote_settings.dart';
import 'package:app/screens/settings/services.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'model.g.dart';

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
