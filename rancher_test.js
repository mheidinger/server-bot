require("dotenv").config();

const Rancher = require("rancher-client");
 
const client = new Rancher.Client({ url: process.env.RANCHER_URL, access_key: process.env.RANCHER_KEY, secret_key: process.env.RANCHER_SECRET });

client.getStacks().then(checkStacks, (reason) => { console.error("Getting stacks failed: ", reason); });

function checkStacks(stacks) {
	stacks = stacks.filter(stack => !stack.system);
	for (var it in stacks) {
		console.log(stacks[it].name + " [" + stacks[it].description + "]: " + stacks[it].healthState);
	}
}