
var name = "livejazz-admin";
var mysql = "livejazz-mysql-" + (process.env.RUNN_ENV||'development');

module.exports = {
    name: name,
    link: {
        $staging: {
            mysql: mysql
        }
    },
    $development: {
        publish: "3000:3000"
    },
    envFile: "",
    env: {

    },
    postrun: function () {

    },
    build: {
        dockerfile: "."
    },
    dependencies: [
        {
            name: mysql,
            phase: ["development", "staging"],
            $darwin: {
                publish: '3306:3306'
            },
            prebuild: function () {

            },
            build: {

            },
            dependencies: [
                {name: 'nginx'}
            ]
        }, {
            name: "nginx",
            phase: ["production", "staging"],
            build: {
                dockerfile: "./Dockerfile"
            },
            volumne: [
                "./:/etc/confd"
            ],
            publish: ["80:80","443:443"]

        }
    ]
}