import '../../data/constants.dart';
import '../../interop/notification_service_stub.dart';
import '../../network/interfaces.dart';
import '../../screens/home/service.dart';
import '../../screens/quiz/service.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:hive_ce_flutter/adapters.dart';
import 'package:http/http.dart' as http;

import 'fn/fn.dart';
import 'network/network.dart';

@pragma('vm:entry-point')
void printHello() async {
  await Hive.initFlutter();

  final box = await Hive.openBox<dynamic>(localSettingsHiveBox);
  String? token = box.get(localSettingsHiveLoginToken);
  String? remoteHost = box.get(localSettingsHiveRemoteHost);

  if (token != null && remoteHost != null) {
    FetchDataWithToken f = RemoteServer(
        remoteHost: remoteHost, token: token, client: http.Client());

    final tagsResult =
        await (await (await getAllTags(f)).doFlatMap((tags) async {
      return await getNextQuestion(f, tags.$1.map((e) => e.tag).toISet());
    }))
            .doFlatMap((f) async {
      if (f.dueCards > 0) {
        await NotificationService().showNotification(
            title: 'Flashcard reminder',
            body: 'There are ${f.dueCards} cards that are due!');
      }
      return Ok(());
    });
    tagsResult.doMatch(
        onOk: (_) async {},
        onErr: (e) async {
          await NotificationService().showNotification(
            title: 'Error getting data',
            body: 'Error: $e',
          );
        });
  }
}
