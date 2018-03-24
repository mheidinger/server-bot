require("dotenv").config();

const Rancher = require("rancher-client");
 
const client = new Rancher.Client({ url: "http://max-heidinger.de:8089/v2-beta", access_key: process.env.RANCHER_KEY, secret_key: process.env.RANCHER_SECRET });

var stack;

client.getStacks().then((info) => {
	stack = info[0];
	console.log(stack.healthState);

	client.getStackServices(stack.id).then((info) => {
		console.log(info);
	}).catch((err) => {
		console.error(" ERROR: ", err);
	});
}).catch((err)=>{
  console.error(" ERROR: ", err);
});