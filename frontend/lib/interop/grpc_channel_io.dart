import 'package:grpc/grpc.dart';
import 'package:grpc/grpc_connection_interface.dart';

ClientChannelBase createGrpcChannel(Uri uri) {
  final isSecure = uri.scheme == 'https';
  return ClientChannel(
    uri.host,
    port: uri.hasPort ? uri.port : (isSecure ? 443 : 80),
    options: ChannelOptions(
      credentials: isSecure
          ? const ChannelCredentials.secure()
          : const ChannelCredentials.insecure(),
    ),
  );
}
