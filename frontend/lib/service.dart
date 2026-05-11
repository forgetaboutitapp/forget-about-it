import '../../data/constants.dart';
import '../../interop/notification_service_stub.dart';
import '../../screens/home/service.dart';
import '../../screens/quiz/service.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:hive_ce_flutter/adapters.dart';

import 'fn/fn.dart';
import 'screens/login/mdns_lookup.dart';

@pragma('vm:entry-point')
void printHello() async {
  await Hive.initFlutter();

  final box = await Hive.openBox<dynamic>(localSettingsHiveBox);
  String? token = box.get(localSettingsHiveLoginToken);
  String? remoteHost = box.get(localSettingsHiveRemoteHost);

  if (remoteHost != null && remoteHost.contains('.local')) {
    final uri = Uri.tryParse(remoteHost);
    if (uri != null) {
      await resolveAndCacheMdnsHost(uri.host);
    }
  }

  if (token != null && remoteHost != null) {
    final tagsResult =
        await (await (await getAllTags(remoteHost, token, () => {}))
                .doFlatMap((tags) async {
      return await getNextQuestion(remoteHost, token, () => {},
          tags.$1.map((e) => e.tag).toISet(), false);
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
