package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
		funcs := []string{
		"getAllQuestions",
    "postAllQuestions",
    "generateNewToken",
    "checkNewToken",
    "deleteNewToken",
    "getRemoteSettings",
    "getAllTags",
    "getNextQuestion",
    "gradeQuestion",
    "uploadAlgorithm",
    "setDefaultAlgorithm",
    "removeLogin",
    "removeAlgorithm",
    "getStats",}
	dI:= `
part of 'interfaces.dart';
mixin _$FetchDataWithToken {
`

dC:= `
part of 'network.dart';

  Future<Result<List<int>>> getData(String remoteHost, List<int> data) async =>
      Result.doSafe(() async {
        final res = await client.post(
          Uri.parse('$remoteHost/api/v1/do'),
          headers: <String, String>{
            'Content-Type': 'application/octet-stream',
          },
          body: data,
        );
        if (res.statusCode != 200) {
          throw ServerException(code: res.statusCode);
        }
        return res.bodyBytes as List<int>;
      });

mixin _$RemoteServer {
  String get remoteHost;
	String get token;
` 

	for _, f := range funcs {
		uf := strings.Title(f)
		dI += fmt.Sprintf(`
Future<Result<server_to_client.%s>> %s(client_to_server.%s q);
		`,uf, f, uf)
		dC += fmt.Sprintf(`

  Future<Result<server_to_client.%s>> %s(client_to_server.%s q) async =>
      (await getData(remoteHost,
        client_to_server.Message(
                secureMessage:
                    client_to_server.SecureMessage(token: token, %s: q))
            .writeToBuffer(),
      ))
          .flatMap((v) => Result.safe(()=>server_to_client.Message.fromBuffer(v))
					.flatMap((msg) {
          
          if (msg.hasOkMessage() && msg.okMessage.has%s()) {
            return Ok(msg.okMessage.%s);
          } else if (msg.hasErrorMessage()) {
            return Err(
              MyException(
                shouldLogout: msg.errorMessage.shouldLogOut,
                inner: Exception(msg.errorMessage.error),
              ),
            );
          } else {
           return Ok(msg.okMessage.%s);
        }
      },),);
		`,uf, f, uf, f, uf, f, f)
	}


	dI += "}\n"
	dC += "}\n"
	err := os.WriteFile("frontend/lib/network/fetch_data_with_token.g.dart", []byte(dI), 0777)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("frontend/lib/network/network.g.dart", []byte(dC), 0777)
	if err != nil {
		panic(err)
	}
}
