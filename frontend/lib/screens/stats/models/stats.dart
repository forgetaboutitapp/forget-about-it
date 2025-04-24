import 'package:app/network/interfaces.dart';

import '../service.dart' as services;
import 'package:riverpod_annotation/riverpod_annotation.dart';

import 'stats_data.dart';

part 'stats.g.dart';

@riverpod
class Stats extends _$Stats {
  StatsData? _cache;
  DateTime _originalDate = DateTime.now();
  DateTime _finalDate = DateTime.now();
  @override
  StatsData? build() {
    return _cache;
  }

  Future<void> getStats(
      FetchData fd, DateTime initTime, DateTime endTime) async {
    bool reset = false;
    if (_originalDate.isAfter(initTime)) {
      _originalDate = initTime;
      reset = true;
    }
    if (_finalDate.isBefore(endTime)) {
      _finalDate = endTime;
      reset = true;
    }
    if (reset) {
      _cache = await services.getStats(fd, initTime, endTime);
      ref.invalidateSelf();
    }
  }
}
