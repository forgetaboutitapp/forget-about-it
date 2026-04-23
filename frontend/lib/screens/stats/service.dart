import 'dart:developer' as developer;

import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';
import 'package:grpc/grpc_web.dart';

import '../../fn/fn.dart';
import '../../protobufs-build/google/protobuf/timestamp.pb.dart';
import '../../screens/stats/models/stats_data.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:fixnum/fixnum.dart';

Future<Result<StatsData>> getStats(String remoteHost, String token,
    Function logOut, DateTime startingDate, DateTime endTime) async {
  final tzOffset = Int64(DateTime.now().timeZoneOffset.inSeconds);
  developer.log(
      'tzOffset: $tzOffset Timestamp.now: ${DateTime.now()} Timestamp.now in utc: ${DateTime.now().toUtc()} Timestamp.now in local: ${DateTime.now().toLocal()}');
  developer.log(
    'tzOffset: $tzOffset start time: ${Timestamp.fromDateTime(toDay(startingDate).toUtc())}, ${Timestamp.fromDateTime(toDay(startingDate).toLocal())}, ${Timestamp.fromDateTime(toDay(startingDate).toUtc()).seconds - (Timestamp.fromDateTime(toDay(startingDate).toLocal())).seconds}',
  );
  final client = await ForgetAboutItServiceClient(
          GrpcWebClientChannel.xhr(Uri.parse(remoteHost)))
      .getStats(GetStatsRequest(token: token));
  if (client.hasError()) {
    return Err(Exception(client.error));
  }

  if (client.hasOk()) {
    return Ok(StatsData(
      heatmapData: client.ok.pastUsage
          .map(
            (timestamp, count) => MapEntry(
              DateTime.fromMillisecondsSinceEpoch(
                (timestamp * 1000).toInt(),
              ).toLocal(),
              count,
            ),
          )
          .toIMap(),
      pastResults: client.ok.pastResults
          .map(
            (timestamp, count) => MapEntry(
              DateTime.fromMillisecondsSinceEpoch(
                (timestamp * 1000).toInt(),
              ).toLocal(),
              (count * 100).toInt(),
            ),
          )
          .toIMap(),
      futureUsage: client.ok.futureUsage
          .map(
            (timestamp, count) => MapEntry(
              DateTime.fromMillisecondsSinceEpoch(
                (timestamp * 1000).toInt(),
              ).toLocal(),
              (count * 100).toInt(),
            ),
          )
          .toIMap(),
    ));
  } else {
    return Err(Exception('Unreachable case'));
  }
}
