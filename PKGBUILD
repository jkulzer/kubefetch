# Maintainer: jkulzer <kulzer dot johannes at tutanota dot com>
pkgname=kubefetch-bin
pkgver=0.5.4
pkgrel=1
pkgdesc="A tool to display information about a Kubernetes cluster in the style of Neofetch"
arch=(x86_64)
url="https://github.com/jkulzer/kubefetch"
license=('GPL')
groups=()
depends=()
makedepends=('go')
optdepends=()
provides=()
conflicts=()
replaces=()
backup=()
options=()
install=
changelog=
source=(https://github.com/jkulzer/kubefetch/releases/download/$pkgver/kubefetch)
noextract=()

package() {
	mkdir -p $pkgdir/usr/bin
	cp $srcdir/kubefetch $pkgdir/usr/bin/kubefetch
}
sha256sums=('0c2f5973e391f90b5134dd1ab981ec3a9c21fe71fa51fe946f5ab0d6110cf6a1')
