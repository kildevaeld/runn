


function S3 (options) {
    this._client = __private_s3(options);
}

exports.S3 = S3;

var proto = S3.prototype;

proto.get = function (source, target) {
    var self = this;
    return new Promise(function (resolve, reject) {
        self._client.Get(source, target, function (err) {
            if (err) return reject(err);
            return resolve();
        });
    });
}

proto.list = function (prefix) {
    var self = this;
    /*if (prefix.length > 0 && prefix[0] === '/') {
        prefix = prefix.substr(1)
    }*/

    return new Promise(function (resolve, reject) {
        self._client.List(prefix, function (err, keys) {
            if (err) return reject(err);
            resolve(keys);
        });
    })
}