
var S3 = require('s3').S3;


/*s3.get('/croppings/0006a0ac155dd2bb4091140c3ca1608fa58ac65a.jpg','image.jpg').then(function () {
	console.log('done')
}).catch(console.error);*/

s3.list('croppings/').then(function (list) {
	console.log('list', JSON.stringify(list, null, 2));
}).catch(console.error);


var argv = require('minimist')(process.argv.slice(1))
console.log('test', JSON.stringify(argv, null, 2))

Promise.delay(1000).then(function () {
	console.log('delay');
})


console.log(process.env.NODE_ENV)