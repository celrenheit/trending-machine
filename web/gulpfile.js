var gulp = require("gulp");
var $ = require('gulp-load-plugins')();
var browserSync = require('browser-sync').create()
var reload      = browserSync.reload;
var es = require('event-stream');
var streamqueue = require('streamqueue');
var pkg = require('./package.json');

var baseRoot = "./public/",
	appRoot = baseRoot + "modules/",
	destRoot = baseRoot + "dist/",
	vendorRoot = baseRoot + "components/";

var config = {
	paths: {
		scripts: {
			source: [
				appRoot + '**/*.js'
			],
			vendor: [
				vendorRoot + 'lodash/lodash.js',
				vendorRoot + 'angular/angular.js',
				vendorRoot + 'angular-aria/angular-aria.js',
				vendorRoot + 'angular-animate/angular-animate.js',
				vendorRoot + 'angular-material/angular-material.js',
				vendorRoot + 'angular-ui-router/release/angular-ui-router.js',
				vendorRoot + 'moment/moment.js',
				vendorRoot + 'angular-hotkeys/build/hotkeys.js',
				vendorRoot + 'is_js/is.js',
				vendorRoot + 'mdPickers/dist/mdPickers.js',
				vendorRoot + 'angular-social-links/angular-social-links.js',
				baseRoot   + 'vendor/**/*.js'
			],
			output: destRoot + 'js'
		},
		styles: {
			source: [appRoot + '/index.scss'],
			watchSources: [appRoot + '/**/*.scss'],
			vendor: [
				vendorRoot + 'animate.css/animate.css',
				vendorRoot + 'angular-material/angular-material.css',
				vendorRoot + 'angular-hotkeys/build/hotkeys.css',
				vendorRoot + 'mdPickers/dist/mdPickers.css',
				baseRoot   + 'vendor/**/*.css'
			],
			output: destRoot + 'css'
		}
	}
}

var today = new Date();
function dateToYMD(date) {
    var d = date.getDate();
    var m = date.getMonth() + 1;
    var y = date.getFullYear();
    return '' + y + '-' + (m<=9 ? '0' + m : m) + '-' + (d <= 9 ? '0' + d : d);
}

var banner = ['/**',
  ' * <%= pkg.title || pkg.name %> - v<%= pkg.version %> - <%= todayDate %>',
  ' * @description <%= pkg.description %>',
  ' * @version v<%= pkg.version %>',
  ' * @link <%= pkg.homepage %>',
  ' * @license <%= pkg.license %>',
  ' * Copyright (c) <%= todayYear %> <%= pkg.author %>',
  ' */',
  ''].join('\n');


// Static Server + watching scss/html files
gulp.task('serve', ['sass'], function() {

    browserSync.init({
        proxy: "localhost:3000",
        notify: false,
        open: !!$.util.env.o
    });

    gulp.watch(config.paths.scripts.source.concat(baseRoot   + 'vendor/**/*.js'), ['scripts']);
    gulp.watch(config.paths.styles.watchSources, ['sass']);
    gulp.watch(["public/*.html", appRoot + "**/*.tpl.html"]).on('change', reload);
});

gulp.task('sass', function() {
	var vendor = gulp.src(config.paths.styles.vendor)
					 .pipe($.plumber({errorHandler: $.notify.onError("Error SCSS: <%= error.message %>")}))
	var source = gulp.src(config.paths.styles.source)
					 .pipe($.plumber({errorHandler: $.notify.onError("Error SCSS: <%= error.message %>")}))
					 .pipe($.cssGlobbing({
				        // Configure it to use SCSS files
				        extensions: ['.scss']
				    }))
					 .pipe($.sass({style: 'compressed'}))

    return es.concat(vendor, source)
		.pipe($.concat('app.css'))
		.pipe($.autoprefixer('last 4 version'))
		.pipe($.if($.util.env.prod, $.minifyCss({compatibility: 'ie8'})))
		.pipe($.header(banner, {
			pkg : pkg,
			todayDate: dateToYMD(today),
			todayYear: today.getFullYear()
		}))
    .pipe(gulp.dest(config.paths.styles.output))
		.pipe($.size())
		.pipe(reload({stream: true}))
});

gulp.task('scripts', function() {
	var vendor = gulp.src(config.paths.scripts.vendor)
					.pipe($.plumber({errorHandler: $.notify.onError("Error JS: <%= error.message %>"), inherit: false}))
	var source = gulp.src(config.paths.scripts.source)
					.pipe($.plumber({errorHandler: $.notify.onError("Error JS: <%= error.message %>"), inherit: false}))
					.pipe($.babel())
					.pipe($.ngAnnotate());
	return streamqueue({objectMode: true}, vendor, source)
			.pipe($.cached('scripts'))
			.pipe($.remember('scripts'))
			.pipe($.concat('app.js'))
			.pipe($.if($.util.env.prod, $.uglify({compress: {}, mangle: true})))
			.pipe($.header(banner, {
				pkg : pkg,
				todayDate: dateToYMD(today),
				todayYear: today.getFullYear()
			}))
			.pipe(gulp.dest(config.paths.scripts.output))
			.pipe($.size())
			.pipe(reload({stream: true}))
})

gulp.task('build', ['sass', 'scripts'])

gulp.task('watch', ['serve'])
gulp.task('default', ['build', 'watch'])
