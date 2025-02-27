import 'package:app/data/keys.dart';
import 'package:app/bip39/wordlist.dart';
import 'package:app/state/login.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:flutter_zxing/flutter_zxing.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:app/interop/get_url.dart';
import 'package:http/http.dart';

class LoginScreen extends HookConsumerWidget {
  static String location = '/login';
  final Client client;
  const LoginScreen({super.key, required this.client});
  static String? remoteURL = getCurrentLocation();
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final twelveWordHook = useState(List.generate(12, (i) => '').toIList());
    final show12Words = useState(true);
    final urlController = useTextEditingController(text: remoteURL ?? '');
    final scrollController = useScrollController();
    return Scaffold(
      body: Container(
        color: Colors.black87,
        child: Center(
          child: Padding(
            padding: EdgeInsets.fromLTRB(
                64,
                MediaQuery.sizeOf(context).height > 680 ? 64.0 : 4,
                64,
                MediaQuery.sizeOf(context).height > 680 ? 64.0 : 4),
            child: SizedBox(
              width: MediaQuery.sizeOf(context).width > 800
                  ? 800
                  : MediaQuery.sizeOf(context).width,
              height: MediaQuery.sizeOf(context).height,
              child: Card.outlined(
                child: SingleChildScrollView(
                  controller: scrollController,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      Center(
                        child: Text(
                          'Log In',
                          style: TextStyle(
                              fontSize: MediaQuery.sizeOf(context).height > 550
                                  ? 36
                                  : 26),
                        ),
                      ),
                      show12Words.value
                          ? build12Words(
                              context, ref, urlController, twelveWordHook)
                          : getCamera((Code c) {
                              debugPrint('barcodes: ${c.text}');
                            }),
                      Center(
                        child: SegmentedButton(
                          segments: [
                            ButtonSegment(
                                value: 1,
                                label: Text('12 Words'),
                                icon: Icon(Icons.input)),
                            ButtonSegment(
                                value: 2,
                                label: Text('Camera'),
                                icon: Icon(Icons.camera_alt)),
                          ],
                          selected: {show12Words.value ? 1 : 2},
                          onSelectionChanged: kIsWeb
                              ? null
                              : (v) => v.contains(1)
                                  ? show12Words.value = true
                                  : show12Words.value = false,
                        ),
                      )
                    ],
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Column getCamera(void Function(Code) callback) {
    return Column(
      children: [
        Center(
          child: SizedBox(
            width: 400,
            height: 400,
            child: Center(
              child: kIsWeb
                  ? Text('Camera does not work on the browser')
                  : ReaderWidget(
                      onScan: callback,
                    ),
            ),
          ),
        )
      ],
    );
  }

  Column build12Words(
      BuildContext context,
      WidgetRef ref,
      TextEditingController urlController,
      ValueNotifier<IList<String>> twelveWordHook) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Padding(
          padding:
              EdgeInsets.all(MediaQuery.sizeOf(context).height > 550 ? 8.0 : 0),
          child: TextField(
            key: ValueKey(loginURLKey),
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
                            .where((t) => t.contains(v.text.toLowerCase()));
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
          child: ElevatedButton(
            key: ValueKey(twelveWordLoginButton),
            onPressed: twelveWordHook.value
                        .where((w) =>
                            w.trim().isNotEmpty &&
                            WORDLIST.contains(w.toLowerCase()))
                        .length !=
                    12
                ? null
                : () async {
                    Uri remoteURL;
                    try {
                      remoteURL = Uri.parse(urlController.value.text);
                    } catch (e) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(
                          backgroundColor: Colors.red,
                          content: Text(
                              'Cannot parse URL ${urlController.value.text}'),
                        ),
                      );
                      return;
                    }
                    try {
                      bool success = await ref
                          .read(loginProvider.notifier)
                          .update(client, remoteURL, twelveWordHook.value);
                      if (!success && context.mounted) {
                        ScaffoldMessenger.of(context).showSnackBar(
                          SnackBar(
                            backgroundColor: Colors.red,
                            content: Text('Invalid Twelve Words'),
                          ),
                        );
                      }
                    } catch (e) {
                      if (!context.mounted) return;
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(
                          backgroundColor: Colors.red,
                          content: Text('Error Connecting: $e'),
                        ),
                      );
                    }
                  },
            child: Text('Submit'),
          ),
        )
      ],
    );
  }
}
