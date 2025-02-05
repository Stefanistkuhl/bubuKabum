const fs = require('node:fs');
const path = require('node:path');
const dotenv = require('dotenv');
dotenv.config();
const { MessageFlags, Client, Collection, ModalBuilder, Events, GatewayIntentBits, TextInputBuilder, TextInputStyle, ActionRowBuilder, } = require('discord.js');

const client = new Client({
	intents: [
		GatewayIntentBits.Guilds,
		GatewayIntentBits.GuildMessages,
		GatewayIntentBits.GuildExpressions
	],
});

const requestvals = {
	link: "",
	is2Frame: false,
	name: "",
	guildid: "",
}

const apiUrl = "http://localhost:6999/api/emote"
const apiDownloadUrl = "http://localhost:6999/converted"

const regex = /https:\/\/(?:7tv\.app|old\.7tv\.app)\/emotes\/[A-Z0-9]{26}/;

client.commands = new Collection();
const foldersPath = path.join(__dirname, 'commands');
const commandFolders = fs.readdirSync(foldersPath);

for (const folder of commandFolders) {
	const commandsPath = path.join(foldersPath, folder);
	const commandFiles = fs.readdirSync(commandsPath).filter(file => file.endsWith('.js'));
	for (const file of commandFiles) {
		const filePath = path.join(commandsPath, file);
		const command = require(filePath);
		if ('data' in command && 'execute' in command) {
			client.commands.set(command.data.name, command);
		} else {
			console.log(`[WARNING] The command at ${filePath} is missing a required "data" or "execute" property.`);
		}
	}
}

client.once(Events.ClientReady, readyClient => {
	console.log(`Ready! Logged in as ${readyClient.user.tag}`);
});
client.on(Events.InteractionCreate, async interaction => {
	// Command handler
	if (interaction.isChatInputCommand()) {
		const command = interaction.client.commands.get(interaction.commandName);
		if (!command) {
			console.error(`No command matching ${interaction.commandName} was found.`);
			return;
		}
		try {
			await command.execute(interaction);
		} catch (error) {
			console.error(error);
			if (interaction.replied || interaction.deferred) {
				await interaction.followUp({ content: 'There was an error while executing this command!', flags: MessageFlags.Ephemeral });
			} else {
				await interaction.reply({ content: 'There was an error while executing this command!', flags: MessageFlags.Ephemeral });
			}
		}
	}

	// Button handler
	if (interaction.isButton()) {
		if (interaction.customId === 'addMore') {
			const modal = new ModalBuilder()
				.setCustomId('inputModal')
				.setTitle('Add Emote');

			const linkInput = new TextInputBuilder()
				.setCustomId('linkInput')
				.setLabel("Enter the 7tv Emote URL")
				.setRequired(true)
				.setStyle(TextInputStyle.Short);

			const nameInput = new TextInputBuilder()
				.setCustomId('nameInput')
				.setLabel("Emote Name (Optional)")
				.setStyle(TextInputStyle.Short)
				.setMinLength(2)
				.setMaxLength(32)
				.setRequired(false);
			const is2FrameGifInput = new TextInputBuilder()
				.setCustomId('2FrameGif')
				.setLabel("Convert to 2-frame GIF? (Optional, yes/no)")
				.setStyle(TextInputStyle.Short)
				.setRequired(false);

			const firstActionRow = new ActionRowBuilder().addComponents(linkInput);
			const secondActionRow = new ActionRowBuilder().addComponents(nameInput);
			const thirdActionRow = new ActionRowBuilder().addComponents(is2FrameGifInput);

			modal.addComponents(firstActionRow, secondActionRow, thirdActionRow);

			await interaction.showModal(modal);
		} else if (interaction.customId === 'close') {
			await interaction.update({ content: 'Closing', components: [] });
		}
	}

	// Modal submit handler
	if (interaction.isModalSubmit()) {
		if (interaction.customId === 'inputModal') {
			const link = interaction.fields.getTextInputValue('linkInput');
			const name = interaction.fields.getTextInputValue('nameInput');
			const gif = interaction.fields.getTextInputValue('2FrameGif');
			var gif_conv
			if (gif.toLowerCase() === "yes" || gif.toLowerCase() === "y") {
				gif_conv = true
			} else {
				gif_conv = false
			}
			if (regex.test(link)) {
			} else {
				await interaction.reply({ content: `Please enter a valid link`, ephemeral: true });
				return;
			}
			result = await isEmoteNotFound(link)
			if (result) {
				await interaction.reply({ content: `Please enter a valid link`, ephemeral: true });
				return;
			}
			const guildid = interaction.guild.id
			const guild = interaction.guild
			requestvals.link = link
			requestvals.name = name
			requestvals.is2Frame = gif_conv
			requestvals.guildid = guildid
			let data = await sendEmoteRequest(requestvals)
			const path = "./downloads/" + data.emotes[0].guildId;
			fs.access(path, (error) => {
				if (error) {
					fs.mkdir(path, { recursive: true }, (error) => {
						if (error) {
							console.log(error);
						} else {
						}
					});
				} else {
				}
			});
			const emotePath = await downloadEmote(data)
			let emoteName = data.emotes[0].filename.replace(/\.[^/.]+$/, "");
			if (requestvals.emoteName === "") {
				emoteName = requestvals.emoteName;
			} else {
				if (emoteName.length > 32) {
					emoteName = emoteName.substring(0, 32);
				}
				else if (emoteName.length < 2) {
					emoteName = "_" + emoteName;
				}
			}
			await guild.emojis.fetch();
			const animatedCount = guild.emojis.cache.filter(emoji => emoji.animated).size;
			const staticCount = guild.emojis.cache.filter(emoji => !emoji.animated).size;
			let totalSlots;
			switch (guild.premiumTier) {
				case "TIER_1":
					totalSlots = 200;
					break;
				case "TIER_2":
					totalSlots = 350;
					break;
				case "TIER_3":
					totalSlots = 500;
					break;
				default:
					totalSlots = 100;
			}
			const maxPerType = totalSlots / 2;
			const availableAnimatedSlots = maxPerType - animatedCount;
			const availableStaticSlots = maxPerType - staticCount;
			if (emotePath.endsWith(".gif") && availableAnimatedSlots === 0) {
				await interaction.reply({ content: `No animated emote slots available`, ephemeral: true });
				return;
			} else if (emotePath.endsWith(".png") && availableStaticSlots === 0) {
				await interaction.reply({ content: `No static emote slots available`, ephemeral: true });
				return;
			}

			guild.emojis.create({ attachment: emotePath, name: emoteName })
				.then(emoji => console.log(`Created new emoji with name ${emoji.name}!`))
				.catch(console.error);
			await interaction.reply({ content: `Sucessfully added emote: ${emoteName}`, ephemeral: true });
		}
	}
});

async function isEmoteNotFound(link) {
	var link_to_test = "https://7tv.io/v3/emotes/" + link.substr(-26);
	try {
		const response = await fetch(link_to_test);
		const data = await response.json();
		return data.error_code === 12000;
	} catch (error) {
		return false;
	}
}

async function downloadEmote(data) {
	try {
		const downloadUrl = `${apiDownloadUrl}/${data.emotes[0].guildId}/${data.emotes[0].filename}`;
		const response = await fetch(downloadUrl);

		if (!response.ok) {
			throw new Error('Download failed');
		}

		const buffer = await response.arrayBuffer();

		const filePath = path.join('./downloads', data.emotes[0].guildId, data.emotes[0].filename);

		await fs.promises.writeFile(filePath, Buffer.from(buffer));

		return filePath;
	} catch (error) {
		console.error('Error downloading emote:', error);
		throw error;
	}
}

async function sendEmoteRequest(requestvals) {
	const payload = {
		emotes: [
			{
				link: requestvals.link,
				is_2_frame_gif: requestvals.is2Frame,
				desired_name: requestvals.name,
				guild_id: requestvals.guildid
			}
		]
	};
	const requestOptions = {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(payload),
	};
	try {
		const response = await fetch(apiUrl, requestOptions);
		if (!response.ok) {
			throw new Error('Network response was not ok');
		}
		const data = await response.json();
		return data;
	} catch (error) {
		console.error('Error:', error);
		throw error;
	}

}
client.login(process.env.TOKEN);
