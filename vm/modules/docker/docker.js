var sh = require('sh'),
	format = require('util').format

var slice = Array.prototype.slice;


function Docker () {

}

exports.Docker = Docker;

function call(cmd, args) {
	var s = format('docker %s %s', cmd, args.join(' ')) 
	return sh.exec(s).stdout
}

var proto = Docker.prototype

proto.start = function (name, image, options) {
    options = options||{};
    
    var flags = []
    for (var key in options) {
        var val = options[key];
        if (key == 'envFile'|| key == 'env' || key == 'publish') {
            if ('envFile' == key) key = 'env-file'
            if (!Array.isArray(val)) {
                val = [val]
            }
            val = val.map(function (v) { return '--' + key + ' ' + v; });
            flags.push(val.join(' '));
        } else {
            flags.push('--'+key + ' ' + val);
        }

    }
	var cmd;
	if (!this.hasContainer(name)) {
		cmd = format('docker run -d %s --name %s %s', flags.join(' '), name, image);
	} else if (this.isRunning(name)) {
		return sh.exec('docker inspect -f {{.Id}} ' + name).stdout.trim()
	} else {
		cmd = "docker start " + name;
		sh.exec(cmd);
		cmd = 'docker inspect -f {{.Id}} ' + name
	}
	console.log('running docker ', cmd)
	return sh.exec(cmd).stdout.trim();
}

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
	var reg = new RegExp('\\s*('+name+')\\s+')
	return reg.test(sh.exec("docker ps").stdout)
}	

proto.hasContainer = function (name) {
	var reg = new RegExp('\\s*('+name+')\\s+')
	return reg.test(sh.exec("docker ps -a").stdout)
}

proto.hasImage = function (name) {
	var reg = new RegExp('\\s*('+name+')\\s+')
	var o = reg.test(sh.exec("docker images").stdout)
	return o;
}

proto.inspect = function (name) {
	try {
		var json = call('inspect', [name]);
		return JSON.parse(json);
	} catch (e) {
        console.log(e)
		return null;
	}
}

