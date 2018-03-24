require("dotenv").config();

const low = require("lowdb");
const FileSync = require("lowdb/adapters/FileSync");
 
const adapter = new FileSync("db.json");
const db = low(adapter);
db.defaults({ user: [] }).write();
const userDB = db.get("user");

const telegraf = require("telegraf");
const rancherBot = new telegraf(process.env.BOT_TOKEN);

const noAuthCommands = [ "/start", "/help" ];

rancherBot.use((ctx, next) => {
	var user = userDB.find({ id: ctx.chat.id }).value();
	if (!user && !noAuthCommands.includes(ctx.message.text)) {
		if (ctx.message.text === process.env.CHAT_SECRET) {
			userDB.push({ id: ctx.chat.id }).write();
			ctx.reply("Correct password! Let the spam begin 😅");
			return;
		} else {
			ctx.reply("⛔ You are not authorized! ⛔\nEnter the correct password to gain access!");
			return;
		}
	}

	return next();
});

rancherBot.start((ctx) => ctx.reply("Welcome to the rancher alert bot! 🎉\nFirst unlock the bot with the correct password and then try /help for all commands 😁"));

rancherBot.command("/help", (ctx) => ctx.reply("Nothing here yet 😢"));

rancherBot.on("message", (ctx) => {
	if (!ctx.chat.id.toString().startsWith("-")) {
			ctx.reply("Unknown command 😱\nTry /help to list the best features 🐬");
	}
});

rancherBot.startPolling();