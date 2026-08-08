package fileutil

import (
	"github.com/usesnipet/snipet/pkg/driver/knowledge"
)

// The lists below enumerate the MIME types (and their known aliases) from
// https://github.com/gabriel-vasile/mimetype's supported_mimes.md, grouped
// by SourceItemKind. Formats with no clear content kind (archives,
// executables, fonts, CAD/3D binaries, installers, etc.) are intentionally
// left out and fall through to SourceItemKindUnknown.

var textMediaTypes = []string{
	"text/plain",
	"text/html",
	"application/xhtml+xml",
	"text/x-php",
	"text/javascript", "application/x-javascript", "application/javascript",
	"text/x-lua",
	"text/x-perl",
	"text/x-python", "text/x-script.python", "application/x-python",
	"text/x-ruby", "application/x-ruby",
	"text/rtf", "application/rtf",
	"text/x-tcl", "application/x-tcl",
	"text/x-shellscript", "text/x-sh", "application/x-shellscript", "application/x-sh",
	"message/rfc822",
	"text/vcard",
	"text/calendar",
	"text/vnd.familysearch.gedcom",
	"application/warc",
	"text/vtt",
	"application/x-subrip", "application/x-srt", "text/x-srt",
}

var documentMediaTypes = []string{
	"application/pdf", "application/x-pdf",
	"application/vnd.fdf",
	"application/msword", "application/vnd.ms-word",
	"application/vnd.ms-excel", "application/msexcel",
	"application/vnd.ms-powerpoint", "application/mspowerpoint",
	"application/vnd.ms-publisher",
	"application/vnd.ms-outlook",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"application/vnd.oasis.opendocument.text", "application/x-vnd.oasis.opendocument.text",
	"application/vnd.oasis.opendocument.text-template", "application/x-vnd.oasis.opendocument.text-template",
	"application/vnd.oasis.opendocument.spreadsheet", "application/x-vnd.oasis.opendocument.spreadsheet",
	"application/vnd.oasis.opendocument.spreadsheet-template", "application/x-vnd.oasis.opendocument.spreadsheet-template",
	"application/vnd.oasis.opendocument.presentation", "application/x-vnd.oasis.opendocument.presentation",
	"application/vnd.oasis.opendocument.presentation-template", "application/x-vnd.oasis.opendocument.presentation-template",
	"application/vnd.oasis.opendocument.graphics", "application/x-vnd.oasis.opendocument.graphics",
	"application/vnd.oasis.opendocument.graphics-template", "application/x-vnd.oasis.opendocument.graphics-template",
	"application/vnd.oasis.opendocument.formula", "application/x-vnd.oasis.opendocument.formula",
	"application/vnd.oasis.opendocument.chart", "application/x-vnd.oasis.opendocument.chart",
	"application/vnd.sun.xml.calc",
	"application/epub+zip",
	"application/x-mobipocket-ebook",
	"application/x-ms-reader",
	"application/vnd.wordperfect",
	"application/vnd.framemaker",
	"application/onenote",
	"application/vnd.ms-htmlhelp",
	"application/vnd.ms-visio.drawing.main+xml",
	"application/postscript",
	"application/vnd.lotus-1-2-3",
}

var imageMediaTypes = []string{
	"image/x-xpixmap",
	"image/png", "image/apng", "image/vnd.mozilla.apng",
	"image/jpeg",
	"image/jxl",
	"image/jp2", "image/jpx", "image/jpm", "video/jpm", "image/jxs",
	"image/gif",
	"image/webp",
	"image/tiff",
	"image/bmp", "image/x-bmp", "image/x-ms-bmp",
	"image/x-icon",
	"image/avif",
	"image/heic", "image/heic-sequence", "image/heif", "image/heif-sequence",
	"image/svg+xml",
	"image/vnd.adobe.photoshop", "image/x-psd", "application/photoshop",
	"image/vnd.dwg", "image/x-dwg", "application/acad", "application/x-acad", "application/autocad_dwg", "application/dwg", "application/x-autocad", "drawing/dwg",
	"image/vnd.dxf",
	"image/x-xcf",
	"image/x-gimp-pat", "image/x-gimp-gbr",
	"image/bpg",
	"image/vnd.djvu",
	"image/x-icns",
	"image/vnd.radiance",
	"image/x-portable-bitmap", "image/x-portable-graymap", "image/x-portable-pixmap", "image/x-portable-arbitrarymap",
	"image/fits", "application/fits",
	"image/jxr", "image/vnd.ms-photo",
	"application/dicom",
}

var audioMediaTypes = []string{
	"audio/ogg",
	"audio/flac",
	"audio/midi", "audio/mid", "audio/sp-midi", "audio/x-mid", "audio/x-midi",
	"audio/ape",
	"audio/musepack",
	"audio/amr", "audio/amr-nb",
	"audio/wav", "audio/x-wav", "audio/vnd.wave", "audio/wave",
	"audio/aiff", "audio/x-aiff",
	"audio/basic",
	"audio/mp4", "audio/x-mp4a",
	"audio/x-m4a",
	"audio/aac",
	"audio/x-unknown",
	"application/vnd.apple.mpegurl", "audio/mpegurl", "application/x-mpegurl",
	"audio/qcelp",
	"audio/mpeg", "audio/x-mpeg", "audio/mp3",
}

var videoMediaTypes = []string{
	"video/mpeg",
	"video/quicktime",
	"video/mp4",
	"video/3gpp", "video/3gp", "audio/3gpp",
	"video/3gpp2", "video/3g2", "audio/3gpp2",
	"video/x-m4v",
	"video/mj2",
	"video/vnd.dvb.file",
	"video/webm", "audio/webm",
	"video/x-msvideo", "video/avi", "video/msvideo",
	"video/x-flv",
	"video/matroska", "video/x-matroska",
	"video/x-ms-asf", "video/asf", "video/x-ms-wmv",
	"application/vnd.rn-realmedia-vbr",
	"video/ogg",
}

var structuredMediaTypes = []string{
	"application/json",
	"application/geo+json",
	"model/gltf+json",
	"application/vnd.cyclonedx+json",
	"application/x-ndjson",
	"text/xml", "application/xml",
	"application/rss+xml", "text/rss",
	"application/atom+xml",
	"model/x3d+xml",
	"application/vnd.google-earth.kml+xml",
	"application/x-xliff+xml",
	"model/vnd.collada+xml",
	"application/gml+xml",
	"application/gpx+xml",
	"application/vnd.garmin.tcx+xml",
	"application/x-amf",
	"application/vnd.ms-package.3dmanufacturing-3dmodel+xml",
	"application/vnd.adobe.xfdf",
	"application/owl+xml",
	"application/vnd.cyclonedx+xml",
	"text/csv", "text/tab-separated-values",
	"application/cbor",
	"application/vnd.sqlite3", "application/x-sqlite3",
	"application/vnd.apache.parquet", "application/x-parquet",
	"application/x-dbf",
	"application/vnd.shx", "application/vnd.shp",
	"application/marc",
	"application/grib", "application/bufr",
	"application/x-msaccess",
	// Not detected from content by the underlying library (no distinct magic
	// bytes), kept here so extension-based fallback still classifies them.
	"application/yaml", "application/toml", "application/hcl", "application/sql",
}

var kindByMediaType = buildKindByMediaType()

func buildKindByMediaType() map[string]knowledge.SourceItemKind {
	m := make(map[string]knowledge.SourceItemKind)

	assign := func(mediaTypes []string, kind knowledge.SourceItemKind) {
		for _, mediaType := range mediaTypes {
			m[mediaType] = kind
		}
	}

	assign(textMediaTypes, knowledge.SourceItemKindText)
	assign(documentMediaTypes, knowledge.SourceItemKindDocument)
	assign(imageMediaTypes, knowledge.SourceItemKindImage)
	assign(audioMediaTypes, knowledge.SourceItemKindAudio)
	assign(videoMediaTypes, knowledge.SourceItemKindVideo)
	assign(structuredMediaTypes, knowledge.SourceItemKindStructured)

	return m
}

// MapKind classifies an already-normalized media type into a SourceItemKind.
func MapKind(mediaType string) knowledge.SourceItemKind {
	if kind, ok := kindByMediaType[mediaType]; ok {
		return kind
	}
	return knowledge.SourceItemKindUnknown
}
