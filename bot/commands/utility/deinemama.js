const { SlashCommandBuilder } = require('discord.js');

module.exports = {
	data: new SlashCommandBuilder()
		.setName('deine')
		.setDescription('dei mama'),
	async execute(interaction) {
		await interaction.reply('mama');
	},
};
