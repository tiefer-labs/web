# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
"""Build the WOFF2 web fonts from the variable TTF sources.

Run from the repository root: python3 scripts/fonts.py
Needs fontTools and brotli (pip install fonttools brotli). This lives in a
script because Go has no font tooling; its output is checked by the Go
budget tests.

- Subset to Basic Latin, Latin-1 Supplement, Latin Extended-A and the
  punctuation the site uses.
- Keep the wght axis (200 to 700).
- Pin the wdth axis of Mozilla Headline to 100.
"""
import io

from fontTools import subset
from fontTools.ttLib import TTFont
from fontTools.varLib import instancer

FONTS = "web/static/fonts/"
UNICODES = (
    list(range(0x20, 0x7F))        # Basic Latin
    + list(range(0xA0, 0x100))     # Latin-1 Supplement
    + list(range(0x100, 0x180))    # Latin Extended-A
    + [0x2018, 0x2019, 0x201C, 0x201D, 0x2022, 0x2026, 0x2032, 0x2033, 0x20AC, 0x2122, 0x2212]
)


def build(src, dst, pin=None):
    font = TTFont(FONTS + src, lazy=False)
    if pin:
        font = instancer.instantiateVariableFont(font, pin)
        # Reload, so the subsetter sees fully decompiled tables.
        buf = io.BytesIO()
        font.save(buf)
        buf.seek(0)
        font = TTFont(buf, lazy=False)
    opts = subset.Options()
    opts.flavor = "woff2"
    opts.layout_features = ["kern", "liga", "calt", "tnum", "lnum", "case", "ccmp", "locl", "mark", "mkmk"]
    opts.name_IDs = ["*"]
    opts.name_legacy = True
    opts.name_languages = ["*"]
    opts.notdef_outline = True
    opts.hinting = False
    opts.desubroutinize = True
    s = subset.Subsetter(opts)
    s.populate(unicodes=UNICODES)
    s.subset(font)
    font.flavor = "woff2"
    font.save(FONTS + dst)


build("MozillaHeadline-VF.ttf", "MozillaHeadline.woff2", pin={"wdth": 100})
build("MozillaText-VF.ttf", "MozillaText.woff2")
