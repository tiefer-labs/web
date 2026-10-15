// SPDX-License-Identifier: MPL-2.0
"use client";

import { useEffect, useRef, useState } from "react";
import { buttonClass } from "@/components/ui/button";
import { Link } from "@/i18n/navigation";
import { cn } from "@/lib/cn";

type NavProps = {
  dark: boolean;
  label: string;
  openLabel: string;
  closeLabel: string;
  items: { href: string; label: string }[];
  cta: { href: string; label: string };
};

/** Inline links on desktop; below `lg` a disclosure menu that Escape closes. */
export function SiteNav({ dark, label, openLabel, closeLabel, items, cta }: NavProps) {
  const [open, setOpen] = useState(false);
  const toggle = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    if (!open) return undefined;
    function onKeyDown(event: KeyboardEvent) {
      if (event.key !== "Escape") return;
      setOpen(false);
      toggle.current?.focus();
    }
    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, [open]);

  const close = () => setOpen(false);

  return (
    <nav aria-label={label}>
      <button
        ref={toggle}
        type="button"
        aria-expanded={open}
        aria-controls="nav-list"
        onClick={() => setOpen((value) => !value)}
        className="hidden cursor-pointer items-center gap-3 border-0 bg-transparent py-2.5 font-text text-small leading-none font-medium text-inherit max-lg:inline-flex"
      >
        <span>{open ? closeLabel : openLabel}</span>
        <svg viewBox="0 0 20 12" width="20" height="12" aria-hidden="true" focusable="false">
          <path d="M0 1h20M0 6h20M0 11h20" className="fill-none stroke-current" strokeWidth="1.5" />
        </svg>
      </button>
      <ul
        id="nav-list"
        role="list"
        className={cn(
          "flex items-center gap-7 text-small font-medium",
          "max-lg:absolute max-lg:inset-x-0 max-lg:top-18 max-lg:flex-col max-lg:items-stretch max-lg:gap-0 max-lg:border-b max-lg:px-(--pad) max-lg:pt-2 max-lg:pb-6 max-lg:text-[1.0625rem]",
          dark
            ? "max-lg:border-line-on-ink max-lg:bg-ink max-lg:text-on-ink max-lg:*:border-line-on-ink"
            : "max-lg:border-line max-lg:bg-paper max-lg:text-ink max-lg:*:border-line",
          !open && "max-lg:hidden",
        )}
      >
        {items.map((item) => (
          <li key={item.href} className="max-lg:border-t">
            <Link
              href={item.href}
              onClick={close}
              className="no-underline hover:underline hover:decoration-1 max-lg:block max-lg:py-3.5"
            >
              {item.label}
            </Link>
          </li>
        ))}
        <li className="max-lg:pt-4">
          <Link
            href={cta.href}
            onClick={close}
            className={buttonClass({ variant: "outline", tone: dark ? "dark" : "light", size: "sm" })}
          >
            {cta.label}
          </Link>
        </li>
      </ul>
    </nav>
  );
}
