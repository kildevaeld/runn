
var _ = require('underscore'),
    format = require('util').format,
    fsm = require('fsm'),
    Docker = require('docker.async').Docker,
    EventEmitter = require('events').EventEmitter;


function isObject(obj) {
    return obj === Object(obj);
}

var eachAsync = function (list, cb) {
    var i = 0, l = list.length, result = [];

    return next().then(function () {
        return result;
    });

    function next() {
        if (i == l) return Promise.resolve(result);
        return Promise.resolve(cb(list[i++])).then(function (ret) {
            result.push(ret);
            return next()
        });
    };
}

var Builder = (function (__super) {

    function Builder(modules, env, debug) {
        __super.call(this);
        this.modules = modules;
        this.env = env //||'development';
        //this.startup();
        this.docker = new Docker({debug:!!debug});
    }

    function getCreateOptions(mod, env) {
        var out = {}
        var exclude = ['name', 'postrun', 'prerun', 'prebuild', 'postbuild', 
            'build', 'dependencies', 'phase', 'prestart', 'poststart', 'initialize']
        if (mod.phase) {

            if (!Array.isArray) mod.phase = [mod.phase];
            if (!!!~mod.phase.indexOf(env)) {
                self.trigger('notification', n.skipping, mod);
                return out;
            }
        }
        for (var key in mod) {

            var value = mod[key];
            if (key[0] == "$") {
                if (key.substr(1) === env) {
                    out = _.extend(out, mod[key]);
                } else if (key.substr(1) === process.platform) {
                    out = _.extend(out, mod[key]);
                }
            } else if (!!~exclude.indexOf(key)) {

                continue;
            } else {
                if (isObject(value) && !Array.isArray(value)) {
                    value = getCreateOptions(value, env);
                }

                if (isFunction(value)) continue;
                out[key] = value
            }
        }
        return out;
    }

    function runHook(hook, mod) {
        if (typeof mod[hook] === 'function') return Promise.resolve(mod[hook].call(mod, mod));
        return Promise.resolve();
    }

    _.extend(Builder.prototype, EventEmitter.prototype, {

        _modName: function(mod) {
            return mod.name + "-" + this.env
        },

        build: function () {
            var self = this;
            var builds = this.modules.map(function (step) {
                if (!step.build) return null;
                return step; //_.pick(step, ['build', 'prebuild', 'postbuild'])
            }).filter(function (step) { return step != null })
                .map(function (mod) {
                    return self._buildModule.call(self, mod);
                });



            return Promise.all(builds);


        },

        _buildModule: function (mod) {
            var self = this;
            self.trigger('notification', n.building, mod);
            var build = mod.build;
            
            if (build == undefined) return Promise.resolve();
            if (build.dockerfile == null && build.image == null) {  
                self.trigger('notification', n.build, mod);
                return Promise.resolve()
            } 

            return runHook('prebuild', mod).then(function () {
                return self.docker.build(build.dockerfile, mod.image||mod.name);
            })
                .then(function (out) {
                    return runHook('postbuild', mod)
                        .then(function () {
                            self.trigger('notification', n.build, mod);
                            return out;
                        })


                })
        },

        start: function (autoBuild) {
            var self = this;

            return eachAsync(this.modules, function (mod) {
                if (mod.phase) {
                    if (!Array.isArray(mod.phase)) mod.phase = [mod.phase];

                    if (!!!~mod.phase.indexOf(self.env)) {
                        self.trigger('notification', n.skipping, mod);
                        return;
                    }
                }
                var name = mod.name
                var promises = [
                    self.docker.hasContainer(name),
                    self.docker.isRunning(name),
                    self.docker.hasImage(name)
                ];

                return Promise.all(promises)
                    .then(function (ret) {
                        
                        self.trigger('notification', n.starting, mod)
                        if (ret[1] == true) {
                            self.trigger("notification", n.alreadyStarted, mod);;
                            return false;
                        } else if (ret[0] == true) {

                            return runHook('prestart', mod)
                                .then(function () {
                                    
                                    return self.docker.start(mod.name);
                                })
                                .then(function () {
                                    return true;
                                })
                        } else if (ret[2] != true && autoBuild == true) {
                            return self._buildModule.call(self, mod)
                                .then(function () {
                                    return runHook('prestart', mod)
                                }).then(function () {
                                    var options = getCreateOptions(mod, self.env);
                                    return self.docker.create(name, mod.image||name, options)
                                })/*.then(function () {
                                    return true;
                                })*/
                        } else {
                            var options = getCreateOptions(mod, self.env);
                            return runHook('prestart', mod)
                                .then(function () {
                                    return self.docker.create(name, mod.image||name, options)
                                        .then(function () {
                                            return true;
                                        })
                                })

                        }
                    }).then(function (started) {
                        if (started) {
                            self.trigger('notification', n.started, mod);
                            return runHook('poststart', mod)
                        }
                    });

            }).then(function (results) {
                return results;
            })
        },

        stop: function () {
            var self = this;
            return eachAsync(this.modules.reverse(), function (mod) {
            
                return self.docker.isRunning(mod.name)
                .then(function (ok) {
                    if (!ok) return;

                    self.trigger('notification', n.stopping, mod);

                    return runHook('prestop', mod)
                    .then(function () {
                        return self.docker.stop(mod.name);
                    });
                }).then(function () {
                    self.trigger('notification', n.stopped, mod);
                    return runHook('poststop', mod);
                })

            }).then(function () {
                self.modules.reverse();
            }).catch(function (e) {
                self.modules.reverse();
                return Promise.reject(e);
            })
        },

        remove: function (removeImages) {
            var self = this;

            return eachAsync(this.modules.reverse(), function (mod) {
                var modName = mod.name,
                    imageName = mod.image||mod.name;
                return self.docker.hasContainer(mod.name)
                .then(function (has) {
                    self.trigger('notification', n.removing, mod)
                    if (has) {
                        return runHook("preremove", mod)
                        .then(function () {
                            return self.docker.remove(modName).then(function () { return has; });
                        })
                    } else {
                        self.trigger('notification', n.skipping, mod);
                    }
                    return false;
                }).then(function (has) {
                    var p = Promise.resolve()
                    if (has) {
                        self.trigger('notification', n.removing, mod)
                        p = runHook('postremove', mod);
                    }

                    if (removeImages) {
                        return p.then(function () { return self.docker.hasImage(imageName); });
                    }
                    return p.then(function () { return false; })
                }).then(function (has) {
                    if (has) {
                        self.trigger('notification', n.removingimage, mod)
                        return runHook("preremoveimage", mod)
                        .then(function () {
                            return self.docker.rmi(imageName).then(function () { return has; });
                        })
                    }
                    return false;
                }).then(function (has) {
                    if (has)  {
                        self.trigger('notification', n.removedimage, mod)
                        return runHook('postremoveimage', mod);
                    }
                })
               
            }).then(function () {
                self.modules.reverse();
            })
        }


    });

    /*fsm.create({
        target: Builder.prototype,
        events: [
            { name: "startup", from: 'none', to: 'build' },
            { name: "build", from: "ready", to: "building"},
            { name: 'start', from: ['building', 'stopping'], to: 'starting'},
            { name: 'stop', from: 'starting', to: 'stopping'}
        ],
        callbacks: {
            onbuild: function (event, from, to) {

                return fsm.ASYNC;
            },

            onstart: function (event, from, to) {

                return fsm.ASYNC;
            },

            onstop: function (event, from, to) {

                return fsm.ASYNC

            }
        }
    })*/

    return Builder

})(EventEmitter);


function isFunction(a) {
    return a && typeof a === 'function';
}


exports.createBuilder = function (a, env, debug) {
    if (isFunction(a)) {
        a = a();
    }

    var known_modules = {};

    return Promise.resolve(a)
        .then(function (options) {
            if (typeof options.initialize === 'function') {
                return Promise.resolve(options.initialize(options))
                .then(function () { return options })
            }
            return options;
        })
        .then(function (options) {

            parseModule(options, known_modules)

            var out = [];
            if (options.dependencies != null) {
                resolveDependencies(options.dependencies, known_modules, out);
            }
            
            out.push(options);

            var builds = out.map(function (step) {
                if (!step.build) return null;
                return _.pick(step, ['build', 'prebuild', 'postbuild'])
            }).filter(function (step) { return step != null });

            return new Builder(out, env, debug);

        });
}

var n = exports.notifications = {
    starting: 'starting',
    stopping: 'stopping',
    stopped: 'stopped',
    started: 'started',
    alreadyStarted: 'alreadystarted',
    startError: 'starterror',
    building: 'building',
    build: 'build',
    buildError: 'builderror',
    skipping: "skipping",
    removing: "removing",
    removed: "removed",
    removingimage: 'removingimage',
    removedimage: "removedimage"
}

function parseModule(options, known_modules) {
    var name = options.name;

    if (!known_modules[name]) {
        known_modules[name] = options;
    } else if (Object.keys(known_modules[name]).length < Object.keys(options).length) {
        known_modules[name] = options;
    }
    if (options.dependencies) {
        var deps = options.dependencies
        for (var i = 0, ii = deps.length; i < ii; i++) {
            var mod = deps[i]
            parseModule(mod, known_modules)
        }
    }
}


function resolveDependencies(dependencies, known_modules, out) {
    var first = true;
    for (var i = 0, ii = dependencies.length; i < ii; i++) {
        var mod = dependencies[i];
        mod = known_modules[mod.name];
        var deps = mod.dependencies;

        var found = _.find(out, function (v) {
            return v.name == mod.name;
        })

        if (found != null) continue;

        if (!deps) {
            var found = _.find(out, function (v) {
                return v.name == mod.name;
            })

            if (!found) out.push(known_modules[mod.name]);
            continue;
        }

        var sdep = _.find(deps, function (v) {
            var m = _.find(known_modules, function (vv) {
                return vv.name == v.name
            })
            if (m && m.dependencies) {
                return !!_.find(m.dependencies, function (vvv) {
                    return vvv.name == mod.name
                })
            }
            return false;
        })

        if (sdep) {
            throw new Error(format('circle dep: %s and %s depends on eachother', sdep.name, mod.name))
        }

        resolveDependencies(deps, known_modules, out)

        var found = _.find(out, function (v) {
            return v.name == mod.name;
        })

        if (!found) out.push(mod);

    }

}

