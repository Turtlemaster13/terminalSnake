# Maintainer: Oliver <olivergiordano2@gmail.com>
pkgname=terminalSnake
pkgver=1.0.0
pkgrel=1
pkgdesc="snake game in the terminal"
arch=('x86_64')
url='https://github.com/Turtlemaster13/terminalSnake'
license=('CC-BY')
depends=('go')

#sha256sum=("2a6134805ff111ac383df132cf292d6cbdd70aca841faeaf4ffc9fd57d53aa1a")
#source=("archive.tar.gz")

build() {
  tar -xf "$srcdir/archive.tar.gz"
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