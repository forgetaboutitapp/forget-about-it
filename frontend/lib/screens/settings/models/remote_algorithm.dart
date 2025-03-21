import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'remote_algorithm.freezed.dart';

@freezed
class RemoteAlgorithm with _$RemoteAlgorithm {
  const RemoteAlgorithm({
    required this.algorithmID,
    required this.authorName,
    required this.license,
    required this.remoteURL,
    required this.downloadURL,
    required this.timestampAdded,
    required this.initializationFunctions,
    required this.allocatingFunction,
    required this.freeingFunction,
    required this.algorithm,
    required this.moduleName,
    required this.version,
    required this.algorithmName,
  });
  @override
  final int algorithmID;
  @override
  final String authorName;
  @override
  final String license;
  @override
  final String remoteURL;
  @override
  final String downloadURL;
  @override
  final int timestampAdded;
  @override
  final IList<String> initializationFunctions;
  @override
  final String allocatingFunction;
  @override
  final String freeingFunction;
  @override
  final String? algorithm;
  @override
  final String moduleName;
  @override
  final int version;
  @override
  final String algorithmName;

  static RemoteAlgorithm fromJSON(dynamic m) {
    return RemoteAlgorithm(
      algorithmID: m['id'],
      authorName: m['name'],
      license: m['license'],
      remoteURL: m['remote-url'],
      downloadURL: m['download-url'],
      timestampAdded: m['timestamp-added'],
      initializationFunctions: m['init-functions'],
      allocatingFunction: m['allocating-functions'],
      freeingFunction: m['freeing-functions'],
      algorithm: m['algorithm'],
      moduleName: m['module-name'],
      version: m['version'],
      algorithmName: m['algorithm-name'],
    );
  }
}
