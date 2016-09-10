var sh = require('sh'),
	format = require('util').format,
	_ = require('underscore');

var slice = Array.prototype.slice;


function Docker(options) {
	options = options||{}
	this.debug = options.debug
}

exports.Docker = Docker;

function call(cmd, args, debug) {
	
	var s = "docker " + cmd.trim()
	if (args) {
		s += " " + (Array.isArray(args) ? args.join(' ') : args);
	}
	s = s.split(' ').map(function (m) { return m.trim(); }).join(' ').trim();


	//var s = format('%s %s', cmd, args.join(' '))
	return new Promise(function (resolve, reject) {
		if (debug)
			console.log('call:', s);

		module.__private_docker(s, function (err, result) {
			if (err) return reject(err);
			resolve(result.trim());
		});
	});
	
}


var proto = Docker.prototype

var multiProps = ['env', 'envFile', 'publish', 'volume', 'link']


proto.create = function (name, image, options) {
	options = options || {};

    var flags = []
	
    for (var key in options) {
		if (!!~['image'].indexOf(key)) continue;
        var val = options[key];
        if (!!~multiProps.indexOf(key)) {
            if ('envFile' == key) key = 'env-file'
            if (typeof val === 'string') {
                val = ["--" + key + "=" + val];
            } else if (Array.isArray(val)) {
				val = val.map(function (v) { return '--' + key + '=' + v; });
			} else  {
				var out = [];
				for (var k in val) {
					var v = val[k]
					if (/\d+/.test(k)) {
						k = ""
					} else {
						k += key == 'link' ? ":" : "="
					}
					out.push('--' + key + "=" + k + v);
				
				}
				val = out;
			}
            flags.push(val.join(' '));
        } else {
            flags.push('--' + key + '=' + val);
        }
    }
	
	var args = "";
	if (flags.length > 0) {
		args = format("%s ", flags.join(' '));
	}
	args += format('-d --name=%s %s', name, image);
	return call("run", args, this.debug)
}

proto.start = function (name) {
	return call("start", name, this.debug);
}
/*proto.start = function (name, image, options) {
    options = options || {};

    var flags = []
    for (var key in options) {
        var val = options[key];
        if (key == 'envFile' || key == 'env' || key == 'publish') {
            if ('envFile' == key) key = 'env-file'
            if (!Array.isArray(val)) {
                val = [val]
            }
            val = val.map(function (v) { return '--' + key + ' ' + v; });
            flags.push(val.join(' '));
        } else {
            flags.push('--' + key + ' ' + val);
        }

    }

	var cmd, args;
	if (!this.hasContainer(name)) {
		cmd = "run"
		args = format('-d %s --name %s %s', flags.join(' '), name, image);
	} else if (this.isRunning(name)) {
		cmd = "inspect"
		args = 
		cmd = sh.exec('docker inspect -f {{.Id}} ' + name).stdout.trim()
	} else {
		cmd = "docker start " + name;
		sh.exec(cmd);
		cmd = 'docker inspect -f {{.Id}} ' + name
	}
	return sh.exec(cmd).stdout.trim();
}*/

proto.stop = function () {
	return call('stop', slice.call(arguments));
}

proto.rm = function (name) {
	return call('rm -f', [name]);
}

proto.rmi = function (name) {
	return call('rmi', [name]);
}

proto.build = function (path, tag) {
	return call('build', ['--tag', tag, path]);
}

proto.isRunning = function (name) {
	var reg = new RegExp('\\s*(' + name + ')\\s+')
	return call('ps').then(function (out) {
		return reg.test(out);
	});
}

proto.hasContainer = function (name) {
	var reg = new RegExp('\\s*(' + name + ')\\s+')
	return call('ps', '-a').then(function (out) {
		return reg.test(out);
	});
}

proto.list = function (all) {
	return call("ps", !!all ? ["-a"] : []);
}

proto.hasImage = function (name) {
	var reg = new RegExp('\\s*(' + name + ')\\s+')
	
	return call('images').then(function (out) {
		return reg.test(out);
	});
}

proto.inspect = function (name) {
	return call('inspect', name)
	.then(function (out) {
		return JSON.parse(out);
	})
}

