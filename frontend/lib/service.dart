import 'package:app/data/constants.dart';
import 'package:app/interop/notification_service_stub.dart';
import 'package:app/network/interfaces.dart';
import 'package:app/screens/home/service.dart';
import 'package:app/screens/quiz/service.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:hive_ce_flutter/adapters.dart';
import 'package:http/http.dart' as http;

import 'network/network.dart';

@pragma('vm:entry-point')
void printHello() async {
  await Hive.initFlutter();

  final box = await Hive.openBox<dynamic>(localSettingsHiveBox);
  String? token = box.get(localSettingsHiveLoginToken);
  String? remoteHost = box.get(localSettingsHiveRemoteHost);

  if (token != null && remoteHost != null) {
    FetchData f = RemoteServer(
        remoteHost: remoteHost, token: token, client: http.Client());
    try {
      final tags = (await getAllTags(f)).$1;
      final q = await getNextQuestion(f, tags.map((e) => e.tag).toISet());
      if (q.dueCards > 0) {
        await NotificationService().showNotification(
            title: 'Flashcard reminder',
            body: 'There are ${q.dueCards} cards that are due!');
      }
    } on Exception catch (e) {
      await NotificationService().showNotification(
        title: 'Error getting data',
        body: 'Error: $e',
      );
    }
  }
}
