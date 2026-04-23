import '../../../screens/settings/models/remote_settings.dart';
import '../../../screens/settings/services.dart';
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

  Future<Exception?> getData(
      String remoteServer, String token, Function logOut) async {
    final newState = await getRemoteSettings(remoteServer, token, logOut);
    Exception? returnVal;
    newState.match(onErr: (err) {
      returnVal = err;
    }, onOk: (ok) {
      if (!_init || state != ok) {
        state = ok;
        _init = true;
      }
    });

    return returnVal;
  }
}
