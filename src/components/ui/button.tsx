// SPDX-License-Identifier: MPL-2.0
import type { ComponentProps } from "react";
import { Link } from "@/i18n/navigation";
import { cn } from "@/lib/cn";

type ButtonStyle = {
  variant?: "primary" | "outline";
  /** `dark` for buttons on the ink background. */
  tone?: "light" | "dark";
  size?: "md" | "sm";
};

const variants = {
  primary: {
    light: "border-ink bg-ink text-on-ink hover:border-ink-80 hover:bg-ink-80",
    dark: "border-on-ink bg-on-ink text-ink hover:border-on-ink-muted hover:bg-on-ink-muted",
  },
  outline: {
    light: "border-ink bg-transparent text-ink hover:bg-ink hover:text-on-ink",
    dark: "border-on-ink bg-transparent text-on-ink hover:bg-on-ink hover:text-ink",
  },
};

const sizes = {
  md: "px-[22px] py-[13px] text-base",
  sm: "px-4 py-[9px] text-small",
};

export function buttonClass({ variant = "primary", tone = "light", size = "md" }: ButtonStyle = {}) {
  return cn(
    "inline-block cursor-pointer rounded-[2px] border font-text leading-[1.25] font-medium no-underline transition-colors duration-150 forced-colors:border-[ButtonText]",
    variants[variant][tone],
    sizes[size],
  );
}

type ButtonLinkProps = ButtonStyle & Omit<ComponentProps<"a">, "href"> & { href: string };

/** Internal paths go through the locale-aware Link; mailto and other schemes stay plain anchors. */
export function ButtonLink({ href, variant, tone, size, className, children, ...props }: ButtonLinkProps) {
  const classes = cn(buttonClass({ variant, tone, size }), className);
  if (href.startsWith("/")) {
    return (
      <Link href={href} className={classes} {...props}>
        {children}
      </Link>
    );
  }
  return (
    <a href={href} className={classes} {...props}>
      {children}
    </a>
  );
}
