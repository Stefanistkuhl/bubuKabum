const { ActionRowBuilder, ButtonBuilder, ButtonStyle, SlashCommandBuilder } = require('discord.js');

module.exports = {
	data: new SlashCommandBuilder()
		.setName('emoteclone')
		.setDescription('Clone an Emote from 7tv'),
	async execute(interaction) {
		const close = new ButtonBuilder()
			.setCustomId('close')
			.setLabel('Close')
			.setStyle(ButtonStyle.Danger);
		const addMore = new ButtonBuilder()
			.setCustomId('addMore')
			.setLabel('Add emote')
			.setStyle(ButtonStyle.Success);

		const row = new ActionRowBuilder()
			.addComponents(addMore, close);

		await interaction.reply({
			content: `Emote cloning initiated. Please select an option below:`,
			components: [row],
			ephemeral: true,
		});
	},
};
