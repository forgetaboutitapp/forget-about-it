import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'stats_data.freezed.dart';

@freezed
class StatsData with _$StatsData {
  const StatsData({
    required this.heatmapData,
    required this.pastResults,
    required this.futureUsage,
  });
  @override
  final IMap<DateTime, int> heatmapData;
  @override
  final IMap<DateTime, int> pastResults;

  @override
  final IMap<DateTime, int> futureUsage;
}

DateTime toDay(DateTime dateTime) =>
    DateTime(dateTime.year, dateTime.month, dateTime.day);
