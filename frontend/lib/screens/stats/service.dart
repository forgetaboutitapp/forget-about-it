import 'dart:developer' as developer;

import '../../fn/fn.dart';
import '../../protobufs-build/client_to_server.pb.dart' as client_to_server;
import '../../protobufs-build/google/protobuf/timestamp.pb.dart';
import '../../screens/stats/models/stats_data.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:fixnum/fixnum.dart';

import '../../network/interfaces.dart';

Future<Result<StatsData>> getStats(FetchDataWithToken remoteServer,
    DateTime startingDate, DateTime endTime) async {
  final tzOffset = Int64(DateTime.now().timeZoneOffset.inSeconds);
  developer.log(
      'tzOffset: $tzOffset Timestamp.now: ${DateTime.now()} Timestamp.now in utc: ${DateTime.now().toUtc()} Timestamp.now in local: ${DateTime.now().toLocal()}');
  developer.log(
    'tzOffset: $tzOffset start time: ${Timestamp.fromDateTime(toDay(startingDate).toUtc())}, ${Timestamp.fromDateTime(toDay(startingDate).toLocal())}, ${Timestamp.fromDateTime(toDay(startingDate).toUtc()).seconds - (Timestamp.fromDateTime(toDay(startingDate).toLocal())).seconds}',
  );
  return (await remoteServer.getStats(
    client_to_server.GetStats(
      tzOffset: tzOffset,
      startTime: Timestamp.fromDateTime(toDay(startingDate)),
      endTime: Timestamp.fromDateTime(
        toDay(endTime),
      ),
    ),
  ))
      .map((origVal) {
    developer.log('origVal: $origVal');
    return StatsData(
      heatmapData: origVal.pastUsage
          .map(
            (timestamp, count) => MapEntry(
              DateTime.fromMillisecondsSinceEpoch(
                (timestamp * 1000).toInt(),
              ).toLocal(),
              count,
            ),
          )
          .toIMap(),
      pastResults: origVal.pastResults
          .map(
            (timestamp, count) => MapEntry(
              DateTime.fromMillisecondsSinceEpoch(
                (timestamp * 1000).toInt(),
              ).toLocal(),
              (count * 100).toInt(),
            ),
          )
          .toIMap(),
      futureUsage: origVal.futureUsage
          .map(
            (timestamp, count) => MapEntry(
              DateTime.fromMillisecondsSinceEpoch(
                (timestamp * 1000).toInt(),
              ).toLocal(),
              (count * 100).toInt(),
            ),
          )
          .toIMap(),
    );
  });
}
