// This file is automatically generated by qtc from "article.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line article.qtpl:1
package templates

//line article.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line article.qtpl:1
import "fmt"

//line article.qtpl:2
import "strconv"

//line article.qtpl:3
import "github.com/bakape/meguca/common"

//line article.qtpl:4
import "github.com/bakape/meguca/lang"

//line article.qtpl:5
import "github.com/bakape/meguca/imager/assets"

//line article.qtpl:7
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line article.qtpl:7
func streamrenderArticle(qw422016 *qt422016.Writer, p common.Post, op uint64, omit, imageOmit int, subject, root string) {
	//line article.qtpl:8
	id := strconv.FormatUint(p.ID, 10)

	//line article.qtpl:8
	qw422016.N().S(`<article id="p`)
	//line article.qtpl:9
	qw422016.N().S(id)
	//line article.qtpl:9
	qw422016.N().S(`" class="glass`)
	//line article.qtpl:9
	if p.Editing {
		//line article.qtpl:9
		qw422016.N().S(` `)
		//line article.qtpl:9
		qw422016.N().S(`editing`)
		//line article.qtpl:9
	}
	//line article.qtpl:9
	qw422016.N().S(`"><header class="spaced"><input type="checkbox" class="mod-checkbox hidden">`)
	//line article.qtpl:12
	if subject != "" {
		//line article.qtpl:12
		qw422016.N().S(`<h3>「`)
		//line article.qtpl:14
		qw422016.E().S(subject)
		//line article.qtpl:14
		qw422016.N().S(`」</h3>`)
		//line article.qtpl:16
	}
	//line article.qtpl:16
	qw422016.N().S(`<b class="name`)
	//line article.qtpl:17
	if p.Auth != "" {
		//line article.qtpl:17
		qw422016.N().S(` `)
		//line article.qtpl:17
		qw422016.N().S(`admin`)
		//line article.qtpl:17
	}
	//line article.qtpl:17
	qw422016.N().S(`">`)
	//line article.qtpl:18
	if p.Name != "" || p.Trip == "" {
		//line article.qtpl:19
		if p.Name != "" {
			//line article.qtpl:20
			qw422016.E().S(p.Name)
			//line article.qtpl:21
		} else {
			//line article.qtpl:21
			qw422016.N().S(`Anonymous`)
			//line article.qtpl:23
		}
		//line article.qtpl:24
		if p.Trip != "" {
			//line article.qtpl:25
			qw422016.N().S(` `)
			//line article.qtpl:26
		}
		//line article.qtpl:27
	}
	//line article.qtpl:28
	if p.Trip != "" {
		//line article.qtpl:28
		qw422016.N().S(`<code>!`)
		//line article.qtpl:30
		qw422016.E().S(p.Trip)
		//line article.qtpl:30
		qw422016.N().S(`</code>`)
		//line article.qtpl:32
	}
	//line article.qtpl:33
	if p.Auth != "" {
		//line article.qtpl:34
		qw422016.N().S(` `)
		//line article.qtpl:34
		qw422016.N().S(`##`)
		//line article.qtpl:34
		qw422016.N().S(` `)
		//line article.qtpl:34
		qw422016.N().S(lang.Packs["en_GB"].Common.Posts[p.Auth])
		//line article.qtpl:35
	}
	//line article.qtpl:35
	qw422016.N().S(`</b><time>`)
	//line article.qtpl:38
	qw422016.N().S(formatTime(p.Time))
	//line article.qtpl:38
	qw422016.N().S(`</time><nav><a href="#p`)
	//line article.qtpl:41
	qw422016.N().S(id)
	//line article.qtpl:41
	qw422016.N().S(`">No.</a><a class="quote">`)
	//line article.qtpl:45
	qw422016.N().S(id)
	//line article.qtpl:45
	qw422016.N().S(`</a></nav><a class="control"><svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" viewBox="0 0 8 8"><path d="M1.5 0l-1.5 1.5 4 4 4-4-1.5-1.5-2.5 2.5-2.5-2.5z" transform="translate(0 1)" /></svg></a></header>`)
	//line article.qtpl:54
	var src string

	//line article.qtpl:55
	if p.Image != nil {
		//line article.qtpl:56
		img := *p.Image

		//line article.qtpl:57
		src = assets.SourcePath(img.FileType, img.SHA1)

		//line article.qtpl:58
		ISSrc := root + assets.RelativeSourcePath(img.FileType, img.SHA1)

		//line article.qtpl:58
		qw422016.N().S(`<figcaption class="spaced"><a class="image-toggle act" hidden></a><span class="spaced image-search-container"><a class="image-search google" target="_blank" rel="nofollow" href="https://www.google.com/searchbyimage?image_url=`)
		//line article.qtpl:62
		qw422016.N().S(ISSrc)
		//line article.qtpl:62
		qw422016.N().S(`">G</a><a class="image-search iqdb" target="_blank" rel="nofollow" href="http://iqdb.org/?url=`)
		//line article.qtpl:65
		qw422016.N().S(ISSrc)
		//line article.qtpl:65
		qw422016.N().S(`">Iq</a><a class="image-search saucenao" target="_blank" rel="nofollow" href="http://saucenao.com/search.php?db=999&url=`)
		//line article.qtpl:68
		qw422016.N().S(ISSrc)
		//line article.qtpl:68
		qw422016.N().S(`">Sn</a><a class="image-search whatAnime" target="_blank" rel="nofollow" href="https://whatanime.ga/?url=`)
		//line article.qtpl:71
		qw422016.N().S(ISSrc)
		//line article.qtpl:71
		qw422016.N().S(`">WA</a><a class="image-search desustorage" target="_blank" rel="nofollow" href="https://desuarchive.org/_/search/image/`)
		//line article.qtpl:74
		qw422016.N().S(img.MD5)
		//line article.qtpl:74
		qw422016.N().S(`">Ds</a><a class="image-search exhentai" target="_blank" rel="nofollow" href="http://exhentai.org/?fs_similar=1&fs_exp=1&f_shash=`)
		//line article.qtpl:77
		qw422016.N().S(img.SHA1)
		//line article.qtpl:77
		qw422016.N().S(`">Ex</a></span><span>(`)
		//line article.qtpl:83
		if img.Audio {
			//line article.qtpl:83
			qw422016.N().S(`♫,`)
			//line article.qtpl:84
			qw422016.N().S(` `)
			//line article.qtpl:85
		}
		//line article.qtpl:86
		if img.Length != 0 {
			//line article.qtpl:87
			l := img.Length

			//line article.qtpl:88
			if l < 60 {
				//line article.qtpl:89
				qw422016.N().S(fmt.Sprintf("0:%02d", l))
				//line article.qtpl:90
			} else {
				//line article.qtpl:91
				min := l / 6

				//line article.qtpl:92
				qw422016.N().S(fmt.Sprintf("%02d:%02d", min, l-min))
				//line article.qtpl:93
			}
			//line article.qtpl:93
			qw422016.N().S(`,`)
			//line article.qtpl:94
			qw422016.N().S(` `)
			//line article.qtpl:95
		}
		//line article.qtpl:96
		if img.APNG {
			//line article.qtpl:96
			qw422016.N().S(`APNG,`)
			//line article.qtpl:97
			qw422016.N().S(` `)
			//line article.qtpl:98
		}
		//line article.qtpl:99
		qw422016.N().S(readableFileSize(img.Size))
		//line article.qtpl:99
		qw422016.N().S(`,`)
		//line article.qtpl:99
		qw422016.N().S(` `)
		//line article.qtpl:100
		qw422016.N().S(strconv.FormatUint(uint64(img.Dims[0]), 10))
		//line article.qtpl:100
		qw422016.N().S(`x`)
		//line article.qtpl:102
		qw422016.N().S(strconv.FormatUint(uint64(img.Dims[1]), 10))
		//line article.qtpl:102
		qw422016.N().S(`)</span>`)
		//line article.qtpl:105
		name := imageName(img.FileType, img.Name)

		//line article.qtpl:105
		qw422016.N().S(`<a href="`)
		//line article.qtpl:106
		qw422016.N().S(src)
		//line article.qtpl:106
		qw422016.N().S(`" download="`)
		//line article.qtpl:106
		qw422016.N().S(name)
		//line article.qtpl:106
		qw422016.N().S(`">`)
		//line article.qtpl:107
		qw422016.N().S(name)
		//line article.qtpl:107
		qw422016.N().S(`</a></figcaption>`)
		//line article.qtpl:110
	}
	//line article.qtpl:110
	qw422016.N().S(`<div class="post-container">`)
	//line article.qtpl:112
	if p.Image != nil {
		//line article.qtpl:113
		img := *p.Image

		//line article.qtpl:113
		qw422016.N().S(`<figure><a target="_blank" href="`)
		//line article.qtpl:115
		qw422016.N().S(src)
		//line article.qtpl:115
		qw422016.N().S(`">`)
		//line article.qtpl:116
		if img.Spoiler {
			//line article.qtpl:116
			qw422016.N().S(`<!-- TODO: board-specific server-side spoiler rendering --><img src="/assets/spoil/default.jpg" width="125" height="125">`)
			//line article.qtpl:119
		} else {
			//line article.qtpl:120
			w, h := correctDims(subject != "", img.Dims[2], img.Dims[3])

			//line article.qtpl:120
			qw422016.N().S(`<img src="`)
			//line article.qtpl:121
			qw422016.N().S(assets.ThumbPath(img.ThumbType, img.SHA1))
			//line article.qtpl:121
			qw422016.N().S(`" width="`)
			//line article.qtpl:121
			qw422016.N().S(w)
			//line article.qtpl:121
			qw422016.N().S(`" height="`)
			//line article.qtpl:121
			qw422016.N().S(h)
			//line article.qtpl:121
			qw422016.N().S(`">`)
			//line article.qtpl:122
		}
		//line article.qtpl:122
		qw422016.N().S(`</a></figure>`)
		//line article.qtpl:125
	}
	//line article.qtpl:125
	qw422016.N().S(`<blockquote>`)
	//line article.qtpl:127
	streambody(qw422016, p, op)
	//line article.qtpl:127
	qw422016.N().S(`</blockquote>`)
	//line article.qtpl:129
	if p.Banned {
		//line article.qtpl:129
		qw422016.N().S(`<b class="admin banned">`)
		//line article.qtpl:131
		qw422016.N().S(lang.Packs["en_GB"].Common.Posts["banned"])
		//line article.qtpl:131
		qw422016.N().S(`</b>`)
		//line article.qtpl:133
	}
	//line article.qtpl:133
	qw422016.N().S(`</div>`)
	//line article.qtpl:135
	if omit != 0 {
		//line article.qtpl:135
		qw422016.N().S(`<span class="omit" data-omit="`)
		//line article.qtpl:136
		qw422016.N().D(omit)
		//line article.qtpl:136
		qw422016.N().S(`" data-image-omit="`)
		//line article.qtpl:136
		qw422016.N().D(imageOmit)
		//line article.qtpl:136
		qw422016.N().S(`">`)
		//line article.qtpl:137
		qw422016.N().D(omit)
		//line article.qtpl:137
		qw422016.N().S(` `)
		//line article.qtpl:137
		qw422016.N().S(`post`)
		//line article.qtpl:137
		if omit > 1 {
			//line article.qtpl:137
			qw422016.N().S(`s`)
			//line article.qtpl:137
		}
		//line article.qtpl:138
		qw422016.N().S(` `)
		//line article.qtpl:138
		qw422016.N().S(`and`)
		//line article.qtpl:138
		qw422016.N().S(` `)
		//line article.qtpl:138
		qw422016.N().D(imageOmit)
		//line article.qtpl:139
		qw422016.N().S(` `)
		//line article.qtpl:139
		qw422016.N().S(`image`)
		//line article.qtpl:139
		if imageOmit > 1 {
			//line article.qtpl:139
			qw422016.N().S(`s`)
			//line article.qtpl:139
		}
		//line article.qtpl:140
		qw422016.N().S(` `)
		//line article.qtpl:140
		qw422016.N().S(`omitted`)
		//line article.qtpl:140
		qw422016.N().S(` `)
		//line article.qtpl:140
		qw422016.N().S(`<span class="act"><a href="`)
		//line article.qtpl:142
		qw422016.N().S(strconv.FormatUint(op, 10))
		//line article.qtpl:142
		qw422016.N().S(`" class="history">See All</a></span></span>`)
		//line article.qtpl:147
	}
	//line article.qtpl:148
	if p.Backlinks != nil {
		//line article.qtpl:148
		qw422016.N().S(`<span class="backlinks spaced">`)
		//line article.qtpl:150
		for id, link := range p.Backlinks {
			//line article.qtpl:150
			qw422016.N().S(`<em>`)
			//line article.qtpl:152
			streampostLink(qw422016, id, link.OP, link.Board, link.OP != op)
			//line article.qtpl:152
			qw422016.N().S(`</em>`)
			//line article.qtpl:154
		}
		//line article.qtpl:154
		qw422016.N().S(`</span>`)
		//line article.qtpl:156
	}
	//line article.qtpl:156
	qw422016.N().S(`</article>`)
//line article.qtpl:158
}

//line article.qtpl:158
func writerenderArticle(qq422016 qtio422016.Writer, p common.Post, op uint64, omit, imageOmit int, subject, root string) {
	//line article.qtpl:158
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line article.qtpl:158
	streamrenderArticle(qw422016, p, op, omit, imageOmit, subject, root)
	//line article.qtpl:158
	qt422016.ReleaseWriter(qw422016)
//line article.qtpl:158
}

//line article.qtpl:158
func renderArticle(p common.Post, op uint64, omit, imageOmit int, subject, root string) string {
	//line article.qtpl:158
	qb422016 := qt422016.AcquireByteBuffer()
	//line article.qtpl:158
	writerenderArticle(qb422016, p, op, omit, imageOmit, subject, root)
	//line article.qtpl:158
	qs422016 := string(qb422016.B)
	//line article.qtpl:158
	qt422016.ReleaseByteBuffer(qb422016)
	//line article.qtpl:158
	return qs422016
//line article.qtpl:158
}
