import 'dart:developer' as developer;

import 'package:app/network/interfaces.dart';
import 'package:app/screens/general-display/show_error.dart';
import 'package:app/screens/stats/models/stats.dart';
import 'package:app/screens/stats/models/stats_data.dart';
import 'package:flutter/material.dart';
import 'package:flutter_heatmap_calendar/flutter_heatmap_calendar.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../fn/fn.dart';

class Stats extends HookConsumerWidget {
  static String location = '/stats';
  final FetchDataWithToken remoteServer;
  const Stats({super.key, required this.remoteServer});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final statsData = ref.watch(statsProvider);
    final isError = useState(false);
    final useDate = useState(DateTime.now());
    final pastDate = useState(DateTime.now());
    final futureDate = useState(DateTime.now());

    Future<void> getDataForYear() async {
      final useYear = useDate.value;
      final pastYear = pastDate.value;
      final futureYear = futureDate.value;
      isError.value = false;
      await ref.read(statsProvider.notifier).getStats(
            remoteServer,
            DateTime(
              useYear.year,
            ),
            DateTime(
              useYear.year,
              12,
              31,
            ),
          );
      await ref.read(statsProvider.notifier).getStats(
            remoteServer,
            DateTime(
              pastYear.year,
            ),
            DateTime(
              pastYear.year,
              12,
              31,
            ),
          );
      await ref.read(statsProvider.notifier).getStats(
            remoteServer,
            DateTime(
              futureYear.year,
            ),
            DateTime(
              futureYear.year,
              12,
              31,
            ),
          );
    }

    useEffect(() {
      getDataForYear();
      return null;
    }, []);
    return Scaffold(
        appBar: AppBar(
          title: Text('Statistics'),
        ),
        body: statsData == null
            ? isError.value
                ? Container()
                : Center(
                    child: CircularProgressIndicator(),
                  )
            : switch (statsData) {
                Ok(:final value) => StatsView(
                    statsData: value,
                    remoteServer: remoteServer,
                    getNewUseDate: (newDate) {
                      useDate.value = newDate;
                      getDataForYear();
                    },
                    getNewPastDate: (newDate) {
                      pastDate.value = newDate;
                      getDataForYear();
                    },
                    getNewFutureDate: (newDate) {
                      futureDate.value = newDate;
                      getDataForYear();
                    },
                  ),
                Err(:final value) => localShowError(context, value),
              });
  }

  Widget localShowError(BuildContext context, Exception error) {
    WidgetsBinding.instance.addPostFrameCallback(
      (_) => showError(
        context,
        error.toString(),
      ),
    );
    return Container();
  }
}

class StatsView extends HookConsumerWidget {
  final FetchDataWithToken remoteServer;
  final StatsData statsData;
  final Function(DateTime newDate) getNewUseDate;
  final Function(DateTime newDate) getNewPastDate;
  final Function(DateTime newDate) getNewFutureDate;

  const StatsView({
    super.key,
    required this.getNewUseDate,
    required this.getNewPastDate,
    required this.getNewFutureDate,
    required this.remoteServer,
    required this.statsData,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final pastResults = statsData.pastResults;
    final useResults = statsData.heatmapData;
    developer.log('pastResults: $pastResults, useResults: $useResults');
    return SingleChildScrollView(
      child: Center(
        child: Column(
          children: [
            Card(
              child: Column(
                children: [
                  Text('App Usage'),
                  Padding(
                    padding: EdgeInsets.all(8),
                    child: HeatMapCalendar(
                      onClick: (v) {
                        developer
                            .log('v: $v, pastResults[v]: ${pastResults[v]}');
                        pastResults[v] == null
                            ? null
                            : ScaffoldMessenger.of(context).showSnackBar(
                                SnackBar(
                                  content: Text(
                                      'You answered ${useResults[v]} questions on ${v.year}/${v.month}/${v.day}'),
                                ),
                              );
                      },
                      onMonthChange: (newDate) => getNewUseDate(newDate),
                      colorsets: {0: Colors.green},
                      datasets: useResults.unlock,
                    ),
                  ),
                ],
              ),
            ),
            Card(
              child: Column(
                children: [
                  Text('Past Correct Values'),
                  Padding(
                    padding: EdgeInsets.all(8),
                    child: HeatMapCalendar(
                      onClick: (v) => pastResults[v] == null
                          ? null
                          : ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(
                                content: Text(
                                    'You answered ${pastResults[v]}% of the cards offered on ${v.year}/${v.month}/${v.day} correctly'),
                              ),
                            ),
                      onMonthChange: (newDate) => getNewPastDate(newDate),
                      colorsets: {0: Colors.green},
                      datasets: pastResults.isEmpty ||
                              pastResults.values
                                  .where((e) => e != pastResults.values.first)
                                  .isEmpty
                          ? null
                          : pastResults.unlock,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
