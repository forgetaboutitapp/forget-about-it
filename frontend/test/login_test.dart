import 'dart:convert';

import 'package:app/data/keys.dart';
import 'package:app/screens/login/view.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

import 'package:app/main.dart';
import 'package:http/testing.dart';

void main() {
  testWidgets('Counter increments smoke test', (WidgetTester tester) async {
    SharedPreferences.setMockInitialValues({});
    sharedPreferences = await SharedPreferences.getInstance();
    expect(sharedPreferences.getKeys(), <Map>{});
    final client = MockClient((request) async {
      if (request.url.host == 'localhost' &&
          request.url.path == '/api/v0/get-token/by-twelve-words') {
        Map<String, dynamic> v = jsonDecode(request.body);
        assert(v.containsKey('twelve-words'));
        IList<String> l = (v['twelve-words'] as List<dynamic>)
            .map(
              (e) => e.toString(),
            )
            .toIList();
        if (l ==
            IList([
              'assume',
              'assume',
              'assume',
              'assume',
              'assume',
              'assume',
              'assume',
              'assume',
              'assume',
              'assume',
              'assume',
              'assume'
            ])) {
          return http.Response('{"token": "1234567890"}', 200);
        }
        return http.Response('{"token": ""}', 200);
      }
      return http.Response('{"token": ""}', 404);
    });
    await tester.pumpWidget(WrapTestWidget(
        child: LoginScreen(
      client: client,
    )));
    expect(find.text('Submit'), findsOneWidget);
    expect(
        tester
            .widget<ElevatedButton>(find.byKey(ValueKey(twelveWordLoginButton)))
            .enabled,
        false);
    await tester.pumpAndSettle();
    for (int row = 0; row < 6; row++) {
      for (int col = 0; col < 2; col++) {
        var newKey = ValueKey(TwelveWordKeys(col: col, row: row));
        expect(find.byKey(newKey), findsOneWidget);
        await tester.enterText(find.byKey(newKey), 'assume');
      }
    }
    await tester.enterText(
        find.byKey(ValueKey(loginURLKey)), 'http://localhost');
    await tester.pumpAndSettle();
    expect(
        tester
            .widget<ElevatedButton>(find.byKey(ValueKey(twelveWordLoginButton)))
            .enabled,
        true);
    await tester.tap(find.byKey(ValueKey(twelveWordLoginButton)));
    await tester.pump(Duration(seconds: 1));
    expect(sharedPreferences.getString('LOGGED_IN_KEY'), '1234567890');
  });
}

class WrapTestWidget extends StatelessWidget {
  final Widget child;
  const WrapTestWidget({super.key, required this.child});

  @override
  Widget build(BuildContext context) => ProviderScope(
      child: MaterialApp(home: child, debugShowCheckedModeBanner: false));
}
