const { ActionRowBuilder, ButtonBuilder, ButtonStyle, SlashCommandBuilder } = require('discord.js');

module.exports = {
	data: new SlashCommandBuilder()
		.setName('deine')
		.setDescription('dei mama'),
	async execute(interaction) {
		const close = new ButtonBuilder()
			.setCustomId('close')
			.setLabel('Close')
			.setStyle(ButtonStyle.Danger);
		const submit = new ButtonBuilder()
			.setCustomId('submit')
			.setLabel('Submit')
			.setStyle(ButtonStyle.Success);

		const addMore = new ButtonBuilder()
			.setCustomId('addMore')
			.setLabel('Add more emotes')
			.setStyle(ButtonStyle.Success);

		const row = new ActionRowBuilder()
			.addComponents(addMore, close);

		await interaction.reply({
			content: `erm`,
			components: [row],
			ephemeral: true,
		});
	},
};
