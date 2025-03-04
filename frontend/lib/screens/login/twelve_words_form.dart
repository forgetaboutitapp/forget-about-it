import 'package:app/bip39/wordlist.dart';
import 'package:app/screens/login/submit_type.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../data/keys.dart';
import '../../state/login.dart';
import 'login_button.dart';

class TwelveWordsForm extends HookConsumerWidget {
  final String? remoteURL;
  final http.Client client;
  const TwelveWordsForm(
      {super.key, required this.remoteURL, required this.client});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final urlController = useTextEditingController(text: remoteURL ?? '');
    final urlText = useState(remoteURL ?? '');
    final twelveWordHook = useState(List.generate(12, (i) => '').toIList());

    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Padding(
          padding:
              EdgeInsets.all(MediaQuery.sizeOf(context).height > 550 ? 8.0 : 0),
          child: TextField(
            key: ValueKey(loginURLKey),
            onChanged: (value) => urlText.value = value,
            controller: urlController,
            readOnly: kReleaseMode && remoteURL != null,
            decoration: InputDecoration(
              border: OutlineInputBorder(),
              labelText: 'Server Host',
            ),
          ),
        ),
        Table(
          children: List.generate(
            6,
            (i) => TableRow(
              children: List.generate(
                2,
                (j) => Padding(
                  padding: EdgeInsets.fromLTRB(
                      8,
                      MediaQuery.sizeOf(context).height > 550 ? 4 : 0,
                      8,
                      MediaQuery.sizeOf(context).height > 550 ? 4 : 0),
                  child: Autocomplete<String>(
                    //https://github.com/flutter/flutter/issues/98728
                    optionsViewBuilder: (BuildContext context,
                        AutocompleteOnSelected<String> onSelected,
                        Iterable<String> options) {
                      return Align(
                        alignment: Alignment.topLeft,
                        child: Material(
                          elevation: 4.0,
                          child: ConstrainedBox(
                            constraints: const BoxConstraints(
                                maxHeight: 200, maxWidth: 200),
                            child: ListView.builder(
                              padding: EdgeInsets.zero,
                              shrinkWrap: true,
                              itemCount: options.length,
                              itemBuilder: (BuildContext context, int index) {
                                final String option = options.elementAt(index);
                                return InkWell(
                                  onTap: () {
                                    onSelected(option);
                                  },
                                  child: Container(
                                    color: Theme.of(context).focusColor,
                                    padding: const EdgeInsets.all(16.0),
                                    child: Text(option),
                                  ),
                                );
                              },
                            ),
                          ),
                        ),
                      );
                    },
                    fieldViewBuilder: (context, controller, focusNode, f) {
                      return TextField(
                        key: ValueKey(
                          TwelveWordKeys(row: i, col: j),
                        ),
                        onChanged: (s) => twelveWordHook.value = twelveWordHook
                            .value
                            .replace(i * 2 + j, s.toString()),
                        onEditingComplete: f,
                        controller: controller,
                        focusNode: focusNode,
                        autofocus: true,
                        decoration: InputDecoration(
                          border: OutlineInputBorder(),
                          labelText: 'Word ${i * 2 + j + 1}',
                        ),
                      );
                    },
                    initialValue:
                        TextEditingValue(text: twelveWordHook.value[i * 2 + j]),
                    optionsBuilder: (v) {
                      if (v.text == '') {
                        return [];
                      } else {
                        return WORDLIST
                            .where((t) => t.startsWith(v.text.toLowerCase()));
                      }
                    },
                    onSelected: (s) => twelveWordHook.value =
                        twelveWordHook.value.replace(i * 2 + j, s.toString()),
                  ),
                ),
              ),
            ),
          ),
        ),
        Padding(
          padding: EdgeInsets.fromLTRB(
              8,
              MediaQuery.sizeOf(context).height > 550 ? 4 : 2,
              8,
              MediaQuery.sizeOf(context).height > 550 ? 16 : 2),
          child: LoginButton(
            shouldEnable: twelveWordHook.value
                    .where(
                      (w) =>
                          w.trim().isNotEmpty &&
                          WORDLIST.contains(
                            w.toLowerCase(),
                          ),
                    )
                    .length ==
                12,
            remoteURLString: urlController.value.text,
            toRun: (uri) async => await ref.read(loginProvider.notifier).update(
                  client,
                  uri,
                  TwelveWords(twelveWords: twelveWordHook.value),
                ),
          ),
        )
      ],
    );
  }
}
