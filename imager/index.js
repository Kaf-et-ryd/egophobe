const async = require('async'),
    config = require('../config'),
    child_process = require('child_process'),
    db = require('./db'),
    etc = require('../util/etc'),
    fs = require('fs'),
    hooks = require('../util/hooks'),
    path = require('path'),
    winston = require('winston');

exports.Onegai = db.Onegai;
exports.config = config;

const image_attrs = ('src thumb ext dims size MD5 SHA1 hash imgnm spoiler'
	+ ' apng mid audio length').split(' ');
exports.image_attrs = image_attrs;

function nestImageProps(post) {
	if (!is_image(post))
		return;

	// Restructures the flat hash from redis to have image attributes in an
	// embeded hash. Better structure.
	let image = {};
	for (let key of image_attrs) {
		if (key in post) {
			image[key] = post[key];
			delete post[key];
		}
	}
	if (image.dims.split)
		image.dims = image.dims.split(',').map(parse_number);
	image.size = parse_number(image.size);

	// Hashes are only used for image duplicate detection and are useless
	// client-side
	delete image.hash;
	post.image = image;
}
exports.nestImageProps = nestImageProps;

function deleteImageProps(post) {
	if (!is_image(post))
		return;
	for (let key of image_attrs) {
		delete post[key];
	}
}
exports.deleteImageProps = deleteImageProps;

function parse_number(n) {
	return parseInt(n, 10);
}

hooks.hook_sync('inlinePost', function (info) {
	let post = info.dest;
	const image = info.src.image;
	if (!image)
		return;
	for (let i = 0, l = image_attrs.length; i < l; i++) {
		let key = image_attrs[i];
		if (key in image)
			post[key] = image[key];
	}
});

function publish(alloc, cb) {
	let mvs = [];
	for (let kind in alloc.tmps) {
		mvs.push(etc.cpx.bind(etc,
			media_path('tmp', alloc.tmps[kind]),
			media_path(kind, alloc.image[kind])
		));
	}
	async.parallel(mvs, cb);
}

function validate_alloc(alloc) {
	if (!alloc || !alloc.image || !alloc.tmps)
		return;
	for (let dir in alloc.tmps) {
		const fnm = alloc.tmps[dir];
		if (!/^[\w_]+$/.test(fnm)) {
			winston.warn("Suspicious filename: " + JSON.stringify(fnm));
			return false;
		}
	}
	return true;
}

function is_image(image) {
	return image && image.src;
}

function media_path(dir, filename) {
	return path.join(config.MEDIA_DIRS[dir], filename);
}
exports.media_path = media_path;

function make_media_dirs (cb) {
	const keys = ['src', 'thumb', 'tmp'];
	if (config.EXTRA_MID_THUMBNAILS)
		keys.push('mid');
	async.forEach(keys, 
		(key, cb) =>
			fs.mkdir(config.MEDIA_DIRS[key], err =>
				cb(err && err.code == 'EEXIST' ? null : err)),
		err => cb(err)
	);
}
exports.make_media_dirs = make_media_dirs;

function squish_MD5 (hash) {
	if (typeof hash == 'string')
		hash = new Buffer(hash, 'hex');
	return hash.toString('base64').replace(/\//g, '_').replace(/=*$/, '');
}
exports.squish_MD5 = squish_MD5;

async function obtain_image_alloc (id) {
	const onegai = new db.Onegai
	onegai.obtain_image_alloc(id, (err, alloc) => {
		if (err)
			return cb(err)

		if (validate_alloc(alloc))
			cb(null, alloc)
		else
			cb("Invalid image alloc")
	})
}
exports.obtain_image_alloc = obtain_image_alloc;

async function commit_image_alloc (alloc) {
	for (let kind in alloc.tmps) {
		await etc.copyAsync(media_path('tmp', alloc.tmps[kind]),
			media_path(kind, alloc.image[kind]))
	}
	await new db.Onegai.commit_image_alloc(alloc)
}
exports.commit_image_alloc = commit_image_alloc;
