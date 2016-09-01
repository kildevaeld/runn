var argv = require('minimist')(process.argv.slice(1))
console.log('test', JSON.stringify(argv, null, 2))

Promise.delay(1000).then(function () {
	console.log('delay');
})


console.log(process.env.NODE_ENV)