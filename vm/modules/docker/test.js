
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
    prestart: function () {
        console.log('prestart')
    },
    build: {
        dockerfile: process.cwd() + "/admin/Dockerfile"
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
                dockerfile: process.cwd() + "/mysql/Dockerfile"
            },
            dependencies: [
                {name: 'nginx'}
            ]
        }, {
            name: "nginx",
            phase: ["production", "staging"],
            volumne: [
                "./:/etc/confd"
            ],
            prestart: function () {
                console.log('prestart')
            },
            publish: ["80:80","443:443"]

        }
    ]
}