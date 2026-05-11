import 'package:grpc/grpc.dart';
import 'package:grpc/grpc_connection_interface.dart';
import '../screens/login/mdns_lookup.dart';

ClientChannelBase createGrpcChannel(Uri uri) {
  final isSecure = uri.scheme == 'https';
  var host = uri.host;
  if (resolvedMdnsHosts.containsKey(host)) {
    host = resolvedMdnsHosts[host]!;
  }
  return ClientChannel(
    host,
    port: uri.hasPort ? uri.port : (isSecure ? 443 : 80),
    options: ChannelOptions(
      credentials: isSecure
          ? const ChannelCredentials.secure()
          : const ChannelCredentials.insecure(),
    ),
  );
}
