// SPDX-License-Identifier: MPL-2.0
"use client";

import type { ReactNode } from "react";
import { Link, usePathname } from "@/i18n/navigation";

const classes = "text-ink no-underline hover:underline aria-[current=page]:underline";

export function FooterLink({ href, children }: { href: string; children: ReactNode }) {
  const pathname = usePathname();
  return (
    <Link href={href} aria-current={pathname === href ? "page" : undefined} className={classes}>
      {children}
    </Link>
  );
}

export function ExternalFooterLink({ href, children }: { href: string; children: ReactNode }) {
  return (
    <a href={href} rel="noopener noreferrer" className={classes}>
      {children}
    </a>
  );
}
