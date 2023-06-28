# Maintainer: Oliver <olivergiordano2@gmail.com>
pkgname=terminalSnake
pkgver=1.0.0
pkgrel=1
pkgdesc="snake game in the terminal"
arch=('x86_64')
url='https://github.com/Turtlemaster13/terminalSnake'
license=('CC-BY')
depends=('go')


build() {
  cd "$srcdir/$pkgname-$pkgver"
  go build -o "snakeGame"
}

package() {
  username=$LOGNAME
  cd "$srcdir/$pkgname-$pkgver"
  install -Dm755 snakeGame "$pkgdir/usr/bin/terminalSnake"
  install -Dm644 snakeGame.csv "$pkgdir/usr/share/terminalSnake.csv"
  chown -R "$username":"$username" "$pkgdir/usr/share/terminalSnake.csv"
}
