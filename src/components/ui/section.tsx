// SPDX-License-Identifier: MPL-2.0
import type { ReactNode } from "react";
import { cn } from "@/lib/cn";

type Tone = "paper" | "tint" | "ink";

const tones: Record<Tone, string> = {
  paper: "",
  tint: "bg-tint",
  ink: "bg-ink text-on-ink-muted",
};

type SectionProps = {
  id?: string;
  labelledBy: string;
  tone?: Tone;
  /** Vertical padding; the default is the responsive section rhythm. */
  padding?: string;
  className?: string;
  /** Decorative layers drawn behind the content (SVG art, marks). */
  decoration?: ReactNode;
  as?: "section" | "article";
  children: ReactNode;
};

/** A page section: full-bleed background, centred 12-column grid. */
export function Section({
  id,
  labelledBy,
  tone = "paper",
  padding = "py-(--section)",
  className,
  decoration,
  as: Tag = "section",
  children,
}: SectionProps) {
  return (
    <Tag
      id={id}
      aria-labelledby={labelledBy}
      data-tone={tone === "ink" ? "dark" : undefined}
      className={cn("relative", padding, tones[tone], className)}
    >
      {decoration}
      <div className="relative wrap grid-12">{children}</div>
    </Tag>
  );
}

export function SectionLabel({ className, children }: { className?: string; children: ReactNode }) {
  return (
    <p className={cn("mb-6 pt-[0.6em] label-caps text-ink-60 in-data-[tone=dark]:text-on-ink-muted", className)}>
      {children}
    </p>
  );
}
