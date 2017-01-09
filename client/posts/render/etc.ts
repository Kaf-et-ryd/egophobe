// Miscellaneous post component rendering functions

import { page, mine } from '../../state'
import lang from '../../lang'
import { PostLinks } from '../../common'

// Render a link to other posts
export function renderPostLink(num: number, board: string, op: number): string {
    let html = `<a class="history post-link" data-id="${num}" href="`
    const cross = op !== page.thread

    if (cross) {
        html += `/${board}/${op}`
    }
    html += `#p${num}">>>`

    if (cross) {
        html += `>/${board}/`
    }
    html += num

    if (mine.has(num)) { // Post, I made
        html += ' ' + lang.posts["you"]
    }

    html += `</a><a class="hash-link"> #</a>`

    return html
}

// TODO: Reimplement, when moderation done

// Render USER WAS BANNED FOR THIS POST message, or similar
// export function renderBanned(): string {
// 	return `<b class="admin banMessage">${lang.mod.banMessage}</b>`
// }
//
// Render moderation information. Only exposed to authenticated staff.
// export function renderModInfo(info) {
// 	let html = '<b class="modLog admin">'
// 	for (let action of info) {
// 		html += `${lang.mod.formatLog(action)}<br>`
// 	}
// 	html += '</b>'
// 	return html
// }

// Render links to posts that are linking to the target post
export function renderBacklinks(post: DocumentFragment, links: PostLinks) {
    if (!links) {
        return
    }

    let el = post.querySelector(".backlinks")
    if (!el) {
        el = document.createElement("span")
        el.classList.add("spaced", "backlinks")
        post.append(el)
    }

    let html = ''
    for (let id in links) {
        const {board, op} = links[id]
        html += "<em>" + renderPostLink(parseInt(id), board, op) + "</em>"
    }

    el.innerHTML = html
}
