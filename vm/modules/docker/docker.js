
var sh = require('shell');

function Docker(options) {
	this.options = options;
}

var proto = Docker.prototype;

proto.build = function (options) {
	if (options == null) throw new Error('no options');
	__private_docker(this.options, 'build', options, function (err, result) {

	});
}

proto.run = function (options) {
	__private_docker(this.options, 'build', options, function (err, result) {
		
	});
}

proto.stop = function (options) {

}

proto.start = function (options) {

}

module.export = Docker;

