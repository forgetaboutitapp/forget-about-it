import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import '../../interop/get_url.dart';

import 'token_login.dart';
import 'twelve_words_form.dart';
import 'qr_login.dart';

enum LoginMethod { twelveWords, token, qrCode }

class LoginScreen extends HookConsumerWidget {
  static String location = '/login';
  const LoginScreen({super.key});
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
                MediaQuery.sizeOf(context).width > 500 ? 64 : 4,
                MediaQuery.sizeOf(context).height > 760 ? 64.0 : 4,
                MediaQuery.sizeOf(context).width > 500 ? 64 : 4,
                MediaQuery.sizeOf(context).height > 760 ? 64.0 : 4),
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
                            remoteURL: remoteURL,
                          ),
                        LoginMethod.token => TokenLogin(
                            remoteURL: remoteURL,
                          ),
                        LoginMethod.qrCode => QrLogin(
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
                                  value: 2,
                                  label: Text('Token'),
                                  icon: Icon(Icons.text_fields)),
                              if (defaultTargetPlatform == TargetPlatform.android)
                                ButtonSegment(
                                    value: 3,
                                    label: Text('QR Code'),
                                    icon: Icon(Icons.qr_code_scanner)),
                            ],
                            selected: switch (show12Words.value) {
                              LoginMethod.twelveWords => {1},
                              LoginMethod.token => {2},
                              LoginMethod.qrCode => {3},
                            },
                            onSelectionChanged: (v) =>
                                switch (v.map((e) => e).toList()[0]) {
                                  1 => show12Words.value =
                                      LoginMethod.twelveWords,
                                  2 => show12Words.value = LoginMethod.token,
                                  3 => show12Words.value = LoginMethod.qrCode,
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
}
