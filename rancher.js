const Rancher = require("rancher-client");
 
const client = new Rancher.Client({ url: process.env.RANCHER_URL, access_key: process.env.RANCHER_KEY, secret_key: process.env.RANCHER_SECRET });

function getStackStatus() {
	return new Promise((resolve, reject) => {
		client.getStacks().then((stacks) => {
			stacks = stacks.filter(stack => !stack.system);
			var returnString = "";
			for (var it in stacks) {
				returnString += stacks[it].name + " [" + stacks[it].description + "]: " + stacks[it].healthState + "\n";
			}
			resolve(returnString);
		}, (reason) => {
			reject(Error("Getting stacks failed: " + reason));
		});
	});
}

exports.getStackStatus = getStackStatus;