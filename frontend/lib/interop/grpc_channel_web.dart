import 'package:grpc/grpc_web.dart';

GrpcWebClientChannel createGrpcChannel(Uri uri) => GrpcWebClientChannel.xhr(uri);
