// SPDX-License-Identifier: MPL-2.0
import { createCn } from "cn/config";

// The type scale tokens from globals.css are font sizes; without this, cn would merge them as colours.
export const cn = createCn({
  extend: {
    classGroups: {
      "font-size": [{ text: ["h1", "h2", "name", "statement", "h3", "lead", "body", "small", "label", "tag"] }],
    },
  },
});
