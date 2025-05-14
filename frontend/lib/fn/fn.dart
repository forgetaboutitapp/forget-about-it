sealed class Result<T> {
  Result<S> map<S>(S Function(T origVal) f);

  Result<S> flatMap<S>(Result<S> Function(T origVal) f);

  Future<Result<S>> doMap<S>(Future<S> Function(T origVal) f);

  Future<Result<S>> doFlatMap<S>(Future<Result<S>> Function(T origVal) f);

  void match({
    required void Function(T val) onOk,
    required void Function(Exception err) onErr,
  });

  Future<void> doMatch({
    required Future<void> Function(T val) onOk,
    required Future<void> Function(Exception err) onErr,
  });

  static Result<T> safe<T>(T Function() f) {
    try {
      return Ok(f());
    } on Exception catch (e) {
      return Err(e);
    }
  }

  static Future<Result<T>> doSafe<T>(Future<T> Function() f) async {
    try {
      return Ok(await f());
    } on Exception catch (e) {
      return Err(e);
    }
  }

  @override
  bool operator ==(Object other) {
    if (other is Ok<T> && this is Ok<T>) {
      return other.value == (this as Ok<T>).value;
    } else if (other is Err && this is Err) {
      return other.value.toString() == (this as Err).value.toString();
    }
    return false;
  }

  @override
  int get hashCode => switch (this) {
        Ok(:final value) => value.hashCode,
        Err(:final value) => value.hashCode,
      };
}

class Ok<T> extends Result<T> {
  final T value;

  Ok(this.value);

  @override
  Result<S> map<S>(S Function(T origVal) f) => Ok(f(value));

  @override
  Future<Result<S>> doMap<S>(Future<S> Function(T origVal) f) async =>
      Ok(await f(value));

  @override
  Future<Result<S>> doFlatMap<S>(
      Future<Result<S>> Function(T origVal) f) async {
    final r = await f(value);
    return switch (r) {
      Ok(:final value) => Ok(value),
      Err(:final value) => Err(value),
    };
  }

  @override
  String toString() => 'Ok($value)';

  @override
  Result<S> flatMap<S>(Result<S> Function(T origVal) f) {
    final r = f(value);
    return switch (r) {
      Ok(:final value) => Ok(value),
      Err(:final value) => Err(value),
    };
  }

  @override
  Future<void> doMatch({
    required Future<void> Function(T val) onOk,
    required Future<void> Function(Exception err) onErr,
  }) async =>
      await onOk(value);

  @override
  void match({
    required void Function(T val) onOk,
    required void Function(Exception err) onErr,
  }) =>
      onOk(value);
}

class Err<T> extends Result<T> {
  final Exception value;
  Err(this.value);

  @override
  Result<S> map<S>(S Function(T origVal) f) => Err(this.value);

  @override
  Future<Result<S>> doMap<S>(Future<S> Function(T origVal) f) async =>
      Err(this.value);

  @override
  Future<void> doMatch({
    required Future<void> Function(T val) onOk,
    required Future<void> Function(Exception err) onErr,
  }) async =>
      await onErr(value);

  @override
  void match({
    required void Function(T val) onOk,
    required void Function(Exception err) onErr,
  }) =>
      onErr(value);

  @override
  Future<Result<S>> doFlatMap<S>(
          Future<Result<S>> Function(T origVal) f) async =>
      Err(value);

  @override
  Result<S> flatMap<S>(Result<S> Function(T origVal) f) => Err(value);
  @override
  String toString() => 'Err($value)';
}
