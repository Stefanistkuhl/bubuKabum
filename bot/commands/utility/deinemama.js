const { ActionRowBuilder, ButtonBuilder, ButtonStyle, SlashCommandBuilder } = require('discord.js');

module.exports = {
	data: new SlashCommandBuilder()
		.setName('deine')
		.setDescription('dei mama'),
	async execute(interaction) {
		const confirm = new ButtonBuilder()
			.setCustomId('close')
			.setLabel('close')
			.setStyle(ButtonStyle.Success);

		const addMore = new ButtonBuilder()
			.setCustomId('addMore')
			.setLabel('Add more emotes')
			.setStyle(ButtonStyle.Primary);

		const row = new ActionRowBuilder()
			.addComponents(addMore, confirm);

		await interaction.reply({
			content: `erm`,
			components: [row],
		});
	},
};
