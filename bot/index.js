const fs = require('node:fs');
const path = require('node:path');
const dotenv = require('dotenv');
dotenv.config();
const { MessageFlags, Client, Collection, ModalBuilder, Events, ButtonStyle, GatewayIntentBits, StringSelectMenuOptionBuilder, StringSelectMenuBuilder, TextInputBuilder, TextInputStyle, ActionRowBuilder, ButtonBuilder } = require('discord.js');

const client = new Client({
	intents: [
		GatewayIntentBits.Guilds,
		GatewayIntentBits.GuildMessages,
		GatewayIntentBits.GuildExpressions
	],
});

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
				.setLabel("Enter a 7tv link")
				.setRequired(true)
				.setStyle(TextInputStyle.Short);

			const nameInput = new TextInputBuilder()
				.setCustomId('nameInput')
				.setLabel("Desired Emote Name")
				.setStyle(TextInputStyle.Short)
				.setRequired(false);
			const is2FrameGifInput = new TextInputBuilder()
				.setCustomId('2FrameGif')
				.setLabel("make this a 2 frame gif?")
				.setStyle(TextInputStyle.Short)
				.setRequired(false);

			const firstActionRow = new ActionRowBuilder().addComponents(linkInput);
			const secondActionRow = new ActionRowBuilder().addComponents(nameInput);
			const thirdActionRow = new ActionRowBuilder().addComponents(is2FrameGifInput);

			modal.addComponents(firstActionRow, secondActionRow, thirdActionRow);

			await interaction.showModal(modal);
		} else if (interaction.customId === 'close') {
			await interaction.update({ content: 'leck eier', components: [] });
		}
	}

	// Modal submit handler
	if (interaction.isModalSubmit()) {
		if (interaction.customId === 'inputModal') {
			const link = interaction.fields.getTextInputValue('linkInput');
			const name = interaction.fields.getTextInputValue('nameInput');
			const gif = interaction.fields.getTextInputValue('2FrameGif');

			await interaction.reply({ content: `Received link: ${link}\nName: ${name} \n gif: ${gif}`, ephemeral: true });
			// Handle your emote addition logic here
		}
	}
});


client.login(process.env.TOKEN);
