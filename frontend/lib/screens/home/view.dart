import 'package:app/network/interfaces.dart';
import 'package:app/screens/bulk-edit/view.dart';
import 'package:app/screens/home/model.dart';
import 'package:app/screens/quiz/view.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../settings/view.dart';
import 'service.dart';

class HomeScreen extends HookConsumerWidget {
  static String location = '/';
  final FetchData remoteServer;
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
            ],
          )
        ],
      ),
      body: FutureBuilder(
        future: getAllTags(remoteServer),
        builder: (context, snapshot) => switch (snapshot.connectionState) {
          ConnectionState.none => Center(child: Text('Error')),
          ConnectionState.waiting ||
          ConnectionState.active =>
            Center(child: CircularProgressIndicator()),
          ConnectionState.done => buildDisplay(context, snapshot),
        },
      ),
    );
  }

  buildDisplay(BuildContext context, AsyncSnapshot<IList<Tag>> snapshot) {
    if (snapshot.hasError) {
      WidgetsBinding.instance.addPostFrameCallback(
        (_) => ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            backgroundColor: Colors.red,
            content: Text(
              'Error ${snapshot.error}',
            ),
          ),
        ),
      );
      return Center(child: Text(''));
    } else if (snapshot.hasData) {
      return TagsView(
        tagList: snapshot.data!,
        remoteServer: remoteServer,
      );
    } else {
      return Center(child: Text('There is no error or data'));
    }
  }
}

class TagsView extends HookWidget {
  final FetchData remoteServer;
  final IList<Tag> tagList;
  const TagsView({
    super.key,
    required this.remoteServer,
    required this.tagList,
  });

  @override
  Widget build(BuildContext context) {
    final tagListResponsive = useState(tagList);
    final selectedTags = useState(ISet<int>());
    return Column(
      children: [
        Expanded(
          child: RefreshIndicator(
            onRefresh: () async {
              tagListResponsive.value = await getAllTags(remoteServer);
              selectedTags.value = ISet();
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
                        context.go(
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
                        );
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
