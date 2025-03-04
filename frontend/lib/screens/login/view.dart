import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:flutter_zxing/flutter_zxing.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:app/interop/get_url.dart';
import 'package:http/http.dart';

import 'token_login.dart';
import 'twelve_words_form.dart';

enum LoginMethod { twelveWords, camera, token }

class LoginScreen extends HookConsumerWidget {
  static String location = '/login';
  final Client client;
  const LoginScreen({super.key, required this.client});
  static String? remoteURL = getCurrentLocation();
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final show12Words = useState(LoginMethod.twelveWords);
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
                      switch (show12Words.value) {
                        LoginMethod.twelveWords => TwelveWordsForm(
                            client: client,
                            remoteURL: remoteURL,
                          ),
                        LoginMethod.camera => getCamera((Code c) {
                            debugPrint('barcodes: ${c.text}');
                          }),
                        LoginMethod.token => TokenLogin(
                            client: client,
                            remoteURL: remoteURL,
                          ),
                      },
                      Center(
                        child: SegmentedButton(
                            segments: [
                              ButtonSegment(
                                  value: 1,
                                  label: Text('12 Words'),
                                  icon: Icon(Icons.input)),
                              ButtonSegment(
                                  enabled: !kIsWeb,
                                  value: 2,
                                  label: Text('Camera'),
                                  icon: Icon(Icons.camera_alt)),
                              ButtonSegment(
                                  value: 3,
                                  label: Text('Token'),
                                  icon: Icon(Icons.text_fields)),
                            ],
                            selected: switch (show12Words.value) {
                              LoginMethod.twelveWords => {1},
                              LoginMethod.camera => {2},
                              LoginMethod.token => {3},
                            },
                            onSelectionChanged: (v) =>
                                switch (v.map((e) => e).toList()[0]) {
                                  1 => show12Words.value =
                                      LoginMethod.twelveWords,
                                  2 => show12Words.value = LoginMethod.camera,
                                  3 => show12Words.value = LoginMethod.token,
                                  _ => throw AssertionError(
                                      'selected an invalid state $v')
                                }),
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
}
