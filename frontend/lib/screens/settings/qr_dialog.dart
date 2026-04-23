import 'dart:convert';

import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';
import 'package:grpc/grpc_web.dart';

import '../../screens/login/view.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:qr_flutter/qr_flutter.dart';
import '../../screens/general-display/show_error.dart';

class QRDialog extends HookWidget {
  final String remoteHost;
  final String token;
  final Function() logOut;
  const QRDialog({
    super.key,
    required this.remoteHost,
    required this.token,
    required this.logOut,
  });

  @override
  Widget build(BuildContext context) {
    final ValueNotifier<String?> qrCodeData = useState(null);
    final show12Words = useState(LoginMethod.twelveWords);
    useEffect(() {
      bool isCancelled = false;

      ForgetAboutItServiceClient(
              GrpcWebClientChannel.xhr(Uri.parse(remoteHost)))
          .generateNewToken(GenerateNewTokenRequest(
        token: token,
      ))
          .then((e) async {
        if (e.hasError()) {
          showError(context, e.error.error);
          return;
        }
        final res = e.ok;
        qrCodeData.value = res.newUuid;
        while (!isCancelled) {
          final a = await ForgetAboutItServiceClient(
                  GrpcWebClientChannel.xhr(Uri.parse(remoteHost)))
              .checkNewToken(CheckNewTokenRequest());
          if (a.hasError()) {
            if (context.mounted) {
              showError(context, a.error.error);
            }
          } else if (a.hasOk()) {
            if (a.ok.done) {
              isCancelled = true;
              if (context.mounted) {
                Navigator.of(context).pop();
              }
            }
          }
          await Future.delayed(Duration(seconds: 1));
        }
      });
      return () async {
        final p = await ForgetAboutItServiceClient(
                GrpcWebClientChannel.xhr(Uri.parse(remoteHost)))
            .deleteNewToken(DeleteNewTokenRequest(
          token: token,
        ));
        if (p.hasError() && context.mounted) {
          showError(context, p.error.error);
        }
      };
    }, []);
    final qrCodeDataValue = qrCodeData.value;
    if (qrCodeDataValue == null) {
      return Center(child: CircularProgressIndicator());
    }
    final decode = jsonDecode(qrCodeDataValue);
    return Dialog(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(8.0, 8.0, 8.0, 8.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(0, 0, 0, 16.0),
              child: switch (show12Words.value) {
                LoginMethod.twelveWords =>
                  TwelveWordView(twelveWords: decode['mnemonic']),
                LoginMethod.token => TokenView(
                    remoteHost: remoteHost,
                    uuid: decode['new-uuid'],
                  ),
              },
            ),
            SegmentedButton(
                segments: [
                  ButtonSegment(
                      value: 1,
                      label: Text('12 Words'),
                      icon: Icon(Icons.input)),
                  ButtonSegment(
                      value: 2,
                      label: Text('Token'),
                      icon: Icon(Icons.text_fields)),
                ],
                selected: switch (show12Words.value) {
                  LoginMethod.twelveWords => {1},
                  LoginMethod.token => {2},
                },
                onSelectionChanged: (v) =>
                    switch (v.map((e) => e).toList()[0]) {
                      1 => show12Words.value = LoginMethod.twelveWords,
                      2 => show12Words.value = LoginMethod.token,
                      _ => throw AssertionError('selected an invalid state $v')
                    }),
          ],
        ),
      ),
    );
  }
}

class TokenView extends StatelessWidget {
  final String uuid;
  final String remoteHost;
  const TokenView({super.key, required this.uuid, required this.remoteHost});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Table(
        defaultColumnWidth: IntrinsicColumnWidth(),
        children: [
          TableRow(
            children: [
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: Text('URL'),
              ),
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: SelectableText(remoteHost),
              )
            ],
          ),
          TableRow(children: [
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: Text('Token'),
            ),
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: SelectableText(uuid),
            )
          ]),
        ],
      ),
    );
  }
}

class TwelveWordView extends StatelessWidget {
  final List<dynamic> twelveWords;

  const TwelveWordView({super.key, required this.twelveWords});
  @override
  Widget build(BuildContext context) {
    return Table(
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
              child: Center(
                child: Text(
                  twelveWords[i * 2 + j].toString(),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class QrImageViewer extends StatelessWidget {
  final String uuid;
  final String remoteServer;
  final String token;
  final Function() logOut;
  const QrImageViewer(
      {super.key,
      required this.uuid,
      required this.remoteServer,
      required this.token,
      required this.logOut});
  @override
  Widget build(BuildContext context) => QrImageView(
        data: '$remoteServer;$uuid',
        version: QrVersions.auto,
        size: 400,
        gapless: false,
      );
}
