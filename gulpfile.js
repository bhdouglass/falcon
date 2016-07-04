var gulp = require('gulp');
var shell = require('gulp-shell');
var del = require('del');
var fs = require('fs');
var path = require('path');

var sdk = 'ubuntu-sdk-15.04';
var paths = {
    src: {
        click: ['click/manifest.json', 'click/falcon.apparmor'],
        scope: [
            'images/icon.png',
            'images/logo.png',
            'src/falcon.bhdouglass_falcon.ini',
            'src/falcon.bhdouglass_falcon-settings.ini',
        ],
        go: 'src/*.go',
    },
    dist: {
        click: 'dist',
        scope: 'dist/falcon/',
        go: 'dist/falcon/falcon.bhdouglass_falcon',
    }
};

gulp.task('clean', function() {
    del.sync(paths.dist.click);
});

gulp.task('move-click', function() {
    return gulp.src(paths.src.click)
        .pipe(gulp.dest(paths.dist.click));
});

gulp.task('move-scope', function() {
    return gulp.src(paths.src.scope)
        .pipe(gulp.dest(paths.dist.scope));
});

gulp.task('build-go', ['clean', 'move-click', 'move-scope'], shell.task('GOPATH=`pwd`/go go build -o ' + paths.dist.go + ' ' + paths.src.go));

gulp.task('build-go-armhf', ['clean', 'move-click', 'move-scope'], shell.task(
    'CGO_ENABLED=1 ' +
    'GOPATH=`pwd`/go ' +
    'GOARCH=arm ' +
    'GOARM=7 ' +
    'CXX=arm-linux-gnueabihf-g++ ' +
    'PKG_CONFIG_LIBDIR=/usr/lib/arm-linux-gnueabihf/pkgconfig:/usr/lib/pkgconfig:/usr/share/pkgconfig ' +
    'CC=arm-linux-gnueabihf-gcc ' +
    'go build -o ' + paths.dist.go + ' ' +
    '-ldflags \'-extld=arm-linux-gnueabihf-g++\' ' + paths.src.go
));

gulp.task('run', ['build-go'], shell.task(
    'unity-scope-tool ' + paths.dist.scope + '/falcon.bhdouglass_falcon.ini'
));

gulp.task('default', ['run']);

gulp.task('build-chroot', shell.task('sudo click chroot -a armhf -f ' + sdk + ' run gulp build-go-armhf'));

gulp.task('build-click', ['build-chroot'], shell.task('cd dist && click build .'));

function findClick() {
    var click = null;

    var dir = fs.readdirSync('./dist');
    dir.forEach(function(file) {
        if (path.extname(file) == '.click') {
            click = file;
        }
    });

    return click;
}

gulp.task('push-click', ['build-click'], shell.task(
    'adb push dist/' + findClick() + ' /home/phablet/')
);

gulp.task('install-click', ['push-click'], shell.task(
    'echo "pkcon install-local --allow-untrusted ./' + findClick() + '" | phablet-shell'
));
