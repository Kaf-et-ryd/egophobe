// Specs for individual option models

import { config } from '../state'
import { makeEl, HTML, setCookie } from "../util"
import { render as renderBG } from "./background"
import { render as renderMascot } from "./mascot"
import options from "."

// Types of option models
export const enum optionType {
	checkbox, number, image, shortcut, menu, textarea,
}

// Full schema of the option interface
export type OptionSpec = {
	// Type of option. Determines storage and rendering method. Defaults to
	// 'checkbox', if omitted.
	type?: optionType

	// Default value. false, if omitted.
	default?: any

	// Function to execute on option change
	exec?: (val?: any) => void

	// Should the function not be executed on model population?
	noExecOnStart?: boolean
}

function renderBackground(_: boolean) {
	renderBG()
}

// Specifications of option behavior, where needed. Some properties defined as
// getters to prevent race with "state" module
export const specs: { [id: string]: OptionSpec } = {
	// Thumbnail inline expansion mode
	inlineFit: {
		type: optionType.menu,
		default: "largeur",
	},
	// Hide thumbnails
	hideThumbs: {},
	// Boss key toggle
	workModeToggle: {
		type: optionType.checkbox,
		default: false,
		exec: toggleHeadStyle("work-mode", ".image-banner{display: none;}"),
	},
	// Image hover expansion
	imageHover: {
		default: true,
	},
	// WebM hover expansion
	webmHover: {},
	// Animated GIF thumbnails
	autogif: {},
	// Enable thumbnail spoilers
	spoilers: {
		default: true,
	},
	// Desktop Notifications
	notification: {
		default: true,
		exec(enabled: boolean) {
			const req = enabled
				&& typeof Notification === "function"
				&& (Notification as any).permission !== "granted"
			if (req) {
				Notification.requestPermission()
			}
		},
	},
	// Anonymise all poster names
	anonymise: {},
	// Hide posts that linked to a hidden post
	hideRecursively: {},
	// Expand post links inline
	postInlineExpand: {
		default: true,
		exec: toggleHeadStyle(
			"postInlineExpand",
			".hash-link{ display: inline; }"
		),
	},
	// Relative post timestamps
	relativeTime: {},
	// Tile posts horizontally too
	horizontalPosting: {
		exec: toggleHeadStyle(
			'horizontal',
			HTML`#thread-container {
				display:flex;
				flex-direction: row;
				flex-wrap: wrap;
				align-items: baseline;
			}`,
		)
	},
	// Move [Reply] to the right side of the screen
	replyRight: {
		exec: toggleHeadStyle(
			'reply-at-right',
			'aside.posting{margin: -26px 0 2px auto;}'
		)
	},
	// Change theme
	theme: {
		type: optionType.menu,
		get default() {
			return config.defaultCSS
		},
		noExecOnStart: true,
		exec(theme: string) {
			if (!theme) {
				return
			}
			document
				.getElementById('theme-css')
				.setAttribute('href', `/assets/css/${theme}.css`)
			setCookie("theme", theme, 365 * 10)
		},
	},
	// Custom user-set background
	userBG: {
		noExecOnStart: true,
		exec: renderBackground,
	},
	// Upload field for the custom background image
	userBGImage: {
		type: optionType.image,
	},
	// Mascot in the corner
	mascot: {
		noExecOnStart: true,
		exec: renderMascot,
	},
	mascotImage: {
		type: optionType.image,
	},
	// User-set CSS rules
	customCSSToggle: {
		noExecOnStart: true,
		exec(on: boolean) {
			let el = document
				.getElementById("custom-CSS-style") as HTMLStyleElement
			if (!el) {
				el = document.createElement("style")
				el.id = "custom-CSS-style"
				document.head.append(el)
				// The disabled property only exists on elements in the DOM,
				// so we do another query
				el = document
					.getElementById("custom-CSS-style") as HTMLStyleElement
			}
			el.innerHTML = options.customCSS
			el.disabled = !on
		},
	},
	customCSS: {
		type: optionType.textarea,
	},
	// Lock thread scrolling to bottom, when bottom in view, even when the
	// tab is hidden
	alwaysLock: {},
	// Shortcut keys
	newPost: {
		default: 78,
		type: optionType.shortcut,
	},
	done: {
		default: 83,
		type: optionType.shortcut,
	},
	toggleSpoiler: {
		default: 73,
		type: optionType.shortcut,
	},
	expandAll: {
		default: 69,
		type: optionType.shortcut,
	},
	workMode: {
		default: 66,
		type: optionType.shortcut,
	},
	galleryMode: {
		default: 71,
		type: optionType.shortcut,
	},
}

// Toggle an optional style element in the head
function toggleHeadStyle(
	name: string,
	css: string,
): (toggle: boolean) => void {
	return toggle => {
		const id = name + "-toggle"
		if (!document.getElementById(id)) {
			const html = `<style id="${id}">${css}</style>`
			document.head.append(makeEl(html))
		}

		// The disabled property only exists on elements in the DOM, so we do
		// another query
		(document.getElementById(id) as HTMLStyleElement).disabled = !toggle
	}
}
