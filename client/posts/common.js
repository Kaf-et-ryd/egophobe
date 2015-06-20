/*
 * Common methods to both OP and regular posts
 */

let main = require('../main'),
	imager = require('./imager'),
	{$, _, common, lang, oneeSama, options, state} = main;

module.exports = {
	initCommon() {
		this.$blockquote = this.$el.children('blockquote');
		this.listenTo(this.model, 'dispatch', this.redirect)
			.fun();
		return this;
	},
	// Extra initialisation logic for posts renderred client-side
	clientInit() {
		if (options.get('anonymise'))
			this.anonymise();
		return this;
	},
	// Proxy to the appropriate method
	redirect(command, ...args) {
		this[command](...args);
	},
	updateBody(update) {
		oneeSama.dice = update.dice;
		oneeSama.links = update.links;
		oneeSama.callback = this.inject;
		oneeSama.$buffer = this.$blockquote;
		oneeSama.state = update.state;
		oneeSama.fragment(update.frag);
	},
	// Inject various tags into the blockqoute
	inject(frag) {
		var $dest = this.$buffer;
		for (var i = 0; i < this.state[1]; i++)
			$dest = $dest.children('del').last();
		if (this.state[0] == common.S_QUOTE)
			$dest = $dest.children('em').last();
		if (this.strong)
			$dest = $dest.children('strong').last();
		var out = null;
		if (frag.safe) {
			var m = frag.safe.match(/^<(\w+)>$/);
			if (m)
				out = document.createElement(m[1]);
			else if (/^<\/\w+>$/.test(frag.safe))
				out = '';
		}
		if (out === null) {
			if (Array.isArray(frag))
				out = $(common.flatten(frag).join(''));
			else
				out = common.escape_fragment(frag);
		}
		if (out)
			$dest.append(out);
		return out;
	},
	renderTime() {
		let el = this.el.getElementsByTagName('time')[0];
		el.outerHTML = oneeSama.time(this.model.get('time'));
	},
	renderBacklinks(links) {
		let el = this.el.getElementsByTagName('small')[0];
		main.follow(() => el.innerHTML = oneeSama.backlinks(links));
	},
	// Admin JS injections
	fun() {
		// Fun goes here
	},
	// Self-delusion tripfag filter
	anonymise() {
		this.el
			.getElementsByClassName('name')[0]
			.innerHTML = `<b class="name">${lang.anon}<b>`;
	},
	// Restore regular name
	renderName() {
		this.el
			.getElementsByClassName('name')[0]
			.outerHTML = oneeSama.name(this.model.attributes);
	}
};

_.extend(module.exports, imager.Hidamari);
