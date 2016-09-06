

exports.archive = function (format, path, target) {
    return new Promise(function (resolve, reject) {
        __private_archive('pack', {
            format: format,
            source: path,
            target: target,
        }, function (err) {
            if (err) return reject(err);
            resolve();
        });
    });
}

exports.unarchive = function(format, path, target) {
    return new Promise(function (resolve, reject) {
        __private_archive('unpack', {
            format: format,
            source: path,
            target: target,
        }, function (err) {
            if (err) return reject(err);
            resolve();
        });
    });
}