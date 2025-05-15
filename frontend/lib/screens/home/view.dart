import '../../future_widget/future_widget.dart';
import '../../network/interfaces.dart';
import '../../screens/bulk-edit/view.dart';
import '../../screens/general-display/show_error.dart';
import '../../screens/home/model.dart';
import '../../screens/quiz/view.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../fn/fn.dart';
import '../settings/view.dart';
import '../stats/view.dart';
import 'service.dart';

class HomeScreen extends HookConsumerWidget {
  static String location = '/';
  final FetchDataWithToken remoteServer;

  const HomeScreen({super.key, required this.remoteServer});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Forget About It'),
        actions: [
          PopupMenuButton(
            // https://stackoverflow.com/a/75375750 for the strategy to use a function as value
            onSelected: (fn) => fn(),
            itemBuilder: (context) => [
              PopupMenuItem(
                value: () => context.go(BulkEditScreen.location),
                child: Text(
                  'Bulk Edit',
                ),
              ),
              PopupMenuItem(
                value: () => context.go(SettingsScreen.location),
                child: Text(
                  'Settings',
                ),
              ),
              PopupMenuItem(
                value: () => context.go(Stats.location),
                child: Text(
                  'Statistics',
                ),
              ),
            ],
          )
        ],
      ),
      body: FutureWidget(
        future: getAllTags(remoteServer),
        built: (context, r) => buildDisplay(context, r),
        waiting: (context) => Center(child: CircularProgressIndicator()),
      ),
    );
  }

  Widget buildDisplay(BuildContext context, Result<(IList<Tag>, bool)> data) {
    return switch (data) {
      Ok(value: final d) =>
        TagsView(remoteServer: remoteServer, canRun: d.$2, tagList: d.$1),
      Err(value: final error) => localShowError(context, error),
    };
  }

  Widget localShowError(BuildContext context, Exception error) {
    showErrorDelayed(context, error.toString());
    return Container();
  }
}

class TagsView extends HookWidget {
  final FetchDataWithToken remoteServer;
  final IList<Tag> tagList;
  final bool canRun;

  const TagsView({
    super.key,
    required this.remoteServer,
    required this.canRun,
    required this.tagList,
  });

  @override
  Widget build(BuildContext context) {
    final tagListResponsive = useState(tagList);
    final canRunResponsive = useState(canRun);
    final selectedTags = useState(ISet<int>());
    return Column(
      children: [
        Expanded(
          child: RefreshIndicator(
            onRefresh: () async {
              (await getAllTags(remoteServer)).match(onErr: (e) {
                showError(context, e.toString());
              }, onOk: (v) {
                final (tags, canRun) = v;
                tagListResponsive.value = tags;
                canRunResponsive.value = canRun;
                selectedTags.value = ISet();
              });
            },
            child: ScrollConfiguration(
              behavior: ScrollConfiguration.of(context).copyWith(
                dragDevices: {
                  // https://github.com/flutter/flutter/issues/142529
                  PointerDeviceKind.mouse,
                  PointerDeviceKind.touch,
                  PointerDeviceKind.stylus,
                  PointerDeviceKind.unknown,
                },
              ),
              child: ListView.builder(
                physics: const AlwaysScrollableScrollPhysics(),
                itemCount: tagListResponsive.value.length,
                itemBuilder: (BuildContext context, int count) => ListTile(
                  leading: Checkbox(
                      value: selectedTags.value.contains(count),
                      onChanged: (q) {
                        if (q == true) {
                          selectedTags.value = selectedTags.value.add(count);
                        } else {
                          selectedTags.value = selectedTags.value.remove(count);
                        }
                      }),
                  title: Text(tagListResponsive.value[count].tag),
                  subtitle: Text(
                      'Total questions: ${tagListResponsive.value[count].totalQuestions}'),
                ),
              ),
            ),
          ),
        ),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: TextButton(
                onPressed: selectedTags.value.isEmpty
                    ? null
                    : () {
                        canRunResponsive.value
                            ? context.go(
                                Uri(
                                  path: QuizView.location,
                                  queryParameters: {
                                    'tags': selectedTags.value
                                        .map(
                                          (e) => tagList[e].tag,
                                        )
                                        .toList()
                                  },
                                ).toString(),
                              )
                            : showError(context, 'You need to add a scheduler');
                      },
                child: Text('Quiz Me!'),
              ),
            )
          ],
        ),
      ],
    );
  }
}
