// SPDX-License-Identifier: MPL-2.0
import { Mozilla_Headline, Mozilla_Text } from "next/font/google";

export const headline = Mozilla_Headline({
  subsets: ["latin", "latin-ext"],
  weight: ["500", "600"],
  variable: "--font-mozilla-headline",
  display: "swap",
  adjustFontFallback: false,
});

export const text = Mozilla_Text({
  subsets: ["latin", "latin-ext"],
  weight: ["400", "500"],
  variable: "--font-mozilla-text",
  display: "swap",
  adjustFontFallback: false,
});
