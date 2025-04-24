import 'dart:convert';

import 'package:app/screens/stats/models/stats_data.dart';

import '../../network/interfaces.dart';

Future<StatsData> getStats(FetchData remoteServer, DateTime startingDate,
        DateTime endTime) async =>
    dbg(StatsData.fromJSON(
        jsonDecode(await remoteServer.getStats(startingDate, endTime))));

T dbg<T>(T a) {
  print(a);
  return a;
}
