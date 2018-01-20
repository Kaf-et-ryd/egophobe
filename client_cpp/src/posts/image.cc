#include "../../brunhild/events.hh"
#include "../../brunhild/mutations.hh"
#include "../../brunhild/util.hh"
#include "../lang.hh"
#include "../options/options.hh"
#include "../state.hh"
#include "../util.hh"
#include "etc.hh"
#include "models.hh"
#include <emscripten.h>
#include <emscripten/bind.h>
#include <iomanip>
#include <sstream>

using brunhild::escape;
using std::optional;
using std::ostringstream;
using std::string;

// TODO: Expand all images and automatic image expansion

Node Post::render_figcaption()
{
    auto& img = *image;

    Node n = { "figcaption", { { "class", "spaced" } } };
    n.children.reserve(4);

    if (options->hide_thumbs || options->work_mode_toggle) {
        n.children.push_back({
            "a",
            {
                { "class", "image-toggle act" },
                { "data-id", std::to_string(id) },
            },
            lang->posts.at(img.reveal_thumbnail ? "hide" : "show"),
        });
    }
    n.children.push_back(render_file_info());

    // File name + download link
    auto& ext = file_extentions.at(img.file_type);
    ostringstream name, url;
    name << escape(img.name) << '.' << ext;
    url << "/assets/images/src/" << img.SHA1 << '.' << ext;
    n.children.push_back({ "a",
        { { "href", url.str() }, { "download", name.str() } }, name.str() });

    n.stringify_subtree();
    return n;
}

Node Post::render_file_info()
{
    using std::setw;

    auto& img = *image;
    ostringstream s;
    bool first = true;
    s << '(';

// Appends a comma and a space after the first invocation
#define COMMA                                                                  \
    if (!first) {                                                              \
        s << ", ";                                                             \
    } else {                                                                   \
        first = false;                                                         \
    }

    if (img.artist) {
        COMMA
        s << escape(*img.artist);
    }
    if (img.title) {
        COMMA
        s << escape(*img.title);
    }
    if (img.audio) {
        COMMA
        s << "♫";
    }
    if (img.length) {
        COMMA
        if (img.length < 60) {
            s << "0:" << setw(2) << img.length;
        } else {
            s << setw(2) << img.length / 60 << ':' << setw(2)
              << img.length % 60;
        }
    }
    if (img.apng) {
        COMMA
        s << "APNG";
    }

    // Readable file size
    COMMA
    if (img.size < 1 << 10) {
        s << img.size << " B";
    } else if (img.size < 1 << 20) {
        s << img.size / (1 << 10) << " KB";
    } else {
        s << std::setprecision(1) << std::fixed << (float)img.size / (1 << 20)
          << " MB";
    }

    // Media dimensions
    if (const auto[w, h, _, __] = img.dims; w && h) {
        COMMA
        s << w << 'x' << h;
    }

    s << ')';
    return Node("span", s.str());
}

// Render unexpanded file thumbnail image
static Node render_thumbnail(const Image& img)
{
    string thumb;
    uint16_t h, w;

    if (img.thumb_type == FileType::no_file) {
        // No thumbnail exists. Assign default.
        string file;
        switch (img.file_type) {
        case FileType::mp4:
        case FileType::mp3:
        case FileType::ogg:
        case FileType::flac:
            file = "audio";
            break;
        default:
            file = "file";
        }
        thumb = "/assets/" + file + ".png";
        h = w = 150;
    } else if (img.spoiler) {
        thumb = "/assets/spoil/default.jpg";
        h = w = 150;
    } else {
        thumb = img.thumb_path();
        w = img.dims[2];
        h = img.dims[3];
    }

    return {
        "img",
        {
            { "src", thumb }, { "width", std::to_string(w) },
            { "height", std::to_string(h) },
        },
    };
}

// Render expanded file image, video or audio
static void render_expanded(
    const Image& img, Node& inner, optional<Node>& audio)
{
    const auto src = img.source_path();

    switch (img.file_type) {
    case FileType::ogg:
    case FileType::mp4:
        // Can have only audio
        if (img.video) {
            goto render_video;
        }
    case FileType::flac:
    case FileType::mp3:
        // Audio controls are rendered outside the figure. Keep the
        // thumbnail.
        audio = {
            {
                "audio",
                {
                    { "autoplay", "" }, { "controls", "" }, { "loop`", "" },
                    { "src", src },
                },
            },
        };
        inner = render_thumbnail(img);
        return;
    case FileType::webm:
    render_video:
        inner = {
            "video",
            {
                { "autoplay", "" }, { "controls", "" }, { "loop`", "" },
            },
        };
        break;
    default:
        inner = { "img" };
    }

    inner.attrs["class"] = options->inline_fit == Options::FittingMode::width
        ? "fit-to-width"
        : "fit-to-screen";
    inner.attrs["src"] = src;
}

std::tuple<Node, optional<Node>> Post::render_image()
{
    auto& img = *image;
    Node inner;
    optional<Node> audio;

    if (img.expanded) {
        render_expanded(img, inner, audio);
    } else {
        inner = render_thumbnail(img);
    }

    const string id_str = std::to_string(id);
    inner.attrs["data-id"] = id_str;
    Node n({
        "figure", {},
        {
            {
                "a",
                {
                    { "href", img.source_path() }, { "target", "_blank" },
                    { "data-id", id_str },
                },
                { inner },
            },
        },
    });
    n.stringify_subtree();
    return { n, audio };
}

void handle_image_click(const brunhild::EventTarget& target)
{
    // Identify and validate parent post
    if (page->catalog) {
        return;
    }
    auto p = match_post(target.attrs);
    if (!p || !p->image) {
        return;
    }
    auto& img = *p->image;

    // Simply download the file
    switch (img.file_type) {
    case FileType::pdf:
    case FileType::zip:
    case FileType::_7z:
    case FileType::targz:
    case FileType::tarxz:
    case FileType::txt:
        EM_ASM_INT(
            {
                if (!document.querySelector) {
                    // Really old browser. Fuck it!
                    return;
                }
                document.getElementById('p' + $0.toString())
                    .querySelector('figcaption a[download]')
                    .click();
            },
            p->id);
        return;
    }

    img.expanded = !img.expanded;
    if (options->inline_fit == Options::FittingMode::width
        && img.dims[1] > emscripten::val::global("window")["innerHeight"]
                             .as<unsigned>()) {
        brunhild::scroll_into_view('p' + target.attrs.at("data-id"));
    }
    p->patch();
}

void toggle_hidden_thumbnail(const brunhild::EventTarget& target)
{
    auto p = match_post(target.attrs);
    if (!p || !p->image) {
        return;
    }
    p->image->reveal_thumbnail = !p->image->reveal_thumbnail;
    p->patch();
}
