import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'stats_data.freezed.dart';

@freezed
class StatsData with _$StatsData {
  const StatsData({
    required this.heatmapData,
    required this.pastResults,
  });
  @override
  final IMap<DateTime, int> heatmapData;
  @override
  final IMap<DateTime, int> pastResults;

  static StatsData fromJSON(Map<String, dynamic> m) {
    Map<String, dynamic> heatmapDataMap =
        m['heatmap-data'] as Map<String, dynamic>;
    Map<String, dynamic> pastDataMap =
        m['past-results'] as Map<String, dynamic>;
    Map<DateTime, int> heatmapData = heatmapDataMap.map(
      (i, e) => MapEntry(
          toDay(DateTime.fromMillisecondsSinceEpoch(int.parse(i) * 1000)), e),
    );
    Map<DateTime, int> pastResults = pastDataMap.map(
      (i, e) => MapEntry(
          toDay(DateTime.fromMillisecondsSinceEpoch(int.parse(i) * 1000)), e),
    );
    return StatsData(
      heatmapData: heatmapData.toIMap(),
      pastResults: pastResults.toIMap(),
    );
  }
}

DateTime toDay(DateTime dateTime) =>
    DateTime(dateTime.year, dateTime.month, dateTime.day);
