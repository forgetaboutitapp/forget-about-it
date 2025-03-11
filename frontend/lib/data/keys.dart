class TwelveWordKeys {
  const TwelveWordKeys({required this.row, required this.col});
  final int col;
  final int row;

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) {
      return true;
    }
    if (other.runtimeType != runtimeType) {
      return false;
    }
    return other is TwelveWordKeys && other.col == col && other.row == row;
  }

  @override
  int get hashCode => Object.hash(row, col);
}

const twelveWordLoginButton = 10001;
const loginURLKey = 10002;
const loginURLKeyInToken = 10003;
const loginTokenKeyInToken = 10004;
