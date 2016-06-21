/*
 Renders the HTML of the options panel
*/

import {
	filter, extend, groupBy, HTML, makeAttrs, ElementAttributes
} from '../util'
import {opts as lang, OptLabel} from '../lang'
import {specs, OptionSpec, optionType} from './specs'

// Render the inner HTML of the options panel
export default function (): string {
	let html = '<ul class="option_tab_sel">'
	const {tabs} = lang,
		byTab = groupBy(specs, 'tab'),
		opts: {[key: number]: OptionSpec[]} = []

	// Render tab butts
	for (let i = 0; i < tabs.length; i++) {
		// No options in this tab
		if (!byTab[i]) {
			continue
		}

		// Pick the options for this specific tab, according to current
		// template and server configuration
		opts[i] = filter(byTab[i], spec =>
			!spec.noLoad && !spec.hidden)

		// All options disaled
		if (!opts[i].length) {
			continue
		}

		const attrs: ElementAttributes = {
			'data-content': `tab-${i}`,
			class: 'tab_link'
		}

		// Highlight the first tabButt by default
		if (i === 0) {
			attrs['class'] += ' tab_sel'
		}
		html += HTML
			`<li>
				<a ${makeAttrs(attrs)}>
					${tabs[i]}
				</a>
			</li>`
	}

	html += '</ul><ul class="option_tab_cont">'
	for (let tabNumber in opts) {
		html += renderTab(opts[tabNumber], parseInt(tabNumber))
	}
	html += '</ul>'

	return html
}

// Render tab contents
function renderTab(opts: OptionSpec[], i: number): string {
	if (!opts.length) {
		return ''
	}
	let html = ''
	html += `<li class="tab-${i}`

	// Show the first tab by default
	if (i === 0) {
		html += ' tab_sel'
	}
	html += '">'

	// Render the actual options
	for (let opt of opts) {
		html += renderOption(opt)
	}

	if (i === 0) {
		html += renderExtras()
	}
	html += '</li>'

	return html
}

// Render a single option from it's schema
function renderOption(spec: OptionSpec): string {
	switch (spec.type) {
		case optionType.shortcut:
			return 'Alt+' + renderInput(spec.id, {
				maxlength: '1',
				class: 'shortcut'
			})
		case optionType.checkbox:
			return renderInput(spec.id, {type: 'checkbox'})
		case optionType.number:
			return renderInput(spec.id, {
				style: 'width: 4em;',
				maxlength: '4'
			})
		case optionType.image:
			return renderInput(spec.id, {type: 'file'})
		case optionType.menu:
			return renderMenu(spec)
	}
}

// Common input field render logic
function renderInput(id: string, attrs: ElementAttributes): string {
	const [label, title] = lang.labels[id]
	extend(attrs, {id, title})
	return `<input ${makeAttrs(attrs)}>` + renderLabel(id, title, label)
}

// Render the description label to the right of the option
function renderLabel(id: string, title: string, label: string): string {
	return HTML
		`<label for="${id}" title="${title}">
			${label}
		</label>
		<br>`
}

// Render drop down selection menu
function renderMenu({id, list}: OptionSpec): string {
	const [label, title] = lang.labels[id]
	let html = `<select id="${id}" title="${title}">`
	for (let item of list) {
		html += HTML
			`<option value="${item}">
				${lang.modes[item] || item}
			</option>`
	}
	html += '</select>' + renderLabel(id, title, label)
	return html
}

// Hidden post reset, Export and Import links to first tab
function renderExtras(): string {
	let html = '<br>'
	const links = ['export', 'import', 'hidden']
	for (let id of links) {
		const [label, title] = lang.labels[id]
		html += HTML
			`<a id="${id}" title="${title}">
				${label}
			</a> `
	}

	// Hidden file input for uploading the JSON
	const attrs: ElementAttributes = {
		type: 'file',
		id: 'importSettings',
		name: "Import Settings"
	}
	html += `<input ${makeAttrs(attrs)}>`

	return html
}
