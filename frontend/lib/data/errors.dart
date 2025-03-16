class ServerException implements Exception {
  final int code;

  ServerException({required this.code});
  @override
  String toString() => switch (code) {
        -1 => 'The server is not responding',
        401 => 'The token is invalid',
        400 ||
        500 ||
        404 ||
        405 =>
          'Either the host is incorrect or the application has a bug. Please file a bug report on https://github.com/forgetaboutitapp/forget-about-it',
        _ =>
          'Either the host is incorrect or the application has a bug (the server returned a $code). Please file a bug report on https://github.com/forgetaboutitapp/forget-about-it'
      };
}
