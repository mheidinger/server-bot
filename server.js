require('dotenv').config()

const telegraf = require('telegraf')
const rancherBot = new telegraf(process.env.BOT_TOKEN)

rancherBot.start((ctx) => ctx.reply('Welcome to the rancher alert bot! 🎉\nTry /help for all commands.'))

rancherBot.command('/help', (ctx) => ctx.replyWithMarkdown('Nothing here yet 😢'))

rancherBot.on("message", (ctx) => {
	if (!ctx.chat.id.toString().startsWith("-")) {
			ctx.reply("Unknown command 😱\nTry /help to list the best features 🐬");
	}
});

rancherBot.startPolling()